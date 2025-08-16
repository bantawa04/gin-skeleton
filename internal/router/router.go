package router

import (
	"fmt"
	exceptions "gin/internal/api/exception"
	handlers "gin/internal/api/handler"
	middleware "gin/internal/middleware"
	response "gin/internal/response"
	"gin/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type Router interface {
	Run(addr ...string) error
}

func NewRouter(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	jwtManager *utils.JWTManager,
) *gin.Engine {
	router := gin.Default()

	// Add middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CaseConverterMiddleware())
	router.Use(exceptions.ErrorHandler())

	// Add custom logger middleware
	router.Use(func(c *gin.Context) {
		// Start timer
		t := time.Now()

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(t)
		status := c.Writer.Status()
		fmt.Printf("[GIN] %s | %3d | %s | %s\n",
			c.Request.Method,
			status,
			latency,
			c.Request.URL.Path,
		)
	})

	router.NoRoute(func(c *gin.Context) {
		desc := "The requested endpoint does not exist"
		appErr := exceptions.NotFoundError("Route not found", &desc)
		_ = c.Error(appErr)
	})

	// Register routes
	router.GET("/ping", func(c *gin.Context) {
		response.SendResponse(c, "pong", "pong")
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.POST("", userHandler.CreateUser)

			// Protected routes - require JWT authentication
			protected := users.Group("/")
			protected.Use(middleware.JWTAuthMiddleware(jwtManager))
			{
				protected.PUT("/:id", userHandler.UpdateUser)
				protected.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}
	return router
}
