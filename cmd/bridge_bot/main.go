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
	swapServer := &bridge_swap.SwapServer{DB: db}

	log.Info("Starting server on :8080...")
	e := api.NewServer(swapServer)
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Infof("%#v", bridge_swap.QuoteReq{})
	log.Infof("%#v", bridge_swap.QuoteRes{})
}

// TODO:
// *  1. Fix the bridge quote from Token Amount field (decimal issue)  
// !  2. Update the quote migration and update the quote struct
// ?  3. Add token , chain struct and chains variable in the constants package
// *  4. Add the rest of  bridgers APIs for tracking state of the transaction
// !  5. Fix duplcate swap for the same quote ID
// ?  6. Add proper state management for the swap
// *  7. Add getQuoteDetails API
// !  8. add a worker wating for the transaction to be verified
