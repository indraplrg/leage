package dtos

import (
	"time"

	"github.com/google/uuid"
)

// request
type BaseRequest struct {
	Username string `json:"username" binding:"required,min=6"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserRegisterRequest struct {
	BaseRequest
	Email string `json:"email" binding:"required,email"`
}

type UserLoginRequest struct  {
	BaseRequest
}

// response
type BaseResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

type RegisterData struct {
	ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
}

type RegisterResponse struct {
	BaseResponse
	Data *RegisterData `json:"data,omitempty"`
}

type LoginData struct {
	AccessToken string `json:"access_token"`
	RefreshToken string
}

type LoginResponse struct {
	BaseResponse
	Data *LoginData `json:"data,omitempty"`
}

type UserResponse struct {
	ID uuid.UUID `json:"id"`
	Username string	`json:"username"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

