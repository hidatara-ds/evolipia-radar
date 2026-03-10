package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/db"
	"github.com/hidatara-ds/evolipia-radar/internal/services"
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

type NewsData struct {
	Items       []NewsItem `json:"items"`
	LastUpdated time.Time  `json:"last_updated"`
	TotalCount  int        `json:"total_count"`
}

func main() {
	cfg := config.Load()

	// Connect to database
	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run ingestion
	w := services.NewWorker(database, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	log.Println("Starting ingestion...")
	if err := w.RunIngestion(ctx); err != nil {
		log.Printf("Ingestion error: %v", err)
	} else {
		log.Println("Ingestion completed successfully")
	}

	// Fetch latest news from database
	log.Println("Fetching latest news...")
	items, err := fetchLatestNews(ctx, database)
	if err != nil {
		log.Fatalf("Failed to fetch news: %v", err)
	}

	// Write to JSON file
	outputPath := os.Getenv("JSON_OUTPUT_PATH")
	if outputPath == "" {
		outputPath = "data/news.json"
	}

	newsData := NewsData{
		Items:       items,
		LastUpdated: time.Now(),
		TotalCount:  len(items),
	}

	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(newsData); err != nil {
		log.Fatalf("Failed to write JSON: %v", err)
	}

	log.Printf("Successfully wrote %d items to %s", len(items), outputPath)
}

func fetchLatestNews(ctx context.Context, database *db.DB) ([]NewsItem, error) {
	query := `
		SELECT 
			i.id,
			i.title,
			i.url,
			i.domain,
			i.published_at,
			i.category,
			COALESCE(s.final, 0) as score,
			COALESCE(sm.tldr, '') as tldr,
			COALESCE(sm.why_it_matters, '') as why_it_matters,
			COALESCE(sm.tags, '[]'::jsonb) as tags
		FROM items i
		LEFT JOIN scores s ON i.id = s.item_id
		LEFT JOIN summaries sm ON i.id = sm.item_id
		WHERE i.published_at >= NOW() - INTERVAL '7 days'
		ORDER BY COALESCE(s.final, 0) DESC, i.published_at DESC
		LIMIT 100
	`

	rows, err := database.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Parse tags JSON
		if len(tagsJSON) > 0 {
			if err := json.Unmarshal(tagsJSON, &item.Tags); err != nil {
				log.Printf("Error parsing tags: %v", err)
				item.Tags = []string{}
			}
		}

		items = append(items, item)
	}

	return items, rows.Err()
}
