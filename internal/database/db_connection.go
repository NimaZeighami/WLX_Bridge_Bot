package database

import (
	log "bridgebot/internal/utils/logger"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"bridgebot/internal/database/models"
)

// Connect initializes the connection to the MySQL server using GORM
func Connect(config models.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}
	return db, nil
}

// CheckDatabaseExists checks if the database exists using GORM
func CheckDatabaseExists(db *gorm.DB) (bool, error) {
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = ?", "bridgebot_core").Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("error checking database existence: %v", err)
	}
	return count > 0, nil
}

// CreateDatabase attempts to create a database named "bridgebot_core" if it does not already exist
func CreateDatabase(db *gorm.DB) error {
	return db.Exec("CREATE DATABASE IF NOT EXISTS bridgebot_core").Error
}

// UseDatabase switches to the 'bridgebot_core' database using GORM
func UseDatabase(db *gorm.DB) error {
	err := db.Exec("USE bridgebot_core").Error
	if err != nil {
		return fmt.Errorf("failed to switch to 'bridgebot_core' database: %v", err)
	}
	return nil
}



// InitializeDatabase connects to the database, creates the database, switches to it, and creates the table using GORM
func InitializeDatabase(config models.DBConfig) error {
	db, err := Connect(config)
	if err != nil {
		return err
	}

	dbExists, err := CheckDatabaseExists(db)
	if err != nil {
		return err
	}
	if !dbExists {
		err = CreateDatabase(db)
		if err != nil {
			return err
		}
	}
	log.Info("Database is created successfully!")
	log.Warn("But creating tables and initializing data is on command line !")
	return nil
}