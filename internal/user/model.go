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
	// TODO: consider using `gorm.Model`
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey"`
	Username     string       `gorm:"uniqueIndex;not null"`
	Email        string       `gorm:"uniqueIndex;not null"`
	Password     *string      `gorm:"type:varchar(255)"`
	GoogleID     *string      `gorm:"type:varchar(255);uniqueIndex"`
	AvatarURL    *string      `gorm:"type:varchar(500)"`
	AuthProvider AuthProvider `gorm:"type:varchar(20);not null;default:'email'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
