package main

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/internal/crawler"
	"github.com/hidatara-ds/evolipia-radar/pkg/config"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
	"github.com/hidatara-ds/evolipia-radar/pkg/services"
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
	cfg := config.Load()
	outputPath := getOutputPath()

	// 1. Load existing items from disk to prevent data loss (including 30-07-2026 data)
	existingData := loadExistingNews(outputPath)
	log.Printf("Loaded %d existing articles from disk", len(existingData.Items))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var newDiscovered []NewsItem

	// 2. Determine if database should be used
	if shouldSkipDatabaseWork() {
		log.Println("CI or SKIP_DB enabled: Running standalone multi-agent feed crawler...")
		newDiscovered = crawlStandaloneFeeds(ctx)
	} else {
		database, err := db.New(cfg)
		if err != nil {
			log.Printf("Database connection failed (%v). Falling back to standalone multi-agent feed crawler...", err)
			newDiscovered = crawlStandaloneFeeds(ctx)
		} else {
			defer database.Close()

			log.Println("Starting database ingestion worker...")
			w := services.NewWorker(database, cfg)
			if err := w.RunIngestion(ctx); err != nil {
				log.Printf("Ingestion warning: %v", err)
			} else {
				log.Println("Ingestion completed successfully")
			}

			log.Println("Fetching latest news from database...")
			dbItems, fetchErr := fetchLatestNews(ctx, database)
			if fetchErr != nil {
				log.Printf("Failed to fetch news from database: %v. Running feed crawler...", fetchErr)
				newDiscovered = crawlStandaloneFeeds(ctx)
			} else {
				newDiscovered = dbItems
			}
		}
	}

	// 3. Smart Merge: deduplicate and prepend new items while preserving existing historical items
	mergedItems := mergeNewsItems(existingData.Items, newDiscovered)
	log.Printf("Merged results: %d new signals added. Total signals: %d", len(mergedItems)-len(existingData.Items), len(mergedItems))

	finalData := NewsData{
		Items:       mergedItems,
		LastUpdated: time.Now(),
		TotalCount:  len(mergedItems),
	}

	// 4. Save to both outputPath and api/news.json
	writeJSONOutput(outputPath, finalData)
	syncToAPIPath(finalData)
}

func getOutputPath() string {
	outputPath := os.Getenv("JSON_OUTPUT_PATH")
	if outputPath == "" {
		return "data/news.json"
	}
	return outputPath
}

func shouldSkipDatabaseWork() bool {
	return os.Getenv("SKIP_DB") == "true" || os.Getenv("CI") == "true"
}

func loadExistingNews(primaryPath string) NewsData {
	paths := []string{primaryPath, "api/news.json", "data/news.json"}
	for _, p := range paths {
		// #nosec G304
		data, err := os.ReadFile(filepath.Clean(p))
		if err == nil {
			var parsed NewsData
			if jsonErr := json.Unmarshal(data, &parsed); jsonErr == nil && len(parsed.Items) > 0 {
				return parsed
			}
		}
	}
	return NewsData{Items: []NewsItem{}, LastUpdated: time.Now(), TotalCount: 0}
}

func normalizeURL(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return strings.ToLower(strings.TrimSpace(raw))
	}
	q := parsed.Query()
	for k := range q {
		if strings.HasPrefix(k, "utm_") || k == "ref" || k == "source" || k == "fbclid" {
			q.Del(k)
		}
	}
	parsed.RawQuery = q.Encode()
	parsed.Fragment = ""
	res := parsed.String()
	return strings.ToLower(strings.TrimRight(res, "/"))
}

func extractDomain(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "external"
	}
	dom := strings.TrimPrefix(parsed.Hostname(), "www.")
	if dom == "" {
		return "external"
	}
	return dom
}

func inferCategoryAndTags(title, domain string) (string, []string) {
	lower := strings.ToLower(title + " " + domain)
	category := "llm"
	tags := []string{}

	if strings.Contains(lower, "arxiv") || strings.Contains(lower, "paper") || strings.Contains(lower, "benchmark") {
		category = "research"
		tags = append(tags, "research")
	}
	if strings.Contains(lower, "agent") || strings.Contains(lower, "autonomous") || strings.Contains(lower, "workflow") {
		category = "agents"
		tags = append(tags, "agents")
	}
	if strings.Contains(lower, "vision") || strings.Contains(lower, "image") || strings.Contains(lower, "diffusion") {
		category = "vision"
		tags = append(tags, "vision")
	}
	if strings.Contains(lower, "robot") || strings.Contains(lower, "robotics") || strings.Contains(lower, "atlas") {
		category = "robotics"
		tags = append(tags, "robotics")
	}
	if strings.Contains(lower, "open-source") || strings.Contains(lower, "github") || strings.Contains(lower, "weights") {
		category = "open-source"
		tags = append(tags, "open-source")
	}
	if strings.Contains(lower, "infra") || strings.Contains(lower, "vector") || strings.Contains(lower, "rag") || strings.Contains(lower, "cuda") {
		category = "infra"
		tags = append(tags, "infra")
	}

	tags = append(tags, category)
	return category, tags
}

