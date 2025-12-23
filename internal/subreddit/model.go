package subreddit

import (
	"time"

	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subreddit struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"uniqueIndex;not null;size:21"`
	DisplayName string    `gorm:"not null;size:255"`
	Description *string   `gorm:"size:500"`
	IconURL     *string   `gorm:"size:500"`

	CreatorID uuid.UUID   `gorm:"type:uuid;not null;index"`
	Creator   user.User   `gorm:"foreignKey:CreatorID;references:ID"`
	Members   []user.User `gorm:"many2many:subreddit_members"`

	// Update on every action(sub/unsub)
	MemberCount int `gorm:"default:0;not null"`
	PostCount   int `gorm:"default:0;not null"`

	IsPublic bool `gorm:"default:true;not null"`
	IsNSFW   bool `gorm:"default:false;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SubredditMember struct {
	SubredditID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt   time.Time `gorm:"not null"`
}
