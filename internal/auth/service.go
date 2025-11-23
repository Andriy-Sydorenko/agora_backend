package auth

import (
	"context"
	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
)

type Service struct {
	userService *user.Service
}

func NewService(userService *user.Service) *Service {
	return &Service{
		userService: userService,
	}
}

func (s *Service) Register(ctx context.Context, email, username, password string) (*user.User, error) {
	return s.userService.CreateUser(ctx, email, username, password)
}
