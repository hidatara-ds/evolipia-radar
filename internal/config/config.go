// Package config provides application configuration loading and validation.
package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPort              = "8080"
	defaultDatabaseURL       = "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable"
	defaultCacheTTL          = 60
	defaultCrawlInterval     = "@every 6h"
	defaultMinRelevance      = 30
	defaultMaxRetries        = 3
	defaultFetchTimeout      = 8
	defaultMaxFetchBytes     = 2000000
	defaultTopicKeywords     = "llm,agents,vision,open source,infra,robotics,security,ai,machine learning"
	defaultFallbackLLMModels = "anthropic/claude-3.5-sonnet,meta-llama/llama-3.1-70b-instruct"
)

// Config holding all environment configuration options for Evolipia Radar.
type Config struct {
	DatabaseURL         string
	Port                string
	CacheTTLSeconds     int
	WorkerCron          string
	CrawlInterval       string
	MinRelevanceScore   int
	MaxCrawlRetries     int
	TopicKeywords       []string
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

// Load populates Config from environment variables with safe default values.
func Load() *Config {
	fallbackModelsStr := getEnv("LLM_FALLBACK_MODELS", defaultFallbackLLMModels)
	fallbackModels := splitString(fallbackModelsStr, ",")

	topicsStr := getEnv("TOPICS_KEYWORDS", defaultTopicKeywords)
	topics := splitString(topicsStr, ",")

	cfg := &Config{
		DatabaseURL:         getEnv("DATABASE_URL", defaultDatabaseURL),
		Port:                getEnv("PORT", defaultPort),
		CacheTTLSeconds:     getEnvInt("CACHE_TTL_SECONDS", defaultCacheTTL),
		WorkerCron:          getEnv("WORKER_CRON", "*/10 * * * *"),
		CrawlInterval:       getEnv("CRAWL_INTERVAL", defaultCrawlInterval),
		MinRelevanceScore:   getEnvInt("MIN_RELEVANCE_SCORE", defaultMinRelevance),
		MaxCrawlRetries:     getEnvInt("MAX_CRAWL_RETRIES", defaultMaxRetries),
		TopicKeywords:       topics,
		MaxFetchBytes:       int64(getEnvInt("MAX_FETCH_BYTES", defaultMaxFetchBytes)),
		FetchTimeoutSeconds: getEnvInt("FETCH_TIMEOUT_SECONDS", defaultFetchTimeout),

		// LLM Configuration
		LLMProvider:       getEnv("LLM_PROVIDER", "openrouter"),
		LLMModel:          getEnv("LLM_MODEL", "google/gemini-flash-1.5"),
		LLMFallbackModels: fallbackModels,
		LLMAPIKey:         getEnv("LLM_API_KEY", ""),
		LLMMaxTokens:      getEnvInt("LLM_MAX_TOKENS", 500),
		LLMTemperature:    getEnvFloat("LLM_TEMPERATURE", 0.7),
		LLMEnabled:        getEnvBool("LLM_ENABLED", false),
	}

	cfg.Validate()
	return cfg
}

// Validate checks critical config parameters and logs warnings for suspicious values.
func (c *Config) Validate() {
	if c.MinRelevanceScore < 0 || c.MinRelevanceScore > 100 {
		slog.Warn("MIN_RELEVANCE_SCORE out of bounds [0-100], defaulting to 30", "val", c.MinRelevanceScore)
		c.MinRelevanceScore = defaultMinRelevance
	}
	if c.MaxCrawlRetries <= 0 {
		slog.Warn("MAX_CRAWL_RETRIES must be positive, defaulting to 3", "val", c.MaxCrawlRetries)
		c.MaxCrawlRetries = defaultMaxRetries
	}
}

// CacheTTL returns duration for cache expiry.
func (c *Config) CacheTTL() time.Duration {
	return time.Duration(c.CacheTTLSeconds) * time.Second
}

// FetchTimeout returns HTTP client timeout duration.
func (c *Config) FetchTimeout() time.Duration {
	return time.Duration(c.FetchTimeoutSeconds) * time.Second
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
