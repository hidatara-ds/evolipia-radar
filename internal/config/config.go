package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL         string
	Port                string
	CacheTTLSeconds     int
	WorkerCron          string
	MaxFetchBytes       int64
	FetchTimeoutSeconds int
}

func Load() *Config {
	return &Config{
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable"),
		Port:                getEnv("PORT", "8080"),
		CacheTTLSeconds:     getEnvInt("CACHE_TTL_SECONDS", 60),
		WorkerCron:          getEnv("WORKER_CRON", "*/10 * * * *"),
		MaxFetchBytes:       int64(getEnvInt("MAX_FETCH_BYTES", 2000000)),
		FetchTimeoutSeconds: getEnvInt("FETCH_TIMEOUT_SECONDS", 8),
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
