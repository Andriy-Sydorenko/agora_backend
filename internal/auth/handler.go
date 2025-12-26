package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"

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
			c.JSON(
				http.StatusBadRequest, gin.H{
					"error":   "Validation failed",
					"details": validationErrs,
				},
			)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	c.JSON(
		http.StatusCreated, gin.H{
			"message": "Registration successful",
		},
	)
}

func (h *Handler) Login(c *gin.Context) {
	var req BasicLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tokenPair, err := h.service.Login(c.Request.Context(), h.config.JWT, req.Email, req.Password)

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

		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		if errors.Is(err, ErrOAuthAccountNoPassword) {
			c.JSON(
				http.StatusBadRequest, gin.H{
					"error": "This account uses Google Sign-In. Please login with Google.",
				},
			)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	h.setTokenCookies(c, tokenPair)

	c.JSON(
		http.StatusOK, gin.H{
			"message": "Login successful",
		},
	)
}

func (h *Handler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(h.config.JWT.RefreshTokenCookieKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}
	err = h.service.blacklistToken(
		c.Request.Context(),
		&h.config.JWT,
		refreshToken,
		utils.TokenTypeRefresh,
	)
	if err != nil {
		// TODO: Implement logging instead of builtin logic
		log.Println("Failed to blacklist token")
	}
	c.SetCookie(
		h.config.JWT.AccessTokenCookieKey,
		"",
		-1,
		"/",
		"",
		h.config.Project.IsProduction,
		true,
	)
	c.SetCookie(
		h.config.JWT.RefreshTokenCookieKey,
		"",
		-1,
		"/",
		"",
		h.config.Project.IsProduction,
		true,
	)
	c.JSON(
		http.StatusOK, gin.H{
			"message": "Logout successful",
		},
	)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(h.config.JWT.RefreshTokenCookieKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	tokenPair, err := h.service.refreshTokens(c.Request.Context(), refreshToken, &h.config.JWT)

	if err != nil {
		if errors.Is(err, utils.ErrInvalidRefreshToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	h.setTokenCookies(c, tokenPair)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed successfully"})
}

func (h *Handler) GoogleURL(c *gin.Context) {
	googleURL, err := h.service.CreateGoogleURL(h.config)

	if err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": "Problem generating google auth url",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"url": googleURL,
		},
	)
}

func (h *Handler) HandleGoogleCallback(c *gin.Context) {
	googleAuthCode := c.Query("code")
	googleAuthState := c.Query("state")
	if googleAuthCode == "" || googleAuthState == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state"})
		return
	}

	tokenPair, err := h.service.HandleGoogleCallback(
		c.Request.Context(),
		&h.config.JWT,
		googleAuthCode,
		googleAuthState,
	)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized, gin.H{
				"error": "OAuth authentication failed",
			},
		)
	}

	h.setTokenCookies(c, tokenPair)
	c.Redirect(http.StatusTemporaryRedirect, h.config.Project.FrontendURL)
}

func (h *Handler) setTokenCookies(c *gin.Context, tokenPair *utils.TokenPair) {
	c.SetCookie(
		h.config.JWT.AccessTokenCookieKey,
		tokenPair.AccessToken,
		int(h.config.JWT.AccessLifetime.Seconds()),
		"/",
		"",
		h.config.Project.IsProduction,
		true,
	)

	c.SetCookie(
		h.config.JWT.RefreshTokenCookieKey,
		tokenPair.RefreshToken,
		int(h.config.JWT.RefreshLifetime.Seconds()),
		"/",
		"",
		h.config.Project.IsProduction,
		true,
	)
}
