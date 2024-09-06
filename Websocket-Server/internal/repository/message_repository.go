package repository

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"marketplace_websocket/internal/models"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, message *models.Message) error
	GetMessagesFromChatRoom(ctx context.Context, chatRoomID gocql.UUID) ([]models.MessageWrap, error)
	GetLastMessageFromChatRoom(ctx context.Context, chatRoomID gocql.UUID) (string, error)
}

type messageRepository struct {
	session *gocql.Session
}

func NewMessageRepository(session *gocql.Session) MessageRepository {
	return &messageRepository{session: session}
}

func (r *messageRepository) SaveMessage(ctx context.Context, message *models.Message) error {
	query := "INSERT INTO messenger_keyspace.message(id, chatID, content, senderid, timestamp) VALUES (?, ?, ?, ?, ?)"
	err := r.session.Query(query,
		message.ID,
		message.ChatRoomID,
		message.Content,
		message.SenderID,
		message.Timestamp,
	).WithContext(ctx).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (r *messageRepository) GetMessagesFromChatRoom(ctx context.Context, chatRoomID gocql.UUID) ([]models.MessageWrap, error) {
	var messages []models.MessageWrap
	var message models.MessageWrap
	query := "SELECT id, content, senderid, chatid, timestamp FROM messenger_keyspace.message WHERE chatid = ?"
	iter := r.session.Query(query, chatRoomID).WithContext(ctx).Iter()

	for iter.Scan(&message.ID, &message.Content, &message.SenderID, &message.ChatRoomID, &message.Timestamp) {
		messages = append(messages, models.MessageWrap{
			ID:         message.ID,
			Content:    message.Content,
			SenderID:   message.SenderID,
			ChatRoomID: message.ChatRoomID,
			Timestamp:  message.Timestamp,
		})
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *messageRepository) GetLastMessageFromChatRoom(ctx context.Context, chatRoomID gocql.UUID) (string, error) {
	query := "SELECT content FROM messenger_keyspace.message WHERE chatID = ? LIMIT 1"

	// Execute the query
	var message string
	if err := r.session.Query(query, chatRoomID).WithContext(ctx).Scan(&message); err != nil {
		if err == gocql.ErrNotFound {
			return "", fmt.Errorf("no messages found for chat room %s", chatRoomID)
		}
		log.Printf("Error executing query: %v", err)
		return "", err
	}

	return message, nil
}
