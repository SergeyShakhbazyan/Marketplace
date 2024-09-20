package websocket

import (
	"github.com/gocql/gocql"
	"marketplace_websocket/internal/models"
	"sync"
)

type Hub struct {
	clients    map[gocql.UUID]*Client
	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

type Notification struct {
	Type    string     `json:"type"`
	Message string     `json:"message"`
	UserID  gocql.UUID `json:"userID"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[gocql.UUID]*Client),
		broadcast:  make(chan models.Message, 1000),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		//messageHandler: messageHandler,
	}
}

func (h *Hub) IsUserConnected(userID gocql.UUID) bool {
	_, connected := h.clients[userID]
	return connected
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.id] = client
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			recipient, ok := h.clients[message.RecipientID]
			h.mu.Unlock()
			if ok {
				select {
				case recipient.send <- message:
				default:
					close(recipient.send)
					h.mu.Lock()
					delete(h.clients, recipient.id)
					h.mu.Unlock()
				}
			}
		}
	}
}
