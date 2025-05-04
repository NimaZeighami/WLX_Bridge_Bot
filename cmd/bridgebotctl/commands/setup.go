// This is going to create the database if it doesn't exist
package commands

import (
	"bridgebot/configs" // or wherever your GetEnv is
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Create database if it doesn't exist",
	Run: func(cmd *cobra.Command, args []string) {
		configs.LoadEnv(".env")

		user := configs.GetEnv("DB_USER", "root")
		pass := configs.GetEnv("DB_PASS", "@Nima8228")
		host := configs.GetEnv("DB_HOST", "localhost")
		port := configs.GetEnv("DB_PORT", "3306")
		dbName := configs.GetEnv("DB_NAME", "bridgebot_config")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pass, host, port)
		fmt.Println("Connecting to DSN:", dsn)


		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer db.Close()

		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
		if err != nil {
			log.Fatalf("Failed to create DB: %v", err)
		}

		fmt.Printf("✅ Database '%s' is ready.\n", dbName)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
