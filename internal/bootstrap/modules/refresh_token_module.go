package modules

import (
	refreshTokenRepository "gin/internal/refresh_token/repository"
	refreshTokenService "gin/internal/refresh_token/service"

	"go.uber.org/fx"
)

// RefreshTokenModule provides refresh token-related dependencies (repository, service)
var RefreshTokenModule = fx.Options(
	fx.Provide(refreshTokenRepository.NewRefreshTokenRepository),
	fx.Provide(refreshTokenService.NewRefreshTokenService),
)
