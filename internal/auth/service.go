package auth

import (
	"context"
	"fmt"
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Service struct {
	userService *user.Service
	validator   *Validator
}

func NewService(userService *user.Service) *Service {
	return &Service{
		userService: userService,
		validator:   NewValidator(userService),
	}
}

func (s *Service) Register(ctx context.Context, email, username, password string) error {
	if errs := s.validator.ValidateRegistrationInput(ctx, email, username, password); len(errs) > 0 {
		return errs
	}

	_, err := s.userService.CreateUser(ctx, email, username, password)
	return err
}

func (s *Service) Login(ctx context.Context, cfg config.JWTConfig, email, password string) (string, error) {
	if errs := s.validator.ValidateLoginInput(ctx, email, password); len(errs) > 0 {
		return "", errs
	}
	userObj, _ := s.userService.GetByEmail(ctx, email)
	jwtToken, err := utils.GenerateJWT(cfg.Secret, cfg.AccessLifetime, userObj.ID.String())
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (s *Service) GoogleURL(cfg *config.Config) (string, error) {
	state, err := GenerateState(cfg.JWT.Secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate state token: %w", err)
	}

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		RedirectURL:  cfg.Google.ClientRedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
	authURL := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.SetAuthURLParam("prompt", "select_account"))
	return authURL, nil
}
