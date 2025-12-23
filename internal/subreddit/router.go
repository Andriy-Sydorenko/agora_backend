package subreddit

import (
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, h *Handler) {
	subredditRouter := router.Group("/subreddits")
	{
		subredditRouter.GET("", h.GetSubredditList)
		subredditRouter.GET(":id", h.GetSubreddit)

		subredditRouter.POST("", utils.JWTAuthMiddleware(&h.config.JWT), h.CreateSubreddit)
		subredditRouter.PATCH(":id", utils.JWTAuthMiddleware(&h.config.JWT), h.UpdateSubreddit)
		subredditRouter.DELETE(":id", utils.JWTAuthMiddleware(&h.config.JWT), h.DeleteSubreddit)
	}
}
