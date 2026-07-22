package scoring

import (
	"math"
	"time"

	"github.com/hidatara-ds/evolipia-radar/pkg/models"
)

type Weights struct {
	W1 float64 // popularity
	W2 float64 // relevance/impact
	W3 float64 // credibility
	W4 float64 // engineering value
	W5 float64 // novelty
}

var DefaultWeights = Weights{
	W1: 0.0, // Hot/popularity (handled by gravity decay now, but keep field)
	W2: 0.3, // Impact
	W3: 0.1, // Credibility
	W4: 0.4, // Engineering Value
	W5: 0.2, // Novelty
}

var (
	defaultCredibilityConfig = DefaultCredibilityConfig()
	defaultRelevanceKeywords = DefaultRelevanceKeywords()
)

func ComputeScore(item *models.Item, signal *models.Signal, summary *models.Summary, existingScore *models.Score, weights Weights) *models.Score {
	hot := computeHotScore(signal, item.PublishedAt)
	credibility := computeCredibilityScore(item.Domain)
	
	relevance := computeRelevanceScore(item, summary) // Fallback for impact if no LLM
	
	// Default to heuristic novelty
	novelty := computeNoveltyScore(item.PublishedAt)
	impact := relevance
	engineeringValue := relevance
	reasoning := ""
	
	// Use LLM scores if available
	if existingScore != nil && existingScore.Impact > 0 {
		// LLM scores are 1-10. Normalize to 0-1 for internal math if needed, 
		// but wait, let's keep them as 0-1 or 1-10.
		// AnalyzeArticle returns 1-10. Let's normalize them here to 0-1.
		novelty = existingScore.Novelty / 10.0
		impact = existingScore.Impact / 10.0
		engineeringValue = existingScore.EngineeringValue / 10.0
		reasoning = existingScore.Reasoning
	}

	final := (weights.W2 * impact) + (weights.W3 * credibility) + (weights.W4 * engineeringValue) + (weights.W5 * novelty)

	return &models.Score{
		ItemID:           item.ID,
		Hot:              hot,
		Relevance:        relevance,
		Credibility:      credibility,
		Novelty:          novelty,
		Impact:           impact,
		EngineeringValue: engineeringValue,
		Reasoning:        reasoning,
		Final:            final,
		ComputedAt:       time.Now(),
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
	if val, ok := config.Tier1[domain]; ok {
		return val
	}
	if val, ok := config.Tier2[domain]; ok {
		return val
	}
	if val, ok := config.Tier3[domain]; ok {
		return val
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

// ConvertToScale10 converts a 0-1 score to 1-10 scale for better UX
// 0.0-0.1 -> 1
// 0.1-0.2 -> 2
// ...
// 0.9-1.0 -> 10
func ConvertToScale10(score float64) float64 {
	if score <= 0 {
		return 1.0
	}
	if score >= 1.0 {
		return 10.0
	}

	// Convert 0-1 to 1-10
	scaled := (score * 9.0) + 1.0

	// Round to 1 decimal place
	return math.Round(scaled*10) / 10
}

// ConvertScoreToScale10 converts all score components to 1-10 scale
func ConvertScoreToScale10(score *models.Score) map[string]float64 {
	return map[string]float64{
		"final":       ConvertToScale10(score.Final),
		"hot":         ConvertToScale10(score.Hot),
		"relevance":   ConvertToScale10(score.Relevance),
		"credibility": ConvertToScale10(score.Credibility),
		"novelty":     ConvertToScale10(score.Novelty),
	}
}
