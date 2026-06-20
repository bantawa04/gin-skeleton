package router

import (
	"crypto/subtle"
	_ "gin/docs" // Swagger documentation
	authhandler "gin/internal/domain/auth/handler"
	healthhandler "gin/internal/domain/health/handler"
	"gin/internal/infra/config"
	middleware "gin/internal/infra/middleware"
	exceptions "gin/internal/shared/exception"
	response "gin/internal/shared/response"
	"gin/internal/shared/utils"
	userhandler "gin/internal/domain/user/handler"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Router interface {
	Run(addr ...string) error
}

type routerDeps struct {
	userHandler   *userhandler.UserHandler
	authHandler   *authhandler.AuthHandler
	healthHandler *healthhandler.HealthHandler
	jwtManager    *utils.JWTManager
	db            *gorm.DB
}

func NewRouter(
	userHandler *userhandler.UserHandler,
	authHandler *authhandler.AuthHandler,
	healthHandler *healthhandler.HealthHandler,
	jwtManager *utils.JWTManager,
	cfg *config.Config,
	db *gorm.DB,
) *gin.Engine {
	router := gin.Default()

	// Add global middleware (order matters)
	router.Use(middleware.CORSMiddleware(cfg.CORS().AllowedOrigins)) // CORS should be first
	router.Use(middleware.RequestIDMiddleware())     // Request ID for tracing
	router.Use(middleware.LoggingMiddleware())       // Structured logging
	router.Use(middleware.SanitizeMiddleware())      // Input sanitization (XSS prevention)
	router.Use(middleware.CaseConverterMiddleware()) // Case conversion
	router.Use(exceptions.ErrorHandler())            // Error handling

	router.NoRoute(func(c *gin.Context) {
		desc := "The requested endpoint does not exist"
		appErr := exceptions.NotFoundError("Route not found", &desc)
		_ = c.Error(appErr)
	})

	// Register routes
	router.GET("/ping", func(c *gin.Context) {
		response.SendResponse(c, "pong", "pong")
	})

	router.GET("/health", healthHandler.Health)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Gin Skeleton API is running",
			"health":  "/api/health",
			"docs":    "/swagger/index.html",
		})
	})

	registerSwaggerRoutes(router, cfg)

	// API routes
	api := router.Group("/api")
	deps := &routerDeps{
		userHandler:   userHandler,
		authHandler:   authHandler,
		healthHandler: healthHandler,
		jwtManager:    jwtManager,
		db:            db,
	}

	registerWebRoutes(api, deps)
	return router
}

func registerSwaggerRoutes(router *gin.Engine, cfg *config.Config) {
	swaggerAuth := func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		expected := cfg.Swagger()

		if !ok ||
			subtle.ConstantTimeCompare([]byte(username), []byte(expected.Username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(password), []byte(expected.Password)) != 1 {
			c.Header("WWW-Authenticate", `Basic realm="Gin Skeleton Swagger"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}

	router.GET("/docs/swagger.yaml", swaggerAuth, func(c *gin.Context) {
		spec, err := os.ReadFile("docs/swagger.yaml")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "spec unavailable"})
			return
		}
		host := c.Request.Host
		if host == "" {
			host = "localhost:" + cfg.Server().Port
		}
		content := strings.Replace(string(spec), "host: localhost", "host: "+host, 1)
		c.Header("Content-Type", "application/yaml")
		c.String(http.StatusOK, content)
	})
	router.GET("/swagger/*any", swaggerAuth, ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.yaml"), ginSwagger.PersistAuthorization(true)))
}
