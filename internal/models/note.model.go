package models

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;type:uuid;not null"`
	Title string `gorm:"not null"`
	Content string
	IsPublic bool
	CreatedAt time.Time
	UpdateAt time.Time
}