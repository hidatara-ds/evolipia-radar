package scoring

import (
	"math"
	"time"

	"github.com/hidatara-ds/evolipia-radar/internal/models"
)

type Weights struct {
	W1 float64 // popularity
	W2 float64 // relevance
	W3 float64 // credibility
	W4 float64 // novelty
}

var DefaultWeights = Weights{
	W1: 0.55,
	W2: 0.25,
	W3: 0.15,
	W4: 0.05,
}

var (
	defaultCredibilityConfig = DefaultCredibilityConfig()
	defaultRelevanceKeywords = DefaultRelevanceKeywords()
)

func ComputeScore(item *models.Item, signal *models.Signal, summary *models.Summary, weights Weights) *models.Score {
	hot := computeHotScore(signal, item.PublishedAt)
	relevance := computeRelevanceScore(item, summary)
	credibility := computeCredibilityScore(item.Domain)
	novelty := computeNoveltyScore(item.PublishedAt)

	final := weights.W1*hot + weights.W2*relevance + weights.W3*credibility + weights.W4*novelty

	return &models.Score{
		ItemID:      item.ID,
		Hot:         hot,
		Relevance:   relevance,
		Credibility: credibility,
		Novelty:     novelty,
		Final:       final,
		ComputedAt:  time.Now(),
	}
}

func computeHotScore(signal *models.Signal, publishedAt time.Time) float64 {
	if signal == nil {
		return 0.0
	}

	points := 0
	comments := 0
	if signal.Points != nil {
		points = *signal.Points
	}
	if signal.Comments != nil {
		comments = *signal.Comments
	}

	// Recency decay: older items get lower hot score
	ageHours := time.Since(publishedAt).Hours()
	decayFactor := math.Exp(-ageHours / 48.0) // Half-life of 48 hours

	// Simple scoring: points * 10 + comments * 5
	rawScore := float64(points*10 + comments*5)

	// Normalize to 0-1 range (assuming max 1000 points, 500 comments)
	maxScore := 1000.0*10 + 500.0*5
	normalized := rawScore / maxScore
	if normalized > 1.0 {
		normalized = 1.0
	}

	return normalized * decayFactor
}

func computeRelevanceScore(item *models.Item, summary *models.Summary) float64 {
	return computeRelevanceScoreWithConfig(item, summary, defaultRelevanceKeywords)
}

func computeRelevanceScoreWithConfig(item *models.Item, summary *models.Summary, keywords RelevanceKeywords) float64 {
	// Check title and excerpt for AI/ML keywords
	text := item.Title
	if item.RawExcerpt != nil {
		text += " " + *item.RawExcerpt
	}
	textLower := toLower(text)

	score := 0.0
	matches := 0

	// LLM keywords
	for _, kw := range keywords.LLM {
		if contains(textLower, kw) {
			score += keywords.Weights["llm"]
			matches++
		}
	}

	// MLOps keywords
	for _, kw := range keywords.MLOps {
		if contains(textLower, kw) {
			score += keywords.Weights["mlops"]
			matches++
		}
	}

	// CV keywords
	for _, kw := range keywords.CV {
		if contains(textLower, kw) {
			score += keywords.Weights["cv"]
			matches++
		}
	}

	// Check tags from summary
	if summary != nil {
		for _, tag := range summary.Tags {
			tagLower := toLower(tag)
			if contains(tagLower, "llm") || contains(tagLower, "ml") || contains(tagLower, "ai") {
				score += keywords.Weights["tag"]
				matches++
			}
		}
	}

	// Normalize: max score is around 1.0-2.0, cap at 1.0
	if score > 1.0 {
		score = 1.0
	}

	// If no matches, give baseline 0.1 (might still be relevant)
	if matches == 0 {
		score = 0.1
	}

	return score
}

func computeCredibilityScore(domain string) float64 {
	return computeCredibilityScoreWithConfig(domain, defaultCredibilityConfig)
}

func computeCredibilityScoreWithConfig(domain string, config CredibilityConfig) float64 {
	if config.Whitelist[domain] {
		return 1.0
	}
	if config.Blacklist[domain] {
		return 0.2
	}
	return 0.5 // Baseline
}

func computeNoveltyScore(publishedAt time.Time) float64 {
	// Newer items get higher novelty
	ageHours := time.Since(publishedAt).Hours()

	// Decay: items older than 7 days get very low novelty
	if ageHours > 168 { // 7 days
		return 0.1
	}

	// Linear decay from 1.0 to 0.1 over 7 days
	novelty := 1.0 - (ageHours / 168.0 * 0.9)
	if novelty < 0.1 {
		novelty = 0.1
	}

	return novelty
}

func toLower(s string) string {
	// Simple lowercase conversion
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			result[i] = s[i] + 32
		} else {
			result[i] = s[i]
		}
	}
	return string(result)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			indexOf(s, substr) != -1)))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
