package config

import "os"

// Config holds runtime settings loaded from env vars.
type Config struct {
	HTTPAddr string
	DSN      string // postgres DSN
}

func Load() Config {
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		DSN:      getEnv("DSN", "host=localhost user=postgres password=postgres dbname=consultations port=5432 sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
