package repositories

import (
	"context"
	"share-notes-app/internal/models"

	"gorm.io/gorm"
)

type NoteRepositories interface {
	CreateNote(ctx context.Context, entity *models.Note) error
	FilteringGetAllNotes(ctx context.Context, status bool, offset, limit int) ([]models.Note, int64, error)
	GetUserNotes(ctx context.Context, userID string, limit, offset int) ([]models.Note, int64, error)
	GetOneNote(ctx context.Context, noteID string) (*models.Note, error)
	UpdateNote(ctx context.Context, updateField *models.Note) error
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

func (r *noteRepository) FilteringGetAllNotes(ctx context.Context, status bool, offset, limit int) ([]models.Note, int64, error) {
	var notes []models.Note
	var total int64

	tx := r.db.WithContext(ctx)
	
	// query data 
	if err := tx.Joins("User").Where("is_public = ?", status).Order("notes.id ASC").Limit(limit).Offset(offset).Find(&notes).Error; err != nil {
		return nil, 0, err
	}

	// query total
	if err := tx.Model(&models.Note{}).Where("is_public = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

func (r *noteRepository) GetUserNotes(ctx context.Context, userID string, limit, offset int) ([]models.Note, int64, error) {
	var notes []models.Note
	var total int64

	tx := r.db.WithContext(ctx)
	
	// query data 
	if err := tx.Joins("User").Where("user_id = ?", userID).Order("notes.id ASC").Limit(limit).Offset(offset).Find(&notes).Error; err != nil {
		return nil, 0, err
	}

	// query total
	if err := tx.Model(&models.Note{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return notes, total, nil	
}

func (r *noteRepository) GetOneNote(ctx context.Context, noteID string) (*models.Note, error) {
	var note *models.Note
	tx := r.db.WithContext(ctx)

	if err := tx.Preload("User").Where("id = ?", noteID).First(&note).Error; err != nil {
		return nil, err
	}

	return note, nil
}

func (r *noteRepository) UpdateNote(ctx context.Context, updateField *models.Note) error {
	return r.db.WithContext(ctx).Save(updateField).Error;
}
