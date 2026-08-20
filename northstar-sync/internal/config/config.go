package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	WarehouseAPIURL string
	PollInterval    time.Duration
	ServerPort      string
	RequestTimeout  time.Duration
	MaxRetries      int
}

func Load() (Config, error) {
	cfg := Config{
		WarehouseAPIURL: getEnv("WAREHOUSE_API_URL", "http://localhost:9090/stock"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		PollInterval:    5 * time.Minute,
		RequestTimeout:  10 * time.Second,
		MaxRetries:      3,
	}

	if v := os.Getenv("POLL_INTERVAL_SECONDS"); v != "" {
		secs, err := strconv.Atoi(v)
		if err != nil || secs <= 0 {
			return cfg, fmt.Errorf("POLL_INTERVAL_SECONDS must be a positive integer, got %q", v)
		}
		cfg.PollInterval = time.Duration(secs) * time.Second
	}

	if v := os.Getenv("MAX_RETRIES"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return cfg, fmt.Errorf("MAX_RETRIES must be a non-negative integer, got %q", v)
		}
		cfg.MaxRetries = n
	}

	if cfg.WarehouseAPIURL == "" {
		return cfg, fmt.Errorf("WAREHOUSE_API_URL cannot be empty")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
