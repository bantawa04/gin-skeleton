package modules

import (
	"gin/internal/config"

	"go.uber.org/fx"
)

// ConfigModule provides configuration dependencies
var ConfigModule = fx.Options(
	fx.Provide(config.LoadConfig),
	fx.Provide(config.InitDatabase),
)
