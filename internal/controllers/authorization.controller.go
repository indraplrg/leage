package controllers

import (
	"share-notes-app/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthorizationController struct {
	service services.AuthorizationService
}

func NewAuthorizationsController(service services.AuthorizationService) *AuthorizationController {
	return &AuthorizationController{service:service}
}

func (s *AuthorizationController) VerifyEmail(ctx *gin.Context) {}