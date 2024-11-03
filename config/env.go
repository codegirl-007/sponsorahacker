package config

import (
	"github.com/joho/godotenv"
	"os"
)

func GetEnvVar(key string) string {
	if os.Getenv("ENVIRONMENT") == "DEV" {
		err := godotenv.Load()
		if err != nil {
			// we're going to panic because without the env variables, our app couldn't run anyways
			panic("Error loading environment variable")
		}
	}

	return os.Getenv(key)
}
