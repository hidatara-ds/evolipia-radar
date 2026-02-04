package services

import (
	"context"
	"time"

	"github.com/evolipia/radar/internal/db"
	"github.com/evolipia/radar/internal/models"
)

type FeedService struct {
	itemRepo    *db.ItemRepository
	signalRepo  *db.SignalRepository
	scoreRepo   *db.ScoreRepository
	summaryRepo *db.SummaryRepository
}

func NewFeedService(database *db.DB) *FeedService {
	return &FeedService{
		itemRepo:    db.NewItemRepository(database),
		signalRepo:  db.NewSignalRepository(database),
		scoreRepo:   db.NewScoreRepository(database),
		summaryRepo: db.NewSummaryRepository(database),
	}
}

func (s *FeedService) BuildFeedResponse(ctx context.Context, items []models.Item, date time.Time, topic *string) map[string]interface{} {
	responseItems := make([]map[string]interface{}, 0, len(items))

	for rank, item := range items {
		score, _ := s.scoreRepo.GetByItemID(ctx, item.ID)
		summary, _ := s.summaryRepo.GetByItemID(ctx, item.ID)

		itemResp := map[string]interface{}{
			"id":           item.ID,
			"rank":         rank + 1,
			"title":        item.Title,
			"url":          item.URL,
			"domain":       item.Domain,
			"published_at": item.PublishedAt.Format(time.RFC3339),
			"scores": map[string]float64{
				"final":       0.0,
				"hot":         0.0,
				"relevance":   0.0,
				"credibility": 0.0,
				"novelty":     0.0,
			},
			"summary": map[string]interface{}{
				"tldr":           "",
				"why_it_matters": "",
				"tags":           []string{},
				"method":         "extractive",
			},
		}

		if score != nil {
			itemResp["scores"] = map[string]float64{
				"final":       score.Final,
				"hot":         score.Hot,
				"relevance":   score.Relevance,
				"credibility": score.Credibility,
				"novelty":     score.Novelty,
			}
		}

		if summary != nil {
			itemResp["summary"] = map[string]interface{}{
				"tldr":           summary.TLDR,
				"why_it_matters": summary.WhyItMatters,
				"tags":           summary.Tags,
				"method":         summary.Method,
			}
		}

		responseItems = append(responseItems, itemResp)
	}

	dateStr := date.Format("2006-01-02")
	topicStr := ""
	if topic != nil {
		topicStr = *topic
	}

	return map[string]interface{}{
		"date":  dateStr,
		"topic": topicStr,
		"items": responseItems,
	}
}

func (s *FeedService) BuildRisingResponse(ctx context.Context, items []models.Item, window time.Duration) map[string]interface{} {
	responseItems := make([]map[string]interface{}, 0, len(items))

	for rank, item := range items {
		signals, _ := s.signalRepo.GetRisingSignals(ctx, item.ID, window)

		pointsDelta := 0
		commentsDelta := 0

		if len(signals) >= 2 {
			first := signals[0]
			last := signals[len(signals)-1]

			if first.Points != nil && last.Points != nil {
				pointsDelta = *last.Points - *first.Points
			}
			if first.Comments != nil && last.Comments != nil {
				commentsDelta = *last.Comments - *first.Comments
			}
		}

		risingScore := float64(pointsDelta*10 + commentsDelta*5)

		itemResp := map[string]interface{}{
			"id":           item.ID,
			"rank":         rank + 1,
			"title":        item.Title,
			"url":          item.URL,
			"domain":       item.Domain,
			"published_at": item.PublishedAt.Format(time.RFC3339),
			"rising_score": risingScore,
			"signals": map[string]int{
				"points_delta":   pointsDelta,
				"comments_delta": commentsDelta,
			},
		}

		responseItems = append(responseItems, itemResp)
	}

	return map[string]interface{}{
		"window": window.String(),
		"items":  responseItems,
	}
}
