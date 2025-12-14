package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

type Service struct {
	userService *user.Service
	validator   *Validator
	oauthConfig *oauth2.Config
}

var (
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrOAuthAccountNoPassword = errors.New("account uses OAuth, no password set")
)

func NewService(userService *user.Service, googleCfg config.GoogleConfig) *Service {
	oauthConfig := &oauth2.Config{
		ClientID:     googleCfg.ClientID,
		ClientSecret: googleCfg.ClientSecret,
		RedirectURL:  googleCfg.ClientRedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return &Service{
		userService: userService,
		validator:   NewValidator(userService),
		oauthConfig: oauthConfig,
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
	userObj, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if userObj.Password == nil {
		return "", ErrOAuthAccountNoPassword
	}

	if !utils.VerifyPassword(password, *userObj.Password) {
		return "", ErrInvalidCredentials
	}

	jwtToken, err := utils.GenerateJWT(cfg.Secret, cfg.AccessLifetime, userObj.ID.String())
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (s *Service) CreateGoogleURL(cfg *config.Config) (string, error) {
	state, err := GenerateState(cfg.JWT.Secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate state token: %w", err)
	}
	authURL := s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.SetAuthURLParam("prompt", "select_account"))
	return authURL, nil
}

func (s *Service) HandleGoogleCallback(ctx context.Context, jwtCfg *config.JWTConfig, code, state string) (string, error) {
	if err := ValidateState(state, jwtCfg.Secret); err != nil {
		return "", fmt.Errorf("invalid state: %w", err)
	}

	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("code exchange failed: %w", err)
	}

	userInfo, err := s.fetchGoogleUserInfo(ctx, token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user info: %w", err)
	}

	userObj, err := s.userService.FindOrCreateByGoogle(ctx, userInfo.Email, userInfo.ID, userInfo.AvatarURL)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return utils.GenerateJWT(jwtCfg.Secret, jwtCfg.AccessLifetime, userObj.ID.String())
}

func (s *Service) fetchGoogleUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	// TODO: replace hardcoded values
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}
