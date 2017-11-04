package utils

import (
	"fmt"
	"os"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func RequireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf(
			"Failed to retrieve required environment variable %s", key,
		))
	}
	return value
}
