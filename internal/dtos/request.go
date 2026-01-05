package dtos

type BaseRequest struct {
	Username string `json:"username" binding:"required,min=6"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserRequest struct {
	BaseRequest
	Email string `json:"email" binding:"required,email"`
}

type LoginRequest struct  {
	BaseRequest
}

type NoteRequest struct {
	Title string `json:"title" binding:"required,min=6"`
	Content string `json:"content" binding:"required,min=12"` 
}

type ResendTokenRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type AuthPayload struct {
	UserID string 
	Username string
}