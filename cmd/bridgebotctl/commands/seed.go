package commands

import (
	"database/sql"
	log "bridgebot/internal/utils/logger"

	"github.com/spf13/cobra"
	_ "github.com/go-sql-driver/mysql"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/bridgebot_config?parseTime=true")
		if err != nil {
			log.Fatalf("Error opening DB: %v", err)
		}
		defer db.Close()

		_, err = db.Exec(`
			INSERT INTO bridge_configurations (token, network, token_contract_address, token_decimals, bridgers_smart_contract_address, is_enabled)
			VALUES ('USDT', 'POLYGON', '0xc2132D05D31c914a87C6611C10748AEb04B58e8F', 6, '0x1234567890abcdef', true)
		`)
		if err != nil {
			log.Fatalf("Error seeding DB: %v", err)
		}

		log.Info("✅ Database seeded.")
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
