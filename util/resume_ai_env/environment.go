package resume_ai_env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const FetchaEnvKey = "RESUME_AI_ENV"

func IsProd() bool {
	return os.Getenv(FetchaEnvKey) == "production"
}

func IsDev() bool {
	return os.Getenv(FetchaEnvKey) == "development"
}

func IsTest() bool {
	return os.Getenv(FetchaEnvKey) == "test"
}

func LoadEnvironmentIfNeeded() {
	if IsDev() {
		err := godotenv.Load(".env.development")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
