package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitMiddleware creates a rate limiting middleware
// rate: "10-M" means 10 requests per minute, "100-H" means 100 requests per hour
func RateLimitMiddleware(rate string) gin.HandlerFunc {
	// Create a rate limiter with in-memory store
	store := memory.NewStore()

	// Parse rate string (format: "N-M" where N is limit and M is period)
	rateLimiter := limiter.New(store, limiter.Rate{
		Period: parsePeriod(rate),
		Limit:  parseLimit(rate),
	})

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Get rate limit context
		context, err := rateLimiter.Get(c, clientIP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limit error",
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		// Check if limit exceeded
		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// parsePeriod extracts the period from rate string (e.g., "100-M" -> minute)
func parsePeriod(rate string) time.Duration {
	// Simple parser for common formats
	// "100-M" = 100 requests per minute
	// "1000-H" = 1000 requests per hour
	// "10000-D" = 10000 requests per day
	if len(rate) < 3 {
		return time.Minute
	}

	lastChar := rate[len(rate)-1]
	switch lastChar {
	case 'S', 's':
		return time.Second
	case 'M', 'm':
		return time.Minute
	case 'H', 'h':
		return time.Hour
	case 'D', 'd':
		return 24 * time.Hour
	default:
		return time.Minute
	}
}

// parseLimit extracts the limit number from rate string (e.g., "100-M" -> 100)
func parseLimit(rate string) int64 {
	// Simple parser - extract number before the last character
	if len(rate) < 2 {
		return 100 // Default limit
	}

	// Find where the number ends
	var num int64
	for i := 0; i < len(rate)-1; i++ {
		if rate[i] >= '0' && rate[i] <= '9' {
			num = num*10 + int64(rate[i]-'0')
		} else {
			break
		}
	}

	if num == 0 {
		return 100 // Default limit
	}

	return num
}
