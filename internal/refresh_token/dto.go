package refreshtoken

import (
	"time"
)

// RefreshTokenDTO represents the data transfer object for RefreshToken
type RefreshTokenDTO struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	ExpiresAt time.Time  `json:"expiresAt"`
	Revoked   bool       `json:"revoked"`
	RevokedAt *time.Time `json:"revokedAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// FromRefreshTokenModel converts a RefreshToken model to a RefreshTokenDTO
func FromRefreshTokenModel(token RefreshToken) RefreshTokenDTO {
	return RefreshTokenDTO{
		ID:        token.ID,
		UserID:    token.UserID,
		ExpiresAt: token.ExpiresAt,
		Revoked:   token.Revoked,
		RevokedAt: token.RevokedAt,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}
}

// ToRefreshTokenModel converts a RefreshTokenDTO to a RefreshToken model
func (dto RefreshTokenDTO) ToRefreshTokenModel() RefreshToken {
	return RefreshToken{
		ID:        dto.ID,
		UserID:    dto.UserID,
		ExpiresAt: dto.ExpiresAt,
		Revoked:   dto.Revoked,
		RevokedAt: dto.RevokedAt,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
