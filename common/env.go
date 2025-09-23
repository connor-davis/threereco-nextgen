package common

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func EnvString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
