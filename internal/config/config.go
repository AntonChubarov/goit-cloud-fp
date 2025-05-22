package config

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	DB_DSN   string
	Port     string
	LogLevel slog.Level
}

func FromEnv() Config {
	return Config{
		DB_DSN:   envOr("DB_DSN", "postgres://user:pass@localhost:5432/shortlink?sslmode=disable"),
		Port:     envOr("PORT", "8080"),
		LogLevel: parseLogLevel(envOr("LOG_LEVEL", "INFO")),
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
