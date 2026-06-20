package modules

import (
	"gin/internal/domain/user/handler"
	userRepository "gin/internal/domain/user/repository"
	userService "gin/internal/domain/user/service"

	"go.uber.org/fx"
)

// UserModule provides user-related dependencies (repository, service, handler)
var UserModule = fx.Options(
	fx.Provide(userRepository.NewUserRepository),
	fx.Provide(userService.NewUserService),
	fx.Provide(handler.NewUserHandler),
)
