package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evolipia/radar/internal/config"
	"github.com/evolipia/radar/internal/db"
	"github.com/evolipia/radar/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	router := gin.Default()

	// Health check
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	v1 := router.Group("/v1")
	{
		h := handlers.New(database)
		v1.GET("/feed", h.GetFeed)
		v1.GET("/rising", h.GetRising)
		v1.GET("/items/:id", h.GetItem)
		v1.GET("/search", h.Search)
		v1.GET("/sources", h.ListSources)
		v1.POST("/sources", h.CreateSource)
		v1.POST("/sources/test", h.TestSource)
		v1.PATCH("/sources/:id/enable", h.EnableSource)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("API server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
