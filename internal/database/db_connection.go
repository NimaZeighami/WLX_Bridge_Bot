package database

import (
	"fmt"
	log "bridgebot/internal/utils/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type BridgeConfiguration struct {
	ID                           uint   `gorm:"primaryKey"`
	Network                      string `gorm:"not null"`
	ChainID                      int
	Token                        string `gorm:"not null"`
	TokenContractAddress         string `gorm:"not null"`
	TokenDecimals                int    `gorm:"not null"`
	BridgersSmartContractAddress string `gorm:"not null"`
	IsEnabled                    bool   `gorm:"default:true"`
	CreatedAt                    string `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt                    string `gorm:"default:CURRENT_TIMESTAMP"`
}

// TokenInfo holds the relevant details for a token on a given network
type TokenInfo struct {
	ChainID					 int
	TokenContractAddress        string
	TokenDecimals               int
	BridgersSmartContractAddress string
	IsEnabled                   bool
}

var DB *gorm.DB

// Connect initializes the connection to the MySQL server using GORM
func Connect(config DBConfig) (*gorm.DB, error) {
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

// CheckTableExists checks if the bridge_configurations table exists using GORM
func CheckTableExists(db *gorm.DB) (bool, error) {
    // Make sure we're using the correct database first
    err := UseDatabase(db)
    if err != nil {
        return false, fmt.Errorf("error using database: %v", err)
    }

    // Now, check if the table exists
    var tableName string
    err = db.Raw("SHOW TABLES LIKE 'bridge_configurations'").Scan(&tableName).Error
    if err != nil {
        return false, fmt.Errorf("error checking table existence: %v", err)
    }

    // If tableName is empty, it means the table doesn't exist
    return tableName == "bridge_configurations", nil
}

// CreateTable creates the bridge_configurations table using GORM
func CreateTable(db *gorm.DB) error {
	type BridgeConfiguration struct {
		ID                           uint   `gorm:"primaryKey"`
		Network                      string `gorm:"not null"`
		ChainID                      int
		Token                        string `gorm:"not null"`
		TokenContractAddress         string `gorm:"not null"`
		TokenDecimals                int    `gorm:"not null"`
		BridgersSmartContractAddress string `gorm:"not null"`
		IsEnabled                    bool   `gorm:"default:true"`
		CreatedAt                    string `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt                    string `gorm:"default:CURRENT_TIMESTAMP"`
	}

	if err := db.AutoMigrate(&BridgeConfiguration{}); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	return nil
}

// HasInitialData checks if the bridge_configurations table contains any data using GORM
func HasInitialData(db *gorm.DB) (bool, error) {
	var count int64
	err := db.Model(&BridgeConfiguration{}).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if initial data exists: %v", err)
	}
	return count > 0, nil
}

// InsertInitialData fills the table with initial data using GORM
func InsertInitialData(db *gorm.DB) error {
	initialData := []BridgeConfiguration{
		{
			Network:                      "BSC",
			Token:                        "USDT",
			ChainID:                      56,
			TokenContractAddress:         "0x55d398326f99059ff775485246999027b3197955",
			TokenDecimals:                18,
			BridgersSmartContractAddress: "0xb685760ebd368a891f27ae547391f4e2a289895b",
			IsEnabled:                    true,
		},
		{
			Network:                      "TRON",
			Token:                        "USDT",
			ChainID:                      728126428,
			TokenContractAddress:         "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
			TokenDecimals:                6,
			BridgersSmartContractAddress: "0xb685760ebd368a891f27ae547391f4e2a289895b",
			IsEnabled:                    true,
		},
		{
			Network:                      "POL",
			Token:                        "USDT",
			ChainID:                      137,
			TokenContractAddress:         "0xc2132d05d31c914a87c6611c10748aeb04b58e8f",
			TokenDecimals:                6,
			BridgersSmartContractAddress: "0xb685760ebd368a891f27ae547391f4e2a289895b",
			IsEnabled:                    true,
		},
	}

	if err := db.Create(&initialData).Error; err != nil {
		return fmt.Errorf("failed to insert initial data: %v", err)
	}
	return nil
}

// InitializeDatabase connects to the database, creates the database, switches to it, and creates the table using GORM
func InitializeDatabase(config DBConfig) error {
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

	tableExists, err := CheckTableExists(db)
	if err != nil {
		return err
	}
	if !tableExists {
		err = CreateTable(db)
		if err != nil {
			return err
		}
	}

	dataExists, err := HasInitialData(db)
	if err != nil {
		return err
	}
	if !dataExists {
		err = InsertInitialData(db)
		if err != nil {
			return err
		}
	}

	log.Info("Database, table, and initial data are ready.")
	return nil
}
// ? Do we need to have this functionality with Pure Sql?
// ? Or we can use GORM for all the database operations?

// package database

// import (
// 	log "bridgebot/internal/utils/logger"
// 	"database/sql"
// 	"fmt"
// 	_ "github.com/go-sql-driver/mysql"
// )

// // DBConfig holds MySQL database connection details
// type DBConfig struct {
// 	Username string
// 	Password string
// 	Host     string
// 	Port     string
// 	Database string
// }

