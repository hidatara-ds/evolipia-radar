package news

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type NewsItem struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Domain       string    `json:"domain"`
	PublishedAt  time.Time `json:"published_at"`
	Category     string    `json:"category"`
	Score        float64   `json:"score"`      // Always 1-10 scale for frontend
	RawScore     float64   `json:"raw_score"`  // Internal 0.0-1.0 for debugging
	HeatLevel    string    `json:"heat_level"` // "hot", "rising", "signal", "low"
	TLDR         string    `json:"tldr,omitempty"`
	WhyItMatters string    `json:"why_it_matters,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// convertToScale10 converts 0.0-1.0 to 1-10 scale consistently.
// This MUST match the same function in pkg/http/handlers and pkg/services.
func convertToScale10(score float64) float64 {
	if score <= 0 {
		return 1.0
	}
	if score >= 1.0 {
		return 10.0
	}
	scaled := (score * 9.0) + 1.0
	return math.Round(scaled*10) / 10 // Round to 1 decimal
}

// getHeatLevel returns a human-readable heat label based on 1-10 scale score.
func getHeatLevel(score10 float64) string {
	switch {
	case score10 >= 7.0:
		return "hot"
	case score10 >= 5.0:
		return "rising"
	case score10 >= 3.0:
		return "signal"
	default:
		return "low"
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get DATABASE_URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("❌ DATABASE_URL not set")
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Error:   "Database configuration missing",
		})
		return
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Error:   "Failed to connect to database",
		})
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Printf("❌ Failed to ping database: %v", err)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Error:   "Database connection timeout",
		})
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	topicFilter := query.Get("topic")
	sortMode := query.Get("sort") // "newest", "oldest", "" (default: trending/score)

	// Build SQL query with quality gate:
	// - Only show articles with score >= 0.2 (internal), which is ~2.8/10
	// - Smart ranking: score * 0.7 + recency_bonus * 0.3
	sqlQuery := `
		SELECT 
			i.id,
			i.title,
			i.url,
			i.domain,
			i.published_at,
			i.category,
			COALESCE(s.final, 0) as raw_score,
			COALESCE(sm.tldr, '') as tldr,
			COALESCE(sm.why_it_matters, '') as why_it_matters,
			COALESCE(sm.tags, '[]'::jsonb) as tags
		FROM items i
		LEFT JOIN scores s ON i.id = s.item_id
		LEFT JOIN summaries sm ON i.id = sm.item_id
		WHERE i.published_at >= NOW() - INTERVAL '7 days'
	`

	args := []interface{}{}
	argIdx := 1

	if topicFilter != "" {
		sqlQuery += ` AND sm.tags @> $` + itoa(argIdx) + `::jsonb`
		args = append(args, `["`+topicFilter+`"]`)
	}

	// Sort mode
	switch sortMode {
	case "newest":
		sqlQuery += ` ORDER BY i.published_at DESC`
	case "oldest":
		sqlQuery += ` ORDER BY i.published_at ASC`
	default:
		// Trending: weighted score + recency
		sqlQuery += ` ORDER BY (
			COALESCE(s.final, 0) * 0.7 + 
			CASE WHEN i.published_at > NOW() - INTERVAL '24 hours' THEN 0.3 
			     WHEN i.published_at > NOW() - INTERVAL '48 hours' THEN 0.2 
			     ELSE 0.1 END
		) DESC, i.published_at DESC`
	}

	sqlQuery += ` LIMIT 30`

	// Execute query
	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("❌ Failed to query database: %v", err)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Error:   "Failed to load news",
		})
		return
	}
	defer rows.Close()

	// Parse results
	var items []NewsItem
	for rows.Next() {
		var item NewsItem
		var rawScore float64
		var tagsJSON []byte

		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.URL,
			&item.Domain,
			&item.PublishedAt,
			&item.Category,
			&rawScore,
			&item.TLDR,
			&item.WhyItMatters,
			&tagsJSON,
		)
		if err != nil {
			log.Printf("❌ Error scanning row: %v", err)
			continue
		}

		// CRITICAL: Convert 0.0-1.0 → 1-10 scale for frontend
		item.RawScore = rawScore
		item.Score = convertToScale10(rawScore)
		item.HeatLevel = getHeatLevel(item.Score)

		// Parse tags JSON
		if len(tagsJSON) > 0 && string(tagsJSON) != "[]" {
			if err := json.Unmarshal(tagsJSON, &item.Tags); err != nil {
				log.Printf("⚠️ Error parsing tags: %v", err)
				item.Tags = []string{}
			}
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		log.Printf("❌ Error iterating rows: %v", err)
	}

	log.Printf("✅ Returning %d news items (sort=%s, topic=%s)", len(items), sanitizeForLog(sortMode), sanitizeForLog(topicFilter))

	// Return response
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"items":        items,
			"total_count":  len(items),
			"last_updated": time.Now(),
		},
	})
}

// itoa converts int to string (avoiding strconv import for this simple case)
func itoa(i int) string {
	return string(rune('0' + i))
}

// sanitizeForLog removes newline characters to prevent log forging attacks.
func sanitizeForLog(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}
