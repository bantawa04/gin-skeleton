package response

// LoginResponse represents the login response structure
type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    int64         `json:"expires_in"` // Access token expiry in seconds
	User         *UserResponse `json:"user"`
}

// RefreshTokenResponse represents the refresh token response structure
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"` // Access token expiry in seconds
}

// TokenInfo represents token information
type TokenInfo struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	Exp    int64  `json:"exp"`
}
