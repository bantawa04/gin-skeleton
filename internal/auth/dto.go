package auth

import "gin/internal/user"

// LoginResponseDTO represents the login response with tokens
type LoginResponseDTO struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	TokenType    string       `json:"tokenType"`
	ExpiresIn    int64        `json:"expiresIn"`
	User         user.UserDTO `json:"user"`
}

// RefreshTokenResponseDTO represents the refresh token response
type RefreshTokenResponseDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"` // New refresh token (rotation)
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64  `json:"expiresIn"`
}
