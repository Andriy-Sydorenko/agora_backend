package router

import (
	"github.com/Andriy-Sydorenko/agora_backend/internal/auth"
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/Andriy-Sydorenko/agora_backend/internal/database"
	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	// Infrastructure layer - Database
	db := database.Connect(&cfg.Database)

	// Data layer - Repositories
	userRepo := user.NewRepository(db)

	// Domain layer - Services
	userService := user.NewService(userRepo)
	authService := auth.NewService(userService)

	// Presentation layer - Handlers
	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService)

	//Router setup
	router := gin.Default()
	router.SetTrustedProxies(cfg.Server.Cors.AllowedOrigins)

	// Register domain routes
	user.RegisterRoutes(router, userHandler)
	auth.RegisterRoutes(router, authHandler)

	return router
}
