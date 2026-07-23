package crawler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// RetryRunner executes tasks with exponential backoff and tracks source health.
type RetryRunner struct {
	maxRetries       int
	initialBackoff   time.Duration
	unhealthySources map[string]time.Time
	mu               sync.RWMutex
}

// NewRetryRunner initializes a RetryRunner.
func NewRetryRunner(maxRetries int, initialBackoff time.Duration) *RetryRunner {
	if maxRetries <= 0 {
		maxRetries = 3
	}
	if initialBackoff <= 0 {
		initialBackoff = 1 * time.Second
	}
	return &RetryRunner{
		maxRetries:       maxRetries,
		initialBackoff:   initialBackoff,
		unhealthySources: make(map[string]time.Time),
	}
}

// IsSourceHealthy checks if a source is currently healthy or cooling down.
func (r *RetryRunner) IsSourceHealthy(sourceName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	failedTime, exists := r.unhealthySources[sourceName]
	if !exists {
		return true
	}

	// Unhealthy status expires after 1 hour cooldown
	if time.Since(failedTime) > 1*time.Hour {
		return true
	}
	return false
}

// ExecuteWithRetry attempts fn up to maxRetries with exponential backoff (1s, 2s, 4s).
func (r *RetryRunner) ExecuteWithRetry(ctx context.Context, sourceName string, fn func(ctx context.Context) error) error {
	if !r.IsSourceHealthy(sourceName) {
		slog.Warn("Skipping source marked as unhealthy", "source", sourceName)
		return fmt.Errorf("source %s is currently unhealthy and skipped", sourceName)
	}

	backoff := r.initialBackoff
	var lastErr error

	for attempt := 1; attempt <= r.maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := fn(ctx)
		if err == nil {
			// Mark as healthy if it succeeded
			r.mu.Lock()
			delete(r.unhealthySources, sourceName)
			r.mu.Unlock()
			return nil
		}

		lastErr = err
		slog.Warn("Source fetch failed, retrying...", "source", sourceName, "attempt", attempt, "max", r.maxRetries, "err", err)

		if attempt < r.maxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff *= 2 // Exponential backoff (1s -> 2s -> 4s)
			}
		}
	}

	// Mark source as unhealthy after all retries fail
	r.mu.Lock()
	r.unhealthySources[sourceName] = time.Now()
	r.mu.Unlock()

	slog.Error("Source marked as unhealthy after max retries exceeded", "source", sourceName, "err", lastErr)
	return fmt.Errorf("source %s failed after %d retries: %w", sourceName, r.maxRetries, lastErr)
}
