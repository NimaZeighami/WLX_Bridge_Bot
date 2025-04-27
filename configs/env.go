package configs

import (
	log "bridgebot/internal/utils/logger"
	"github.com/joho/godotenv"
	"os"
)

func GetPrivateKeyHex() string {
	return os.Getenv("POLYGON_PRIVATE_KEY")
}

func GetRpcURL() string {
	return os.Getenv("POLYGON_RPC_URL")
}

// LoadEnv loads the .env file from the configs directory
func LoadEnv(envFilePath string) {
	err := godotenv.Load(envFilePath) 
	if err != nil {
		log.Errorf("Error loading .env file: %v", err)
	}
	if GetPrivateKeyHex() == "" {
		log.Fatal("POLYGON_PRIVATE_KEY not set in environment")
	}

	if GetRpcURL() == "" {
		log.Fatal("POLYGON_RPC_URL not set in environment")
	}
}

// GetEnv fetches an environment variable or returns a default value if missing
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
