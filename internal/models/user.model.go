package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username string `gorm:"unique;not null"`
	Email string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	CreatedAt time.Time
	IsVerified bool `gorm:"default:false"`
	Notes []Note `gorm:"foreignKey:UserID"`
	RefreshToken *Token `gorm:"foreignKey:UserID"`
	Verifications []EmailVerification `gorm:"foreignKey:UserID"`
}