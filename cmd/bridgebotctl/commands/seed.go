package commands

import (
	"bridgebot/configs"
	log "bridgebot/internal/utils/logger"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	_ "github.com/go-sql-driver/mysql"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Apply Goose seed migrations (SQL files)",
	Run: func(cmd *cobra.Command, args []string) {
		configs.LoadEnv(".env")

		user := configs.GetEnv("DB_USER", "root")
		pass := configs.GetEnv("DB_PASS", "")
		host := configs.GetEnv("DB_HOST", "localhost")
		port := configs.GetEnv("DB_PORT", "3306")
		name := configs.GetEnv("DB_NAME", "bridgebot_core")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}
		defer db.Close()

		if err := goose.SetDialect("mysql"); err != nil {
			log.Fatalf("Failed to set Goose dialect: %v", err)
		}

		if err := goose.Up(db, "internal/database/migrations"); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}

		log.Info("✅ Seed migrations applied successfully.")
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
