package auth

import (
	"errors"
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := h.service.Register(c.Request.Context(), req.Email, req.Username, req.Password)

	if err != nil {
		var validationErrs ValidationErrors
		if errors.As(err, &validationErrs) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationErrs,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req BasicLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	jwtToken, err := h.service.Login(c.Request.Context(), h.config.JWT, req.Email, req.Password)

	if err != nil {
		var validationErrs ValidationErrors
		if errors.As(err, &validationErrs) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationErrs,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	c.SetCookie(h.config.JWT.JwtTokenCookieKey, jwtToken, int(h.config.JWT.AccessLifetime.Seconds()), "/", "", h.config.Project.IsProduction, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}
