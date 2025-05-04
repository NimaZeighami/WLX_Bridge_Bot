package commands

import (
	"bridgebot/configs"
	"database/sql"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	_ "github.com/go-sql-driver/mysql"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run Goose migrations from SQL files",
	Run: func(cmd *cobra.Command, args []string) {
		configs.LoadEnv("../../.env")

		user := configs.GetEnv("DB_USER", "root")
		pass := configs.GetEnv("DB_PASS", "@Nima8228")
		host := configs.GetEnv("DB_HOST", "localhost")
		port := configs.GetEnv("DB_PORT", "3306")
		name := configs.GetEnv("DB_NAME", "bridgebot_config")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}
		defer db.Close()

		if err := goose.Up(db, "internal/database/migrations"); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		log.Println("✅ Migrations complete.")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
