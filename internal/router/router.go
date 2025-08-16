package router

import (
	"fmt"
	exceptions "gin/internal/api/exception"
	handlers "gin/internal/api/handler"
	"time"

	"github.com/gin-gonic/gin"
)

// Router interface defines methods for the router
type Router interface {
	Run(addr ...string) error
}

// NewRouter creates and configures a Gin router
// Update the NewRouter function to include brand and category routes
// NewRouter creates and configures a new router
func NewRouter(

	userHandler *handlers.UserHandler,
) *gin.Engine {
	router := gin.Default()

	// Add middleware
	// router.Use(exceptions.CaseConverterMiddleware())
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

	// Handle 404 Not Found
	router.NoRoute(func(c *gin.Context) {
		// Add error to context instead of handling directly
		// appErr := exceptions.NewNotFoundError("Route not found", "The requested endpoint does not exist")
		// _ = c.Error(appErr)
	})

	// Register routes
	// router.GET("/ping", healthHandler.HealthCheck)

	// API routes
	api := router.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

	}
	return router
}
