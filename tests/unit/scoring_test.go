package scoring_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/pkg/models"
	"github.com/hidatara-ds/evolipia-radar/pkg/scoring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeScore_Basic(t *testing.T) {
	item := &models.Item{
		ID:          uuid.New(),
		Title:       "Popular AI Article",
		URL:         "https://example.com/ai",
		Domain:      "openai.com",
		PublishedAt: time.Now().Add(-1 * time.Hour),
	}

	points := 1000
	comments := 100
	signal := &models.Signal{
		ID:        uuid.New(),
		ItemID:    item.ID,
		Points:    &points,
		Comments:  &comments,
		FetchedAt: time.Now(),
	}

	summary := &models.Summary{
		ItemID:       item.ID,
		TLDR:         "AI breakthrough",
		WhyItMatters: "Important for ML engineers",
		Tags:         []string{"llm", "mlops"},
		Method:       "extractive",
		CreatedAt:    time.Now(),
	}

	score := scoring.ComputeScore(item, signal, summary, nil, scoring.DefaultWeights)

	require.NotNil(t, score)
	assert.Greater(t, score.Final, 0.0)
	assert.LessOrEqual(t, score.Final, 1.0)
}



func TestComputeScore_Relevance(t *testing.T) {
	now := time.Now()

	itemRelevant := &models.Item{
		ID:          uuid.New(),
		Title:       "LLM inference optimization techniques",
		PublishedAt: now,
	}
	itemIrrelevant := &models.Item{
		ID:          uuid.New(),
		Title:       "Best Pasta Recipes for Dinner",
		PublishedAt: now,
	}

	signal := &models.Signal{
		ID:        uuid.New(),
		ItemID:    itemRelevant.ID,
		FetchedAt: now,
	}

	relevant := scoring.ComputeScore(itemRelevant, signal, nil, nil, scoring.DefaultWeights)
	irrelevant := scoring.ComputeScore(itemIrrelevant, signal, nil, nil, scoring.DefaultWeights)

	assert.Greater(t, relevant.Final, irrelevant.Final, "ML-related content should have higher score")
}

func TestComputeScore_Credibility(t *testing.T) {
	now := time.Now()

	itemHigh := &models.Item{
		ID:          uuid.New(),
		Title:       "Research Paper",
		URL:         "https://arxiv.org/abs/1234.5678",
		Domain:      "arxiv.org",
		PublishedAt: now,
	}
	itemMedium := &models.Item{
		ID:          uuid.New(),
		Title:       "Blog Post",
		URL:         "https://techcrunch.com/blog",
		Domain:      "techcrunch.com",
		PublishedAt: now,
	}
	itemLow := &models.Item{
		ID:          uuid.New(),
		Title:       "Opinion Piece",
		URL:         "https://unknown-blog.com/article",
		Domain:      "unknown-blog.com",
		PublishedAt: now,
	}

	signal := &models.Signal{
		ID:        uuid.New(),
		ItemID:    itemHigh.ID,
		FetchedAt: now,
	}

	high := scoring.ComputeScore(itemHigh, signal, nil, nil, scoring.DefaultWeights)
	medium := scoring.ComputeScore(itemMedium, signal, nil, nil, scoring.DefaultWeights)
	low := scoring.ComputeScore(itemLow, signal, nil, nil, scoring.DefaultWeights)

	assert.Greater(t, high.Final, medium.Final, "whitelisted domains should have highest credibility/score")
	assert.Greater(t, medium.Final, low.Final, "blacklisted domains should have lowest credibility/score")
}

func TestComputeScore_Novelty(t *testing.T) {
	now := time.Now()

	itemNew := &models.Item{
		ID:          uuid.New(),
		Title:       "Breaking News",
		PublishedAt: now,
	}
	itemDayOld := &models.Item{
		ID:          uuid.New(),
		Title:       "Yesterday News",
		PublishedAt: now.Add(-24 * time.Hour),
	}
	itemWeekOld := &models.Item{
		ID:          uuid.New(),
		Title:       "Last Week News",
		PublishedAt: now.Add(-7 * 24 * time.Hour),
	}

	signal := &models.Signal{
		ID:        uuid.New(),
		ItemID:    itemNew.ID,
		FetchedAt: now,
	}

	scoreNew := scoring.ComputeScore(itemNew, signal, nil, nil, scoring.DefaultWeights)
	scoreDayOld := scoring.ComputeScore(itemDayOld, signal, nil, nil, scoring.DefaultWeights)
	scoreWeekOld := scoring.ComputeScore(itemWeekOld, signal, nil, nil, scoring.DefaultWeights)

	assert.Greater(t, scoreNew.Final, scoreDayOld.Final)
	assert.Greater(t, scoreDayOld.Final, scoreWeekOld.Final)
}

func BenchmarkComputeScore(b *testing.B) {
	item := &models.Item{
		ID:          uuid.New(),
		Title:       "Benchmark Article",
		URL:         "https://example.com/benchmark",
		Domain:      "openai.com",
		PublishedAt: time.Now(),
	}

	points := 100
	comments := 10
	signal := &models.Signal{
		ID:        uuid.New(),
		ItemID:    item.ID,
		Points:    &points,
		Comments:  &comments,
		FetchedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = scoring.ComputeScore(item, signal, nil, nil, scoring.DefaultWeights)
	}
}

func FuzzComputeScore(f *testing.F) {
	now := time.Now()

	// Seed corpus
	f.Add("Title", "example.com", now.Unix())
	f.Add("Another Title", "openai.com", now.Add(-time.Hour).Unix())

	f.Fuzz(func(t *testing.T, title, domain string, publishedUnix int64) {
		published := time.Unix(publishedUnix, 0)
		item := &models.Item{
			ID:          uuid.New(),
			Title:       title,
			Domain:      domain,
			PublishedAt: published,
		}

	score := scoring.ComputeScore(item, nil, nil, nil, scoring.DefaultWeights)
		require.GreaterOrEqual(t, score.Final, 0.0)
		require.LessOrEqual(t, score.Final, 1.0)
	})
}
