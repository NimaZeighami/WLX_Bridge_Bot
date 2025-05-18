package services

import (
	"bridgebot/internal/client/http/bridgers"
	// "bridgebot/internal/database/models"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
)

func BuildQuoteRequest(userAddr string, fromTokenAddr, toTokenAddr string) bridgers.QuoteRequest {
	return bridgers.QuoteRequest{
		FromTokenAddress: fromTokenAddr,
		ToTokenAddress:   toTokenAddr,
		FromTokenAmount:  fmt.Sprintf("%d", BridgingAmount),
		FromTokenChain:   "POLYGON",
		ToTokenChain:     "BSC",
		UserAddr:         userAddr,
		EquipmentNo:      GenerateEquipmentNo(userAddr),
		SourceFlag:       "WLXBridgeApp",
		SourceType:       "",
	}
}

func RequestQuote(ctx context.Context, req bridgers.QuoteRequest) *bridgers.QuoteResponse {
	resp, err := bridgers.RequestQuote(ctx, req)
	if err != nil {
		log.Errorf("Error fetching quote: %v", err)
		panic(err)
	}
	return resp
}
