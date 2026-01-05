package subreddit

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo      *Repository
	validator *Validator
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo:      repo,
		validator: NewValidator(repo),
	}
}

var (
	ErrNotAuthorized      = errors.New("not authorized to perform this action")
	ErrCreatorCannotLeave = errors.New("creator cannot leave subreddit, delete it instead")
)

func (s *Service) GetSubredditList(ctx context.Context) ([]Subreddit, error) {
	return s.repo.GetList(ctx)
}

func (s *Service) GetSubredditById(ctx context.Context, id uuid.UUID) (*Subreddit, error) {
	return s.repo.GetByID(ctx, id, true)
}

func (s *Service) CreateSubreddit(
	ctx context.Context, creatorID uuid.UUID, name string, displayName string, description *string,
	iconURL *string, isPublic bool, isNSFW bool,
) (*Subreddit, error) {

	if errs := s.validator.ValidateCreateSubredditInput(
		ctx, name, displayName, description, iconURL,
	); len(errs) > 0 {
		return nil, errs
	}

	subreddit := &Subreddit{
		ID:          uuid.New(),
		Name:        name,
		DisplayName: displayName,
		Description: description,
		IconURL:     iconURL,
		CreatorID:   creatorID,
		MemberCount: 1,
		PostCount:   0,
		IsPublic:    isPublic,
		IsNSFW:      isNSFW,
	}

	err := s.repo.WithTx(
		ctx, func(txRepo Repository) error {
			if err := txRepo.Create(ctx, subreddit); err != nil {
				return err
			}
			if err := txRepo.AddMember(ctx, subreddit.ID, creatorID); err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return subreddit, nil
}

func (s *Service) UpdateSubreddit(
	ctx context.Context,
	subredditID, userID uuid.UUID,
	req UpdateSubredditRequest,
) (
	*Subreddit,
	error,
) {

	_, err := s.ensureCreator(ctx, subredditID, userID)
	if err != nil {
		return nil, err
	}

	if errs := s.validator.ValidateUpdateSubredditInput(req); len(errs) > 0 {
		return nil, errs
	}

	// TODO: Implement better fields mapping
	updates := make(map[string]interface{})

	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Description != nil {
		updates["description"] = req.Description
	}
	if req.IconURL != nil {
		updates["icon_url"] = req.IconURL
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.IsNSFW != nil {
		updates["is_nsfw"] = *req.IsNSFW
	}

	err = s.repo.Update(ctx, subredditID, updates)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, subredditID, false)
}

func (s *Service) DeleteSubreddit(
	ctx context.Context,
	subredditID,
	userID uuid.UUID,
) error {
	_, err := s.ensureCreator(ctx, subredditID, userID)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, subredditID)
}

func (s *Service) ensureCreator(ctx context.Context, subredditID, userID uuid.UUID) (
	*Subreddit,
	error,
) {
	subreddit, err := s.repo.GetByID(ctx, subredditID, false)
	if err != nil {
		return nil, err
	}

	if subreddit.CreatorID != userID {
		return nil, ErrNotAuthorized
	}

	return subreddit, nil
}

func (s *Service) JoinSubreddit(ctx context.Context, subredditID, userID uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, subredditID, false)
	if err != nil {
		return err
	}

	// TODO: add request logic to join subreddit if it's private
	return s.repo.WithTx(
		ctx, func(txRepo Repository) error {
			return txRepo.AddMember(ctx, subredditID, userID)
		},
	)
}

func (s *Service) LeaveSubreddit(ctx context.Context, subredditID, userID uuid.UUID) error {
	subreddit, err := s.repo.GetByID(ctx, subredditID, false)
	if err != nil {
		return err
	}

	if subreddit.CreatorID == userID {
		return ErrCreatorCannotLeave
	}

	return s.repo.WithTx(
		ctx, func(txRepo Repository) error {
			return txRepo.RemoveMember(ctx, subredditID, userID)
		},
	)
}
