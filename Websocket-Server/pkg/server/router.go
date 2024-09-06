package server

import (
	"github.com/gin-gonic/gin"
	"marketplace_websocket/internal/handlers"
	"marketplace_websocket/internal/service"
	"marketplace_websocket/internal/websocket"
)

func (a *App) setupRouterSocket(messageService *service.MessageService) {

	a.router.GET("/ws", func(c *gin.Context) {
		websocket.HandleConnection(a.hub, c.Writer, c.Request, messageService)
	})

}

func (a *App) setupRouterChat(messageHandler *handlers.MessageHandler, chatRoomHandler *handlers.ChatRoomHandler) {
	a.router.GET("/getChatMessages", chatRoomHandler.GetMessagesFromChatID)
	a.router.GET("/getChatRoomID", chatRoomHandler.GetChatIDByUsers)
	a.router.GET("/getUserChats", chatRoomHandler.GetUserChats)
}
