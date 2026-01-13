package router

import (
	"github.com/Andriy-Sydorenko/agora_backend/internal/auth"
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/Andriy-Sydorenko/agora_backend/internal/database"
	"github.com/Andriy-Sydorenko/agora_backend/internal/email"
	"github.com/Andriy-Sydorenko/agora_backend/internal/subreddit"
	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	// Infrastructure layer - Database
	db := database.Connect(&cfg.Database)
	// Infrastructure layer - Redis (singleton)
	redisClient := database.ConnectRedisClient(&cfg.Redis)

	// Data layer - Repositories
	userRepo := user.NewRepository(db)
	subredditRepo := subreddit.NewRepository(db)

	// Domain layer - Services
	emailService := email.NewService(cfg.SMTP)
	userService := user.NewService(userRepo)
	authService := auth.NewService(
		userService,
		emailService,
		&cfg.Google,
		&cfg.Project,
		redisClient,
	)
	subredditService := subreddit.NewService(subredditRepo)

	// Presentation layer - Handlers
	userHandler := user.NewHandler(userService, cfg)
	authHandler := auth.NewHandler(authService, cfg)
	subredditHandler := subreddit.NewHandler(subredditService, cfg)

	// Router setup
	router := gin.Default()
	router.Use(utils.CORS(&cfg.Server.Cors))

	// Register domain routes
	user.RegisterRoutes(router, userHandler)
	auth.RegisterRoutes(router, authHandler)
	subreddit.RegisterRoutes(router, subredditHandler)

	return router
}
