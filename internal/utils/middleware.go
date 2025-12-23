package utils

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func JWTAuthMiddleware(cfgJWT *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(cfgJWT.AccessTokenCookieKey)
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
		userID, _, err := DecryptJWT(tokenString, cfgJWT.Secret, TokenTypeAccess)
		if err != nil {
			if errors.Is(err, ErrExpiredToken) {
				c.JSON(
					http.StatusUnauthorized, gin.H{
						"error": "Token expired",
					},
				)
				c.Abort()
				return
			}

			if errors.Is(err, ErrInvalidTokenType) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}

func CORS(cfgCors *config.CorsConfig) gin.HandlerFunc {
	allowedOriginsSet := make(map[string]struct{}, len(cfgCors.AllowedOrigins))
	for _, origin := range cfgCors.AllowedOrigins {
		allowedOriginsSet[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		if slices.Equal(cfgCors.AllowedOrigins, []string{"*"}) {
			c.Next()
			return
		}
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if _, ok := allowedOriginsSet[origin]; !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Header("Access-Control-Allow-Origin", origin) // Permits the requesting origin
		c.Header(
			"Vary",
			"Origin",
		) // Tells caches response varies by Origin header
		c.Header(
			"Access-Control-Allow-Credentials",
			"true",
		) // Allows cookies/auth headers in cross-origin requests
		c.Header(
			"Access-Control-Allow-Methods",
			"GET,POST,PUT,PATCH,DELETE,OPTIONS",
		) // HTTP methods allowed in cross-origin requests
		c.Header(
			"Access-Control-Allow-Headers",
			"Authorization,Content-Type",
		) // Request headers allowed in cross-origin requests

		if strings.EqualFold(c.Request.Method, http.MethodOptions) {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// GetUserIDFromContext extracts and parses user ID from gin context
// Returns the user ID or an error response is sent and false is returned
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDString := c.GetString("user_id")
	if userIDString == "" {
		c.JSON(
			http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			},
		)
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": ErrInvalidAccessToken.Error(),
			},
		)
		return uuid.Nil, false
	}

	return userID, true
}
