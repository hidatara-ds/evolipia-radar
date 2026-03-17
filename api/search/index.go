package search

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/hidatara-ds/evolipia-radar/pkg/api"
)

// Handler for /api/search - Search news
func Handler(w http.ResponseWriter, r *http.Request) {
	api.EnableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		if err := json.NewEncoder(w).Encode(api.Response{
			Success: false,
			Error:   "Query parameter 'q' is required",
		}); err != nil {
			log.Printf("Error encoding error response: %v", err)
		}
		return
	}

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

	// Simple search in title and tags
	queryLower := strings.ToLower(query)
	var results []api.NewsItem

	for _, item := range newsData.Items {
		titleMatch := strings.Contains(strings.ToLower(item.Title), queryLower)
		tagMatch := false
		for _, tag := range item.Tags {
			if strings.Contains(strings.ToLower(tag), queryLower) {
				tagMatch = true
				break
			}
		}

		if titleMatch || tagMatch {
			results = append(results, item)
		}
	}

	// Limit to 20 items
	if len(results) > 20 {
		results = results[:20]
	}

	if err := json.NewEncoder(w).Encode(api.Response{
		Success: true,
		Data: map[string]interface{}{
			"items":       results,
			"total_count": len(results),
			"query":       query,
		},
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
