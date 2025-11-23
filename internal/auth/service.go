package auth

import (
	"context"
	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
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
