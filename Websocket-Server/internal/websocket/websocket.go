package websocket

import (
	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
	"log"
	"marketplace_websocket/internal/models"
	"marketplace_websocket/internal/service"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleConnection(hub *Hub, w http.ResponseWriter, r *http.Request, messageService *service.MessageService) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return
	}

	clientIDStr := r.URL.Query().Get("id")
	if clientIDStr == "" {
		log.Println("Client ID is missing")
		conn.Close()
		return
	}

	clientID, err := gocql.ParseUUID(clientIDStr)
	if err != nil {
		log.Println("Invalid Client ID:", err)
		conn.Close()
		return
	}

	client := &Client{
		id:   clientID,
		hub:  hub,
		conn: conn,
		send: make(chan interface{}),
	}

	hub.NotifyUserOnline(client.id)

	hub.register <- client
	go client.ReadMessages(messageService)
	go client.WriteMessages()
}

func (h *Hub) NotifyChatRoomCreation(chatRoom *models.ChatRoom) {
	h.mu.Lock()
	defer h.mu.Unlock()

	notification := models.ChatRoomNotification{
		Type:        "chatRoomCreated",
		SenderID:    chatRoom.CurrentUserID,
		RecipientID: chatRoom.TargetUserID,
		ChatRoomID:  chatRoom.ID,
	}

	if client, ok := h.clients[chatRoom.CurrentUserID]; ok {
		select {
		case client.send <- notification:
		default:
			log.Printf("Failed to send notification to sender %s, channel might be blocked or closed", chatRoom.CurrentUserID)
			close(client.send)
			delete(h.clients, chatRoom.CurrentUserID)
		}
	}
}

func (h *Hub) NotifyUserOnline(userID gocql.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	notification := Notification{
		Type:    "status",
		UserID:  userID,
		Message: "online",
	}

	for _, client := range h.clients {
		select {
		case client.send <- notification:
		default:
			log.Printf("Failed to send online notification to user %s, channel might be blocked or closed", userID)
			close(client.send) // Ensure you close the channel
			delete(h.clients, client.id)
		}
	}
}

func (h *Hub) NotifyUserOffline(userID gocql.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	notification := Notification{
		Type:    "status",
		UserID:  userID,
		Message: "offline",
	}

	for _, client := range h.clients {
		select {
		case client.send <- notification:
		default:
			log.Printf("Failed to send online notification to user %s, channel might be blocked or closed", userID)
			close(client.send)
			delete(h.clients, client.id)
		}
	}
}
