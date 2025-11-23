package user

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	// INFO: making db unexported so that service layer can't bypass repository and use ORM directly(uncoupling service from gorm)
	db *gorm.DB
	// TODO: consider adding cache(cache *redis.Client) if needed later
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create inserts a new user into DB
func (repo *Repository) Create(ctx context.Context, user *User) error {
	return repo.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves user by ID
func (repo *Repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	// TODO: why we're using query result assignment by pointer as destination, instead of user:=...?
	var currentUser User                                          // INFO: using value in this case instead of pointer to prevent nil pointer dereferencing in orm method
	err := repo.db.WithContext(ctx).First(&currentUser, id).Error // INFO: using only First() without Where() here because of the default search by primary key
	if err != nil {
		return nil, err
	}
	return &currentUser, nil
}

// GetByEmail retrieves user by email
func (repo *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var currentUser User
	err := repo.db.WithContext(ctx).Where("email = ?", email).First(&currentUser).Error // INFO: using Where() instead of plain First() for better method chaining and readability
	if err != nil {
		return nil, err
	}
	return &currentUser, nil
}

// GetByUsername retrieves user by username
func (repo *Repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var currentUser User
	err := repo.db.WithContext(ctx).Where("username = ?", username).First(&currentUser).Error
	if err != nil {
		return nil, err
	}
	return &currentUser, nil
}

// TODO: consider unifying ExistsBy methods into 1, or use abstract helper method

// ExistsByEmail checks if user with given email exists
func (repo *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByUsername checks if user with given username exists
func (repo *Repository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).Model(&User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
