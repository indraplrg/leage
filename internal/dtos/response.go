package dtos

import "github.com/google/uuid"

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
}

type LoginResponse struct {
	BaseResponse
	Data *LoginData `json:"data,omitempty"`
}

type NoteResponse struct {
	Title string `json:"title"`
	Content string `json:"content"`
}