func crawlStandaloneFeeds(ctx context.Context) []NewsItem {
	rssAgent := crawler.NewRSSAgent()
	trendingAgent := crawler.NewTrendingAgent()

	var discovered []NewsItem

	// 1. RSS Agent (HN RSS, ArXiv AI, TechCrunch AI, OpenAI Blog)
	rssArticles, err := rssAgent.Crawl(ctx, 25)
	if err != nil {
		log.Printf("RSS crawl error: %v", err)
	} else {
		log.Printf("RSS agent fetched %d articles", len(rssArticles))
		for _, art := range rssArticles {
			discovered = append(discovered, articleToNewsItem(art))
		}
	}

	// 2. Trending Agent (Algolia HN Front-page)
	trendArticles, err := trendingAgent.Crawl(ctx, 20)
	if err != nil {
		log.Printf("Trending crawl error: %v", err)
	} else {
		log.Printf("Trending agent fetched %d articles", len(trendArticles))
		for _, art := range trendArticles {
			discovered = append(discovered, articleToNewsItem(art))
		}
	}

	return discovered
}

func articleToNewsItem(art crawler.Article) NewsItem {
	dom := extractDomain(art.Link)
	cat, tags := inferCategoryAndTags(art.Title, dom)

	pubDate := art.PublishedAt
	if pubDate.IsZero() {
		pubDate = time.Now()
	}

	tldr := art.Content
	if len(tldr) > 200 {
		tldr = tldr[:197] + "..."
	}
	if tldr == "" {
		tldr = art.Title + ". Discovered from " + art.Source + "."
	}

	return NewsItem{
		ID:           uuid.New().String(),
		Title:        art.Title,
		URL:          art.Link,
		Domain:       dom,
		PublishedAt:  pubDate,
		Category:     cat,
		Score:        0.88,
		TLDR:         tldr,
		WhyItMatters: "High impact artificial intelligence signal and community insight.",
		Tags:         tags,
	}
}

func mergeNewsItems(existing []NewsItem, incoming []NewsItem) []NewsItem {
	seenUrls := make(map[string]bool)
	seenTitles := make(map[string]bool)

	for _, it := range existing {
		norm := normalizeURL(it.URL)
		titleKey := strings.ToLower(strings.TrimSpace(it.Title))
		if norm != "" {
			seenUrls[norm] = true
		}
		if titleKey != "" {
			seenTitles[titleKey] = true
		}
	}

	var genuinelyNew []NewsItem
	for _, it := range incoming {
		norm := normalizeURL(it.URL)
		titleKey := strings.ToLower(strings.TrimSpace(it.Title))
		if norm == "" || seenUrls[norm] || seenTitles[titleKey] {
			continue
		}
		seenUrls[norm] = true
		seenTitles[titleKey] = true
		genuinelyNew = append(genuinelyNew, it)
	}

	// Prepend new unique items at top, keep existing below
	merged := append(genuinelyNew, existing...)
	if len(merged) > 150 {
		merged = merged[:150]
	}
	return merged
}

func writeJSONOutput(outputPath string, newsData NewsData) {
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// #nosec G304
	file, err := os.Create(filepath.Clean(outputPath))
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("Error closing file: %v", cerr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(newsData); err != nil {
		log.Fatalf("Failed to write JSON: %v", err)
	}

	log.Printf("Successfully wrote %d items to %s", newsData.TotalCount, outputPath)
}

func syncToAPIPath(newsData NewsData) {
	apiDir := "api"
	apiPath := filepath.Join(apiDir, "news.json")
	if err := os.MkdirAll(apiDir, 0o750); err != nil {
		return
	}

	// #nosec G304
	file, err := os.Create(filepath.Clean(apiPath))
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(newsData)
	log.Printf("Successfully synced %d items to %s", newsData.TotalCount, apiPath)
}

func fetchLatestNews(ctx context.Context, database *db.DB) ([]NewsItem, error) {
	query := `
		SELECT 
			i.id,
			i.title,
			i.url,
			i.domain,
			i.published_at,
			i.category,
			COALESCE(s.final, 0.85) as score,
			COALESCE(sm.tldr, '') as tldr,
			COALESCE(sm.why_it_matters, '') as why_it_matters,
			COALESCE(sm.tags, '[]'::jsonb) as tags
		FROM items i
		LEFT JOIN scores s ON i.id = s.item_id
		LEFT JOIN summaries sm ON i.id = sm.item_id
		ORDER BY i.published_at DESC
		LIMIT 100
	`

	rows, err := database.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []NewsItem
	for rows.Next() {
		var item NewsItem
		var tagsJSON []byte

		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.URL,
			&item.Domain,
			&item.PublishedAt,
			&item.Category,
			&item.Score,
			&item.TLDR,
			&item.WhyItMatters,
			&tagsJSON,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		if len(tagsJSON) > 0 {
			if err := json.Unmarshal(tagsJSON, &item.Tags); err != nil {
				item.Tags = []string{}
			}
		}

		items = append(items, item)
	}

	return items, rows.Err()
}
