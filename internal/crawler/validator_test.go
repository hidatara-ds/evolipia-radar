package crawler_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/internal/crawler"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateItem(t *testing.T) {
	v := crawler.NewValidator(30, []string{"llm", "agents", "vision"})

	excerpt := "This is a sufficiently long excerpt that contains detailed content for validation testing."

	validItem := &models.Item{
		ID:          uuid.New(),
		SourceID:    uuid.New(),
		Title:       "Valid Title With Enough Length",
		URL:         "https://example.com/article/1",
		PublishedAt: time.Now(),
		RawExcerpt:  &excerpt,
	}

	err := v.ValidateItem(validItem)
	assert.NoError(t, err)

	// Short title
	shortTitleItem := *validItem
	shortTitleItem.Title = "Short"
	err = v.ValidateItem(&shortTitleItem)
	assert.Error(t, err)

	// Invalid URL
	invalidURLItem := *validItem
	invalidURLItem.URL = "ftp://invalid-url.com"
	err = v.ValidateItem(&invalidURLItem)
	assert.Error(t, err)

	// Future Date
	futureDateItem := *validItem
	futureDateItem.PublishedAt = time.Now().Add(2 * time.Hour)
	err = v.ValidateItem(&futureDateItem)
	assert.Error(t, err)
}

func TestValidator_ScoreRelevance(t *testing.T) {
	v := crawler.NewValidator(30, []string{"llm", "agents", "vision"})

	score := v.ScoreRelevance("Building LLM Agents in Go", "This paper presents autonomous agents powered by LLM models.")
	assert.GreaterOrEqual(t, score, 30)

	irrelevantScore := v.ScoreRelevance("Unrelated Baking Recipe", "Baking delicious cookies step by step.")
	assert.Equal(t, 0, irrelevantScore)
}
