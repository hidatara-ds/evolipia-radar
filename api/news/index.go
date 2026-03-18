package news

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hidatara-ds/evolipia-radar/pkg/api"
	"github.com/hidatara-ds/evolipia-radar/pkg/config"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
)

const httpMethodOptions = "OPTIONS"

// Handler for /api/news - Get all news from database
func Handler(w http.ResponseWriter, r *http.Request) {
	api.EnableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == httpMethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	topicFilter := query.Get("topic")
	dateFilter := query.Get("date")

	// Connect to database
	cfg := &config.Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	
	if cfg.DatabaseURL == "" {
		log.Println("❌ DATABASE_URL not set")
		if encErr := json.NewEncoder(w).Encode(api.Response{
			Success: false,
			Error:   "Database configuration missing",
		}); encErr != nil {
			log.Printf("Error encoding error response: %v", encErr)
		}
		return
	}

	database, err := db.New(cfg)
	if err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
		if encErr := json.NewEncoder(w).Encode(api.Response{
			Success: false,
			Error:   "Failed to connect to database: " + err.Error(),
		}); encErr != nil {
			log.Printf("Error encoding error response: %v", encErr)
		}
		return
	}
	defer database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get items from database
	itemRepo := db.NewItemRepository(database)
	summaryRepo := db.NewSummaryRepository(database)
	scoreRepo := db.NewScoreRepository(database)

	// Get top items from last 7 days
	var topic *string
	if topicFilter != "" {
		topic = &topicFilter
	}

	items, err := itemRepo.GetTopDaily(ctx, time.Now(), topic, 100)
	if err != nil {
		log.Printf("❌ Failed to get items: %v", err)
		if encErr := json.NewEncoder(w).Encode(api.Response{
			Success: false,
			Error:   "Failed to load news: " + err.Error(),
		}); encErr != nil {
			log.Printf("Error encoding error response: %v", encErr)
		}
		return
	}

	// Convert to API format
	var newsItems []api.NewsItem
	for _, item := range items {
		// Get summary and score
		summary, _ := summaryRepo.GetByItemID(ctx, item.ID)
		score, _ := scoreRepo.GetByItemID(ctx, item.ID)

		newsItem := api.NewsItem{
			ID:          item.ID.String(),
			Title:       item.Title,
			URL:         item.URL,
			Domain:      item.Domain,
			PublishedAt: item.PublishedAt,
			Category:    item.Category,
			Score:       0.5, // default
		}

		if summary != nil {
			newsItem.TLDR = summary.TLDR
			newsItem.WhyItMatters = summary.WhyItMatters
			newsItem.Tags = summary.Tags
		}

		if score != nil {
			newsItem.Score = score.Final
		}

		newsItems = append(newsItems, newsItem)
	}

	// Filter by date if specified
	if dateFilter == "today" {
		now := time.Now()
		var filtered []api.NewsItem
		for _, item := range newsItems {
			if item.PublishedAt.Year() == now.Year() &&
				item.PublishedAt.Month() == now.Month() &&
				item.PublishedAt.Day() == now.Day() {
				filtered = append(filtered, item)
			}
		}
		newsItems = filtered
	}

	// Limit to 20 items
	if len(newsItems) > 20 {
		newsItems = newsItems[:20]
	}

	log.Printf("✅ Returning %d news items", len(newsItems))

	if err := json.NewEncoder(w).Encode(api.Response{
		Success: true,
		Data: map[string]interface{}{
			"items":        newsItems,
			"total_count":  len(newsItems),
			"last_updated": time.Now(),
		},
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
