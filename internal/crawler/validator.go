// Package crawler provides web crawling, validation, scheduling, and retry capabilities.
package crawler

import (
	"fmt"
	"strings"
	"time"

	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/hidatara-ds/evolipia-radar/pkg/utils"
)

const (
	minTitleLength   = 10
	minContentLength = 50
	maxFutureDrift   = 5 * time.Minute
)

// Validator validates candidate items and computes relevance scores.
type Validator struct {
	minRelevanceScore int
	topicKeywords     []string
}

// NewValidator constructs a Validator instance.
func NewValidator(minRelevanceScore int, topicKeywords []string) *Validator {
	if minRelevanceScore < 0 || minRelevanceScore > 100 {
		minRelevanceScore = 30
	}
	return &Validator{
		minRelevanceScore: minRelevanceScore,
		topicKeywords:     topicKeywords,
	}
}

// ValidateItem checks item required fields and returns error if invalid.
func (v *Validator) ValidateItem(item *models.Item) error {
	if item == nil {
		return fmt.Errorf("item is nil")
	}

	trimmedTitle := strings.TrimSpace(item.Title)
	if len(trimmedTitle) < minTitleLength {
		return fmt.Errorf("title must be at least %d characters long", minTitleLength)
	}

	if !utils.IsValidURL(item.URL) {
		return fmt.Errorf("invalid URL format: must start with http or https")
	}

	if item.SourceID.String() == "" && item.SourceName == "" {
		return fmt.Errorf("source identification is required")
	}

	if item.PublishedAt.After(time.Now().Add(maxFutureDrift)) {
		return fmt.Errorf("published_at date cannot be in the future")
	}

	contentLen := 0
	if item.RawExcerpt != nil {
		contentLen = len(strings.TrimSpace(*item.RawExcerpt))
	}
	if contentLen < minContentLength {
		return fmt.Errorf("content excerpt must be at least %d characters long", minContentLength)
	}

	return nil
}

// ScoreRelevance calculates a 0-100 relevance score based on keyword occurrences.
func (v *Validator) ScoreRelevance(title, excerpt string) int {
	combined := strings.ToLower(title + " " + excerpt)
	if combined == "" {
		return 0
	}

	matches := 0
	for _, kw := range v.topicKeywords {
		kwClean := strings.TrimSpace(strings.ToLower(kw))
		if kwClean != "" && strings.Contains(combined, kwClean) {
			matches++
		}
	}

	if len(v.topicKeywords) == 0 {
		return 50 // Default score if no topics defined
	}

	// Calculate percentage match with a base score curve
	score := (matches * 100) / len(v.topicKeywords)
	if matches > 0 && score < 30 {
		score = 30 + matches*5 // Boost items that matched at least 1 keyword
	}

	if score > 100 {
		score = 100
	}
	return score
}

// IsReleasesSufficientlyRelevant returns true if score meets minimum threshold.
func (v *Validator) IsReleasesSufficientlyRelevant(score int) bool {
	return score >= v.minRelevanceScore
}
