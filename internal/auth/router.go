package auth

import (
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, h *Handler) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", h.Register)
		authRouter.POST("/login", h.Login)
		authRouter.POST("/logout", utils.JWTAuthMiddleware(&h.config.JWT), h.Logout)
	}

	registerGoogleAuthRoutes(authRouter, h)
}

func registerGoogleAuthRoutes(baseRouter *gin.RouterGroup, h *Handler) {
	googleAuthRouter := baseRouter.Group("/google")
	{
		googleAuthRouter.GET("/url", h.GoogleURL)
		googleAuthRouter.GET("/callback", h.HandleGoogleCallback)
	}
}
