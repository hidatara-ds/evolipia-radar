package metrics

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// MetricsResponse holds live DB-derived metrics for the dashboard.
type MetricsResponse struct {
	ArticlesProcessed int      `json:"articles_processed"`
	FilteredArticles  int      `json:"filtered_articles"`
	Clusters          int      `json:"clusters"`
	AvgClusterScore   float64  `json:"avg_cluster_score"`
	TopClusterTitles  []string `json:"top_cluster_titles"`
	APIHits           int      `json:"api_hits"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("[METRICS] DATABASE_URL not set")
		json.NewEncoder(w).Encode(&MetricsResponse{})
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("[METRICS] DB Connection failed: %v", err)
		json.NewEncoder(w).Encode(&MetricsResponse{})
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp := &MetricsResponse{}

	// 1. Total articles crawled in last 7 days
	_ = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM items WHERE created_at >= NOW() - INTERVAL '7 days'`,
	).Scan(&resp.ArticlesProcessed)

	// 2. Filtered (scored) articles — those that passed AI analysis
	_ = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM items i 
		 JOIN scores s ON i.id = s.item_id 
		 WHERE i.created_at >= NOW() - INTERVAL '7 days'`,
	).Scan(&resp.FilteredArticles)

	// 3. Clusters — count distinct tag groups from summaries
	_ = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM summaries WHERE created_at >= NOW() - INTERVAL '7 days'`,
	).Scan(&resp.Clusters)

	// 4. Average score (convert to 1-10 scale for consistency)
	var avgRaw float64
	_ = db.QueryRowContext(ctx,
		`SELECT COALESCE(AVG(s.final), 0) FROM scores s 
		 JOIN items i ON s.item_id = i.id 
		 WHERE i.created_at >= NOW() - INTERVAL '7 days'`,
	).Scan(&avgRaw)
	resp.AvgClusterScore = math.Round(convertToScale10(avgRaw)*10) / 10

	// 5. Top trending cluster titles (top 5 by score)
	rows, err := db.QueryContext(ctx, `
		SELECT i.title FROM items i
		JOIN scores s ON i.id = s.item_id
		WHERE i.created_at >= NOW() - INTERVAL '7 days'
		ORDER BY s.final DESC
		LIMIT 5
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var title string
			if rows.Scan(&title) == nil {
				resp.TopClusterTitles = append(resp.TopClusterTitles, title)
			}
		}
	}

	log.Printf("✅ [METRICS] articles=%d, filtered=%d, clusters=%d, avg=%.1f",
		resp.ArticlesProcessed, resp.FilteredArticles, resp.Clusters, resp.AvgClusterScore)

	json.NewEncoder(w).Encode(resp)
}

// convertToScale10 converts 0.0-1.0 to 1-10 scale.
// Must match the same function used across all handlers.
func convertToScale10(score float64) float64 {
	if score <= 0 {
		return 1.0
	}
	if score >= 1.0 {
		return 10.0
	}
	return (score * 9.0) + 1.0
}
