package dtos

type ResendTokenRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type AuthPayload struct {
	UserID string 
	Username string
}