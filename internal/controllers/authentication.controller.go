package controllers

import (
	"net/http"
	"share-notes-app/internal/dtos"
	"share-notes-app/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		ctx.JSON(http.StatusBadRequest, dtos.LoginResponse{
			BaseResponse: dtos.BaseResponse{
				Success: false,
				Message: err.Error(),
			},
			Data: nil,
		})
		return 
	}

	ctx.SetCookieData(&http.Cookie{
		Name : "paseto_token", 
		Value: token, 
		Expires: time.Now().Add(2 * time.Hour),
		MaxAge: 7200, 
		Path: "/", 
		Domain: "localhost", 
		Secure: false, 
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	ctx.JSON(http.StatusOK, dtos.LoginResponse{
		BaseResponse: dtos.BaseResponse{
			Success: true,
			Message: "login successfully",
		},
		Data: &dtos.LoginData{
			AccessToken: token,
		},
	})
}

func (c *AuthenticationController) Logout(ctx *gin.Context) {
	
	// Ambil payload dari context
	auth, ok := ctx.MustGet("auth").(*dtos.AuthPayload)
	logrus.Info(auth, ok)
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

		ctx.SetCookieData(&http.Cookie{
		Name : "paseto_token", 
		Value: "", 
		Expires: time.Now().Add(1 * time.Second),
		MaxAge: -1, 
		Path: "/", 
		Domain: "localhost", 
		Secure: false, 
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	ctx.JSON(http.StatusOK, dtos.BaseResponse{
		Success: true,
		Message: "Logout successfully",
	})
}