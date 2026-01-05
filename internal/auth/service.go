package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Service struct {
	userService *user.Service
	validator   *Validator
	oauthConfig *oauth2.Config
	redis       *redis.Client
}

func NewService(
	userService *user.Service,
	googleCfg config.GoogleConfig,
	redisClient *redis.Client,
) *Service {
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
		redis:       redisClient,
	}
}

var (
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrOAuthAccountNoPassword = errors.New("account uses OAuth, no password set")
)

func (s *Service) Register(ctx context.Context, email, username, password string) error {
	if errs := s.validator.ValidateRegistrationInput(
		ctx,
		email,
		username,
		password,
	); len(errs) > 0 {
		return errs
	}

	_, err := s.userService.CreateUser(ctx, email, username, password)
	return err
}

func (s *Service) Login(
	ctx context.Context,
	cfg config.JWTConfig,
	email, password string,
) (*utils.TokenPair, error) {
	if errs := s.validator.ValidateLoginInput(ctx, email, password); len(errs) > 0 {
		return nil, errs
	}
	userObj, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if userObj.Password == nil {
		return nil, ErrOAuthAccountNoPassword
	}

	if !utils.VerifyPassword(password, *userObj.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokenPair(&cfg, userObj.ID.String())
}

func (s *Service) CreateGoogleURL(cfg *config.Config) (string, error) {
	state, err := GenerateState(cfg.JWT.Secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate state token: %w", err)
	}
	authURL := s.oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOnline,
		oauth2.SetAuthURLParam("prompt", "select_account"),
	)
	return authURL, nil
}

func (s *Service) HandleGoogleCallback(
	ctx context.Context,
	jwtCfg *config.JWTConfig,
	code, state string,
) (*utils.TokenPair, error) {
	if err := ValidateState(state, jwtCfg.Secret); err != nil {
		return nil, fmt.Errorf("invalid state: %w", err)
	}

	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	userInfo, err := s.fetchGoogleUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	userObj, err := s.userService.FindOrCreateByGoogle(
		ctx,
		userInfo.Email,
		userInfo.ID,
		userInfo.AvatarURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.generateTokenPair(jwtCfg, userObj.ID.String())
}

func (s *Service) fetchGoogleUserInfo(ctx context.Context, accessToken string) (
	*GoogleUserInfo,
	error,
) {
	// TODO: replace hardcoded values
	req, _ := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://www.googleapis.com/oauth2/v2/userinfo",
		nil,
	)
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

func (s *Service) refreshTokens(
	ctx context.Context,
	refreshToken string,
	cfg *config.JWTConfig,
) (*utils.TokenPair, error) {
	userID, _, err := utils.DecryptJWT(refreshToken, cfg.Secret, utils.TokenTypeRefresh)
	if err != nil {
		return nil, utils.ErrInvalidRefreshToken
	}
	if s.isTokenBlacklisted(ctx, refreshToken) {
		return nil, utils.ErrInvalidRefreshToken
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, utils.ErrInvalidRefreshToken
	}

	_, err = s.userService.GetUserById(ctx, userUUID)
	if err != nil {
		return nil, utils.ErrInvalidRefreshToken
	}

	err = s.blacklistToken(ctx, cfg, refreshToken, utils.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	return s.generateTokenPair(cfg, userID)

}

func (s *Service) generateTokenPair(cfg *config.JWTConfig, userID string) (
	*utils.TokenPair,
	error,
) {
	accessToken, err := utils.GenerateJWT(
		cfg.Secret,
		utils.TokenTypeAccess,
		cfg.AccessLifetime,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateJWT(
		cfg.Secret,
		utils.TokenTypeRefresh,
		cfg.RefreshLifetime,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &utils.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) isTokenBlacklisted(ctx context.Context, token string) bool {
	if s.redis == nil {
		return false
	}

	key := utils.RefreshTokenBlacklistPrefix + token
	exists, err := s.redis.Exists(ctx, key).Result()
	return err == nil && exists > 0
}

func (s *Service) blacklistToken(
	ctx context.Context,
	cfg *config.JWTConfig,
	token string,
	expectedTokenType string,
) error {
	if s.redis == nil {
		return nil
	}

	_, tokenClaims, err := utils.DecryptJWT(token, cfg.Secret, expectedTokenType)
	if err != nil {
		return err
	}
	exp, ok := tokenClaims["exp"].(float64)
	if !ok {
		return fmt.Errorf("token missing expiry claim")
	}

	ttl := time.Until(time.Unix(int64(exp), 0))
	if ttl <= 0 {
		return fmt.Errorf("token already expired")
	}

	key := utils.RefreshTokenBlacklistPrefix + token
	err = s.redis.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set blacklist key: %w", err)
	}

	return nil
}
