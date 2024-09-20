package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"marketplace_websocket/internal/service"
	"net/http"
)

type MessageHandler struct {
	messageService  *service.MessageService
	chatRoomService *service.ChatRoomService
}

func NewMessageHandler(messageService *service.MessageService, chatRoomService *service.ChatRoomService) *MessageHandler {
	return &MessageHandler{messageService: messageService, chatRoomService: chatRoomService}
}

func (h *MessageHandler) GetMessagesFromChatID(c *gin.Context) {
	chatID, err := gocql.ParseUUID(c.Query("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}
	messages, err := h.messageService.GetMessagesFromChatRoom(context.Background(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
