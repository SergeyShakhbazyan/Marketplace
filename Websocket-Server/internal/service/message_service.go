package service

import (
	"context"
	"github.com/gocql/gocql"
	"marketplace_websocket/internal/models"
	"marketplace_websocket/internal/repository"
)

type MessageService struct {
	messageRepo  repository.MessageRepository
	chatRoomRepo repository.ChatRoomRepository
}

func NewMessageService(messageRepo repository.MessageRepository, chatRoomRepo repository.ChatRoomRepository) *MessageService {
	return &MessageService{
		messageRepo:  messageRepo,
		chatRoomRepo: chatRoomRepo,
	}
}

func (s *MessageService) SaveMessage(ctx context.Context, message *models.Message) (*models.ChatRoom, error) {
	chatRoomID := message.ChatRoomID
	if message.ChatRoomID == (gocql.UUID{}) {
		chatID, err := s.chatRoomRepo.GetChatRoomByUsers(ctx, message.SenderID, message.RecipientID)
		if err != nil {
			return nil, err
		}
		if chatID == nil {
			chatRoomID = gocql.TimeUUID()
			chatRoom := models.ChatRoom{
				ID:            chatRoomID,
				CurrentUserID: message.SenderID,
				TargetUserID:  message.RecipientID,
			}
			err = s.chatRoomRepo.CreateChatRoom(ctx, chatRoom)
			if err != nil {
				return nil, err
			}
			message.ChatRoomID = chatRoomID
			return &chatRoom, s.messageRepo.SaveMessage(ctx, message)
		} else {
			message.ChatRoomID = *chatID
		}
	}
	return nil, s.messageRepo.SaveMessage(ctx, message)
}

func (s *MessageService) GetMessagesFromChatRoom(ctx context.Context, chatRoomID gocql.UUID) ([]models.MessageWrap, error) {
	return s.messageRepo.GetMessagesFromChatRoom(ctx, chatRoomID)
}

func (s *MessageService) GetLastMessage(ctx context.Context, chatRoomID gocql.UUID) (string, error) {
	return s.messageRepo.GetLastMessageFromChatRoom(ctx, chatRoomID)
}

//// SaveMessage saves a new message to the repository.
//func (s *MessageService) SaveMessage(message *models.Message) error {
//	message.Timestamp = time.Now()
//	return s.messageRepo.Save(message)
//}
//
//// GetMessages retrieves messages for a specific chat room or user.
//func (s *MessageService) GetMessages(chatRoomID gocql.UUID, userID gocql.UUID) ([]models.Message, error) {
//	return s.messageRepo.FindByChatRoomAndUser(chatRoomID, userID)
//}
//
//// GetMessagesByRecipient retrieves messages for a specific recipient.
//func (s *MessageService) GetMessagesByRecipient(recipientID gocql.UUID) ([]models.Message, error) {
//	return s.messageRepo.FindByRecipient(recipientID)
//}
