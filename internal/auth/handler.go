package auth

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
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
	err := h.service.Register(c.Request.Context(), req.Email, req.Username, req.Password)

	if err != nil {
		var validationErrs ValidationErrors
		if errors.As(err, &validationErrs) {
			c.JSON(400, gin.H{
				"error":   "Validation failed",
				"details": validationErrs,
			})
			return
		}

		log.Printf("Registration failed: %v", err)
		c.JSON(500, gin.H{"error": "Registration failed"})
		return
	}

	c.JSON(201, gin.H{
		"message": "Registration successful",
	})
}
