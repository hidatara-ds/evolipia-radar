package summarizer

import (
	"strings"

	"github.com/hidatara-ds/evolipia-radar/internal/models"
)

func GenerateExtractiveSummary(item *models.Item) *models.Summary {
	text := item.Title
	if item.RawExcerpt != nil {
		text += " " + *item.RawExcerpt
	}

	// Simple extractive: take first 3 sentences or first 200 chars
	sentences := extractSentences(text)
	tldr := ""
	if len(sentences) >= 3 {
		tldr = strings.Join(sentences[:3], " ")
	} else if len(sentences) > 0 {
		tldr = strings.Join(sentences, " ")
	} else {
		// Fallback: first 200 chars
		if len(text) > 200 {
			tldr = text[:200] + "..."
		} else {
			tldr = text
		}
	}

	// Why it matters: AI/ML engineer angle
	whyItMatters := generateWhyItMatters(item, text)

	// Extract tags
	tags := extractTags(item, text)

	return &models.Summary{
		ItemID:         item.ID,
		TLDR:           tldr,
		WhyItMatters:   whyItMatters,
		Tags:           tags,
		Method:         "extractive",
	}
}

func extractSentences(text string) []string {
	// Simple sentence splitting
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	var sentences []string
	current := ""
	
	for _, char := range text {
		current += string(char)
		if char == '.' || char == '!' || char == '?' {
			sentences = append(sentences, strings.TrimSpace(current))
			current = ""
		}
	}
	
	if strings.TrimSpace(current) != "" {
		sentences = append(sentences, strings.TrimSpace(current))
	}

	// Filter out very short sentences
	var filtered []string
	for _, s := range sentences {
		if len(s) > 20 {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

func generateWhyItMatters(item *models.Item, text string) string {
	textLower := strings.ToLower(text)
	
	// Check for specific topics and generate relevant angle
	if contains(textLower, "llm") || contains(textLower, "transformer") || contains(textLower, "gpt") {
		return "This development could impact how AI engineers build and deploy language models, potentially affecting inference costs, model architecture choices, and RAG system design."
	}
	if contains(textLower, "mlops") || contains(textLower, "deployment") {
		return "For ML engineers, this addresses critical production challenges around model deployment, monitoring, and maintaining model performance in real-world environments."
	}
	if contains(textLower, "computer vision") || contains(textLower, "detection") {
		return "Advances in computer vision directly impact applications in autonomous systems, medical imaging, and industrial automation, requiring engineers to stay updated on state-of-the-art techniques."
	}
	if contains(textLower, "rag") || contains(textLower, "retrieval") {
		return "This could improve how AI systems access and utilize external knowledge, which is crucial for building more capable and accurate AI applications."
	}
	
	// Default
	return "Staying informed about AI/ML developments helps engineers make better technical decisions, adopt new tools and techniques, and understand the evolving landscape of machine learning."
}

func extractTags(item *models.Item, text string) []string {
	return extractTagsWithConfig(item, text, DefaultTopicKeywordsConfig())
}

func extractTagsWithConfig(item *models.Item, text string, config TopicKeywordsConfig) []string {
	textLower := strings.ToLower(text)
	var tags []string

	for topic, keywords := range config.Keywords {
		for _, kw := range keywords {
			if strings.Contains(textLower, kw) {
				tags = append(tags, topic)
				break // Only add topic once
			}
		}
	}

	// If no tags found, add general_ai
	if len(tags) == 0 {
		tags = append(tags, "general_ai")
	}

	return tags
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
