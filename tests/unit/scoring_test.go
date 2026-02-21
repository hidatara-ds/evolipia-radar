package scoring

import (
	"testing"
	"time"

	"github.com/hidatara-ds/evolipia-radar/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScoreCalculator_CalculateScore(t *testing.T) {
	config := DefaultScoringConfig()
	calculator := NewScoreCalculator(config)

	tests := []struct {
		name     string
		item     *dto.ContentItemDTO
		expected float64
	}{
		{
			name: "high popularity item",
			item: &dto.ContentItemDTO{
				Title:        "Popular AI Article",
				URL:          "https://example.com/ai",
				SourceType:   "hackernews",
				Popularity:   1000,
				CommentCount: 100,
				PublishedAt:  time.Now(),
			},
			expected: 0.8, // High score due to popularity
		},
		{
			name: "low popularity recent item",
			item: &dto.ContentItemDTO{
				Title:        "New ML Release",
				URL:          "https://example.com/ml",
				SourceType:   "rss",
				Popularity:   10,
				CommentCount: 0,
				PublishedAt:  time.Now(),
			},
			expected: 0.3, // Lower score due to low popularity
		},
		{
			name: "old item",
			item: &dto.ContentItemDTO{
				Title:        "Old Article",
				URL:          "https://example.com/old",
				SourceType:   "arxiv",
				Popularity:   100,
				CommentCount: 10,
				PublishedAt:  time.Now().Add(-7 * 24 * time.Hour), // 7 days old
			},
			expected: 0.4, // Lower score due to age
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculator.CalculateScore(tt.item)
			assert.InDelta(t, tt.expected, score, 0.1, "score should be approximately equal")
		})
	}
}

func TestScoreCalculator_CalculatePopularityScore(t *testing.T) {
	config := DefaultScoringConfig()
	calculator := NewScoreCalculator(config)

	tests := []struct {
		name     string
		points   int
		comments int
		expected float64
	}{
		{"zero engagement", 0, 0, 0.0},
		{"low engagement", 10, 2, 0.1},
		{"medium engagement", 100, 20, 0.3},
		{"high engagement", 1000, 200, 0.8},
		{"very high engagement", 5000, 1000, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculator.calculatePopularityScore(tt.points, tt.comments)
			assert.InDelta(t, tt.expected, score, 0.1)
		})
	}
}

func TestScoreCalculator_CalculateRelevanceScore(t *testing.T) {
	config := DefaultScoringConfig()
	calculator := NewScoreCalculator(config)

	tests := []struct {
		name     string
		title    string
		summary  string
		expected float64
	}{
		{
			name:     "highly relevant - machine learning",
			title:    "New Machine Learning Model Achieves SOTA",
			summary:  "Researchers develop breakthrough ML algorithm",
			expected: 0.9,
		},
		{
			name:     "moderately relevant - cloud computing",
			title:    "AWS Announces New Cloud Features",
			summary:  "Cloud infrastructure improvements",
			expected: 0.3,
		},
		{
			name:     "not relevant - cooking",
			title:    "Best Pasta Recipes",
			summary:  "Delicious Italian cuisine",
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculator.calculateRelevanceScore(tt.title, tt.summary)
			assert.InDelta(t, tt.expected, score, 0.2)
		})
	}
}

func TestScoreCalculator_CalculateCredibilityScore(t *testing.T) {
	config := DefaultScoringConfig()
	calculator := NewScoreCalculator(config)

	tests := []struct {
		name     string
		url      string
		expected float64
	}{
		{"high credibility - arxiv", "https://arxiv.org/abs/1234.5678", 0.95},
		{"high credibility - github", "https://github.com/user/repo", 0.9},
		{"medium credibility - medium", "https://medium.com/article", 0.6},
		{"low credibility - unknown", "https://unknown-site.com/article", 0.3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculator.calculateCredibilityScore(tt.url)
			assert.InDelta(t, tt.expected, score, 0.1)
		})
	}
}

func TestScoreCalculator_CalculateNoveltyScore(t *testing.T) {
	config := DefaultScoringConfig()
	calculator := NewScoreCalculator(config)

	now := time.Now()

	tests := []struct {
		name      string
		published time.Time
		expected  float64
	}{
		{"just published", now, 1.0},
		{"1 hour ago", now.Add(-1 * time.Hour), 0.9},
		{"1 day ago", now.Add(-24 * time.Hour), 0.7},
		{"3 days ago", now.Add(-72 * time.Hour), 0.4},
		{"1 week ago", now.Add(-168 * time.Hour), 0.1},
		{"old article", now.Add(-720 * time.Hour), 0.0}, // 30 days
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculator.calculateNoveltyScore(tt.published)
			assert.InDelta(t, tt.expected, score, 0.1)
		})
	}
}

func TestScoreCalculator_NormalizeScore(t *testing.T) {
	calculator := &ScoreCalculator{}

	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"within range", 0.5, 0.5},
		{"below zero", -0.5, 0.0},
		{"above one", 1.5, 1.0},
		{"at zero", 0.0, 0.0},
		{"at one", 1.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.normalizeScore(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func BenchmarkScoreCalculator_CalculateScore(b *testing.B) {
	config := DefaultScoringConfig()
	calculator := NewScoreCalculator(config)
	item := &dto.ContentItemDTO{
		Title:        "Benchmark Article",
		URL:          "https://example.com/benchmark",
		SourceType:   "hackernews",
		Popularity:   100,
		CommentCount: 10,
		PublishedAt:  time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculator.CalculateScore(item)
	}
}

func FuzzScoreCalculator_CalculateScore(f *testing.F) {
	config := DefaultScoringConfig()
	calculator := NewScoreCalculator(config)

	// Seed corpus
	f.Add("Test Title", "https://example.com", 100, 10)
	f.Add("Another Title", "https://test.org", 0, 0)

	f.Fuzz(func(t *testing.T, title string, url string, popularity int, comments int) {
		item := &dto.ContentItemDTO{
			Title:        title,
			URL:          url,
			SourceType:   "rss",
			Popularity:   popularity,
			CommentCount: comments,
			PublishedAt:  time.Now(),
		}

		score := calculator.CalculateScore(item)
		require.GreaterOrEqual(t, score, 0.0)
		require.LessOrEqual(t, score, 1.0)
	})
}
