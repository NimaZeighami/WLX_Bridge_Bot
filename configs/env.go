package configs

import (
	log "bridgebot/internal/utils/logger"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func GetPrivateKeyHex() string {
	return os.Getenv("POLYGON_PRIVATE_KEY")
}

func GetRpcURL() string {
	return os.Getenv("POLYGON_RPC_URL")
}


// GetBridgersContractAddr returns the contract address for the given network symbol.
// The network parameter should be provided in uppercase (e.g., "ETH", "BSC").
func GetBridgersContractAddr(networkSymbol string) string {
	return os.Getenv(fmt.Sprintf("THE_BRIDGERS_%s_CONTRACT_ADDRESS", networkSymbol))
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
