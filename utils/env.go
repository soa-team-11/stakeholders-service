package utils

import (
	"os"
	logger "stakeholder-service/utils/logging"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	logger.Init()
	err := godotenv.Load()

	if err != nil {
		log.Warn("No .env file found, using default environment variables")
	}
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
