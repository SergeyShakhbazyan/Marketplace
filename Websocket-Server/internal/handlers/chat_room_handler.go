package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"io/ioutil"
	"marketplace_websocket/internal/models"
	"marketplace_websocket/internal/service"
	"marketplace_websocket/internal/websocket"
	"net/http"
)

type ChatRoomHandler struct {
	chatRoomService *service.ChatRoomService
	messageService  *service.MessageService
	hub             *websocket.Hub
}

func NewChatRoomHandler(chatRoomService *service.ChatRoomService, messageService *service.MessageService, hub *websocket.Hub) *ChatRoomHandler {
	return &ChatRoomHandler{chatRoomService: chatRoomService, messageService: messageService, hub: hub}
}

func (h *ChatRoomHandler) GetMessagesFromChatID(c *gin.Context) {
	chatID, err := gocql.ParseUUID(c.Query("chatRoomID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}
	messages, err := h.messageService.GetMessagesFromChatRoom(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func (h *ChatRoomHandler) GetChatIDByUsers(c *gin.Context) {
	var chatRoom models.ChatRoom
	var err error
	chatRoom.CurrentUserID, err = gocql.ParseUUID(c.Query("currentUserID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CurrentUserID"})
		return
	}

	chatRoom.TargetUserID, err = gocql.ParseUUID(c.Query("targetUserID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TargetUserID"})
		return
	}
	chatID, err := h.chatRoomService.GetChatRoomByUsers(c.Request.Context(), chatRoom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}
	var newChatID gocql.UUID
	if chatID == nil {
		newChatID = gocql.TimeUUID()
		c.JSON(http.StatusOK, newChatID)
		return
	}

	c.JSON(http.StatusOK, chatID)
}

func getProfileData(userID gocql.UUID) (*models.UserProfile, error) {
	// Convert userID to string
	userIDStr := userID.String()

	// Construct the URL with the userID query parameter
	url := fmt.Sprintf("http://localhost:3001/profileData?userID=%s", userIDStr)

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response body: %v", err)
	}

	var userProfile models.UserProfile
	if err := json.Unmarshal(body, &userProfile); err != nil {
		return nil, fmt.Errorf("error unmarshalling the response body: %v", err)
	}

	return &userProfile, nil
}

func (h *ChatRoomHandler) GetUserChats(c *gin.Context) {
	userID, err := gocql.ParseUUID(c.Query("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	chatRooms, err := h.chatRoomService.GetUserChats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat rooms"})
		return
	}

	chatDetails := make([]models.UserProfile, len(chatRooms))

	for i, chatRoom := range chatRooms {
		var userProfile *models.UserProfile
		if chatRoom.TargetUserID == userID {
			userProfile, err = getProfileData(chatRoom.CurrentUserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
				return
			}
		} else {
			userProfile, err = getProfileData(chatRoom.TargetUserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
				return
			}
		}

		lastMessage, err := h.messageService.GetLastMessage(context.Background(), chatRoom.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last message for chat room " + chatRoom.ID.String()})
			return
		}
		chatDetails[i] = models.UserProfile{
			UserID:      userProfile.UserID,
			FirstName:   userProfile.FirstName,
			LastName:    userProfile.LastName,
			Avatar:      userProfile.Avatar,
			Status:      h.hub.IsUserConnected(userProfile.UserID),
			LastMessage: lastMessage,
			ChatID:      chatRoom.ID,
		}
	}
	c.JSON(http.StatusOK, chatDetails)
}
