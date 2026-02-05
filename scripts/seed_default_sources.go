package main

import (
	"context"
	"log"

	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/db"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
)

func main() {
	cfg := config.Load()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	sourceRepo := db.NewSourceRepository(database)
	ctx := context.Background()

	defaultSources := []models.Source{
		{
			Name:     "Hacker News",
			Type:     "hacker_news",
			Category: "news",
			URL:      "https://news.ycombinator.com",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "arXiv AI/ML",
			Type:     "arxiv",
			Category: "news",
			URL:      "http://export.arxiv.org/api/query",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "OpenAI Blog",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://openai.com/blog/rss.xml",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "Google AI Blog",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://ai.googleblog.com/feeds/posts/default",
			Enabled:  true,
			Status:   "active",
		},
	}

	for _, source := range defaultSources {
		// Check if source already exists
		sources, err := sourceRepo.List(ctx)
		if err != nil {
			log.Printf("Error listing sources: %v", err)
			continue
		}

		exists := false
		for _, s := range sources {
			if s.URL == source.URL {
				exists = true
				log.Printf("Source already exists: %s", source.Name)
				break
			}
		}

		if !exists {
			if err := sourceRepo.Create(ctx, &source); err != nil {
				log.Printf("Error creating source %s: %v", source.Name, err)
			} else {
				log.Printf("Created source: %s", source.Name)
			}
		}
	}

	log.Println("Default sources seeded successfully")
}
