package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token      string `json:"token,omitempty"`
	Role       string `json:"role"`
	IsVerified string `json:"is_verified"`
	UserID     int64  `json:"user_id"`
}

type ReturnedURL struct {
	ReturnedURL string `json:"url,omitempty"`
}

type GoogleLoginRequest struct {
	Credential string `json:"credential" binding:"required"`
}

type GoogleLoginCallback struct {
	Email string `json:"email" binding:"required,email"`
}
