package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, h *Handler) {
	userRouter := router.Group("/me")
	{
		userRouter.GET("", h.GetRequestUser)
	}
}
