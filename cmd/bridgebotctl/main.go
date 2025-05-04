package main

import (
	"log"
	"bridgebot/cmd/bridgebotctl/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
