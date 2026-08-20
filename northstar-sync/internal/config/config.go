package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	WarehouseAPIURL string
	PollInterval    time.Duration
	ServerPort      string
}

func Load() Config {
	interval := 5 * time.Minute
	if v := os.Getenv("POLL_INTERVAL_SECONDS"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil {
			interval = time.Duration(secs) * time.Second
		}
	}

	url := os.Getenv("WAREHOUSE_API_URL")
	if url == "" {
		url = "http://localhost:9090/stock"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		WarehouseAPIURL: url,
		PollInterval:    interval,
		ServerPort:      port,
	}
}
