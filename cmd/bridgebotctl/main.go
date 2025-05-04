package main

import (
	log "bridgebot/internal/utils/logger"
	"bridgebot/cmd/bridgebotctl/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
