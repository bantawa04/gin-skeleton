package middlewares

import (
	"context"
	"gin/internal/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const TransactionKey = "db_transaction"

// TransactionMiddleware creates a middleware that wraps requests in a database transaction
// The transaction is committed on success or rolled back on error
// The transaction is stored in both Gin context and request context for repository access
func TransactionMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start transaction
		tx := db.Begin()
		if tx.Error != nil {
			logger.LogError(tx.Error, "Failed to start transaction", nil)
			c.AbortWithStatusJSON(500, gin.H{
				"error": "Failed to start database transaction",
			})
			return
		}

		// Store transaction in Gin context (for middleware access)
		c.Set(TransactionKey, tx)

		// Store transaction in request context (for repository access)
		ctx := context.WithValue(c.Request.Context(), TransactionKey, tx)
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Check if there were any errors
		if len(c.Errors) > 0 || c.Writer.Status() >= 400 {
			// Rollback transaction on error
			if err := tx.Rollback().Error; err != nil {
				logger.LogError(err, "Failed to rollback transaction", nil)
			} else {
				logger.LogInfo("Transaction rolled back due to error", map[string]interface{}{
					"path":   c.Request.URL.Path,
					"status": c.Writer.Status(),
				})
			}
		} else {
			// Commit transaction on success
			if err := tx.Commit().Error; err != nil {
				logger.LogError(err, "Failed to commit transaction", nil)
				c.AbortWithStatusJSON(500, gin.H{
					"error": "Failed to commit database transaction",
				})
			}
		}
	}
}

// GetTransaction retrieves the database transaction from the Gin context
// Returns the transaction if available, otherwise returns the original db connection
func GetTransaction(c *gin.Context, db *gorm.DB) *gorm.DB {
	if tx, exists := c.Get(TransactionKey); exists {
		if transaction, ok := tx.(*gorm.DB); ok {
			return transaction
		}
	}
	return db
}
