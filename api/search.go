package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func loadNewsData() (*NewsData, error) {
	// Try multiple possible paths
	paths := []string{
		"data/news.json",
		"../data/news.json",
		"../../data/news.json",
	}

	var data []byte
	var err error

	for _, path := range paths {
		// #nosec G304 - Path is from a fixed list, not user input
		data, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	var newsData NewsData
	if err := json.Unmarshal(data, &newsData); err != nil {
		return nil, err
	}

	return &newsData, nil
}

// Handler for /api/search - Search news
func Handler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Error:   "Query parameter 'q' is required",
		}); err != nil {
			log.Printf("Error encoding error response: %v", err)
		}
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

	// Simple search in title and tags
	queryLower := strings.ToLower(query)
	var results []NewsItem

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

	if err := json.NewEncoder(w).Encode(Response{
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