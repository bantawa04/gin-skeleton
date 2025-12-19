package middlewares

import (
	exception "gin/internal/exception"
	"gin/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware validates JWT tokens and extracts user information
func JWTAuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appErr := exception.UnauthorizedError("Authorization header is required", nil, nil)
			_ = c.Error(appErr)
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			appErr := exception.UnauthorizedError("Invalid authorization header format. Expected 'Bearer <token>'", nil, nil)
			_ = c.Error(appErr)
			c.Abort()
			return
		}

		// Extract the token (remove "Bearer " prefix)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			appErr := exception.UnauthorizedError("Token is required", nil, nil)
			_ = c.Error(appErr)
			c.Abort()
			return
		}

		// Validate the token using the injected JWT manager
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			appErr := exception.UnauthorizedError("Invalid or expired token", nil, nil)
			_ = c.Error(appErr)
			c.Abort()
			return
		}

		// Ensure this is an access token
		if claims.Type != "access" {
			appErr := exception.UnauthorizedError("Invalid token type. Access token required", nil, nil)
			_ = c.Error(appErr)
			c.Abort()
			return
		}

		// Set user ID in context for later use
		c.Set("user_id", claims.UserID)
		c.Set("user_claims", claims)

		// Continue to the next handler
		c.Next()
	}
}
