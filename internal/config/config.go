package config

import "os"

type Config struct {
	DB_DSN string
	Port   string
}

func FromEnv() Config {
	return Config{
		DB_DSN: envOr("DB_DSN", "postgres://user:pass@localhost:5432/shortlink?sslmode=disable"),
		Port:   envOr("PORT", "8080"),
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
