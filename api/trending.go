package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Handler for /api/trending - Get trending items
func TrendingHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	newsData, err := loadNewsData()
	if err != nil {
		if encErr := json.NewEncoder(w).Encode(Response{
			Success: false,
			Error:   "Failed to load news data: " + err.Error(),
		}); encErr != nil {
			log.Printf("Error encoding error response: %v", encErr)
		}
		return
	}

	// Get items from last 2 hours with high scores
	twoHoursAgo := time.Now().Add(-2 * time.Hour)
	var trending []NewsItem

	for _, item := range newsData.Items {
		if item.PublishedAt.After(twoHoursAgo) && item.Score > 0.5 {
			trending = append(trending, item)
		}
	}

	// Limit to 20 items
	if len(trending) > 20 {
		trending = trending[:20]
	}

	if err := json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"items":       trending,
			"total_count": len(trending),
		},
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}