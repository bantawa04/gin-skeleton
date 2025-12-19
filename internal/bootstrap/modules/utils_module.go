package modules

import (
	"gin/internal/config"
	"gin/internal/utils"
	validators "gin/internal/validator"

	"go.uber.org/fx"
)

// UtilsModule provides utility dependencies (JWT manager, validator)
var UtilsModule = fx.Options(
	fx.Provide(newJWTManager),
	fx.Provide(validators.NewValidator),
	fx.Provide(utils.GeneratePassword),
)

// newJWTManager creates a JWT manager with configuration
func newJWTManager(cfg *config.Config) *utils.JWTManager {
	jwtConfig := cfg.JWT()
	return utils.NewJWTManagerFromConfig(
		jwtConfig.SecretKey,
		jwtConfig.AccessExpiry,
		jwtConfig.RefreshExpiry,
	)
}
