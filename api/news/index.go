package news

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hidatara-ds/evolipia-radar/pkg/api"
)

const httpMethodOptions = "OPTIONS"

// Handler for /api/news - Get all news
func Handler(w http.ResponseWriter, r *http.Request) {
	api.EnableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == httpMethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	topic := query.Get("topic")
	date := query.Get("date")

	newsData, err := api.LoadNewsData()
	if err != nil {
		if encErr := json.NewEncoder(w).Encode(api.Response{
			Success: false,
			Error:   "Failed to load news data: " + err.Error(),
		}); encErr != nil {
			log.Printf("Error encoding error response: %v", encErr)
		}
		return
	}

	// Filter by topic if specified
	filteredItems := newsData.Items
	if topic != "" {
		var filtered []api.NewsItem
		for _, item := range filteredItems {
			for _, tag := range item.Tags {
				if strings.EqualFold(tag, topic) {
					filtered = append(filtered, item)
					break
				}
			}
		}
		filteredItems = filtered
	}

	// Filter by date if specified
	if date == "today" {
		now := time.Now()
		var filtered []api.NewsItem
		for _, item := range filteredItems {
			if item.PublishedAt.Year() == now.Year() &&
				item.PublishedAt.Month() == now.Month() &&
				item.PublishedAt.Day() == now.Day() {
				filtered = append(filtered, item)
			}
		}
		filteredItems = filtered
	}

	// Limit to 20 items
	if len(filteredItems) > 20 {
		filteredItems = filteredItems[:20]
	}

	if err := json.NewEncoder(w).Encode(api.Response{
		Success: true,
		Data: map[string]interface{}{
			"items":        filteredItems,
			"total_count":  len(filteredItems),
			"last_updated": newsData.LastUpdated,
		},
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
