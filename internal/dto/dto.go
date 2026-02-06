package dto

import "time"

// ContentItem is a DTO for transferring content items from connectors
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

// TestResult is a DTO for source connection test results
type TestResult struct {
	Status       string                   `json:"status"`
	PreviewItems []map[string]interface{} `json:"preview_items,omitempty"`
	ErrorCode    string                   `json:"error_code,omitempty"`
	Message      string                   `json:"message,omitempty"`
}

