package dto

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ConfirmResetRequest struct {
	Password string `json:"password" binding:"required,min=8"`
	Token    string `json:"token" binding:"required"`
}
