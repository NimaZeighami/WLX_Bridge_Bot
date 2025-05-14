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


// ! consider that because of relative path, 
// ! the command will not work if you run it from another directory, so you need to run it from the root of the project