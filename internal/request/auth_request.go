package request

// LoginRequest represents the login request structure
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=8"`
}

// RefreshTokenRequest represents the refresh token request structure
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" validate:"required"`
}
