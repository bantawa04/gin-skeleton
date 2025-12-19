package modules

import (
	"gin/internal/health/handler"

	"go.uber.org/fx"
)

// HealthModule provides health check handler
var HealthModule = fx.Options(
	fx.Provide(handler.NewHealthHandler),
)
