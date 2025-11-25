package user

import (
	"context"
	"fmt"
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	//Validator *Validator
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
		//Validator: NewValidator(), //TODO: implement custom validation if needed later
	}
}

func (s *Service) CreateUser(ctx context.Context, email, username, password string) (*User, error) {
	// Using individual field for better separation of concerns (HTTP vs business logic),
	// and maximum reusability (if we'll need service for grpc, jobs, tests)
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}
	user := &User{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
		Password: hashedPassword,
	}

	return user, s.repo.Create(ctx, user)
}

func (s *Service) GetUserById(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *Service) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.repo.ExistsByEmail(ctx, email)
}

func (s *Service) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return s.repo.ExistsByUsername(ctx, username)
}
