// Package orchestration provides functionality for handling and managing
// various operations related to token bridging and orchestration logic.
// This file contains the implementation of the BuildQuoteRequest function,
// which is responsible for constructing a QuoteRequest object to facilitate
// token bridging between different blockchain networks.

package orchestration

import (
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/database"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
)

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

func RequestQuote(ctx context.Context, req bridgers.QuoteRequest) *bridgers.QuoteResponse {
	resp, err := bridgers.RequestQuote(ctx, req)
	if err != nil {
		log.Errorf("Error fetching quote: %v", err)
		panic(err)
	}
	return resp
}
