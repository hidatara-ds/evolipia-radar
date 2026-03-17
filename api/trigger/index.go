package trigger

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/hidatara-ds/evolipia-radar/pkg/ai"
	"github.com/hidatara-ds/evolipia-radar/pkg/api"
	"github.com/hidatara-ds/evolipia-radar/pkg/cluster"
	"github.com/hidatara-ds/evolipia-radar/pkg/config"
	"github.com/hidatara-ds/evolipia-radar/pkg/crawler"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler handles the /v2/crawl/trigger manual webhook route on Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	api.EnableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	dryRunEnv := os.Getenv("DRY_RUN") == "true"
	log.Printf("[VERCEL TRIGGER] Starting crawler cycle. DryRun: %v", dryRunEnv)

	// AI Setup
	aiCfg := config.LoadAIConfig()
	orProvider := ai.NewOpenRouterProvider(ai.OpenRouterProviderConfig{
		APIKey:       aiCfg.APIKey,
		DefaultModel: aiCfg.DefaultModel,
	})
	
	budgetProvider := ai.NewTrackerMiddleware(orProvider, 10000, 300000)
	aiService := ai.NewService(budgetProvider)

	// DB Setup
	var pool *pgxpool.Pool
	if !dryRunEnv {
		cfg := config.Load()
		database, err := db.New(cfg)
		if err != nil {
			log.Printf("[VERCEL TRIGGER] DB Connection failed: %v", err)
			http.Error(w, `{"error":"database connection failed"}`, http.StatusInternalServerError)
			return
		}
		defer database.Close()
		pool = database.Pool
	}

	clusterService := ai.NewClusterService(aiService, pool)
	
	// Phase 5: In-Memory Clustering Routing
	embedder := cluster.NewOpenRouterEmbedder(aiCfg.APIKey)
	inMemClusterSvc := cluster.NewClusterService(embedder)
	
	metricsData := &crawler.Metrics{}
	botOrchestrator := crawler.NewOrchestrator(clusterService, inMemClusterSvc, metricsData, dryRunEnv)

	// Executing the cycle synchronously for Vercel Serverless
	stats := botOrchestrator.RunCycle(context.Background())
	botOrchestrator.UpdateClusterMetrics(context.Background())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "completed",
		"stats":  stats,
	})
}
