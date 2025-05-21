package main

import (
	"bridgebot/internal/api"
	"bridgebot/internal/api/bridge_swap"
	"bridgebot/internal/services"
	log "bridgebot/internal/utils/logger"
)

func main() {
	log.Info("Starting Bridge Bot...")

	db := services.InitDatabase()
	swapServer  := &bridge_swap.SwapServer{DB: db}

	log.Info("Starting server on :8080...")
	e := api.NewServer(swapServer)
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Infof("%#v", bridge_swap.QuoteReq{})
	log.Infof("%#v", bridge_swap.QuoteRes{})
}