// // Connect initializes the connection to the MySQL server
// func Connect(config DBConfig) (*sql.DB, error) {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", config.Username, config.Password, config.Host, config.Port)
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
// 	}
// 	return db, nil
// }

// // CheckDatabaseExists checks if the database exists
// func CheckDatabaseExists(db *sql.DB) (bool, error) {
// 	var name string
// 	err := db.QueryRow("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = 'bridgebot_core'").Scan(&name)
// 	if err != nil && err.Error() != "sql: no rows in result set" {
// 		return false, fmt.Errorf("error checking database existence: %v", err)
// 	}
// 	return name == "bridgebot_core", nil
// }

// // CreateDatabase attempts to create a database named "bridgebot_core" if it does not already exist.
// func CreateDatabase(db *sql.DB) error {
// 	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS bridgebot_core")
// 	if err != nil {
// 		return fmt.Errorf("failed to create database: %v", err)
// 	}
// 	return nil
// }

// // UseDatabase switches to the 'bridgebot_core' database
// func UseDatabase(db *sql.DB) error {
// 	_, err := db.Exec("USE bridgebot_core")
// 	if err != nil {
// 		return fmt.Errorf("failed to switch to 'bridgebot_core' database: %v", err)
// 	}
// 	return nil
// }

// // CheckTableExists checks if the bridge_configurations table exists
// func CheckTableExists(db *sql.DB) (bool, error) {
// 	_, err := db.Exec("USE bridgebot_core")
// 	if err != nil {
// 		return false, fmt.Errorf("failed to select database: %v", err)
// 	}

// 	var name string
// 	err = db.QueryRow("SHOW TABLES LIKE 'bridge_configurations'").Scan(&name)
// 	if err != nil && err.Error() != "sql: no rows in result set" {
// 		return false, fmt.Errorf("error checking table existence: %v", err)
// 	}
// 	return name == "bridge_configurations", nil
// }


// // CreateTable creates the `bridgebot_core.bridge_configurations` table in the database
// func CreateTable(db *sql.DB) error {
// 	createTableSQL := `
// 	CREATE TABLE IF NOT EXISTS bridgebot_core.bridge_configurations (
// 		id INT AUTO_INCREMENT PRIMARY KEY,
// 		network VARCHAR(255) NOT NULL,
// 		chain_id INT ,
// 		token VARCHAR(255) NOT NULL,
// 		token_contract_address VARCHAR(255) NOT NULL,
// 		token_decimals INT NOT NULL,
// 		bridgers_smart_contract_address VARCHAR(255) NOT NULL,
// 		is_enabled BOOLEAN DEFAULT TRUE,
// 		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
// 	);`
// 	_, err := db.Exec(createTableSQL)
// 	if err != nil {
// 		return fmt.Errorf("failed to create table: %v", err)
// 	}
// 	return nil

// }

// // HasInitialData checks if the `bridgebot_core.bridge_configurations` table contains any data
// func HasInitialData(db *sql.DB) (bool, error) {
// 	var recordCount int
// 	err := db.QueryRow("SELECT COUNT(*) FROM bridgebot_core.bridge_configurations").Scan(&recordCount)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check if initial data exists: %v", err)
// 	}
// 	return recordCount > 0, nil
// }

// // InsertInitialData fills the table with initial data after it is created
// func InsertInitialData(db *sql.DB) error {
// 	// Example initial data
// 	insertSQL := `
// 	INSERT INTO bridgebot_core.bridge_configurations 
// 	(network, token, chain_id, token_contract_address, token_decimals, bridgers_smart_contract_address, is_enabled) 
// 	VALUES 
// 	('BSC', 'USDT', 56, '0x55d398326f99059ff775485246999027b3197955', 18, '0xb685760ebd368a891f27ae547391f4e2a289895b', TRUE),
// 	('TRON', 'USDT', 728126428, 'TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t', 6, '0xb685760ebd368a891f27ae547391f4e2a289895b', TRUE),
// 	('POL', 'USDT', 137, '0xc2132d05d31c914a87c6611c10748aeb04b58e8f', 6, '0xb685760ebd368a891f27ae547391f4e2a289895b', TRUE);
// 	`

// 	_, err := db.Exec(insertSQL)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert initial data: %v", err)
// 	}
// 	return nil
// }

// // InitializeDatabase connects to the database, creates the database, switches to it, and creates the table
// func InitializeDatabase(config DBConfig) error {
// 	db, err := Connect(config)
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()

// 	dbExists, err := CheckDatabaseExists(db)
// 	if err != nil {
// 		return err
// 	}
// 	if !dbExists {
// 		err = CreateDatabase(db)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	tableExists, err := CheckTableExists(db)
// 	if err != nil {
// 		return err
// 	}
// 	if !tableExists {
// 		err = CreateTable(db)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	dataExists, err := HasInitialData(db)
// 	if err != nil {
// 		return err
// 	}
// 	if !dataExists {
// 		err = InsertInitialData(db)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	log.Info("Database, table, and initial data are ready.")
// 	return nil
// }