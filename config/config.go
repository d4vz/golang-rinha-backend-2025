package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

const (
	DefaultProcessorName  = "default"
	FallbackProcessorName = "fallback"
)

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func GetEnvOrDefaultInt(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)

	if err != nil {
		return defaultValue
	}

	return parsedValue
}
