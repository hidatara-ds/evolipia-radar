package tagging

import (
	"testing"
)

func TestAutoTagger_AssignTags(t *testing.T) {
	tagger := NewAutoTagger()

	tests := []struct {
		name     string
		title    string
		content  string
		expected []string
	}{
		{
			name:     "LLM article",
			title:    "GPT-4 Turbo Released",
			content:  "OpenAI announces new language model with improved performance",
			expected: []string{"llm"},
		},
		{
			name:     "IDE article",
			title:    "Cursor IDE adds new AI features",
			content:  "The AI-powered code editor Cursor now supports GitHub Copilot",
			expected: []string{"ide"},
		},
		{
			name:     "Free credits article",
			title:    "Anthropic offers $10 free credits for students",
			content:  "Student program provides free API access to Claude",
			expected: []string{"free-credits"},
		},
		{
			name:     "Multiple tags",
			title:    "GitHub Copilot now free for students",
			content:  "GitHub Education program offers free access to AI code assistant",
			expected: []string{"ide", "free-credits"},
		},
		{
			name:     "Vision article",
			title:    "Stable Diffusion 3.0 Released",
			content:  "New text-to-image model with improved quality",
			expected: []string{"vision"},
		},
		{
			name:     "No specific tags",
			title:    "Random tech news",
			content:  "Some general technology update",
			expected: []string{"general_ai"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := tagger.AssignTags(tt.title, tt.content)

			// Check if all expected tags are present
			for _, expectedTag := range tt.expected {
				found := false
				for _, tag := range tags {
					if tag == expectedTag {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected tag '%s' not found in %v", expectedTag, tags)
				}
			}
		})
	}
}

func TestMergeTags(t *testing.T) {
	tests := []struct {
		name         string
		existingTags []string
		newTags      []string
		expectedLen  int
	}{
		{
			name:         "No duplicates",
			existingTags: []string{"llm", "research"},
			newTags:      []string{"tools", "ide"},
			expectedLen:  4,
		},
		{
			name:         "With duplicates",
			existingTags: []string{"llm", "research"},
			newTags:      []string{"llm", "tools"},
			expectedLen:  3,
		},
		{
			name:         "Empty existing",
			existingTags: []string{},
			newTags:      []string{"llm", "tools"},
			expectedLen:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeTags(tt.existingTags, tt.newTags)
			if len(result) != tt.expectedLen {
				t.Errorf("Expected %d tags, got %d: %v", tt.expectedLen, len(result), result)
			}
		})
	}
}
