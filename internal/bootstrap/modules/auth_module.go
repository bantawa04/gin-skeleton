package modules

import (
	"gin/internal/auth/handler"

	"go.uber.org/fx"
)

// AuthModule provides authentication-related dependencies (handler)
// Note: AuthHandler depends on UserService and RefreshTokenService which are provided in other modules
var AuthModule = fx.Options(
	fx.Provide(handler.NewAuthHandler),
)
