// Package api provides HTTP routing, middleware, and Server-Sent Events handling.
package api

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS returns a middleware handling Cross-Origin Resource Sharing.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Logger logs request details cleanly using slog.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		slog.Info("HTTP Request", "method", c.Request.Method, "path", path, "status", status, "latency", latency)
	}
}

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
			return
		case <-ctx.Done():
			if !c.IsAborted() {
				RespondWithError(c, http.StatusGatewayTimeout, ErrCodeTimeout, "AI provider request timed out")
				c.Abort()
			}
		}
	}
}

// AIRecoveryMiddleware safely recovers from panics during AI operations.
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
