package news

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	Score        float64   `json:"score"`
	TLDR         string    `json:"tldr,omitempty"`
	WhyItMatters string    `json:"why_it_matters,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
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

	// Set connection pool settings
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Parse query parameters
	query := r.URL.Query()
	topicFilter := query.Get("topic")

	// Build SQL query
	sqlQuery := `
		SELECT 
			i.id,
			i.title,
			i.url,
			i.domain,
			i.published_at,
			i.category,
			COALESCE(s.final, 0.5) as score,
			COALESCE(sm.tldr, '') as tldr,
			COALESCE(sm.why_it_matters, '') as why_it_matters,
			COALESCE(sm.tags, '[]'::jsonb) as tags
		FROM items i
		LEFT JOIN scores s ON i.id = s.item_id
		LEFT JOIN summaries sm ON i.id = sm.item_id
		WHERE i.published_at >= NOW() - INTERVAL '7 days'
	`

	args := []interface{}{}
	if topicFilter != "" {
		sqlQuery += ` AND sm.tags @> $1::jsonb`
		args = append(args, `["`+topicFilter+`"]`)
	}

	sqlQuery += ` ORDER BY COALESCE(s.final, 0) DESC, i.published_at DESC LIMIT 20`

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
		var tagsJSON []byte

		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.URL,
			&item.Domain,
			&item.PublishedAt,
			&item.Category,
			&item.Score,
			&item.TLDR,
			&item.WhyItMatters,
			&tagsJSON,
		)
		if err != nil {
			log.Printf("❌ Error scanning row: %v", err)
			continue
		}

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

	log.Printf("✅ Returning %d news items", len(items))

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
