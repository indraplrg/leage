package controllers

import (
	"errors"
	"share-notes-app/pkg/apperror"

	"net/http"
	"share-notes-app/internal/dtos"
	"share-notes-app/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NoteController struct {
	service services.NoteService
}

func NewNoteController(service services.NoteService) *NoteController {
	return &NoteController{service:service}
}

func (c *NoteController) CreateNote(ctx *gin.Context) {
	var req dtos.NoteRequest

	// Ambil body request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,dtos.BaseResponse{
			Success: false,
			Message: err.Error(),
		},
	)
	return
	}


	// Ambil payload dari context
	auth, ok := ctx.MustGet("auth").(*dtos.AuthPayload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Success: false,
			Message: "internal server error",
		})
		return
	}

	// Kirim ke service
	note, err := c.service.CreateNote(ctx, req, auth)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dtos.NoteCreatedResponse{
		Title: note.Title,
		Content: note.Content,
	})
}

func (c *NoteController) GetAllNotes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	notes, meta, err := c.service.GetAllNotes(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dtos.BaseResponse{
			Success: false,
			Message: "gagal mengambil note",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NoteListResponse{
		Notes: notes,
		Paginations: meta,
	})
}

func (c *NoteController) GetUserNotes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	auth, ok := ctx.MustGet("auth").(*dtos.AuthPayload)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Success: false,
			Message: "internal server error",
		})
		return
	}


	notes, meta, err := c.service.GetUserNotes(ctx, page, limit, auth)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dtos.BaseResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NoteListResponse{
		Notes: notes,
		Paginations: meta,
	})
}

func (c *NoteController) GetOneNote(ctx *gin.Context) {
	noteID := ctx.Param("id")
	
	note, err := c.service.GetOneNote(ctx, noteID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dtos.BaseResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, dtos.NoteResponse{
		ID: note.ID.String(),
		Title: note.Title,
		Content: note.Content,
		IsPublic: note.IsPublic,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
		User: dtos.UserResponse{
			ID: note.UserID,
			Username: note.User.Username,
			Email: note.User.Email,
			CreatedAt: note.User.CreatedAt,
		},
	})
}

func (c *NoteController) UpdateNote(ctx *gin.Context) {
	var req dtos.UpdateNoteRequest

	// Ambil body request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,dtos.BaseResponse{
			Success: false,
			Message: err.Error(),
		},
	)
	return
}

	// Ambil parameter id
	noteID := ctx.Param("id")

	// update note
	updatedNote, err := c.service.UpdateNote(ctx, noteID, req)
	if err != nil {
		if errors.Is(err, apperror.ErrNoteNotFound) {
			ctx.JSON(http.StatusNotFound, dtos.BaseResponse{
				Success: false,
				Message: "note tidak ditemukan",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Success: false,
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NoteUpdatedResponse{
		NoteCreatedResponse: dtos.NoteCreatedResponse{
			Title: updatedNote.Title,
			Content: updatedNote.Content,
		},
		IsPublic: updatedNote.IsPublic,
	})
}