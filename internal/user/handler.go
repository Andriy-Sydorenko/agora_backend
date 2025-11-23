package user

import (
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

func (h *Handler) GetRequestUser(c *gin.Context) {
	// TODO: implement middleware and extract token/decode uuid from it
	//_, err := h.service.GetUserById(c.Request.Context())

	//if err != nil {
	//	// TODO: add corresponding error handling
	//	if errors.Is(err, ErrUserAlreadyExists) {
	//		c.JSON(400, gin.H{"error": "User already exists!"})
	//	}
	//	log.Printf("Registration failed: %v", err)
	//	c.JSON(500, gin.H{"error": "Registration failed"})
	//	return
	//}
}
