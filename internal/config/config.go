package config

import (
	"os"
)

type Config struct {
	RedisURL    string
	PostgresURL string
	Domain      string
	Port        string
	Environment string
}

func Load() Config {
	return Config{
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		PostgresURL: getEnv("POSTGRES_URL", "postgres://username:password@localhost:5432/goflect"),
		Domain:      getEnv("DOMAIN", "localhost"),
		Port:        getEnv("PORT", "8000"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
