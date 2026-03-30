package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hidatara-ds/evolipia-radar/pkg/ai"
	ai_api "github.com/hidatara-ds/evolipia-radar/pkg/api"
	"github.com/hidatara-ds/evolipia-radar/pkg/cluster"
	"github.com/hidatara-ds/evolipia-radar/pkg/config"
	"github.com/hidatara-ds/evolipia-radar/pkg/crawler"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
	"github.com/hidatara-ds/evolipia-radar/pkg/http/handlers"
)

func main() {
	cfg := config.Load()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	router := gin.Default()

	// Web UI (mobile-first); serve from ./web when running from project root
	router.StaticFile("/", "./web/index.html")
	router.StaticFile("/index.html", "./web/index.html")

	// Serve static folders
	router.Static("/web", "web")
	router.Static("/assets", "assets")

	// ----------------------------------------------------------------------
	// Phase 1-3 AI & Intelligent Discovery Integration Setup
	// ----------------------------------------------------------------------
	// Instantiate the core AI providers. Use standard zero-budget configs.
	aiCfg := config.LoadAIConfig()
	orProvider := ai.NewOpenRouterProvider(ai.OpenRouterProviderConfig{
		APIKey:       aiCfg.APIKey,
		DefaultModel: aiCfg.DefaultModel,
	})

	// Phase 2.9 Budget Control Middleware
	// Give the system 10,000 daily tokens and 300,000 monthly free tier limit.
	budgetGuardedProvider := ai.NewTrackerMiddleware(orProvider, 10000, 300000)
	
	centralAIService := ai.NewService(budgetGuardedProvider)
	clusterService := ai.NewClusterService(centralAIService, database.Pool)

	// Phase 5 In-Memory Service
	embedder := cluster.NewOpenRouterEmbedder(aiCfg.APIKey)
	inMemClusterSvc := cluster.NewClusterService(embedder)

	// Phase 3 Web Crawling Orchestrator
	metricsData := crawler.NewMetrics(database.Pool)
	metricsData.LoadFromDB(context.Background())
	
	summarizer := crawler.NewSummarizer(centralAIService, database)

	dryRunEnv := os.Getenv("DRY_RUN") == "true"
	botOrchestrator := crawler.NewOrchestrator(clusterService, inMemClusterSvc, centralAIService, metricsData, database, dryRunEnv, summarizer)
	
	// Start the intelligent crawling loop in the background (runs every 15 minutes)
	crawlCtx, crawlCancel := context.WithCancel(context.Background())
	defer crawlCancel()
	go botOrchestrator.Start(crawlCtx, 15*time.Minute)

	// Health & Observability
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, metricsData)
	})

	// Register internal AI API components (V2)
	aiHandler := ai_api.NewAIHandler(centralAIService)
	aiHandler.RegisterRoutes(router.Group("/"))
	
	// Register Admin Ops / Manual Triggers
	v2 := router.Group("/v2")
	v2.POST("/crawl/trigger", func(c *gin.Context) {
		// Manual sync triggering of discovery agents (WARNING: Blocks response until cycle completes)
		stats := botOrchestrator.RunCycle(c.Request.Context())
		c.JSON(http.StatusOK, gin.H{
			"status": "completed",
			"stats": stats,
		})
	})
	// ----------------------------------------------------------------------

	// API routes
	v1 := router.Group("/v1")
	{
		h := handlers.New(database, centralAIService)
		v1.GET("/feed", h.GetFeed)
		v1.GET("/rising", h.GetRising)
		v1.GET("/items/:id", h.GetItem)
		v1.GET("/search", h.Search)
		v1.GET("/sources", h.ListSources)
		v1.POST("/sources", h.CreateSource)
		v1.POST("/sources/test", h.TestSource)
		v1.PATCH("/sources/:id/enable", h.EnableSource)

		// Settings API
		settingsHandler := ai_api.NewSettingsHandler(database)
		settingsHandler.RegisterRoutes(v1)
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("API server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
