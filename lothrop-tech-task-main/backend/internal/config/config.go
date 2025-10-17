package config

import "os"

type Config struct {
	Port         string
	PostgresDB   string
	PostgresPass string
	PostgresUser string
	PostgresHost string
	PostgresPort string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		PostgresDB:   getEnv("POSTGRES_DB", "lothrop_db"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "password"),
		PostgresUser: getEnv("POSTGRES_USER", "postgres"),
		PostgresHost: getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort: getEnv("POSTGRES_PORT", "5432"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
