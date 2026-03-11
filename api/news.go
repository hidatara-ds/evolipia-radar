package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const httpMethodOptions = "OPTIONS"

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

// Handler for /api/news - Get all news
func Handler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == httpMethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	topic := query.Get("topic")
	date := query.Get("date")

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

	// Filter by topic if specified
	filteredItems := newsData.Items
	if topic != "" {
		var filtered []NewsItem
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
		var filtered []NewsItem
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

	if err := json.NewEncoder(w).Encode(Response{
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

// ItemHandler for /api/news/[id] - Get single news item
func ItemHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == httpMethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract ID from path
	path := r.URL.Path
	id := filepath.Base(path)

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

	// Find item by ID
	for _, item := range newsData.Items {
		if item.ID == id {
			if encErr := json.NewEncoder(w).Encode(Response{
				Success: true,
				Data:    item,
			}); encErr != nil {
				log.Printf("Error encoding response: %v", encErr)
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   "Item not found",
	}); err != nil {
		log.Printf("Error encoding not found response: %v", err)
	}
}

// TrendingHandler for /api/trending - Get trending items
func TrendingHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == httpMethodOptions {
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

// SearchHandler for /api/search - Search news
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == httpMethodOptions {
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
