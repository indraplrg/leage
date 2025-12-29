package controllers

import (
	"net/http"
	"share-notes-app/internal/dtos"
	"share-notes-app/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	service services.AuthenticationService
}

func NewAuthenticationController(service services.AuthenticationService) *AuthenticationController {
	return &AuthenticationController{service: service}
}


func (c *AuthenticationController) Register(ctx *gin.Context) {
	var req dtos.UserRequest

	// ambil request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Success: false,
			Message: err.Error(),
			Data: nil,
		})
		return
	}

	// buat akun
	status, err := c.service.Register(ctx, req)
	if err != nil {
	ctx.JSON(http.StatusConflict, dtos.Response{
			Success: false,
			Message: err.Error(),
			Data: status, 
			
		})
		return	
	}

	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Message: status,
		Data: nil,
	})
}