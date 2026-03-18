package api

import (
	"encoding/json"
	"fmt"
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
	log.Println("🔍 [LoadNewsData] Starting to load news data...")
	
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
	var successPath string

	for _, path := range paths {
		log.Printf("🔍 [LoadNewsData] Trying path: %s", path)
		// #nosec G304 - Path is from a fixed list
		data, err = os.ReadFile(path)
		if err == nil {
			successPath = path
			log.Printf("✅ [LoadNewsData] Successfully loaded from: %s (size: %d bytes)", path, len(data))
			break
		} else {
			log.Printf("❌ [LoadNewsData] Failed to read %s: %v", path, err)
		}
	}

	if err != nil {
		log.Printf("❌ [LoadNewsData] All paths failed. Last error: %v", err)
		
		// Log current working directory for debugging
		if cwd, cwdErr := os.Getwd(); cwdErr == nil {
			log.Printf("📂 [LoadNewsData] Current working directory: %s", cwd)
		}
		
		// List files in current directory
		if files, lsErr := os.ReadDir("."); lsErr == nil {
			log.Printf("📁 [LoadNewsData] Files in current directory:")
			for _, f := range files {
				log.Printf("   - %s (dir: %v)", f.Name(), f.IsDir())
			}
		}
		
		return nil, fmt.Errorf("failed to load news.json from any path: %w", err)
	}

	var newsData NewsData
	if err := json.Unmarshal(data, &newsData); err != nil {
		log.Printf("❌ [LoadNewsData] Failed to parse JSON: %v", err)
		return nil, fmt.Errorf("failed to parse news.json: %w", err)
	}

	log.Printf("✅ [LoadNewsData] Successfully parsed %d news items from %s", len(newsData.Items), successPath)
	return &newsData, nil
}

func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
