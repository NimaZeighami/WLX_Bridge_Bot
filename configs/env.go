package configs

import (
	log "bridgebot/internal/utils/logger"
	"github.com/joho/godotenv"
	"os"
)

// LoadEnv loads the .env file from the configs directory
func LoadEnv(envFilePath string) {
	err := godotenv.Load(envFilePath) //"C:/Users/Nima/Desktop/WLX_Bridge_Bot/.env"
	if err != nil {
		log.Errorf("Error loading .env file: %v", err)
	}
}

// GetEnv fetches an environment variable or returns a default value if missing
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
