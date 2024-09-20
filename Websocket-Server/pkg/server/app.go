package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"marketplace_websocket/internal/db"
	"marketplace_websocket/internal/handlers"
	"marketplace_websocket/internal/repository"
	"marketplace_websocket/internal/service"
	"marketplace_websocket/internal/websocket"
)

type App struct {
	hub    *websocket.Hub
	router *gin.Engine
}

func (a *App) Initialize() {
	a.hub = websocket.NewHub()
	go a.hub.Run()

	a.router = gin.Default()
	a.router.Use(CORSMiddleware())

	session := db.Connection()
	messageRepo := repository.NewMessageRepository(session)
	chatRoomRepo := repository.NewChatRoomRepository(session)
	messageService := service.NewMessageService(messageRepo, chatRoomRepo)
	chatRoomService := service.NewChatRoomService(chatRoomRepo)
	chatRoomHandler := handlers.NewChatRoomHandler(chatRoomService, messageService, a.hub)
	messageHandler := handlers.NewMessageHandler(messageService, chatRoomService)

	a.setupRouterSocket(messageService)
	a.setupRouterChat(messageHandler, chatRoomHandler)
}

func (a *App) Run(addr string) {
	log.Printf("Server is running on %s", addr)
	log.Fatal(a.router.Run(addr))
}
