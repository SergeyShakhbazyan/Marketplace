package websocket

import (
	"context"
	"encoding/json"
	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
	"log"
	"marketplace_websocket/internal/models"
	"marketplace_websocket/internal/service"
)

type Client struct {
	id   gocql.UUID
	hub  *Hub
	conn *websocket.Conn
	send chan interface{}
}

func (c *Client) ReadMessages(messageService *service.MessageService) {
	defer func() {
		c.hub.NotifyUserOffline(c.id)
		c.hub.unregister <- c
		c.conn.Close()
	}()

	//c.hub.NotifyUserOnline(c.id)

	for {
		_, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			}
			break
		}

		if len(msgBytes) == 0 {
			log.Println("Received empty message")
			continue
		}

		var message models.Message
		err = json.Unmarshal(msgBytes, &message)
		if err != nil {
			log.Println("Error unmarshaling message:", err)
			continue
		}

		message.ID = gocql.TimeUUID()

		chatRoom, err := messageService.SaveMessage(context.Background(), &message)
		if err != nil {
			log.Println("Error saving message:", err)
			continue
		}

		if chatRoom != nil {
			c.hub.NotifyChatRoomCreation(chatRoom)
		}

		if message.RecipientID != c.id {
			select {
			case c.hub.broadcast <- message:
			default:
				log.Println("Error broadcasting message: hub broadcast channel full or closed")
			}
		} else {
			select {
			case c.send <- message:
			default:
				log.Println("Error sending message to client: client send channel full or closed")
			}
		}
	}
}

func (c *Client) WriteMessages() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		err := c.conn.WriteJSON(message)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
