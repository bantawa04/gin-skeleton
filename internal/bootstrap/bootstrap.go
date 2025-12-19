package bootstrap

import (
	"context"
	"gin/internal/bootstrap/modules"
	"gin/internal/config"
	"gin/internal/logger"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module exported for initializing application
// Composes all domain modules in the correct order
var Module = fx.Options(
	// Core infrastructure modules (must be first)
	modules.ConfigModule,
	modules.UtilsModule,

	// Domain modules
	modules.UserModule,
	modules.RefreshTokenModule,
	modules.AuthModule,
	modules.HealthModule,

	// Application modules (must be last)
	modules.RouterModule,

	// Bootstrap lifecycle
	fx.Invoke(bootstrap),
)

// BuildApp constructs the fx application with all dependencies
func BuildApp() *fx.App {
	return fx.New(Module)
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
			sqlDB, err := db.DB()
			if err != nil {
				log.Printf("Failed to get database connection for shutdown: %v", err)
			} else if sqlDB != nil {
				if err := sqlDB.Close(); err != nil {
					log.Printf("Failed to close database connection: %v", err)
				} else {
					log.Println("Database connection closed successfully")
				}
			}

			// Shutdown server with context timeout
			if err := server.Shutdown(ctx); err != nil {
				log.Printf("Server forced to shutdown: %v", err)
				return err
			}

			log.Println("Server shutdown completed")
			return nil
		},
	})
}
