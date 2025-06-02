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
	"gorm.io/gorm"
)

//todo: move config of db to .env
//todo: move this func to database directory
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