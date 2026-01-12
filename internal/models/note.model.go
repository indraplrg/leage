package models

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;type:uuid;not null"`
	User User `gorm:"foreignKey:UserID"`
	Title string `gorm:"not null"`
	Content string
	IsPublic bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}