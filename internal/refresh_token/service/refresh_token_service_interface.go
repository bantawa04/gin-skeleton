package refreshtoken

import (
	"context"
	tokenmodel "gin/internal/refresh_token"
)

type RefreshTokenServiceInterface interface {
	Create(ctx context.Context, token *tokenmodel.RefreshToken) (*tokenmodel.RefreshToken, error)
	FindByToken(ctx context.Context, token string) (*tokenmodel.RefreshToken, error)
	RevokeByToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
}
