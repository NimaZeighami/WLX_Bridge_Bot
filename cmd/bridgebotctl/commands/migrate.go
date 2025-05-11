package commands

import (
	"bridgebot/configs"
	"database/sql"
	"fmt"
	log "bridgebot/internal/utils/logger"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	_ "github.com/go-sql-driver/mysql"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run Goose migrations from SQL files",
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

		// ✅ Set dialect to MySQL before running migrations
		if err := goose.SetDialect("mysql"); err != nil {
			log.Fatalf("Failed to set Goose dialect: %v", err)
		}

		if err := goose.Up(db, "internal/database/migrations"); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		log.Info("✅ Migrations complete.")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}


// ! consider that because of relative path, 
// ! the command will not work if you run it from another directory, so you need to run it from the root of the project