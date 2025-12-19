package refreshtoken

import (
	"context"
	tokenmodel "gin/internal/refresh_token"

	"gorm.io/gorm"
)

// RefreshTokenRepository handles refresh token database operations
type RefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(ctx context.Context, token *tokenmodel.RefreshToken) (*tokenmodel.RefreshToken, error) {
	err := r.db.WithContext(ctx).Create(token).Error
	if err != nil {
		return nil, err
	}
	return token, nil
}

// FindByToken finds a refresh token by token string
func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*tokenmodel.RefreshToken, error) {
	var refreshToken tokenmodel.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

// Revoke revokes a refresh token by setting revoked flag
func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenID string) error {
	return r.db.WithContext(ctx).Model(&tokenmodel.RefreshToken{}).
		Where("id = ?", tokenID).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

// RevokeByToken revokes a refresh token by token string
func (r *RefreshTokenRepository) RevokeByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&tokenmodel.RefreshToken{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (r *RefreshTokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&tokenmodel.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

// DeleteExpiredTokens deletes expired refresh tokens (cleanup job)
func (r *RefreshTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", gorm.Expr("CURRENT_TIMESTAMP")).
		Delete(&tokenmodel.RefreshToken{}).Error
}

// FindByUserID finds all refresh tokens for a user
func (r *RefreshTokenRepository) FindByUserID(ctx context.Context, userID string) ([]*tokenmodel.RefreshToken, error) {
	var tokens []*tokenmodel.RefreshToken
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}
