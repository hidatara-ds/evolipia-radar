package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hidatara-ds/evolipia-radar/pkg/config"
	"github.com/hidatara-ds/evolipia-radar/pkg/connectors"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
	"github.com/hidatara-ds/evolipia-radar/pkg/dto"
	"github.com/hidatara-ds/evolipia-radar/pkg/models"
	"github.com/hidatara-ds/evolipia-radar/pkg/normalizer"
	"github.com/hidatara-ds/evolipia-radar/pkg/scoring"
	"github.com/hidatara-ds/evolipia-radar/pkg/summarizer"
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

	if len(sources) == 0 {
		log.Println("No enabled sources found. Creating default Hacker News source...")
		if err := w.ensureDefaultSource(ctx); err != nil {
			log.Printf("Warning: Failed to create default source: %v", err)
		} else {
			sources, err = w.sourceRepo.GetEnabled(ctx)
			if err != nil {
				return fmt.Errorf("failed to get enabled sources after creating default: %w", err)
			}
		}
	}

	log.Printf("Found %d enabled sources", len(sources))

	if len(sources) == 0 {
		log.Println("No enabled sources to process")
		return nil
	}

	for _, source := range sources {
		if ctx.Err() != nil {
			log.Printf("Context canceled, stopping ingestion")
			return ctx.Err()
		}
		if err := w.processSource(ctx, source); err != nil {
			log.Printf("Error processing source %s: %v", source.Name, err)
		}
	}

	return nil
}

func (w *Worker) ensureDefaultSource(ctx context.Context) error {
	defaultSource := &models.Source{
		Name:     "Hacker News",
		Type:     "hacker_news",
		Category: "news",
		URL:      "https://news.ycombinator.com",
		Enabled:  true,
		Status:   "active",
	}

	allSources, err := w.sourceRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list sources: %w", err)
	}

	for _, s := range allSources {
		if s.URL == defaultSource.URL || s.Name == defaultSource.Name {
			if !s.Enabled {
				if err := w.sourceRepo.SetEnabled(ctx, s.ID, true, "active"); err != nil {
					log.Printf("Warning: Failed to enable existing source %s: %v", s.Name, err)
				} else {
					log.Printf("Enabled existing source: %s", s.Name)
					return nil
				}
			} else {
				return nil
			}
		}
	}

	if err := w.sourceRepo.Create(ctx, defaultSource); err != nil {
		return fmt.Errorf("failed to create default source: %w", err)
	}

	log.Printf("Created default source: %s", defaultSource.Name)
	return nil
}

func (w *Worker) processSource(ctx context.Context, source models.Source) error {
	log.Printf("Processing source: %s (%s)", source.Name, source.Type)

	fetchRun := &models.FetchRun{
		SourceID: source.ID,
		Status:   "success",
	}

	items, err := w.fetchItems(ctx, source, fetchRun)
	if err != nil {
		return err
	}

	fetchRun.ItemsFetched = len(items)
	log.Printf("Fetched %d items from %s", len(items), source.Name)

	inserted := w.processItems(ctx, source, items)
	fetchRun.ItemsInserted = inserted

	if err := w.fetchRunRepo.Create(ctx, fetchRun); err != nil {
		log.Printf("Error creating fetch run: %v", err)
	}

	log.Printf("Inserted %d new items from %s", inserted, source.Name)

	if err := w.computeScores(ctx); err != nil {
		log.Printf("Error computing scores: %v", err)
	}

	return nil
}

func (w *Worker) fetchItems(ctx context.Context, source models.Source, fetchRun *models.FetchRun) ([]dto.ContentItem, error) {
	var items []dto.ContentItem
	var err error

	switch source.Type {
	case "hacker_news", "hackernews":
		items, err = connectors.FetchHackerNews(ctx, w.cfg)
	case "rss_atom":
		items, err = connectors.FetchRSSAtom(ctx, source.URL, w.cfg)
	case "arxiv":
		query := "cat:cs.AI OR cat:cs.LG OR cat:cs.CV OR cat:cs.CL"
		items, err = connectors.FetchArxiv(ctx, query, w.cfg)
	case "huggingface":
		items, err = connectors.FetchHuggingFaceTrending(ctx, w.cfg)
	case "lmsys":
		items, err = connectors.FetchLMSYSArena(ctx, w.cfg)
	case "json_api":
		items, err = w.fetchJSONAPI(ctx, source)
	default:
		return nil, fmt.Errorf("unsupported source type: %s", source.Type)
	}

	if err != nil {
		fetchRun.Status = "failed"
		errorMsg := err.Error()
		if len(errorMsg) > 500 {
			errorMsg = errorMsg[:500] + "..."
		}
		fetchRun.Error = &errorMsg
		if createErr := w.fetchRunRepo.Create(ctx, fetchRun); createErr != nil {
			log.Printf("Warning: Failed to create fetch run record: %v", createErr)
		}
		return nil, err
	}

	return items, nil
}

func (w *Worker) fetchJSONAPI(ctx context.Context, source models.Source) ([]dto.ContentItem, error) {
	var mapping map[string]interface{}
	if source.MappingJSON != nil {
		if err := json.Unmarshal(source.MappingJSON, &mapping); err != nil {
			return nil, fmt.Errorf("invalid mapping_json: %w", err)
		}
	}
	return connectors.FetchJSONAPI(ctx, source.URL, mapping, w.cfg)
}

func (w *Worker) processItems(ctx context.Context, source models.Source, items []dto.ContentItem) int {
	inserted := 0
	for _, contentItem := range items {
		if err := w.processItem(ctx, source, contentItem); err != nil {
			log.Printf("Error processing item %s: %v", contentItem.Title, err)
			continue
		}
		inserted++
	}
	return inserted
}

func (w *Worker) processItem(ctx context.Context, source models.Source, contentItem dto.ContentItem) error {
	normalizedURL, err := normalizer.NormalizeURL(contentItem.URL)
	if err != nil {
		return fmt.Errorf("failed to normalize URL: %w", err)
	}

	contentHash := normalizer.ContentHash(contentItem.Title, normalizedURL)

	existing, err := w.itemRepo.GetByContentHash(ctx, contentHash)
	if err != nil {
		return fmt.Errorf("failed to check duplicate: %w", err)
	}

	var item *models.Item
	if existing != nil {
		item = existing
	} else {
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

		summary := summarizer.GenerateExtractiveSummary(item)
		if err := w.summaryRepo.Upsert(ctx, summary); err != nil {
			log.Printf("Error creating summary: %v", err)
		}
	}

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
	items, err := w.itemRepo.GetItemsNeedingScoring(ctx, 7, 1000)
	if err != nil {
		return fmt.Errorf("failed to get items needing scoring: %w", err)
	}

	log.Printf("Computing scores for %d items", len(items))

	for _, item := range items {
		signal, _ := w.signalRepo.GetLatestByItemID(ctx, item.ID)
		summary, _ := w.summaryRepo.GetByItemID(ctx, item.ID)
		score := scoring.ComputeScore(&item, signal, summary, scoring.DefaultWeights)

		if err := w.scoreRepo.Upsert(ctx, score); err != nil {
			log.Printf("Error upserting score: %v", err)
			continue
		}
	}

	return nil
}
