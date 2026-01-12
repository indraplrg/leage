package repositories

import (
	"context"
	"share-notes-app/internal/models"

	"gorm.io/gorm"
)

type NoteRepositories interface {
	CreateNote(ctx context.Context, entity *models.Note) error
	GetAllNotes(ctx context.Context, status bool, offset, limit int) ([]models.Note, int64, error)
}

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepositories {
	return &noteRepository{db:db}
}

func (r *noteRepository) CreateNote(ctx context.Context, entity *models.Note) error {
	return r.db.WithContext(ctx).Create(entity).Error	
}

func (r *noteRepository) GetAllNotes(ctx context.Context, status bool, offset, limit int) ([]models.Note, int64, error) {
	var notes []models.Note
	var total int64

	tx := r.db.WithContext(ctx)
	
	// query data 
	if err := tx.Joins("User").Where("is_public = ?", true).Order("notes.id ASC").Limit(limit).Offset(offset).Find(&notes).Error; err != nil {
		return nil, 0, err
	}

	// query total
	if err := tx.Model(&models.Note{}).Where("is_public = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}
