package router

import (
	_ "gin/docs" // Swagger documentation
	authhandler "gin/internal/auth/handler"
	exceptions "gin/internal/exception"
	healthhandler "gin/internal/health/handler"
	middleware "gin/internal/middleware"
	response "gin/internal/response"
	userhandler "gin/internal/user/handler"
	"gin/internal/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Router interface {
	Run(addr ...string) error
}

func NewRouter(
	userHandler *userhandler.UserHandler,
	authHandler *authhandler.AuthHandler,
	healthHandler *healthhandler.HealthHandler,
	jwtManager *utils.JWTManager,
	db *gorm.DB,
) *gin.Engine {
	router := gin.Default()

	// Add global middleware (order matters)
	router.Use(middleware.CORSMiddleware())          // CORS should be first
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

	// Health check endpoint (no middleware needed, should be fast)
	router.GET("/health", healthHandler.Health)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public) - with rate limiting and transaction
		auth := api.Group("/auth")
		auth.Use(middleware.RateLimitMiddleware("10-M")) // 10 requests per minute for auth endpoints
		{
			auth.POST("/signup", middleware.TransactionMiddleware(db), authHandler.Signup)
			// Login creates tokens, use transaction
			auth.POST("/login", middleware.TransactionMiddleware(db), authHandler.Login)
			// Refresh token updates tokens, use transaction
			auth.POST("/refresh", middleware.TransactionMiddleware(db), authHandler.RefreshToken)

			// Logout requires authentication - user must be logged in to logout
			auth.POST("/logout", middleware.JWTAuthMiddleware(jwtManager), middleware.TransactionMiddleware(db), authHandler.Logout)
		}

		// User routes - with transaction middleware for write operations
		users := api.Group("/users")
		{
			// Read operations (no transaction needed)
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)

			// Write operations (with transaction)
			// users.POST("", middleware.TransactionMiddleware(db), userHandler.CreateUser)

			// Protected routes - require JWT authentication
			protected := users.Group("/")
			protected.Use(middleware.JWTAuthMiddleware(jwtManager))
			{
				// Write operations (with transaction)
				protected.PUT("/:id", middleware.TransactionMiddleware(db), userHandler.UpdateUser)
				protected.DELETE("/:id", middleware.TransactionMiddleware(db), userHandler.DeleteUser)
			}
		}
	}
	return router
}
