package controllers

import (
	"net/http"
	"share-notes-app/internal/dtos"
	"share-notes-app/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthorizationController struct {
	service services.AuthorizationService
}

func NewAuthorizationsController(service services.AuthorizationService) *AuthorizationController {
	return &AuthorizationController{service:service}
}

func (c *AuthorizationController) VerifyEmail(ctx *gin.Context) {
	token := ctx.Param("token")
	logrus.Info(token, "ini token")
	
	ok, err := c.service.VerifyEmail(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusNotFound, &dtos.BaseResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &dtos.BaseResponse{
		Success: true,
		Message: ok,
	})
}