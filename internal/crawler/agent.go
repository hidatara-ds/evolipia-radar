package crawler

import (
	"context"
	"time"
)

// Article represents a raw discovered source before being processed into a cluster.
type Article struct {
	Title       string
	Content     string
	Link        string
	PublishedAt time.Time
	Source      string
}

// DiscoveryAgent defines the contract for any source crawler (RSS, Trending, Search, DeepCrawl).
type DiscoveryAgent interface {
	Name() string
	// Crawl executes the discovery hook. It respects the crawler budget provided via the context or budget manager.
	Crawl(ctx context.Context, maxItems int) ([]Article, error)
}
