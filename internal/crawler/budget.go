package crawler

import (
	"log"
	"sync"
	"time"
)

// CrawlBudget enforces strict limits on ingestion rates and avoids duplicate crawling.
type CrawlBudget struct {
	mu sync.Mutex

	seenCache map[string]time.Time // Basic deduplication by URL. (Use Redis for scale).
	
	hourlyIngested   int
	maxHourlyIngests int
	lastReset        time.Time

	metrics *Metrics
}

// Metrics monitors health and scale of the ingestion engine
type Metrics struct {
	mu sync.Mutex

	ArticlesProcessed int      `json:"articles_processed"`
	FilteredArticles  int      `json:"filtered_articles"`
	APIHits           int      `json:"api_hits"`
	ClustersCount     int      `json:"clusters"`
	AvgClusterScore   float64  `json:"avg_cluster_score"`
	TopClusterTitles  []string `json:"top_cluster_titles"`
}

// AddProcessed increments safely
func (m *Metrics) AddProcessed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ArticlesProcessed++
}

// AddFiltered increments safely
func (m *Metrics) AddFiltered() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FilteredArticles++
}

// NewCrawlBudget initializes the budget restrictions for crawlers.
func NewCrawlBudget(maxHourly int, m *Metrics) *CrawlBudget {
	return &CrawlBudget{
		seenCache:        make(map[string]time.Time),
		maxHourlyIngests: maxHourly,
		lastReset:        time.Now(),
		metrics:          m,
	}
}

// Consume check if we can ingest a URL. Returns false if deduplicated or budget exhausted.
func (b *CrawlBudget) Consume(url string) bool {
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
		b.metrics.AddFiltered()
		return false // Deduplicated
	}

	// Allowed
	b.seenCache[url] = now
	b.hourlyIngested++
	b.metrics.AddProcessed()
	return true
}

func (b *CrawlBudget) LogStatus() {
	b.mu.Lock()
	defer b.mu.Unlock()
	log.Printf("[BUDGET] Crawl Status: %d / %d allowed this hour. Cache size: %d", b.hourlyIngested, b.maxHourlyIngests, len(b.seenCache))
}
