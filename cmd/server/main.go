// Package main is the entry point for the Evolipia Radar backend server.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/internal/api"
	"github.com/hidatara-ds/evolipia-radar/internal/crawler"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/hidatara-ds/evolipia-radar/pkg/config"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
)

type DiskNewsItem struct {
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

type DiskNewsData struct {
	Items       []DiskNewsItem `json:"items"`
	LastUpdated time.Time      `json:"last_updated"`
	TotalCount  int            `json:"total_count"`
}

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	slog.Info("Starting Evolipia Radar Server...")

	cfg := config.Load()

	// Initialize Database connection (if database URL provided)
	var database *db.DB
	if cfg.DatabaseURL != "" {
		dbInstance, err := db.New(cfg)
		if err != nil {
			slog.Warn("Database connection failed, starting with in-memory mode", "err", err)
		} else {
			database = dbInstance
			defer database.Close()
			slog.Info("Database connected successfully")
		}
	}

	// Initialize SSE Progress Broadcaster
	broadcaster := api.NewProgressBroadcaster()

	// Initialize Validator & Agents
	validator := crawler.NewValidator(cfg.MinRelevanceScore, cfg.TopicKeywords)
	rssAgent := crawler.NewRSSAgent()
	trendingAgent := crawler.NewTrendingAgent()

	// Crawl Task definition with live multi-agent discovery & smart merging
	crawlTaskFunc := func(ctx context.Context, onProgress func(models.CrawlProgressEvent)) (int, error) {
		slog.Info("Starting crawl step 1: Initializing crawler...")
		onProgress(models.CrawlProgressEvent{
			Step:       1,
			Message:    "Initializing multi-source crawler...",
			Progress:   10,
			Timestamp:  time.Now(),
			IsComplete: false,
		})

		var discoveredArticles []crawler.Article

		// Step 2: RSS Feeds (Hacker News, ArXiv AI, TechCrunch AI, OpenAI Blog)
		slog.Info("Starting RSS feeds crawl", "step", 2)
		onProgress(models.CrawlProgressEvent{
			Step:          2,
			Message:       "Scanning RSS feeds (HackerNews, ArXiv AI, TechCrunch)...",
			Progress:      35,
			CurrentSource: "RSS Feeds",
			TotalSources:  2,
			Timestamp:     time.Now(),
		})
		rssItems, err := rssAgent.Crawl(ctx, 25)
		if err != nil {
			slog.Warn("RSS agent crawl warning", "err", err)
		} else {
			discoveredArticles = append(discoveredArticles, rssItems...)
		}

		// Step 3: Trending Feed (Algolia HN)
		slog.Info("Starting Trending feed crawl", "step", 3)
		onProgress(models.CrawlProgressEvent{
			Step:          3,
			Message:       "Scanning Algolia trending stories...",
			Progress:      65,
			CurrentSource: "Algolia HN Trending",
			TotalSources:  2,
			Timestamp:     time.Now(),
		})
		trendItems, err := trendingAgent.Crawl(ctx, 20)
		if err != nil {
			slog.Warn("Trending agent crawl warning", "err", err)
		} else {
			discoveredArticles = append(discoveredArticles, trendItems...)
		}

		// Step 4: Validate and score relevance
		slog.Info("Crawl step 4: Validating and scoring data...", "raw_count", len(discoveredArticles))
		onProgress(models.CrawlProgressEvent{
			Step:       4,
			Message:    "Validating & deduplicating incoming signals...",
			Progress:   80,
			Timestamp:  time.Now(),
			IsComplete: false,
		})

		var validItems []models.Item
		for _, art := range discoveredArticles {
			parsedURL, parseErr := url.Parse(art.Link)
			domain := "external"
			if parseErr == nil {
				domain = strings.TrimPrefix(parsedURL.Hostname(), "www.")
			}

			candidate := &models.Item{
				ID:          uuid.New(),
				Title:       art.Title,
				URL:         art.Link,
				Domain:      domain,
				SourceName:  art.Source,
				PublishedAt: art.PublishedAt,
				RawExcerpt:  stringPtr(art.Content),
				Category:    "llm",
				CrawlStatus: "done",
				CreatedAt:   time.Now(),
			}

			if candidate.PublishedAt.IsZero() {
				candidate.PublishedAt = time.Now()
			}

			if valErr := validator.ValidateItem(candidate); valErr != nil {
				continue
			}

			excerptText := ""
			if candidate.RawExcerpt != nil {
				excerptText = *candidate.RawExcerpt
			}
			score := validator.ScoreRelevance(candidate.Title, excerptText)
			candidate.RelevanceScore = score
			candidate.ScaledScore = float64(score) / 10.0

			validItems = append(validItems, *candidate)
		}

		// Step 5: Merge into data/news.json & api/news.json (and DB if connected)
		slog.Info("Crawl step 5: Merging into feeds and database...", "valid_count", len(validItems))
		onProgress(models.CrawlProgressEvent{
			Step:       5,
			Message:    "Merging signals into JSON feeds & database...",
			Progress:   90,
			Timestamp:  time.Now(),
			IsComplete: false,
		})

		newMergedCount := mergeAndSaveServerNews(validItems)

		// Upsert to PostgreSQL if connected
		if database != nil && len(validItems) > 0 {
			upsertItemsToDB(ctx, database, validItems)
		}

		return newMergedCount, nil
	}

	// Initialize & Start Auto-Scheduler
	scheduler, err := crawler.NewScheduler(cfg.CrawlInterval, crawlTaskFunc, broadcaster.Broadcast)
	if err != nil {
		slog.Error("Failed to initialize scheduler", "err", err)
		os.Exit(1)
	}
	scheduler.Start()

	// Initialize Gin HTTP Router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(api.CORS())
	router.Use(api.Logger())

	// Serve Web UI static files
	router.StaticFile("/", "./web/index.html")
	router.Static("/web", "./web")

	// Handlers
	itemsHandler := api.NewItemsHandler(database)

	// API Routes
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
	})

	router.GET("/api/items", itemsHandler.HandleGetItems)
	router.GET("/api/crawl/progress", broadcaster.HandleSSEProgress())

	router.GET("/api/crawl/status", func(c *gin.Context) {
		lastTime, status, count, errStr, isRunning := scheduler.GetStatus()
		c.JSON(http.StatusOK, gin.H{
			"last_run_time":    lastTime,
			"last_run_status":  status,
			"last_items_count": count,
			"last_error":       errStr,
			"is_running":       isRunning,
		})
	})

	router.POST("/api/crawl", func(c *gin.Context) {
		go func() {
			_, _ = scheduler.RunCrawl(context.Background(), "manual")
		}()
		c.JSON(http.StatusAccepted, gin.H{
			"status":  "triggered",
			"message": "Manual crawl triggered successfully",
		})
	})

	// HTTP Server Configuration
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("API Server running", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server forced to shutdown", "err", err)
		}
	}()

	// Graceful Shutdown on SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server gracefully...")

	// Stop scheduler first
	scheduler.Stop()

	// Shutdown HTTP Server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP Server forced to shutdown", "err", err)
	}

	slog.Info("Server stopped cleanly")
}

