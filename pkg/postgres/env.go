package postgres

import (
	"os"
	"strconv"
)

// ConfigFromEnv builds database config from standard POSTGRES_* environment variables.
func ConfigFromEnv() *Config {
	port, _ := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	return NewConfig().
		WithHost(getEnv("POSTGRES_HOST", "localhost")).
		WithPort(port).
		WithUser(getEnv("POSTGRES_USER", "admin")).
		WithPassword(getEnv("POSTGRES_PASSWORD", "admin1234")).
		WithDatabase(getEnv("POSTGRES_DB", "dragon_market")).
		WithSSLMode(getEnv("POSTGRES_SSLMODE", "disable"))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
