package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPPort        string
	DatabaseURL     string
	CacheTTL        time.Duration
	ShortCodeLength int
}

func Load() Config {
	cfg := Config{
		HTTPPort:        getEnvString("PORT", "8080"),
		DatabaseURL:     getEnvString("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/link_storage?sslmode=disable"),
		CacheTTL:        getEnvDuration("CACHE_TTL_SECONDS", 300*time.Second),
		ShortCodeLength: getEnvInt("SHORT_CODE_LENGTH", 20),
	}
	return cfg
}

func getEnvString(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}

	return d
}
