package database

import (
	"share-notes-app/internal/models"

	"gorm.io/gorm"
)

func CreateMigrationTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Note{}, &models.EmailVerification{})
}