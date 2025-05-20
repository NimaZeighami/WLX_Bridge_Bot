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

























// * NoteBook
// ? Tranccation Consist of
// ? 1. the amount of ether you're transferring,
// ? 2. the gas limit,
// ? 3. the gas price,
// ? 1. a nonce,
// ? 5. the receiving address,
// ? 6. and optionally data.

// ? Signing and Broadcasting the transaction
// 1. connecting to the Ethereum client
// 2. load the private key
// 3. get the nonce
// 4. load the public key of receiving address
// 5. value of ether to be transferred ( in wei),
// 6. gas limit (21000) and get gas price dynamically
// ! for erc20 token transfer we have to set the data field of the transaction
// ? considering that in  erc20 token transfer the value is 0
// 7. create the transaction
// 8. sign the transaction with the private key
// 9. send the transaction to the network
// 10. wait for the transaction to be mined


// USE
// "context"
// ctx, cancel := context.WithCancel(context.Background())