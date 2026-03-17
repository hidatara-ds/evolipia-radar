package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/hidatara-ds/evolipia-radar/internal/api"
	"github.com/hidatara-ds/evolipia-radar/internal/crawler"
)

// MetricsHandler for /metrics - Verifies system ingestion stats.
// In a serverless environment, this will only return cold-start zeroed stats 
// unless backed by a persistent Redis or DB structure.
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	api.EnableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// This is a stub for the Vercel architecture. To keep state between serverless
	// invocations, we would query the database here.
	metricsData := &crawler.Metrics{}

	if dryRun := os.Getenv("DRY_RUN"); dryRun == "true" {
		// Mock data for dry-run verification
		metricsData.ArticlesProcessed = 16
		metricsData.FilteredArticles = 24
	}

	json.NewEncoder(w).Encode(metricsData)
}
