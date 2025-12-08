package user

import (
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, h *Handler) {
	userRouter := router.Group("/me")
	{
		userRouter.GET("", utils.JWTAuthMiddleware(h.config), h.GetRequestUser)
	}
}
