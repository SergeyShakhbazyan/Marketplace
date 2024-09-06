package repository

import (
	"context"
	"errors"
	"github.com/gocql/gocql"
	"marketplace_websocket/internal/models"
)

type ChatRoomRepository interface {
	GetChatRoomByUsers(ctx context.Context, firstUserID gocql.UUID, secondUserID gocql.UUID) (*gocql.UUID, error)
	CreateChatRoom(ctx context.Context, chatRoom models.ChatRoom) error
	GetChatsByUserID(ctx context.Context, userID gocql.UUID) ([]models.ChatRoom, error)
}

type chatRoomRepository struct {
	session *gocql.Session
}

func NewChatRoomRepository(session *gocql.Session) ChatRoomRepository {
	return &chatRoomRepository{session: session}
}

func (r *chatRoomRepository) GetChatRoomByUsers(ctx context.Context, firstUserID, secondUserID gocql.UUID) (*gocql.UUID, error) {
	var chatID gocql.UUID
	query := "SELECT chatid FROM messenger_keyspace.chatroom WHERE user1 = ? AND user2 = ?"
	err := r.session.Query(query, firstUserID, secondUserID).WithContext(ctx).Scan(&chatID)

	if err == nil {
		return &chatID, nil
	}

	query = "SELECT chatid FROM messenger_keyspace.chatroom WHERE user1 = ? AND user2 = ?"
	err = r.session.Query(query, secondUserID, firstUserID).WithContext(ctx).Scan(&chatID)

	if err == nil {
		return &chatID, nil
	}

	if errors.Is(err, gocql.ErrNotFound) {
		return nil, nil
	}

	return nil, err
}

func (r *chatRoomRepository) CreateChatRoom(ctx context.Context, chatRoom models.ChatRoom) error {
	query := "INSERT INTO messenger_keyspace.chatroom (chatid, user1, user2) VALUES (?, ?, ?)"
	if err := r.session.Query(query, chatRoom.ID, chatRoom.CurrentUserID, chatRoom.TargetUserID).WithContext(ctx).Exec(); err != nil {
		return err
	}
	return nil
}

func (r *chatRoomRepository) GetChatsByUserID(ctx context.Context, userID gocql.UUID) ([]models.ChatRoom, error) {
	var chatRooms []models.ChatRoom

	query1 := `SELECT chatid, user1, user2 FROM messenger_keyspace.chatroom WHERE user1 = ?`
	iter1 := r.session.Query(query1, userID).WithContext(ctx).Iter()
	var chatID, user1, user2 gocql.UUID

	for iter1.Scan(&chatID, &user1, &user2) {
		chatRooms = append(chatRooms, models.ChatRoom{
			ID:            chatID,
			CurrentUserID: user1,
			TargetUserID:  user2,
		})
	}

	if err := iter1.Close(); err != nil {
		return nil, err
	}

	query2 := `SELECT chatid, user1, user2 FROM messenger_keyspace.chatroom WHERE user2 = ?`
	iter2 := r.session.Query(query2, userID).WithContext(ctx).Iter()

	for iter2.Scan(&chatID, &user1, &user2) {
		chatRooms = append(chatRooms, models.ChatRoom{
			ID:            chatID,
			CurrentUserID: user1,
			TargetUserID:  user2,
		})
	}

	if err := iter2.Close(); err != nil {
		return nil, err
	}

	return chatRooms, nil
}
