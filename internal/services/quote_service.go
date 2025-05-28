package services

import (
	"bridgebot/internal/client/http/bridgers"
	// "bridgebot/internal/database/models"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
)

func BuildQuoteRequest(userAddr, fromTokenAddr, toTokenAddr, fromTokenChain, toTokenChain string, bridgingAmount uint) bridgers.QuoteRequest {
	return bridgers.QuoteRequest{
		FromTokenAddress: fromTokenAddr,
		ToTokenAddress:   toTokenAddr,
		FromTokenAmount:  fmt.Sprintf("%d", bridgingAmount),
		FromTokenChain:   fromTokenChain,
		ToTokenChain:     toTokenChain,
		UserAddr:         userAddr,
		EquipmentNo:      GenerateEquipmentNo(userAddr),
		SourceFlag:       "WLXBridgeApp",
		SourceType:       "",
	}
}

func RequestQuote(ctx context.Context, req bridgers.QuoteRequest) *bridgers.QuoteResponse {
	resp, err := bridgers.FetchQuote(ctx, req)
	if err != nil {
		log.Errorf("Error fetching quote: %v", err)
		panic(err)
	}
	return resp
}
