package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserId uuid.UUID `gorm:"type:uuid;not null"`
	User User `gorm:"constraint:OnDelete:CASCADE"`
	Token string `gorm:"unique;not null"`
	IsUsed bool `gorm:"default:false"`
	ExpiresAt time.Time
	CreatedAt time.Time
}