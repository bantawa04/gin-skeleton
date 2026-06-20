package router

import (
	middleware "gin/internal/infra/middleware"

	"github.com/gin-gonic/gin"
)

// registerWebRoutes wires the default web/API route surface under /api.
func registerWebRoutes(api *gin.RouterGroup, d *routerDeps) {
	// Health check endpoint should stay fast and unauthenticated.
	api.GET("/health", d.healthHandler.Health)

	auth := api.Group("/auth")
	auth.Use(middleware.RateLimitMiddleware("10-M"))
	{
		auth.POST("/signup", middleware.TransactionMiddleware(d.db), d.authHandler.Signup)
		auth.POST("/login", middleware.TransactionMiddleware(d.db), d.authHandler.Login)
		auth.POST("/refresh", middleware.TransactionMiddleware(d.db), d.authHandler.RefreshToken)
		auth.POST("/logout", middleware.JWTAuthMiddleware(d.jwtManager), middleware.TransactionMiddleware(d.db), d.authHandler.Logout)
	}

	users := api.Group("/users")
	{
		users.GET("", d.userHandler.GetAllUsers)
		users.GET("/:id", d.userHandler.GetUserByID)

		protected := users.Group("/")
		protected.Use(middleware.JWTAuthMiddleware(d.jwtManager))
		{
			protected.PUT("/:id", middleware.TransactionMiddleware(d.db), d.userHandler.UpdateUser)
			protected.DELETE("/:id", middleware.TransactionMiddleware(d.db), d.userHandler.DeleteUser)
		}
	}
}
