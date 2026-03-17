package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AILoggerMiddleware logs AI requests and records total processing time latency.
func AILoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[AI GATEWAY] %s %s | Status: %d | Latency: %s\n",
			c.Request.Method,
			c.Request.URL.Path,
			status,
			latency,
		)
	}
}

// TimeoutMiddleware sets a deadline on the context for AI requests.
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})

		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			// Request naturally finished
			return
		case <-ctx.Done():
			// Context timed out
			if !c.IsAborted() {
				RespondWithError(c, http.StatusGatewayTimeout, ErrCodeTimeout, "AI provider request timed out")
				c.Abort()
			}
		}
	}
}

// AIRecoveryMiddleware safely recovers from panics during AI operations,
// logging them and returning a clean JSON 500.
func AIRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[AI GATEWAY PANIC] %v\n", err)
				if !c.IsAborted() {
					RespondWithError(c, http.StatusInternalServerError, ErrCodeInternal, "Internal server error during AI processing")
				}
			}
		}()
		c.Next()
	}
}
