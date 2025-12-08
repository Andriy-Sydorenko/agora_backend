package auth

import (
	"context"
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"
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
