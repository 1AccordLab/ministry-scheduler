package env

import "os"

// Common environment utilities shared across features

// GetEnvOrDefault provides utilities for environment variables.
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
