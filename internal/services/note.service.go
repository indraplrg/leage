package services

import (
	"context"
	"errors"
	"math"
	"share-notes-app/internal/dtos"
	"share-notes-app/internal/models"
	"share-notes-app/internal/repositories"
	"share-notes-app/pkg/apperror"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type NoteService interface {
	CreateNote(ctx context.Context, dto dtos.NoteRequest, tokenPayload *dtos.AuthPayload) (*models.Note, error)
	GetAllNotes(ctx context.Context, page, limit int) ([]dtos.NoteResponse, dtos.PaginationMeta, error)
	GetUserNotes(ctx context.Context, page, limit int, tokenPayload *dtos.AuthPayload) ([]dtos.NoteResponse, dtos.PaginationMeta, error)
	GetOneNote(ctx context.Context, noteID string) (*models.Note, error)
	UpdateNote(ctx context.Context, noteID string, req dtos.UpdateNoteRequest) (*models.Note, error)
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
		logrus.Info(err)
		return nil, errors.New("gagal parsing string ke uuid")
	}

	// Buat note
	Note:= &models.Note{
		UserID: parseUUID,
		Title: dto.Title,
		Content: dto.Content,
		IsPublic: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	//  Simpan ke database
	err = s.repo.CreateNote(ctx, Note) 
	if err != nil {
		logrus.Info(err)
		return nil, errors.New("gagal membuat note")
	}

	return Note, nil
}

func (s *noteService) GetAllNotes(ctx context.Context, page, limit int) ([]dtos.NoteResponse, dtos.PaginationMeta, error) {
	// validasi 
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}
	
	// hitung offset
	offset := (page - 1) * limit


	notes, total, err := s.repo.FilteringGetAllNotes(ctx, true, offset, limit)
	if err != nil {
		logrus.Info(err)
		return nil, dtos.PaginationMeta{}, errors.New("gagal mengambil notes") 
	}

	noteDTOs := make([]dtos.NoteResponse, 0, len(notes))
	for _, n := range notes {
		noteDTOs = append(noteDTOs, dtos.NoteResponse{
			ID: n.ID.String(),
			Title: n.Title,
			Content: n.Content,
			IsPublic: n.IsPublic,
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
			User: dtos.UserResponse{
				ID: n.UserID,
				Username: n.User.Username,
				Email: n.User.Email,
				CreatedAt: n.User.CreatedAt,
			},
		})
	} 

	// hitung meta
	totalPage := int(math.Ceil(float64(total) / float64(limit)))
	meta := dtos.PaginationMeta{
		Page: page,
		Limit: limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return noteDTOs, meta, nil
} 

func (s *noteService) GetUserNotes(ctx context.Context, page, limit int, tokenPayload *dtos.AuthPayload) ([]dtos.NoteResponse, dtos.PaginationMeta, error) {
	// validasi 
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}
	
	// hitung offset
	offset := (page - 1) * limit

	notes, total, err := s.repo.GetUserNotes(ctx, tokenPayload.UserID, limit, offset)
	if err != nil {
		logrus.Info(err)
		return nil, dtos.PaginationMeta{}, errors.New("gagal mengambil notes") 
	}

	noteDTOs := make([]dtos.NoteResponse, 0, len(notes))
	for _, n := range notes {
		noteDTOs = append(noteDTOs, dtos.NoteResponse{
			ID: n.ID.String(),
			Title: n.Title,
			Content: n.Content,
			IsPublic: n.IsPublic,
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
			User: dtos.UserResponse{
				ID: n.UserID,
				Username: n.User.Username,
				Email: n.User.Email,
				CreatedAt: n.User.CreatedAt,
			},
		})
	}
	
	// hitung meta
	totalPage := int(math.Ceil(float64(total) / float64(limit)))
	meta := dtos.PaginationMeta{
		Page: page,
		Limit: limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return noteDTOs, meta, nil 
}

func (s *noteService) GetOneNote(ctx context.Context, noteID string) (*models.Note, error) {
	note, err := s.repo.GetOneNote(ctx, noteID)
	if err != nil {
		logrus.Info(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("note tidak ditemukan")
		}

		return nil, err
	}

	return note, nil
}
func (s *noteService) UpdateNote(ctx context.Context, noteID string, req dtos.UpdateNoteRequest) (*models.Note, error) {
	note, err := s.repo.GetOneNote(ctx, noteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrNoteNotFound
		}

		logrus.Info(err)
		return nil, err
	}

	// update note
	note.Title = req.Title
	note.Content = req.Content
	note.IsPublic = req.IsPublic

	// simpan 
	if err := s.repo.UpdateNote(ctx, note); err != nil {
		logrus.Info(err)
		return nil, err
	}

	return note, nil	
}