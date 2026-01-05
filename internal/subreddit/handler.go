package subreddit

import (
	"errors"
	"fmt"
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

func (h *Handler) GetSubredditList(c *gin.Context) {
	subreddits, err := h.service.GetSubredditList(c.Request.Context())
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch subreddits",
			},
		)
		return
	}
	response := ToSubredditListResponse(subreddits)
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetSubreddit(c *gin.Context) {
	subredditIDString := c.Param("id")
	subredditID, err := uuid.Parse(subredditIDString)
	if err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": "Invalid subreddit ID",
			},
		)
		return
	}

	subreddit, err := h.service.GetSubredditById(c.Request.Context(), subredditID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound, gin.H{
					"error": "Subreddit not found",
				},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch subreddit",
			},
		)
		return
	}

	response := ToSubredditResponse(subreddit)
	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateSubreddit(c *gin.Context) {
	var req CreateSubredditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			},
		)
		return
	}

	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return // Error response already sent
	}

	// FIXME: is this the best solution for optional/omitted fields?
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	isNSFW := false
	if req.IsNSFW != nil {
		isNSFW = *req.IsNSFW
	}

	subreddit, err := h.service.CreateSubreddit(
		c.Request.Context(),
		userID,
		req.Name,
		req.DisplayName,
		req.Description,
		req.IconURL,
		isPublic,
		isNSFW,
	)

	if err != nil {
		var validationErrs ValidationErrors
		if errors.As(err, &validationErrs) {
			c.JSON(
				http.StatusBadRequest, gin.H{
					"error":   "Validation failed",
					"details": validationErrs,
				},
			)
			return
		}

		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": "Failed to create subreddit",
			},
		)
		return
	}

	response := ToSubredditResponse(subreddit)
	c.JSON(http.StatusCreated, response)
}

func (h *Handler) UpdateSubreddit(c *gin.Context) {
	var req UpdateSubredditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			},
		)
		return
	}

	subredditIDString := c.Param("id")
	subredditID, err := uuid.Parse(subredditIDString)
	if err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": "Invalid subreddit ID",
			},
		)
		return
	}
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return
	}

	subreddit, err := h.service.UpdateSubreddit(c.Request.Context(), subredditID, userID, req)
	if err != nil {
		var validationErrs ValidationErrors
		if errors.As(err, &validationErrs) {
			c.JSON(
				http.StatusBadRequest, gin.H{
					"error":   "Validation failed",
					"details": validationErrs,
				},
			)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound, gin.H{
					"error": "Subreddit not found",
				},
			)
			return
		}
		if errors.Is(err, ErrNotAuthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot perform this action"})
			return
		}
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch subreddits",
			},
		)
		return
	}

	response := ToSubredditResponse(subreddit)
	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteSubreddit(c *gin.Context) {
	subredditIDString := c.Param("id")
	subredditID, err := uuid.Parse(subredditIDString)
	if err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": "Invalid subreddit ID",
			},
		)
		return
	}
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return // Error response already sent
	}

	err = h.service.DeleteSubreddit(c.Request.Context(), subredditID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound, gin.H{
					"error": "Subreddit not found",
				},
			)
			return
		}
		if errors.Is(err, ErrNotAuthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot perform this action"})
			return
		}
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch subreddits",
			},
		)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) JoinSubreddit(c *gin.Context) {
	subredditIDString := c.Param("id")
	subredditID, err := uuid.Parse(subredditIDString)
	if err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": "Invalid subreddit ID",
			},
		)
		return
	}
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return
	}

	err = h.service.JoinSubreddit(c.Request.Context(), subredditID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound, gin.H{
					"error": "Subreddit not found",
				},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to fetch subreddits: %s", err.Error()),
			},
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Joined subreddit successfully"})
}

func (h *Handler) LeaveSubreddit(c *gin.Context) {
	subredditIDString := c.Param("id")
	subredditID, err := uuid.Parse(subredditIDString)
	if err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": "Invalid subreddit ID",
			},
		)
		return
	}
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return
	}

	err = h.service.LeaveSubreddit(c.Request.Context(), subredditID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound, gin.H{
					"error": "Subreddit not found",
				},
			)
			return
		}
		if errors.Is(err, ErrCreatorCannotLeave) {
			c.JSON(
				http.StatusForbidden, gin.H{
					"error": "As the creator of this community, you cannot leave it. If you no longer want to manage it, you can delete the subreddit.",
				},
			)
			return
		}

		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to fetch subreddits: %s", err.Error()),
			},
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Joined subreddit successfully"})
}
