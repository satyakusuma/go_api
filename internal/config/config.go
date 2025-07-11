package config

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		// Ignore if .env file is not found
	}
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}