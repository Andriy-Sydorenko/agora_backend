package user

import (
	"errors"
	"net/http"

	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
	config  *config.Config
}

func NewHandler(service *Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		config:  cfg,
	}
}

func (h *Handler) GetRequestUser(c *gin.Context) {
	userIDString := c.GetString("user_id")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": utils.ErrInvalidAccessToken.Error(),
			},
		)
		return
	}

	user, err := h.service.GetUserById(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(
		http.StatusOK, PublicUserResponse{
			Email:    user.Email,
			Username: user.Username,
		},
	)
}
