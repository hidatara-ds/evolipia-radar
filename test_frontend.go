package main

import (
	"log"
	"net/http"

	"github.com/hidatara-ds/evolipia-radar/api/news"
	"github.com/hidatara-ds/evolipia-radar/api/metrics"
)

func main() {
	// Enable CORS for local testing
	http.HandleFunc("/api/news", news.Handler)
	http.HandleFunc("/metrics", metrics.Handler)
	
	log.Println("🚀 Test API server starting on http://localhost:8080")
	log.Println("📡 Endpoints:")
	log.Println("   - http://localhost:8080/api/news")
	log.Println("   - http://localhost:8080/metrics")
	log.Println("")
	log.Println("🌐 Frontend running on http://localhost:3000")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