func stringPtr(s string) *string {
	return &s
}

func mergeAndSaveServerNews(validItems []models.Item) int {
	existingData := loadServerNews("data/news.json")

	seenUrls := make(map[string]bool)
	seenTitles := make(map[string]bool)

	normalizeURL := func(raw string) string {
		parsed, err := url.Parse(strings.TrimSpace(raw))
		if err != nil {
			return strings.ToLower(strings.TrimSpace(raw))
		}
		q := parsed.Query()
		for k := range q {
			if strings.HasPrefix(k, "utm_") || k == "ref" || k == "source" {
				q.Del(k)
			}
		}
		parsed.RawQuery = q.Encode()
		parsed.Fragment = ""
		res := parsed.String()
		return strings.ToLower(strings.TrimRight(res, "/"))
	}

	for _, it := range existingData.Items {
		norm := normalizeURL(it.URL)
		titleKey := strings.ToLower(strings.TrimSpace(it.Title))
		if norm != "" {
			seenUrls[norm] = true
		}
		if titleKey != "" {
			seenTitles[titleKey] = true
		}
	}

	var genuinelyNew []DiskNewsItem
	for _, it := range validItems {
		norm := normalizeURL(it.URL)
		titleKey := strings.ToLower(strings.TrimSpace(it.Title))
		if norm == "" || seenUrls[norm] || seenTitles[titleKey] {
			continue
		}
		seenUrls[norm] = true
		seenTitles[titleKey] = true

		tldr := ""
		if it.RawExcerpt != nil {
			tldr = *it.RawExcerpt
			if len(tldr) > 200 {
				tldr = tldr[:197] + "..."
			}
		}
		if tldr == "" {
			tldr = it.Title
		}

		score := float64(it.RelevanceScore) / 100.0
		if score < 0.75 {
			score = 0.85
		}

		genuinelyNew = append(genuinelyNew, DiskNewsItem{
			ID:           it.ID.String(),
			Title:        it.Title,
			URL:          it.URL,
			Domain:       it.Domain,
			PublishedAt:  it.PublishedAt,
			Category:     it.Category,
			Score:        score,
			TLDR:         tldr,
			WhyItMatters: "High signal AI development and community discussions.",
			Tags:         []string{it.Category, "feed"},
		})
	}

	merged := append(genuinelyNew, existingData.Items...)
	if len(merged) > 150 {
		merged = merged[:150]
	}

	finalData := DiskNewsData{
		Items:       merged,
		LastUpdated: time.Now(),
		TotalCount:  len(merged),
	}

	payload, err := json.MarshalIndent(finalData, "", "  ")
	if err == nil {
		_ = os.WriteFile("data/news.json", payload, 0o600)
		_ = os.WriteFile("api/news.json", payload, 0o600)
		slog.Info("Saved merged news to disk", "new_added", len(genuinelyNew), "total", len(merged))
	}

	return len(genuinelyNew)
}

func loadServerNews(primaryPath string) DiskNewsData {
	paths := []string{primaryPath, "api/news.json", "data/news.json"}
	for _, p := range paths {
		// #nosec G304
		data, err := os.ReadFile(filepath.Clean(p))
		if err == nil {
			var parsed DiskNewsData
			if jsonErr := json.Unmarshal(data, &parsed); jsonErr == nil && len(parsed.Items) > 0 {
				return parsed
			}
		}
	}
	return DiskNewsData{Items: []DiskNewsItem{}, LastUpdated: time.Now(), TotalCount: 0}
}

func upsertItemsToDB(ctx context.Context, database *db.DB, items []models.Item) {
	for _, it := range items {
		query := `
			INSERT INTO items (id, source_id, title, url, domain, category, published_at, raw_excerpt, content_hash)
			VALUES ($1, $1, $2, $3, $4, $5, $6, $7, MD5($3))
			ON CONFLICT (content_hash) DO NOTHING
		`
		_, _ = database.Pool.Exec(ctx, query, it.ID, it.Title, it.URL, it.Domain, it.Category, it.PublishedAt, it.RawExcerpt)
	}
}
