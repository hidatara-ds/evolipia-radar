package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/evolipia/radar/internal/config"
	"github.com/evolipia/radar/internal/connectors"
	"github.com/evolipia/radar/internal/db"
	"github.com/evolipia/radar/internal/security"
)

type TestResult struct {
	Status       string                   `json:"status"`
	PreviewItems []map[string]interface{} `json:"preview_items,omitempty"`
	ErrorCode    string                   `json:"error_code,omitempty"`
	Message      string                   `json:"message,omitempty"`
}

type SourceService struct {
	db *db.DB
}

func NewSourceService(database *db.DB) *SourceService {
	return &SourceService{db: database}
}

func (s *SourceService) TestConnection(ctx context.Context, sourceType, category, url string, mappingJSON json.RawMessage) (*TestResult, error) {
	cfg := config.Load()

	// SSRF protection
	if err := security.ValidateURL(url); err != nil {
		return &TestResult{
			Status:    "failed",
			ErrorCode: "SSRF_BLOCKED",
			Message:   err.Error(),
		}, nil
	}

	// Fetch and parse
	var items []connectors.ContentItem
	var err error

	switch sourceType {
	case "rss_atom":
		items, err = connectors.FetchRSSAtom(ctx, url, cfg)
	case "json_api":
		var mapping map[string]interface{}
		if mappingJSON != nil {
			if err := json.Unmarshal(mappingJSON, &mapping); err != nil {
				return &TestResult{
					Status:    "failed",
					ErrorCode: "MAPPING_ERROR",
					Message:   fmt.Sprintf("invalid mapping_json: %v", err),
				}, nil
			}
		}
		items, err = connectors.FetchJSONAPI(ctx, url, mapping, cfg)
	default:
		return &TestResult{
			Status:    "failed",
			ErrorCode: "INVALID_FORMAT",
			Message:   fmt.Sprintf("unsupported source type: %s", sourceType),
		}, nil
	}

	if err != nil {
		errorCode := "TEST_ERROR"
		message := err.Error()

		if err == connectors.ErrTimeout {
			errorCode = "TIMEOUT"
		} else if err == connectors.ErrSizeLimit {
			errorCode = "HTTP_403" // or SIZE_LIMIT
		}

		return &TestResult{
			Status:    "failed",
			ErrorCode: errorCode,
			Message:   message,
		}, nil
	}

	if len(items) < 3 {
		return &TestResult{
			Status:    "failed",
			ErrorCode: "INSUFFICIENT_ITEMS",
			Message:   fmt.Sprintf("found only %d items, need at least 3", len(items)),
		}, nil
	}

	// Build preview (max 5 items)
	previewCount := len(items)
	if previewCount > 5 {
		previewCount = 5
	}

	previewItems := make([]map[string]interface{}, 0, previewCount)
	for i := 0; i < previewCount; i++ {
		previewItems = append(previewItems, map[string]interface{}{
			"title":        items[i].Title,
			"url":          items[i].URL,
			"published_at": items[i].PublishedAt.Format(time.RFC3339),
			"source_type":  sourceType,
			"category":     category,
		})
	}

	return &TestResult{
		Status:       "ok",
		PreviewItems: previewItems,
	}, nil
}
