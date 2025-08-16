package bootstrap

import (
	"context"
	"gin/internal/api/handler"
	"gin/internal/config"
	"gin/internal/logger"
	userRepository "gin/internal/repository/user"
	"gin/internal/router"
	userService "gin/internal/service/user"
	"gin/internal/utils"
	validators "gin/internal/validator"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module exported for initializing application
var Module = fx.Options(
	ConfigModule,
	RepositoryModule,
	ServiceModule,
	ValidatorModule,
	JWTModule,
	HandlerModule,
	RouterModule,
	fx.Invoke(bootstrap),
)

// ConfigModule provides configuration dependencies
var ConfigModule = fx.Options(
	fx.Provide(config.LoadConfig),
	fx.Provide(config.InitDatabase),
)

// RepositoryModule provides repository dependencies
var RepositoryModule = fx.Options(
	fx.Provide(userRepository.NewUserRepository),
)

// ServiceModule provides service dependencies
var ServiceModule = fx.Options(
	fx.Provide(userService.NewUserService),
)

// ValidatorModule provides validator dependencies
var ValidatorModule = fx.Options(
	fx.Provide(validators.NewValidator),
)

// HandlerModule provides handler dependencies
var HandlerModule = fx.Options(
	fx.Provide(handler.NewUserHandler),
	fx.Provide(handler.NewAuthHandler),
)

// JWTModule provides JWT manager dependencies
var JWTModule = fx.Options(
	fx.Provide(newJWTManager),
)

// RouterModule provides router dependencies
var RouterModule = fx.Options(
	fx.Provide(router.NewRouter),
	fx.Provide(newHTTPServer),
)

// BuildApp constructs the fx application with all dependencies
func BuildApp() *fx.App {
	return fx.New(Module)
}

// newJWTManager creates a JWT manager with configuration
func newJWTManager(cfg *config.Config) *utils.JWTManager {
	jwtConfig := cfg.JWT()
	return utils.NewJWTManagerFromConfig(
		jwtConfig.SecretKey,
		jwtConfig.AccessExpiry,
		jwtConfig.RefreshExpiry,
	)
}

// newHTTPServer creates an HTTP server with the provided router and configuration
func newHTTPServer(router *gin.Engine, cfg *config.Config) *http.Server {
	serverConfig := cfg.Server()
	return &http.Server{
		Addr:         ":" + serverConfig.Port,
		Handler:      router,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
	}
}

// bootstrap registers lifecycle hooks and initializes the application
func bootstrap(
	lifecycle fx.Lifecycle,
	router *gin.Engine,
	server *http.Server,
	cfg *config.Config,
	db *gorm.DB,
) {
	// Register server start and stop hooks
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Initialize logger
			if err := logger.Init(); err != nil {
				log.Printf("Failed to initialize logger: %v", err)
				// Continue without logger rather than failing
			} else {
				logger.LogInfo("Logger initialized successfully", nil)
			}

			log.Println("Starting Application")
			log.Println("------------------------")
			log.Println("-- Gin API --")
			log.Println("------------------------")

			// Start the server in a goroutine
			go func() {
				serverConfig := cfg.Server()
				log.Printf("Server starting on port %s", serverConfig.Port)
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")

			// Close database connection
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				_ = sqlDB.Close()
			}

			return server.Shutdown(ctx)
		},
	})
}
