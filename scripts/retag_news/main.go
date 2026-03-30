// Script to re-tag existing news.json with auto-tagger
// Run with: go run scripts/retag_news.go
package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/hidatara-ds/evolipia-radar/pkg/tagging"
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

func main() {
	inputPath := "data/news.json"
	outputPath := "data/news.json"

	log.Println("🔄 Re-tagging news.json with auto-tagger...")

	// Read existing news.json
	data, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("❌ Failed to read %s: %v", inputPath, err)
	}

	var newsData NewsData
	if err := json.Unmarshal(data, &newsData); err != nil {
		log.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	log.Printf("📰 Found %d items to re-tag", len(newsData.Items))

	// Initialize auto-tagger
	tagger := tagging.NewAutoTagger()

	// Re-tag each item
	tagStats := make(map[string]int)
	retaggedCount := 0

	for i := range newsData.Items {
		item := &newsData.Items[i]

		// Combine title and summary for better tagging
		content := item.TLDR
		if content == "" {
			content = item.WhyItMatters
		}

		// Generate new tags
		autoTags := tagger.AssignTags(item.Title, content)

		// Merge with existing tags
		item.Tags = tagging.MergeTags(item.Tags, autoTags)

		// Count tags
		for _, tag := range item.Tags {
			tagStats[tag]++
		}

		retaggedCount++

		if retaggedCount%10 == 0 {
			log.Printf("   Processed %d/%d items...", retaggedCount, len(newsData.Items))
		}
	}

	// Update metadata
	newsData.LastUpdated = time.Now()

	// Write back to file
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("❌ Failed to create output file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(newsData); err != nil {
		log.Fatalf("❌ Failed to write JSON: %v", err)
	}

	log.Println("\n✅ Re-tagging complete!")
	log.Printf("📊 Tag distribution:")

	// Sort and display tag stats
	for tag, count := range tagStats {
		log.Printf("   %s: %d articles", tag, count)
	}

	log.Printf("\n💾 Updated file: %s", outputPath)
}
