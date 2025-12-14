package user

import (
	"context"
	"fmt"
	"github.com/Andriy-Sydorenko/agora_backend/internal/utils"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, email, username, password string) (*User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}
	user := &User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		Password:     &hashedPassword,
		AuthProvider: AuthProviderEmail,
	}

	return user, s.repo.Create(ctx, user)
}

func (s *Service) CreateUserByGoogle(ctx context.Context, email, username, googleID, avatarURL string) (*User, error) {
	user := &User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		AuthProvider: AuthProviderGoogle,
		GoogleID:     &googleID,
		AvatarURL:    &avatarURL,
	}

	return user, s.repo.Create(ctx, user)
}

func (s *Service) GetUserById(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) GetByGoogleID(ctx context.Context, googleID string) (*User, error) {
	return s.repo.GetByGoogleID(ctx, googleID)
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

func (s *Service) FindOrCreateByGoogle(ctx context.Context, email, googleID, avatarURL string) (*User, error) {
	if user, err := s.repo.GetByGoogleID(ctx, googleID); err == nil {
		return user, nil
	}

	if user, err := s.repo.GetByEmail(ctx, email); err == nil {
		user.GoogleID = &googleID
		user.AvatarURL = &avatarURL
		return user, s.repo.Update(ctx, user)
	}

	username := GenerateUsernameFromEmail(email)

	return s.CreateUserByGoogle(ctx, email, username, googleID, avatarURL)
}
