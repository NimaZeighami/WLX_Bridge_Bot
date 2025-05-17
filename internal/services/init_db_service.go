// This will handle:
// Database init
// Config loading
// Token map construction
// OS signal handling

package services

import (
	"bridgebot/internal/database"
	"bridgebot/internal/database/models"
	log "bridgebot/internal/utils/logger"
	"context"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
)

func InitDatabase() *gorm.DB {
	config := models.DBConfig{
		Username: "root",
		Password: "@Nima8228",
		Host:     "localhost",
		Port:     "3306",
	}

	err := database.InitializeDatabase(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	db, err := database.Connect(config)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	if err := db.Exec("USE bridgebot_core").Error; err != nil {
		log.Fatalf("Error switching to database: %v", err)
	}
	return db
}

func LoadBridgeConfigs(db *gorm.DB) []models.BridgeConfiguration {
	var bridgeConfigs []models.BridgeConfiguration
	if err := db.Find(&bridgeConfigs).Error; err != nil {
		log.Fatalf("Error fetching bridge configs: %v", err)
	}
	return bridgeConfigs
}


// TODO: type map[string]map[string] reverse to order
func BuildTokenMap(bridgeConfigs []models.BridgeConfiguration) map[string]map[string]models.TokenInfo {
	tokenMap := make(map[string]map[string]models.TokenInfo)
	for _, config := range bridgeConfigs {
		if _, exists := tokenMap[config.Token]; !exists {
			tokenMap[config.Token] = make(map[string]models.TokenInfo)
		}
		tokenMap[config.Token][config.Network] = models.TokenInfo{
			ChainID:                      config.ChainID,
			TokenContractAddress:         config.TokenContractAddress,
			TokenDecimals:                config.TokenDecimals,
			BridgersSmartContractAddress: config.BridgersSmartContractAddress,
			IsEnabled:                    config.IsEnabled,
		}
	}
	return tokenMap
}

func SetupSignalHandler(cancelFunc context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Warn("Received termination signal. Shutting down gracefully...")
		cancelFunc()
		os.Exit(0)
	}()
}
