package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/connectors"
	"github.com/hidatara-ds/evolipia-radar/internal/db"
	"github.com/hidatara-ds/evolipia-radar/internal/dto"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/hidatara-ds/evolipia-radar/internal/normalizer"
	"github.com/hidatara-ds/evolipia-radar/internal/scoring"
	"github.com/hidatara-ds/evolipia-radar/internal/summarizer"
)

type Worker struct {
	db           *db.DB
	cfg          *config.Config
	sourceRepo   *db.SourceRepository
	itemRepo     *db.ItemRepository
	signalRepo   *db.SignalRepository
	scoreRepo    *db.ScoreRepository
	summaryRepo  *db.SummaryRepository
	fetchRunRepo *db.FetchRunRepository
}

func NewWorker(database *db.DB, cfg *config.Config) *Worker {
	return &Worker{
		db:           database,
		cfg:          cfg,
		sourceRepo:   db.NewSourceRepository(database),
		itemRepo:     db.NewItemRepository(database),
		signalRepo:   db.NewSignalRepository(database),
		scoreRepo:    db.NewScoreRepository(database),
		summaryRepo:  db.NewSummaryRepository(database),
		fetchRunRepo: db.NewFetchRunRepository(database),
	}
}

func (w *Worker) RunIngestion(ctx context.Context) error {
	sources, err := w.sourceRepo.GetEnabled(ctx)
	if err != nil {
		return fmt.Errorf("failed to get enabled sources: %w", err)
	}

	log.Printf("Found %d enabled sources", len(sources))

	for _, source := range sources {
		if err := w.processSource(ctx, source); err != nil {
			log.Printf("Error processing source %s: %v", source.Name, err)
			// Continue with other sources
		}
	}

	return nil
}

func (w *Worker) processSource(ctx context.Context, source models.Source) error {
	log.Printf("Processing source: %s (%s)", source.Name, source.Type)

	fetchRun := &models.FetchRun{
		SourceID:      source.ID,
		Status:        "success",
		ItemsFetched:  0,
		ItemsInserted: 0,
	}

	var items []dto.ContentItem
	var err error

	// Fetch items based on source type
	switch source.Type {
	case "hacker_news":
		items, err = connectors.FetchHackerNews(ctx, w.cfg)
	case "rss_atom":
		items, err = connectors.FetchRSSAtom(ctx, source.URL, w.cfg)
	case "arxiv":
		// Default query for AI/ML papers
		query := "cat:cs.AI OR cat:cs.LG OR cat:cs.CV OR cat:cs.CL"
		items, err = connectors.FetchArxiv(ctx, query, w.cfg)
	case "json_api":
		var mapping map[string]interface{}
		if source.MappingJSON != nil {
			if err := json.Unmarshal(source.MappingJSON, &mapping); err != nil {
				return fmt.Errorf("invalid mapping_json: %w", err)
			}
		}
		items, err = connectors.FetchJSONAPI(ctx, source.URL, mapping, w.cfg)
	default:
		return fmt.Errorf("unsupported source type: %s", source.Type)
	}

	if err != nil {
		fetchRun.Status = "failed"
		errorMsg := err.Error()
		fetchRun.Error = &errorMsg
		w.fetchRunRepo.Create(ctx, fetchRun)
		return err
	}

	fetchRun.ItemsFetched = len(items)
	log.Printf("Fetched %d items from %s", len(items), source.Name)

	// Process each item: normalize, dedup, store
	inserted := 0
	for _, contentItem := range items {
		if err := w.processItem(ctx, source, contentItem); err != nil {
			log.Printf("Error processing item %s: %v", contentItem.Title, err)
			continue
		}
		inserted++
	}

	fetchRun.ItemsInserted = inserted
	if err := w.fetchRunRepo.Create(ctx, fetchRun); err != nil {
		log.Printf("Error creating fetch run: %v", err)
	}

	log.Printf("Inserted %d new items from %s", inserted, source.Name)

	// Compute scores for new items
	if err := w.computeScores(ctx); err != nil {
		log.Printf("Error computing scores: %v", err)
	}

	return nil
}

func (w *Worker) processItem(ctx context.Context, source models.Source, contentItem dto.ContentItem) error {
	// Normalize URL and compute content hash
	normalizedURL, err := normalizer.NormalizeURL(contentItem.URL)
	if err != nil {
		return fmt.Errorf("failed to normalize URL: %w", err)
	}

	contentHash := normalizer.ContentHash(contentItem.Title, normalizedURL)

	// Check for duplicates
	existing, err := w.itemRepo.GetByContentHash(ctx, contentHash)
	if err != nil {
		return fmt.Errorf("failed to check duplicate: %w", err)
	}

	var item *models.Item
	if existing != nil {
		// Item already exists, update signals if applicable
		item = existing
	} else {
		// Create new item
		item = &models.Item{
			SourceID:    source.ID,
			Title:       contentItem.Title,
			URL:         normalizedURL,
			PublishedAt: contentItem.PublishedAt,
			ContentHash: contentHash,
			Domain:      contentItem.Domain,
			Category:    source.Category,
		}
		if contentItem.Excerpt != "" {
			item.RawExcerpt = &contentItem.Excerpt
		}

		if err := w.itemRepo.Create(ctx, item); err != nil {
			return fmt.Errorf("failed to create item: %w", err)
		}

		// Generate summary
		summary := summarizer.GenerateExtractiveSummary(item)
		if err := w.summaryRepo.Upsert(ctx, summary); err != nil {
			log.Printf("Error creating summary: %v", err)
		}
	}

	// Store signals (points, comments, rank) if available
	if contentItem.Points != nil || contentItem.Comments != nil || contentItem.RankPos != nil {
		signal := &models.Signal{
			ItemID:   item.ID,
			Points:   contentItem.Points,
			Comments: contentItem.Comments,
			RankPos:  contentItem.RankPos,
		}
		if err := w.signalRepo.Create(ctx, signal); err != nil {
			log.Printf("Error creating signal: %v", err)
		}
	}

	return nil
}

func (w *Worker) computeScores(ctx context.Context) error {
	// Get items from the last 7 days that need scoring
	items, err := w.itemRepo.GetItemsNeedingScoring(ctx, 7, 1000)
	if err != nil {
		return fmt.Errorf("failed to get items needing scoring: %w", err)
	}

	log.Printf("Computing scores for %d items", len(items))

	for _, item := range items {
		// Get latest signal
		signal, _ := w.signalRepo.GetLatestByItemID(ctx, item.ID)

		// Get summary
		summary, _ := w.summaryRepo.GetByItemID(ctx, item.ID)

		// Compute score
		score := scoring.ComputeScore(&item, signal, summary, scoring.DefaultWeights)

		// Store score
		if err := w.scoreRepo.Upsert(ctx, score); err != nil {
			log.Printf("Error upserting score: %v", err)
			continue
		}
	}

	return nil
}
