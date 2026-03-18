package metrics

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/hidatara-ds/evolipia-radar/pkg/api"
	"github.com/hidatara-ds/evolipia-radar/pkg/config"
	"github.com/hidatara-ds/evolipia-radar/pkg/crawler"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
)

// MetricsHandler for /metrics - Verifies system ingestion stats.
// In a serverless environment, this will only return cold-start zeroed stats 
// unless backed by a persistent Redis or DB structure.
func Handler(w http.ResponseWriter, r *http.Request) {
	api.EnableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Load metrics from DB for serverless persistence
	cfg := config.Load()
	database, err := db.New(cfg)
	if err != nil {
		log.Printf("[METRICS] DB Connection failed: %v", err)
		// Fallback to empty metrics if DB is down
		json.NewEncoder(w).Encode(&crawler.Metrics{})
		return
	}
	defer database.Close()

	metricsData := crawler.NewMetrics(database.Pool)
	metricsData.LoadFromDB(r.Context())

	if dryRun := os.Getenv("DRY_RUN"); dryRun == "true" {
		// Mock data for dry-run verification
		metricsData.ArticlesProcessed = 16
		metricsData.FilteredArticles = 24
	}

	if err := json.NewEncoder(w).Encode(metricsData); err != nil {
		log.Printf("[METRICS] Failed to encode response: %v", err)
	}
}
