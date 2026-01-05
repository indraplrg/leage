package controllers

import (
	"net/http"
	"share-notes-app/internal/dtos"
	"share-notes-app/internal/services"

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

	ctx.JSON(http.StatusCreated, dtos.NoteResponse{
		Title: note.Title,
		Content: note.Content,
	})
}