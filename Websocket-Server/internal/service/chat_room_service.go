package service

import (
	"context"
	"github.com/gocql/gocql"
	"marketplace_websocket/internal/models"
	"marketplace_websocket/internal/repository"
)

type ChatRoomService struct {
	chatRoomRepo repository.ChatRoomRepository
}

func NewChatRoomService(chatRoomRepo repository.ChatRoomRepository) *ChatRoomService {
	return &ChatRoomService{chatRoomRepo: chatRoomRepo}
}

func (s *ChatRoomService) GetChatRoomByUsers(ctx context.Context, chatRoom models.ChatRoom) (*gocql.UUID, error) {
	return s.chatRoomRepo.GetChatRoomByUsers(ctx, chatRoom.CurrentUserID, chatRoom.TargetUserID)
}

func (s *ChatRoomService) CreateChatRoom(ctx context.Context, chatRoom models.ChatRoom) error {
	return s.chatRoomRepo.CreateChatRoom(ctx, chatRoom)
}

func (s *ChatRoomService) GetUserChats(ctx context.Context, userID gocql.UUID) ([]models.ChatRoom, error) {
	return s.chatRoomRepo.GetChatsByUserID(ctx, userID)
}
