package models

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserId uuid.UUID `gorm:"type:uuid;not null"`
	User *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title string `gorm:"not null"`
	Content string
	IsPublic bool
	CreatedAt time.Time
	UpdateTime time.Time
}