// Package api provides HTTP routing, middleware, and Server-Sent Events handling.
package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
)

// ProgressBroadcaster manages active SSE client connections and broadcasts progress events.
type ProgressBroadcaster struct {
	clients map[chan models.CrawlProgressEvent]bool
	mu      sync.RWMutex
}

// NewProgressBroadcaster creates a ProgressBroadcaster instance.
func NewProgressBroadcaster() *ProgressBroadcaster {
	return &ProgressBroadcaster{
		clients: make(map[chan models.CrawlProgressEvent]bool),
	}
}

// Subscribe registers a new SSE client channel.
func (b *ProgressBroadcaster) Subscribe() chan models.CrawlProgressEvent {
	ch := make(chan models.CrawlProgressEvent, 10)
	b.mu.Lock()
	b.clients[ch] = true
	b.mu.Unlock()
	slog.Info("New SSE client subscribed to crawl progress")
	return ch
}

// Unsubscribe removes an SSE client channel.
func (b *ProgressBroadcaster) Unsubscribe(ch chan models.CrawlProgressEvent) {
	b.mu.Lock()
	if _, exists := b.clients[ch]; exists {
		delete(b.clients, ch)
		close(ch)
	}
	b.mu.Unlock()
	slog.Info("SSE client unsubscribed from crawl progress")
}

// Broadcast sends a progress event payload to all active clients.
func (b *ProgressBroadcaster) Broadcast(event models.CrawlProgressEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.clients {
		select {
		case ch <- event:
		default:
			// Drop event if client buffer is full to avoid blocking
		}
	}
}

// HandleSSEProgress returns a Gin HandlerFunc streaming SSE updates to clients.
func (b *ProgressBroadcaster) HandleSSEProgress() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		clientChan := b.Subscribe()
		defer b.Unsubscribe(clientChan)

		// Send initial connection event
		initData, _ := json.Marshal(models.CrawlProgressEvent{
			Step:       1,
			Message:    "Connected to crawl progress stream",
			Progress:   0,
			IsComplete: false,
		})
		c.SSEvent("message", string(initData))
		c.Writer.Flush()

		ctx := c.Request.Context()
		c.Stream(func(w io.Writer) bool {
			select {
			case <-ctx.Done():
				return false
			case event, ok := <-clientChan:
				if !ok {
					return false
				}
				dataBytes, err := json.Marshal(event)
				if err != nil {
					slog.Error("Failed to marshal SSE progress event", "err", err)
					return true
				}
				c.SSEvent("progress", string(dataBytes))
				c.Writer.Flush()
				return true
			}
		})
	}
}
