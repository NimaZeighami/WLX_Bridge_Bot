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
// ?  1. fix the worker with crunjob and update the state to verified 
// *  2. Add get Quote Details API
// !  3. Add bridge provider interface (after learning interface and reading oop from designpattern)
// ?  4. Check Approve functions one by one before writing interfaces



// 🧩 ۴. استفاده از این Interface در SwapServer

// در SwapServer یک فیلد اضافه کن:

// type SwapServer struct {
// 	DB             *gorm.DB
// 	BridgeProvider BridgeProvider
// }

// و در تابع NewSwapServer آن را مقداردهی کن:

// func NewSwapServer(db *gorm.DB, provider BridgeProvider) *SwapServer {
// 	return &SwapServer{
// 		DB:             db,
// 		BridgeProvider: provider,
// 	}
// }

// 🔄 ۵. استفاده در ProcessQuote و ProcessSwap

// quoteResp, err := s.BridgeProvider.Quote(...)
// ...
// isNeeded, err := s.BridgeProvider.ApprovalNeeded(...)
// ...
// err = s.BridgeProvider.Approve(...)
// ...
// callData, _ := s.BridgeProvider.CallData()
// signedTx, _ := s.BridgeProvider.Sign(callData)
// txHash, _ := s.BridgeProvider.BroadCast(signedTx)

// 🧪 ۶. راه‌اندازی

// در فایل main.go یا هر کجا که سرور را بالا می‌آوری:

// provider := &BridgersProvider{}
// swapServer := NewSwapServer(db, provider)