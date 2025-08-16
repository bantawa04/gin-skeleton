package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in our JWT tokens
type JWTClaims struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTManager creates a new JWT manager instance
func NewJWTManager(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(secretKey),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// NewJWTManagerFromConfig creates a JWT manager from configuration
func NewJWTManagerFromConfig(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return NewJWTManager(secretKey, accessExpiry, refreshExpiry)
}

// GenerateAccessToken generates a short-lived access token
func (j *JWTManager) GenerateAccessToken(userID string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mitho-api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// GenerateRefreshToken generates a long-lived refresh token
func (h *JWTManager) GenerateRefreshToken(userID string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mitho-api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.secretKey)
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetAccessExpiry returns the access token expiry duration
func (j *JWTManager) GetAccessExpiry() time.Duration {
	return j.accessExpiry
}

// GetRefreshExpiry returns the refresh token expiry duration
func (j *JWTManager) GetRefreshExpiry() time.Duration {
	return j.refreshExpiry
}

// RefreshAccessToken generates a new access token using a valid refresh token
func (j *JWTManager) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Ensure this is a refresh token
	if claims.Type != "refresh" {
		return "", errors.New("invalid token type, refresh token required")
	}

	// Generate new access token
	return j.GenerateAccessToken(claims.UserID)
}
