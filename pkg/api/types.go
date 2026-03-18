package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

func LoadNewsData() (*NewsData, error) {
	// Try multiple possible paths suitable for both local development and Vercel
	paths := []string{
		"data/news.json",           // Local development
		"../data/news.json",        // From api subfolder
		"../../data/news.json",     // From nested api folders
		"api/news.json",            // Copied to api folder
		"./news.json",              // Same directory as Go function
		"/var/task/data/news.json", // Vercel specific (old)
		"/var/task/api/news.json",  // Vercel specific (new)
	}

	var data []byte
	var err error

	for _, path := range paths {
		// #nosec G304 - Path is from a fixed list
		data, err = os.ReadFile(path)
		if err == nil {
			log.Printf("✅ Successfully loaded news from: %s", path)
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

func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
