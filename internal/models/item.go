// Package models defines data transfer objects and database schema models.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Item represents a scraped or ingested content piece stored in database.
type Item struct {
	ID             uuid.UUID  `json:"id"`
	SourceID       uuid.UUID  `json:"source_id"`
	SourceName     string     `json:"source_name,omitempty"`
	Title          string     `json:"title"`
	URL            string     `json:"url"`
	PublishedAt    time.Time  `json:"published_at"`
	ContentHash    string     `json:"content_hash"`
	Domain         string     `json:"domain"`
	Category       string     `json:"category"`
	RawExcerpt     *string    `json:"raw_excerpt,omitempty"`
	CrawlStatus    string     `json:"crawl_status"` // pending, processing, done, failed
	CrawlError     *string    `json:"crawl_error,omitempty"`
	RelevanceScore int        `json:"relevance_score"` // 0 - 100
	ValidatedAt    *time.Time `json:"validated_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`

	// Joined/computed fields
	Score       *Score   `json:"score,omitempty"`
	Summary     *Summary `json:"summary,omitempty"`
	Source      *Source  `json:"source,omitempty"`
	HeatLevel   string   `json:"heat_level,omitempty"`
	ScaledScore float64  `json:"scaled_score,omitempty"` // 1.0 - 10.0 scale
}

// CrawlProgressEvent represents a real-time progress payload broadcasted over SSE.
type CrawlProgressEvent struct {
	Step                   int       `json:"step"`                     // 1 to 6
	Message                string    `json:"message"`                  // Informative status string
	Progress               int       `json:"progress"`                 // Percentage 0-100
	CurrentSource          string    `json:"current_source,omitempty"` // Active source name
	TotalSources           int       `json:"total_sources"`
	ProcessedItems         int       `json:"processed_items"`
	EstimatedRemainingSecs int       `json:"estimated_remaining_secs"`
	IsComplete             bool      `json:"is_complete"`
	HasError               bool      `json:"has_error"`
	Error                  string    `json:"error,omitempty"`
	Timestamp              time.Time `json:"timestamp"`
}

// ItemQueryParams contains query parameters for filtering and paginating items.
type ItemQueryParams struct {
	Search       string   `form:"search" json:"search"`
	DateFrom     string   `form:"date_from" json:"date_from"`
	DateTo       string   `form:"date_to" json:"date_to"`
	Sources      []string `form:"sources[]" json:"sources"`
	Categories   []string `form:"categories[]" json:"categories"`
	MinRelevance int      `form:"min_relevance" json:"min_relevance"`
	Status       string   `form:"status" json:"status"`
	SortBy       string   `form:"sort_by" json:"sort_by"`       // date, relevance, credibility, impact
	SortOrder    string   `form:"sort_order" json:"sort_order"` // asc, desc
	Page         int      `form:"page" json:"page"`
	Limit        int      `form:"limit" json:"limit"`
}

// PaginatedItemsResponse represents the HTTP API response for item listings.
type PaginatedItemsResponse struct {
	Success       bool    `json:"success"`
	Data          []Item  `json:"data"`
	TotalCount    int64   `json:"total_count"`
	FilteredCount int64   `json:"filtered_count"`
	Page          int     `json:"page"`
	TotalPages    int     `json:"total_pages"`
	LastUpdated   string  `json:"last_updated"`
	Error         *string `json:"error,omitempty"`
}
