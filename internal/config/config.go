// Package config loads the service's runtime settings (listen address,
// database DSN) from environment variables, falling back to development
// defaults when a variable is not set.
package config

import "os"

// Config holds runtime settings loaded from env vars.
type Config struct {
	HTTPAddr string // адрес и порт, на котором слушает HTTP-сервер (HTTP_ADDR)
	DSN      string // postgres DSN (DSN)
}

// Load reads HTTP_ADDR and DSN from the environment and returns a Config,
// substituting development-friendly defaults (":8080" and a local
// PostgreSQL DSN) for whichever variable is unset or empty.
func Load() Config {
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		DSN:      getEnv("DSN", "host=localhost user=postgres password=postgres dbname=consultations port=5432 sslmode=disable"),
	}
}

// getEnv returns the value of the env var key, or fallback if it is unset
// or empty.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
