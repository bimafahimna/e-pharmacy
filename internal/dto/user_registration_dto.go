package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type VerifyRequest struct {
	Password string `json:"password" binding:"required,min=8"`
	Token    string `json:"token" binding:"required"`
}

type GoogleRegisterRequest struct {
	Credential string `json:"credential" binding:"required"`
}

type GoogleRegisterCallback struct {
	Email string `json:"email" binding:"required,email"`
}
