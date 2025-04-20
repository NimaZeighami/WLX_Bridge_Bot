//! In  (PowerShell):
//! $env:POLYGON_RPC_URL="https://polygon-rpc.com"
// ! In  (Command Prompt / CMD):
//! set POLYGON_RPC_URL=https://polygon-rpc.com
// ! POLYGON_PRIVATE_KEY = "dde619e9c94141eb5c60cf3c52e812f95db0a593543767a59e6b12e133a40c6d"
// ! set POLYGON_PRIVATE_KEY=dde619e9c94141eb5c60cf3c52e812f95db0a593543767a59e6b12e133a40c6d
package main

import (
	"bridgebot/internal/orchestration"
	log "bridgebot/internal/utils/logger"
	"context"

)

func main() {
	log.Info("Starting Bridge Bot...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orchestration.SetupSignalHandler(cancel)

	db := orchestration.InitDatabase()
	bridgeConfigs := orchestration.LoadBridgeConfigs(db)
	tokenMap := orchestration.BuildTokenMap(bridgeConfigs)

	userAddr := "0x7d0F13148e85A53227c65Ed013E7961A67839858"
	receiverAddr := userAddr

	usdtBSC := tokenMap["USDT"]["BSC"]
	usdtPolygon := tokenMap["USDT"]["POL"]

	quoteReq := orchestration.BuildQuoteRequest(userAddr, usdtPolygon, usdtBSC)
	quoteRes := orchestration.RequestQuote(ctx, quoteReq)

	log.Infof("Quote Response: %s USDT(BSC) for %s USDT(POLYGON)",
		quoteRes.Data.TxData.ToTokenAmount,
		quoteRes.Data.TxData.FromTokenAmount,
	)

	orchestration.CheckPolygonApproval(ctx, userAddr, usdtPolygon.TokenContractAddress)
	orchestration.SubmitPolygonApproval(ctx, userAddr, usdtPolygon.TokenContractAddress, usdtPolygon.BridgersSmartContractAddress)

	callReq := orchestration.BuildCalldataRequest(
		userAddr,
		receiverAddr,
		usdtPolygon,
		usdtBSC,
		quoteRes.Data.TxData.AmountOutMin,
	)
	orchestration.ExecuteBridgeTransaction(ctx, callReq)

	log.Info("BridgeBot execution completed.")
}


