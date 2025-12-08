package utils

import (
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(cfg.JWT.JwtTokenCookieKey)
		if err != nil {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
				c.Abort()
				return
			}
		}
		userID, err := DecryptJWT(tokenString, cfg.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}
