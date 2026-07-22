package crawler

import (
	"context"

	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/pkg/ai"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
	"github.com/hidatara-ds/evolipia-radar/pkg/models"
)

type Summarizer struct {
	aiSvc     *ai.Service
	repo      *db.SummaryRepository
	scoreRepo *db.ScoreRepository
}

func NewSummarizer(aiSvc *ai.Service, pool *db.DB) *Summarizer {
	return &Summarizer{
		aiSvc:     aiSvc,
		repo:      db.NewSummaryRepository(pool),
		scoreRepo: db.NewScoreRepository(pool),
	}
}

func (s *Summarizer) Process(ctx context.Context, itemID uuid.UUID, title, content string) error {
	// Call AI Service for structured analysis
	resp, err := s.aiSvc.AnalyzeArticle(ctx, ai.AnalyzeRequest{
		Title:   title,
		Content: content,
	})

	var summaryContent, whyMatters string
	var tags []string
	method := "llm-openrouter"

	var impact, engineeringValue, novelty float64
	var reasoning string

	if err != nil {
		// Fallback if AI fails or budget exhausted
		summaryContent = "Summary pending system capacity."
		whyMatters = "Research discovery."
		tags = []string{"AI", "Research"}
		method = "fallback"
	} else {
		summaryContent = resp.TLDR
		whyMatters = resp.WhyItMatters
		tags = []string{"AI", "Innovation"} // Keep simple for now

		impact = resp.Impact
		engineeringValue = resp.EngineeringValue
		novelty = resp.Novelty
		reasoning = resp.Reasoning
	}

	summary := &models.Summary{
		ItemID:       itemID,
		TLDR:         summaryContent,
		WhyItMatters: whyMatters,
		Tags:         tags,
		Method:       method,
	}

	if err := s.repo.Upsert(ctx, summary); err != nil {
		return err
	}

	// Update existing score with LLM metrics if AI was successful
	if err == nil {
		existingScore, _ := s.scoreRepo.GetByItemID(ctx, itemID)
		if existingScore == nil {
			existingScore = &models.Score{ItemID: itemID}
		}

		existingScore.Impact = impact
		existingScore.EngineeringValue = engineeringValue
		existingScore.Novelty = novelty
		existingScore.Reasoning = reasoning

		_ = s.scoreRepo.Upsert(ctx, existingScore)
	}

	return nil
}
