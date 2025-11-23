package auth

import (
	"errors"
	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
	"github.com/gin-gonic/gin"
	"log"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	_, err := h.service.Register(c.Request.Context(), req.Email, req.Username, req.Password)

	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			c.JSON(400, gin.H{"error": "User already exists!"})
		} else {
			log.Printf("Registration failed: %v", err)
			c.JSON(500, gin.H{"error": "Registration failed"})
		}
		return
	}

	c.JSON(201, gin.H{
		"message": "Registration successful",
	})
}
