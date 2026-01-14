package controllers

import (
	"net/http"
	"share-notes-app/internal/dtos"
	"share-notes-app/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type NoteController struct {
	service services.NoteService
}

func NewNoteController(service services.NoteService) *NoteController {
	return &NoteController{service:service}
}

func (c *NoteController) CreateNote(ctx *gin.Context) {
	var req dtos.NoteRequest
	logrus.Info(req)

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