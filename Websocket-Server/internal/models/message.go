package models

import (
	"github.com/gocql/gocql"
	"time"
)

type Message struct {
	ID          gocql.UUID `json:"id"`          // Unique identifier for the message
	Content     string     `json:"content"`     // The text or content of the message
	SenderID    gocql.UUID `json:"senderID"`    // ID of the user who sent the message
	RecipientID gocql.UUID `json:"recipientID"` // ID of the user who receives the message
	ChatRoomID  gocql.UUID `json:"chatRoomID"`
	Timestamp   time.Time  `json:"timestamp"`   // Unix timestamp of when the message was sent
	MessageType string     `json:"messageType"` // Type of the message (text, image, etc.)
}

type MessageWrap struct {
	ID         gocql.UUID `json:"id"`
	Content    string     `json:"content"`
	SenderID   gocql.UUID `json:"senderID"`
	ChatRoomID gocql.UUID `json:"chatRoomID"`
	Timestamp  time.Time  `json:"timestamp"`
}

type ChatRoomNotification struct {
	Type        string     `json:"type"`
	SenderID    gocql.UUID `json:"senderID"`
	RecipientID gocql.UUID `json:"recipientID"`
	ChatRoomID  gocql.UUID `json:"chatRoomID"`
}
