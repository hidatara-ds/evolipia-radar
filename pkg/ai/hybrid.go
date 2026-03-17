package ai

import "strings"

// extractKeywords is a super lightweight local keyword extractor to aid in hybrid scoring.
// For production, this could be an NER agent, but we are keeping it simple per constraints.
func extractKeywords(text string) map[string]bool {
	// A small set of high-signal tech keywords
	signals := []string{
		"llm", "gpt", "rag", "gemini", "llama", "mistral", "transformer",
		"openai", "anthropic", "google", "meta", "microsoft", "apple",
		"mlops", "kubernetes", "GPU", "nvidia", "funding", "startup", "release",
	}

	lowerText := strings.ToLower(text)
	found := make(map[string]bool)

	for _, s := range signals {
		if strings.Contains(lowerText, strings.ToLower(s)) {
			found[strings.ToLower(s)] = true
		}
	}

	return found
}

// KeywordOverlap calculates a fast Jaccard-like similarity (0.0 to 1.0)
// based on the overlapping presence of high-signal tech keywords.
func KeywordOverlap(textA, textB string) float64 {
	aTerms := extractKeywords(textA)
	bTerms := extractKeywords(textB)

	if len(aTerms) == 0 && len(bTerms) == 0 {
		return 0.0 // No overlap and no keywords
	}

	intersection := 0
	for k := range aTerms {
		if bTerms[k] {
			intersection++
		}
	}

	union := len(aTerms) + len(bTerms) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}
