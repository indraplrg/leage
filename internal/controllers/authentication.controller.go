package controllers

import (
	"net/http"
	"share-notes-app/helper"
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
		ctx.JSON(http.StatusBadRequest, dtos.RegisterResponse{
		BaseResponse: dtos.BaseResponse{
			Success: false,
			Message: err.Error(),
		},
		Data: nil,
		})
		return
	}

	// buat akun
	user, err := c.service.Register(ctx, req)
	if err != nil {
	ctx.JSON(http.StatusConflict, dtos.RegisterResponse{
		BaseResponse: dtos.BaseResponse{
			Success: false,
			Message: err.Error(),
		},
		Data: nil,
	})
		return	
	}

	ctx.JSON(http.StatusCreated, dtos.RegisterResponse{
		BaseResponse: dtos.BaseResponse{
			Success: true,
			Message: "register successfully",
		},
		Data: &dtos.RegisterData{
			ID: user.ID,
			Username: user.Username,
			Email: user.Email,
		},
	})
}

func (c *AuthenticationController) Login(ctx *gin.Context) {
	var req dtos.LoginRequest
	
	// ambil request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.LoginResponse{
			BaseResponse: dtos.BaseResponse{
				Success: false,
				Message: err.Error(),
			},
			Data: nil,
		})
		return
	}

	// login
	token, err := c.service.Login(ctx, req); if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.BaseResponse{
				Success: false,
				Message: err.Error(),
			},
		)
		return 
	}

	// set cookie
	helper.SetCookie(ctx, "access_paseto_token", token.AccessToken, 30*60)
	helper.SetCookie(ctx, "refresh_paseto_token", token.RefreshToken, 168*3600)

	ctx.JSON(http.StatusOK, dtos.LoginResponse{
		BaseResponse: dtos.BaseResponse{
			Success: true,
			Message: "login successfully",
		},
		Data: &dtos.LoginData{
			AccessToken: token.AccessToken,
		},
	})
}

func (c *AuthenticationController) Logout(ctx *gin.Context) {
	
	// Ambil payload dari context
	auth, ok := ctx.MustGet("auth").(*dtos.AuthPayload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Success: false,
			Message: "internal server error",
		})
		return
	} 
	
	err := c.service.Logout(ctx, auth)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dtos.BaseResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// hapus cookie
	helper.DeleteCookie(ctx, "access_paseto_token")
	helper.DeleteCookie(ctx, "refresh_paseto_token")

	ctx.JSON(http.StatusOK, dtos.BaseResponse{
		Success: true,
		Message: "Logout successfully",
	})
}

func (c *AuthenticationController) ResendToken(ctx *gin.Context) {
	var req dtos.ResendTokenRequest
	
	// ambil request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.BaseResponse{
				Success: false,
				Message: err.Error(),
			},
		)
		return
	}

	
}

func (c *AuthenticationController) VerifyEmail(ctx *gin.Context) {
	token := ctx.Param("token")
	
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