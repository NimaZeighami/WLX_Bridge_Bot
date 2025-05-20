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

// InitDatabase initializes the database connection and switches to the 'bridgebot_core' database, returning the GORM DB instance.
func InitDatabase() *gorm.DB {
	// todo: this fields should move to .env
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

// todo: check this 
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
