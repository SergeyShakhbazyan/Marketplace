package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"marketplace_project/config"
	"marketplace_project/internal/db"
	"marketplace_project/internal/handler"
	"marketplace_project/internal/middleware"
	"marketplace_project/internal/repository"
	"marketplace_project/internal/service"
)

type App struct {
	Router *gin.Engine
	cfg    config.ServerConfig
}

func (a *App) Initialize() {
	a.Router = gin.Default()

	a.Router.Use(middleware.CORSMiddleware())
	a.cfg = config.ServerConfig{Port: ":3001"}

	session := db.Connection()

	productRepo := repository.NewProductRepository(session)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	userRepo := repository.NewUserRepository(session)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	categoryRepo := repository.NewCategoryRepository(session)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	//productRepo := repository.NewProductRepository(session)
	//productService := service.NewProductService(productRepo)
	//productHandler := handler.NewProductHandler(productService)

	sectionService := service.NewSectionsService(productRepo, categoryRepo, userRepo)
	sectionHandler := handler.NewSectionsHandler(sectionService)

	a.setRoutersForUser(userHandler)
	a.setRoutersForCategory(categoryHandler)
	a.setRoutersForProduct(productHandler)
	a.setRoutersForSections(sectionHandler)
}

func (a *App) Run() {
	if err := a.Router.Run(a.cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
