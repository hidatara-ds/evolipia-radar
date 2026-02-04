package models

import (
	"time"

	"github.com/google/uuid"
)

type Source struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`     // hacker_news, rss_atom, arxiv, json_api
	Category        string    `json:"category"` // news, web
	URL             string    `json:"url"`
	MappingJSON     []byte    `json:"mapping_json,omitempty"`
	Enabled         bool      `json:"enabled"`
	Status          string    `json:"status"` // active, pending, failed
	LastTestStatus  *string   `json:"last_test_status,omitempty"`
	LastTestMessage *string   `json:"last_test_message,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Item struct {
	ID          uuid.UUID `json:"id"`
	SourceID    uuid.UUID `json:"source_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	ContentHash string    `json:"content_hash"`
	Domain      string    `json:"domain"`
	Category    string    `json:"category"`
	RawExcerpt  *string   `json:"raw_excerpt,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Signal struct {
	ID        uuid.UUID `json:"id"`
	ItemID    uuid.UUID `json:"item_id"`
	Points    *int      `json:"points,omitempty"`
	Comments  *int      `json:"comments,omitempty"`
	RankPos   *int      `json:"rank_pos,omitempty"`
	FetchedAt time.Time `json:"fetched_at"`
}

type Score struct {
	ItemID      uuid.UUID `json:"item_id"`
	Hot         float64   `json:"hot"`
	Relevance   float64   `json:"relevance"`
	Credibility float64   `json:"credibility"`
	Novelty     float64   `json:"novelty"`
	Final       float64   `json:"final"`
	ComputedAt  time.Time `json:"computed_at"`
}

type Summary struct {
	ItemID       uuid.UUID `json:"item_id"`
	TLDR         string    `json:"tldr"`
	WhyItMatters string    `json:"why_it_matters"`
	Tags         []string  `json:"tags"`
	Method       string    `json:"method"` // extractive, llm
	CreatedAt    time.Time `json:"created_at"`
}

type FetchRun struct {
	ID            uuid.UUID `json:"id"`
	SourceID      uuid.UUID `json:"source_id"`
	FetchedAt     time.Time `json:"fetched_at"`
	Status        string    `json:"status"` // success, failed
	Error         *string   `json:"error,omitempty"`
	ItemsFetched  int       `json:"items_fetched"`
	ItemsInserted int       `json:"items_inserted"`
}

type ContentItem struct {
	Title       string
	URL         string
	PublishedAt time.Time
	Excerpt     string
	Domain      string
	Category    string
	Points      *int
	Comments    *int
	RankPos     *int
	Tags        []string
}
