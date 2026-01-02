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
	Notes *[]Note `gorm:"foreignKey:UserId"`
	CreatedAt time.Time
	IsVerified bool `gorm:"default:false"`
	RefreshToken *Token `gorm:"foreignKey:UserId"`
	Verifications *[]EmailVerification `gorm:"foreignKey:UserId"`
}