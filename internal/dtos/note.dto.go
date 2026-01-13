package dtos

import (
	"time"
)

// request
type NoteRequest struct {
	Title string `json:"title" binding:"required,min=6"`
	Content string `json:"content" binding:"required,min=12"` 
}

type UpdateNoteRequest struct {
	NoteRequest
}

// response
type NoteCreatedResponse struct {
	Title string `json:"title"`
	Content string `json:"content"`
}

type NoteUpdatedResponse struct {
	NoteCreatedResponse
}

type NoteResponse struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	IsPublic bool `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User UserResponse `json:"user"`
}

type NoteListResponse struct {
	Notes []NoteResponse `json:"data"`
	Paginations PaginationMeta `json:"meta"`
}
