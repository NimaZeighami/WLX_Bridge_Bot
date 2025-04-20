// This will include:
// BuildQuoteRequest
// RequestQuote

package orchestration

import (
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/database"
	log "bridgebot/internal/utils/logger"
	"context"
)

func BuildQuoteRequest(userAddr string, from, to database.TokenInfo) bridgers.QuoteRequest {
	return bridgers.QuoteRequest{
		FromTokenAddress: from.TokenContractAddress,
		ToTokenAddress:   to.TokenContractAddress,
		FromTokenAmount:  "3500000",
		FromTokenChain:   "POLYGON",
		ToTokenChain:     "BSC",
		UserAddr:         userAddr,
		EquipmentNo:      generateEquipmentNo(userAddr),
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

