package thebridgers

import (
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/database"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
)

const (
	BridgingAmount = 20_000_000 // Amount to approve, get qoute from bridgers and broadcast transaction in smallest unit (6 decimals for USDT)
)

// GenerateEquipmentNo ensures a fixed 32-character string for equipmentNo field.
func GenerateEquipmentNo(userAddr string) string {
	if len(userAddr) >= 32 {
		return userAddr[:32]
	}
	return fmt.Sprintf("%032s", userAddr)
}

// BuildQuoteRequest creates a QuoteRequest for bridging tokens between chains.
func BuildQuoteRequest(userAddr string, from, to database.TokenInfo) bridgers.QuoteRequest {
	return bridgers.QuoteRequest{
		FromTokenAddress: from.TokenContractAddress,
		ToTokenAddress:   to.TokenContractAddress,
		FromTokenAmount:  fmt.Sprintf("%d", BridgingAmount),
		FromTokenChain:   "POLYGON",
		ToTokenChain:     "BSC",
		UserAddr:         userAddr,
		EquipmentNo:      GenerateEquipmentNo(userAddr),
		SourceFlag:       "WLXBridgeApp",
		SourceType:       "",
	}
}

// RequestQuote fetches a quote for bridging tokens using the bridgers API.
func RequestQuote(ctx context.Context, req bridgers.QuoteRequest) *bridgers.QuoteResponse {
	resp, err := bridgers.RequestQuote(ctx, req)
	if err != nil {
		log.Errorf("Error fetching quote: %v", err)
		panic(err)
	}
	return resp
}


