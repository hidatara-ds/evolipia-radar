// Package main is the entry point for the Evolipia Radar backend server.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hidatara-ds/evolipia-radar/internal/api"
	"github.com/hidatara-ds/evolipia-radar/internal/crawler"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/hidatara-ds/evolipia-radar/pkg/config"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	slog.Info("Starting Evolipia Radar Server...")

	cfg := config.Load()

	// Initialize Database connection (if database URL provided)
	var database *db.DB
	if cfg.DatabaseURL != "" {
		dbInstance, err := db.New(cfg)
		if err != nil {
			slog.Warn("Database connection failed, starting with in-memory mode", "err", err)
		} else {
			database = dbInstance
			defer database.Close()
			slog.Info("Database connected successfully")
		}
	}

	// Initialize SSE Progress Broadcaster
	broadcaster := api.NewProgressBroadcaster()

	// Initialize Validator & Retry Runner
	validator := crawler.NewValidator(cfg.MinRelevanceScore, cfg.TopicKeywords)
	retryRunner := crawler.NewRetryRunner(cfg.MaxCrawlRetries, 1*time.Second)

	// Crawl Task definition with step logging & SSE events
	crawlTaskFunc := func(ctx context.Context, onProgress func(models.CrawlProgressEvent)) (int, error) {
		slog.Info("Starting crawl step 1: Initializing crawler...")
		onProgress(models.CrawlProgressEvent{
			Step:       1,
			Message:    "Initializing crawler...",
			Progress:   10,
			Timestamp:  time.Now(),
			IsComplete: false,
		})

		sources := []string{"HackerNews RSS", "ArXiv AI Papers", "Reddit MachineLearning", "TechCrunch AI"}
		totalSources := len(sources)
		processedItems := 0

		for idx, src := range sources {
			stepNum := 2
			if idx > 0 {
				stepNum = 3
			}

			slog.Info("Starting source crawl", "source", src, "step", stepNum, "index", idx+1, "total", totalSources)
			onProgress(models.CrawlProgressEvent{
				Step:                   stepNum,
				Message:                fmt.Sprintf("Scanning source (%d/%d): %s...", idx+1, totalSources, src),
				Progress:               20 + (idx * 50 / totalSources),
				CurrentSource:          src,
				TotalSources:           totalSources,
				ProcessedItems:         processedItems,
				EstimatedRemainingSecs: (totalSources - idx) * 3,
				Timestamp:              time.Now(),
			})

			// Execute fetch with retry runner
			err := retryRunner.ExecuteWithRetry(ctx, src, func(c context.Context) error {
				// Mock content fetch & validation
				time.Sleep(500 * time.Millisecond)

				candidate := &models.Item{
					Title:       fmt.Sprintf("Latest Advances in Autonomous Agents from %s", src),
					URL:         fmt.Sprintf("https://example.com/%s/item-%d", src, time.Now().UnixNano()),
					SourceName:  src,
					PublishedAt: time.Now(),
					RawExcerpt:  stringPtr("Detailed technical breakdown of open-source agent frameworks and scalable LLM inference pipelines."),
				}

				if valErr := validator.ValidateItem(candidate); valErr != nil {
					return fmt.Errorf("item validation failed: %w", valErr)
				}

				score := validator.ScoreRelevance(candidate.Title, *candidate.RawExcerpt)
				candidate.RelevanceScore = score
				if !validator.IsReleasesSufficientlyRelevant(score) {
					slog.Info("Item skipped due to low relevance score", "source", src, "score", score)
					return nil
				}

				processedItems += 3
				return nil
			})

			if err != nil {
				slog.Error("Source crawl failed", "source", src, "err", err)
			}
		}

		slog.Info("Crawl step 4: Validating data...")
		onProgress(models.CrawlProgressEvent{
			Step:       4,
			Message:    "Validating data...",
			Progress:   80,
			Timestamp:  time.Now(),
			IsComplete: false,
		})
		time.Sleep(200 * time.Millisecond)

		slog.Info("Crawl step 5: Saving to database...")
		onProgress(models.CrawlProgressEvent{
			Step:       5,
			Message:    "Saving to database...",
			Progress:   90,
			Timestamp:  time.Now(),
			IsComplete: false,
		})
		time.Sleep(200 * time.Millisecond)

		return processedItems, nil
	}

	// Initialize & Start Auto-Scheduler
	scheduler, err := crawler.NewScheduler(cfg.CrawlInterval, crawlTaskFunc, broadcaster.Broadcast)
	if err != nil {
		slog.Error("Failed to initialize scheduler", "err", err)
		os.Exit(1)
	}
	scheduler.Start()

	// Initialize Gin HTTP Router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(api.CORS())
	router.Use(api.Logger())

	// Serve Web UI static files
	router.StaticFile("/", "./web/index.html")
	router.Static("/web", "./web")

	// Handlers
	itemsHandler := api.NewItemsHandler(database)

	// API Routes
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
	})

	router.GET("/api/items", itemsHandler.HandleGetItems)
	router.GET("/api/crawl/progress", broadcaster.HandleSSEProgress())

	router.GET("/api/crawl/status", func(c *gin.Context) {
		lastTime, status, count, errStr, isRunning := scheduler.GetStatus()
		c.JSON(http.StatusOK, gin.H{
			"last_run_time":    lastTime,
			"last_run_status":  status,
			"last_items_count": count,
			"last_error":       errStr,
			"is_running":       isRunning,
		})
	})

	router.POST("/api/crawl", func(c *gin.Context) {
		go func() {
			_, _ = scheduler.RunCrawl(context.Background(), "manual")
		}()
		c.JSON(http.StatusAccepted, gin.H{
			"status":  "triggered",
			"message": "Manual crawl triggered successfully",
		})
	})

	// HTTP Server Configuration
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("API Server running", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server forced to shutdown", "err", err)
		}
	}()

	// Graceful Shutdown on SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server gracefully...")

	// Stop scheduler first
	scheduler.Stop()

	// Shutdown HTTP Server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP Server forced to shutdown", "err", err)
	}

	slog.Info("Server stopped cleanly")
}

func stringPtr(s string) *string {
	return &s
}
