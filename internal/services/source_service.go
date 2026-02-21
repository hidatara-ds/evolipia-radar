package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/connectors"
	"github.com/hidatara-ds/evolipia-radar/internal/db"
	"github.com/hidatara-ds/evolipia-radar/internal/dto"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/hidatara-ds/evolipia-radar/internal/security"
)

type SourceService struct {
	db         *db.DB
	sourceRepo *db.SourceRepository
}

func NewSourceService(database *db.DB) *SourceService {
	return &SourceService{
		db:         database,
		sourceRepo: db.NewSourceRepository(database),
	}
}

func (s *SourceService) TestConnection(ctx context.Context, sourceType, category, url string, mappingJSON json.RawMessage) (*dto.TestResult, error) {
	cfg := config.Load()

	// SSRF protection
	if err := security.ValidateURL(url); err != nil {
		return &dto.TestResult{
			Status:    "failed",
			ErrorCode: "SSRF_BLOCKED",
			Message:   err.Error(),
		}, nil
	}

	// Fetch and parse
	var items []dto.ContentItem
	var err error

	switch sourceType {
	case "rss_atom":
		items, err = connectors.FetchRSSAtom(ctx, url, cfg)
	case "json_api":
		var mapping map[string]interface{}
		if mappingJSON != nil {
			if err := json.Unmarshal(mappingJSON, &mapping); err != nil {
				return &dto.TestResult{
					Status:    "failed",
					ErrorCode: "MAPPING_ERROR",
					Message:   fmt.Sprintf("invalid mapping_json: %v", err),
				}, nil
			}
		}
		items, err = connectors.FetchJSONAPI(ctx, url, mapping, cfg)
	default:
		return &dto.TestResult{
			Status:    "failed",
			ErrorCode: "INVALID_FORMAT",
			Message:   fmt.Sprintf("unsupported source type: %s", sourceType),
		}, nil
	}

	if err != nil {
		errorCode := "TEST_ERROR"
		message := err.Error()

		if errors.Is(err, connectors.ErrTimeout) {
			errorCode = "TIMEOUT"
		} else if errors.Is(err, connectors.ErrSizeLimit) {
			errorCode = "HTTP_403" // or SIZE_LIMIT
		}

		return &dto.TestResult{
			Status:    "failed",
			ErrorCode: errorCode,
			Message:   message,
		}, nil
	}

	if len(items) < 3 {
		return &dto.TestResult{
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

	return &dto.TestResult{
		Status:       "ok",
		PreviewItems: previewItems,
	}, nil
}

// SetEnabled handles the business logic for enabling/disabling a source
func (s *SourceService) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) error {
	// Business logic: set status based on enabled flag
	status := "pending"
	if enabled {
		status = "active"
	}
	return s.sourceRepo.SetEnabled(ctx, id, enabled, status)
}

// ListSources returns all sources
func (s *SourceService) ListSources(ctx context.Context) ([]models.Source, error) {
	return s.sourceRepo.List(ctx)
}

// GetSourceByID returns a source by ID
func (s *SourceService) GetSourceByID(ctx context.Context, id uuid.UUID) (*models.Source, error) {
	return s.sourceRepo.GetByID(ctx, id)
}

// CreateSource creates a new source
func (s *SourceService) CreateSource(ctx context.Context, source *models.Source) error {
	return s.sourceRepo.Create(ctx, source)
}

// UpdateTestStatus updates the test status of a source
func (s *SourceService) UpdateTestStatus(ctx context.Context, id uuid.UUID, status, message string) error {
	return s.sourceRepo.UpdateTestStatus(ctx, id, status, message)
}
