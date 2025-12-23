package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AuthProvider string

const (
	AuthProviderEmail  AuthProvider = "email"
	AuthProviderGoogle AuthProvider = "google"
)

type User struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey"`
	Username     string       `gorm:"size:255;uniqueIndex;not null"`
	Email        string       `gorm:"size:255;uniqueIndex;not null"`
	Password     *string      `gorm:"size:255"`
	GoogleID     *string      `gorm:"size:255;uniqueIndex"`
	AvatarURL    *string      `gorm:"size:500"`
	AuthProvider AuthProvider `gorm:"size:20;not null;default:'email'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
