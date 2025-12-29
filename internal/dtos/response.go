package dtos

import "github.com/google/uuid"

type Response struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{}  `json:"data,omitempty"`
}

type RegisterResponse struct {
	ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
}