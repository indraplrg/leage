package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;type:uuid;not null"`
	Token string `gorm:"unique;not null"`
	IsUsed bool `gorm:"default:false"`
	ExpiresAt time.Time
	CreatedAt time.Time
}