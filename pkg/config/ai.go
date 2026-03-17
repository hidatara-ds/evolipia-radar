package config

import (
	"os"
	"strconv"
	"time"
)

// AIConfig holds settings for the AI Gateway.
type AIConfig struct {
	Provider     string        // e.g., "openrouter"
	APIKey       string        // from env: OR_API_KEY
	DefaultModel string        // e.g., "google/gemini-flash-1.5"
	Timeout      time.Duration // API request timeout
}

// LoadAIConfig reads AI configuration from environment variables with safe defaults.
func LoadAIConfig() AIConfig {
	timeoutStr := os.Getenv("AI_TIMEOUT_SECONDS")
	timeoutSecs, err := strconv.Atoi(timeoutStr)
	if err != nil || timeoutSecs <= 0 {
		timeoutSecs = 30 // default 30 seconds
	}

	model := os.Getenv("AI_DEFAULT_MODEL")
	if model == "" {
		model = "google/gemini-flash-1.5"
	}

	provider := os.Getenv("AI_PROVIDER")
	if provider == "" {
		provider = "openrouter"
	}

	return AIConfig{
		Provider:     provider,
		APIKey:       os.Getenv("AI_API_KEY"),
		DefaultModel: model,
		Timeout:      time.Duration(timeoutSecs) * time.Second,
	}
}
