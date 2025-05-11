package commands

import (
	"bridgebot/configs"
	log "bridgebot/internal/utils/logger"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	_ "github.com/go-sql-driver/mysql"
)

var migrateDownCmd = &cobra.Command{
	Use:   "migratedown [version]",
	Short: "Rollback Goose migrations to a specific version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configs.LoadEnv(".env")

		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Errorf("Invalid version format: %v", err)
			os.Exit(1)
		}

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

		if err := goose.DownTo(db, "internal/database/migrations", version); err != nil {
			log.Fatalf("Failed to migrate down to version %d: %v", version, err)
		}

		log.Infof("✅ Successfully rolled back to version %d", version)
	},
}

func init() {
	rootCmd.AddCommand(migrateDownCmd)
}


// ! consider that because of relative path, 
// ! the command will not work if you run it from another directory, so you need to run it from the root of the project