package crawler

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/pkg/ai"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
	"github.com/hidatara-ds/evolipia-radar/pkg/models"
)

type Summarizer struct {
	aiSvc *ai.Service
	repo  *db.SummaryRepository
}

func NewSummarizer(aiSvc *ai.Service, pool *db.DB) *Summarizer {
	return &Summarizer{
		aiSvc: aiSvc,
		repo:  db.NewSummaryRepository(pool),
	}
}

func (s *Summarizer) Process(ctx context.Context, itemID uuid.UUID, title, content string) error {
	// 1. Prepare AI prompt
	prompt := fmt.Sprintf("Summarize this news article in a very short TL;DR (max 2 sentences) and explain 'Why it matters' for AI researchers.\nTitle: %s\nContent: %s", title, content)
	
	// 2. Call AI Service
	resp, err := s.aiSvc.Summarize(ctx, ai.SummarizeRequest{
		Text:        content,
		Instruction: prompt,
	})

	var summaryContent, whyMatters string
	var tags []string
	method := "llm-openrouter"

	if err != nil {
		// Fallback if AI fails or budget exhausted
		summaryContent = "Summary pending system capacity."
		whyMatters = "Research discovery."
		tags = []string{"AI", "Research"}
		method = "fallback"
	} else {
		// Simple parsing: split by 'Why it matters:' or similar if possible, 
		// but since we want it structured, we'll just put it in TLDR for now
		summaryContent = resp.Summary
		whyMatters = "Refer to summary for details."
		tags = []string{"AI", "Innovation"} // Extracting actual tags would need another AI call or parsing
	}

	summary := &models.Summary{
		ItemID:       itemID,
		TLDR:         summaryContent,
		WhyItMatters: whyMatters, // Future: parse from AI output
		Tags:         tags,
		Method:       method,
	}

	return s.repo.Upsert(ctx, summary)
}
