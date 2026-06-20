package refreshtoken

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// RefreshToken represents a refresh token in the system
type RefreshToken struct {
	ID        string     `json:"id" gorm:"primaryKey;type:char(26)"`
	UserID    string     `json:"user_id" gorm:"type:char(26);not null;index"`
	Token     string     `json:"-" gorm:"type:text;not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index"`
	Revoked   bool       `json:"revoked" gorm:"default:false;index"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// BeforeCreate hook for generating ID
func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == "" {
		// Generate a new ULID
		id := ulid.Make()
		rt.ID = id.String()
	}
	return nil
}

// IsExpired checks if the refresh token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid checks if the refresh token is valid (not revoked and not expired)
func (rt *RefreshToken) IsValid() bool {
	return !rt.Revoked && !rt.IsExpired()
}

// Revoke marks the refresh token as revoked
func (rt *RefreshToken) Revoke() {
	rt.Revoked = true
	now := time.Now()
	rt.RevokedAt = &now
}

// TableName specifies the table name for the RefreshToken model
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
