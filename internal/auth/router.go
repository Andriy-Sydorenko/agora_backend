package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, h *Handler) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", h.Register)
		authRouter.POST("/login", h.Login)
	}
}
