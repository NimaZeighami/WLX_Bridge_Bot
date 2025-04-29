package main

import (
	"bridgebot/internal/api"
	"bridgebot/internal/api/bridge_swap"
	log "bridgebot/internal/utils/logger"
)

func main() {
	e := api.NewServer()

	log.Info("Starting server on :8080...")

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Infof("%#v", bridge_swap.GetQuoteRequest{})
	log.Infof("%#v", bridge_swap.Response{})
}

// package main

// import (
// 	"bridgebot/internal/orchestration"
// log "bridgebot/internal/utils/logger"
// 	"context"
// )

// func main() {
// 	log.Info("Starting Bridge Bot...")

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	orchestration.SetupSignalHandler(cancel)

// 	db := orchestration.InitDatabase()
// 	bridgeConfigs := orchestration.LoadBridgeConfigs(db)
// 	tokenMap := orchestration.BuildTokenMap(bridgeConfigs)

// 	userAddr := "0x7d0F13148e85A53227c65Ed013E7961A67839858"
// 	receiverAddr := userAddr

// 	usdtBSC := tokenMap["USDT"]["BSC"]
// 	usdtPolygon := tokenMap["USDT"]["POL"]

// 	quoteReq := orchestration.BuildQuoteRequest(userAddr, usdtPolygon, usdtBSC)
// 	quoteRes := orchestration.RequestQuote(ctx, quoteReq)

// 	log.Infof("Quote Response: %s USDT(BSC) for %s USDT(POLYGON)",
// 		quoteRes.Data.TxData.ToTokenAmount,
// 		quoteRes.Data.TxData.FromTokenAmount,
// 	)

// 	isApprovalNeeded := orchestration.CheckPolygonApproval(ctx, userAddr, usdtPolygon.TokenContractAddress)
// 	if isApprovalNeeded {
// 		orchestration.SubmitPolygonApproval(ctx, userAddr, usdtPolygon.TokenContractAddress, usdtPolygon.BridgersSmartContractAddress)
// 	}
// 	callReq := orchestration.BuildCalldataRequest(
// 		userAddr,
// 		receiverAddr,
// 		usdtPolygon,
// 		usdtBSC,
// 		quoteRes.Data.TxData.AmountOutMin,
// 	)
// 	orchestration.ExecuteBridgeTransaction(ctx, callReq)

// 	log.Info("BridgeBot execution completed.")
// }
