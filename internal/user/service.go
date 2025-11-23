package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: create centralized error handler
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this username/email already exists")
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
	hashedPassword, err := hashPassword(password)
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
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("finding user by ID: %w", err)
	}
	return user, nil
}
