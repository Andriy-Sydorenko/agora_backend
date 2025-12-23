package subreddit

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) WithTx(ctx context.Context, fn func(txRepo Repository) error) error {
	return repo.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			txRepo := Repository{db: tx}
			return fn(txRepo)
		},
	)
}

func (repo *Repository) GetList(ctx context.Context) ([]Subreddit, error) {
	var subreddits []Subreddit

	err := repo.db.WithContext(ctx).
		Preload("Creator").
		Where("is_public = ?", true).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&subreddits).Error

	if err != nil {
		return nil, err
	}

	return subreddits, nil
}

func (repo *Repository) GetByID(ctx context.Context, id uuid.UUID, includeMembers bool) (
	*Subreddit,
	error,
) {
	var subreddit Subreddit
	query := repo.db.WithContext(ctx).
		Preload("Creator").
		Where("id = ?", id)

	if includeMembers {
		query = query.Preload("Members")
	}
	err := query.First(&subreddit).Error
	if err != nil {
		return nil, err
	}

	return &subreddit, nil
}

func (repo *Repository) GetUserSubreddits(ctx context.Context, userID uuid.UUID) (
	[]Subreddit,
	int64,
	error,
) {
	var subreddits []Subreddit
	var total int64

	err := repo.db.WithContext(ctx).
		Table("subreddits").
		Joins("INNER JOIN subreddit_members ON subreddits.id = subreddit_members.subreddit_id").
		Where("subreddit_members.user_id = ?", userID).
		Where("subreddits.deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// TODO: Add pagination
	err = repo.db.WithContext(ctx).
		Preload("Creator").
		Joins("INNER JOIN subreddit_members ON subreddits.id = subreddit_members.subreddit_id").
		Where("subreddit_members.user_id = ?", userID).
		Where("subreddits.deleted_at IS NULL").
		Order("subreddits.created_at DESC").
		Find(&subreddits).Error

	if err != nil {
		return nil, 0, err
	}

	return subreddits, total, nil
}

func (repo *Repository) Create(ctx context.Context, subreddit *Subreddit) error {
	return repo.db.WithContext(ctx).Create(subreddit).Error
}

func (repo *Repository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).Model(&Subreddit{}).
		Unscoped(). // Include soft deleted records to not allow name reusage
		Where("LOWER(name) = ?", strings.ToLower(name)).
		Count(&count).Error
	return count > 0, err
}

func (repo *Repository) Update(
	ctx context.Context,
	subredditID uuid.UUID,
	updates map[string]interface{},
) error {
	result := repo.db.WithContext(ctx).
		Model(&Subreddit{}).
		Where("id = ?", subredditID).
		Updates(updates)

	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result := repo.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&Subreddit{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo *Repository) AddMember(ctx context.Context, subredditID, userID uuid.UUID) error {
	member := SubredditMember{
		SubredditID: subredditID,
		UserID:      userID,
	}

	result := repo.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}). // Idempotent (no error if already member)
		Create(&member)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}

	return repo.db.WithContext(ctx).
		Model(&Subreddit{}).
		Where("id = ?", subredditID).
		UpdateColumn(
			"member_count",
			gorm.Expr("member_count + 1"),
		).Error
}

func (repo *Repository) RemoveMember(ctx context.Context, subredditID, userID uuid.UUID) error {
	result := repo.db.WithContext(ctx).
		Where("subreddit_id = ? AND user_id = ?", subredditID, userID).
		Delete(&SubredditMember{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil // Already not a member, idempotent behavior
	}

	return repo.db.WithContext(ctx).
		Model(&Subreddit{}).
		Where("id = ?", subredditID).
		UpdateColumn(
			"member_count",
			gorm.Expr("member_count - 1"),
		).Error
}
