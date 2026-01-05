package repositories

import (
	"context"

	"gorm.io/gorm"
)

type NoteRepositories interface {
	CreateNote(ctx context.Context, entity any) error
}

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepositories {
	return &noteRepository{db:db}
}

func (r *noteRepository) CreateNote(ctx context.Context, entity any) error {
	return r.db.WithContext(ctx).Create(entity).Error	
}