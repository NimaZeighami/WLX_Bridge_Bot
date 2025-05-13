// package main

// import (
// 	"bridgebot/internal/api"
// 	"bridgebot/internal/api/bridge_swap"
// 	log "bridgebot/internal/utils/logger"
// )

// func main() {
// 	e := api.NewServer()

// 	log.Info("Starting server on :8080...")

// 	if err := e.Start(":8080"); err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// 	log.Infof("%#v", bridge_swap.GetQuoteRequest{})
// 	log.Infof("%#v", bridge_swap.Response{})
// }

// ! Without the API server, with the code below the bot will not be able to process requests and execute transactions!

package main

import (
	"bridgebot/internal/services"
	log "bridgebot/internal/utils/logger"
	"context"
)

func main() {
	log.Info("Starting Bridge Bot...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	services.SetupSignalHandler(cancel)

	db := services.InitDatabase()
	bridgeConfigs := services.LoadBridgeConfigs(db)
	tokenMap := services.BuildTokenMap(bridgeConfigs)

	userAddr := "0x7d0F13148e85A53227c65Ed013E7961A67839858"
	receiverAddr := userAddr

	usdtBSC := tokenMap["USDT"]["BSC"]
	usdtPolygon := tokenMap["USDT"]["POL"]

	quoteReq := services.BuildQuoteRequest(userAddr, usdtPolygon, usdtBSC)
	quoteRes := services.RequestQuote(ctx, quoteReq)

	log.Infof("Quote Response: %s USDT(BSC) for %s USDT(POLYGON)",
		quoteRes.Data.TxData.ToTokenAmount,
		quoteRes.Data.TxData.FromTokenAmount,
	)

	isApprovalNeeded := services.CheckPolygonApproval(ctx, userAddr, usdtPolygon.TokenContractAddress)
	if isApprovalNeeded {
		services.SubmitPolygonApproval(ctx, userAddr, usdtPolygon.TokenContractAddress, usdtPolygon.BridgersSmartContractAddress)
	}


	callReq := services.BuildCalldataRequest(
		userAddr,
		receiverAddr,
		usdtPolygon,
		usdtBSC,
		//quoteRes.Data.TxData.AmountOutMin,  //changing this to lower value helps to get revert later !
	"19000000000000000000000000",
	)
	services.ExecuteBridgeTransaction(ctx, callReq)

	log.Info("BridgeBot execution completed.")



}
// ! If the total exchange fee of single order is less than 0.5 USDT, it will be charged as 0.5 USDT.



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