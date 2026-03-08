package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DatabaseURL         string
	Port                string
	CacheTTLSeconds     int
	WorkerCron          string
	MaxFetchBytes       int64
	FetchTimeoutSeconds int
	
	// LLM Configuration
	LLMProvider       string
	LLMModel          string
	LLMFallbackModels []string
	LLMAPIKey         string
	LLMMaxTokens      int
	LLMTemperature    float64
	LLMEnabled        bool
}

func Load() *Config {
	// Parse fallback models from comma-separated string
	fallbackModelsStr := getEnv("LLM_FALLBACK_MODELS", "anthropic/claude-3.5-sonnet,meta-llama/llama-3.1-70b-instruct")
	var fallbackModels []string
	if fallbackModelsStr != "" {
		for _, model := range splitString(fallbackModelsStr, ",") {
			if model != "" {
				fallbackModels = append(fallbackModels, model)
			}
		}
	}
	
	return &Config{
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable"),
		Port:                getEnv("PORT", "8080"),
		CacheTTLSeconds:     getEnvInt("CACHE_TTL_SECONDS", 60),
		WorkerCron:          getEnv("WORKER_CRON", "*/10 * * * *"),
		MaxFetchBytes:       int64(getEnvInt("MAX_FETCH_BYTES", 2000000)),
		FetchTimeoutSeconds: getEnvInt("FETCH_TIMEOUT_SECONDS", 8),
		
		// LLM Configuration
		LLMProvider:       getEnv("LLM_PROVIDER", "openrouter"),
		LLMModel:          getEnv("LLM_MODEL", "google/gemini-flash-1.5"),
		LLMFallbackModels: fallbackModels,
		LLMAPIKey:         getEnv("LLM_API_KEY", ""),
		LLMMaxTokens:      getEnvInt("LLM_MAX_TOKENS", 500),
		LLMTemperature:    getEnvFloat("LLM_TEMPERATURE", 0.7),
		LLMEnabled:        getEnvBool("LLM_ENABLED", false),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (c *Config) CacheTTL() time.Duration {
	return time.Duration(c.CacheTTLSeconds) * time.Second
}

func (c *Config) FetchTimeout() time.Duration {
	return time.Duration(c.FetchTimeoutSeconds) * time.Second
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	parts := []string{}
	for _, part := range strings.Split(s, sep) {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}
