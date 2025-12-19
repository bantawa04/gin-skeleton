package refreshtoken

import (
	"context"
	tokenmodel "gin/internal/refresh_token"
	refreshTokenRepository "gin/internal/refresh_token/repository"
)

type RefreshTokenService struct {
	refreshTokenRepo *refreshTokenRepository.RefreshTokenRepository
}

func NewRefreshTokenService(refreshTokenRepo *refreshTokenRepository.RefreshTokenRepository) RefreshTokenServiceInterface {
	return &RefreshTokenService{
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *RefreshTokenService) Create(ctx context.Context, token *tokenmodel.RefreshToken) (*tokenmodel.RefreshToken, error) {
	return s.refreshTokenRepo.Create(ctx, token)
}

func (s *RefreshTokenService) FindByToken(ctx context.Context, token string) (*tokenmodel.RefreshToken, error) {
	return s.refreshTokenRepo.FindByToken(ctx, token)
}

func (s *RefreshTokenService) RevokeByToken(ctx context.Context, token string) error {
	return s.refreshTokenRepo.RevokeByToken(ctx, token)
}

func (s *RefreshTokenService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return s.refreshTokenRepo.RevokeAllUserTokens(ctx, userID)
}
