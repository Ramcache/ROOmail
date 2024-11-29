package config

import (
	"log"
	"os"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	JWTSecret     string
}

func LoadConfig() Config {
	return Config{
		ServerAddress: getEnv("SERVER_ADDRESS", "90.156.156.78:8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:password@localhost/dbname?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Environment variable %s not set, using default value", key)
	return fallback
}
