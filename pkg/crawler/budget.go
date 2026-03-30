package crawler

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CrawlBudget enforces strict limits on ingestion rates and avoids duplicate crawling.
type CrawlBudget struct {
	mu sync.Mutex

	seenCache map[string]time.Time // Basic deduplication by URL. (Use Redis for scale).

	hourlyIngested   int
	maxHourlyIngests int
	lastReset        time.Time

	metrics *Metrics
	db      *pgxpool.Pool
}

// Metrics monitors health and scale of the ingestion engine
type Metrics struct {
	mu sync.Mutex
	db *pgxpool.Pool

	ArticlesProcessed int      `json:"articles_processed"`
	FilteredArticles  int      `json:"filtered_articles"`
	APIHits           int      `json:"api_hits"`
	ClustersCount     int      `json:"clusters"`
	AvgClusterScore   float64  `json:"avg_cluster_score"`
	TopClusterTitles  []string `json:"top_cluster_titles"`
}

// NewMetrics creates a new Metrics instance with DB persistence
func NewMetrics(db *pgxpool.Pool) *Metrics {
	return &Metrics{db: db}
}

// LoadFromDB fetches global metrics from the database
func (m *Metrics) LoadFromDB(ctx context.Context) {
	if m.db == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	query := `SELECT articles_processed, filtered_articles, api_hits, clusters_count, avg_cluster_score, top_cluster_titles FROM global_metrics WHERE id = 1`
	var titlesJSON []byte
	err := m.db.QueryRow(ctx, query).Scan(&m.ArticlesProcessed, &m.FilteredArticles, &m.APIHits, &m.ClustersCount, &m.AvgClusterScore, &titlesJSON)
	if err != nil {
		log.Printf("[METRICS] Failed to load from DB: %v", err)
	}
	if titlesJSON != nil {
		_ = json.Unmarshal(titlesJSON, &m.TopClusterTitles)
	}
}

// AddProcessed increments safely and persists to DB
func (m *Metrics) AddProcessed(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ArticlesProcessed++

	if m.db != nil {
		_, err := m.db.Exec(ctx, `UPDATE global_metrics SET articles_processed = articles_processed + 1, updated_at = CURRENT_TIMESTAMP WHERE id = 1`)
		if err != nil {
			log.Printf("[METRICS] Failed to persist processed count: %v", err)
		}
	}
}

// AddFiltered increments safely and persists to DB
func (m *Metrics) AddFiltered(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FilteredArticles++

	if m.db != nil {
		_, err := m.db.Exec(ctx, `UPDATE global_metrics SET filtered_articles = filtered_articles + 1, updated_at = CURRENT_TIMESTAMP WHERE id = 1`)
		if err != nil {
			log.Printf("[METRICS] Failed to persist filtered count: %v", err)
		}
	}
}

// NewCrawlBudget initializes the budget restrictions for crawlers.
func NewCrawlBudget(maxHourly int, m *Metrics, db *pgxpool.Pool) *CrawlBudget {
	return &CrawlBudget{
		seenCache:        make(map[string]time.Time),
		maxHourlyIngests: maxHourly,
		lastReset:        time.Now(),
		metrics:          m,
		db:               db,
	}
}

// Consume check if we can ingest a URL. Returns false if deduplicated or budget exhausted.
func (b *CrawlBudget) Consume(ctx context.Context, url string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	// Reset hourly budget
	if now.Sub(b.lastReset) >= time.Hour {
		b.hourlyIngested = 0
		b.lastReset = now
		// Simple cache cleanup
		for k, v := range b.seenCache {
			if now.Sub(v) > 24*time.Hour {
				delete(b.seenCache, k)
			}
		}
	}

	if b.hourlyIngested >= b.maxHourlyIngests {
		return false // Rate limit hit
	}

	if _, exists := b.seenCache[url]; exists {
		b.metrics.AddFiltered(ctx)
		return false // Deduplicated
	}

	// Allowed
	b.seenCache[url] = now
	b.hourlyIngested++
	b.metrics.AddProcessed(ctx)
	return true
}

func (b *CrawlBudget) LogStatus() {
	b.mu.Lock()
	defer b.mu.Unlock()
	log.Printf("[BUDGET] Crawl Status: %d / %d allowed this hour. Cache size: %d", b.hourlyIngested, b.maxHourlyIngests, len(b.seenCache))
}

// UpdateClusterStats persists clustering stats to DB
func (m *Metrics) UpdateClusterStats(ctx context.Context, count int, avgScore float64, titles []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ClustersCount = count
	m.AvgClusterScore = avgScore
	m.TopClusterTitles = titles

	if m.db != nil {
		titlesJSON, _ := json.Marshal(titles)
		_, err := m.db.Exec(ctx, `
			UPDATE global_metrics 
			SET clusters_count = $1, avg_cluster_score = $2, top_cluster_titles = $3, updated_at = CURRENT_TIMESTAMP 
			WHERE id = 1
		`, count, avgScore, titlesJSON)
		if err != nil {
			log.Printf("[METRICS] Failed to persist cluster stats: %v", err)
		}
	}
}
