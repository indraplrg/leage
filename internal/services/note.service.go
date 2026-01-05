package services

import (
	"context"
	"errors"
	"share-notes-app/internal/dtos"
	"share-notes-app/internal/models"
	"share-notes-app/internal/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type NoteService interface {
	CreateNote(ctx context.Context, dto dtos.NoteRequest, tokenPayload *dtos.AuthPayload) (*models.Note, error)
}

type noteService struct {
	repo repositories.NoteRepositories
}

func NewNoteService(repo repositories.NoteRepositories) NoteService {
	return &noteService{repo: repo}
}

func (s *noteService) CreateNote(ctx context.Context, dto dtos.NoteRequest, tokenPayload *dtos.AuthPayload) (*models.Note, error) {
	
	userIDString := tokenPayload.UserID
	parseUUID, err := uuid.Parse(userIDString) 
	if err != nil {
		logrus.WithError(err)
		return nil, errors.New("gagal parsing string ke uuid")
	}

	// Buat note
	Note:= &models.Note{
		UserID: parseUUID,
		Title: dto.Title,
		Content: dto.Content,
		IsPublic: true,
		CreatedAt: time.Now(),
		UpdateAt: time.Now(),
	}

	//  Simpan ke database
	err = s.repo.CreateNote(ctx, Note) 
	if err != nil {
		logrus.WithError(err)
		return nil, errors.New("gagal membuat note")
	}

	return Note, nil
}