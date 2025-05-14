package bridge_swap

import (
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/services"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
)

const (
	UsdtPolygonTokenAddress = "0xc2132d05d31c914a87c6611c10748aeb04b58e8f" // USDT(POLYGON)
	UsdtBscTokenAddress     = "0x55d398326f99059ff775485246999027b3197955" // USDT(BSC)
)

type SwapService struct{}

type BridgeProvider interface{}

func (s *SwapService) ProcessQuote(ctx context.Context, req QuoteReq) (*bridgers.QuoteResponse, error) {
	equipmentNo := services.GenerateEquipmentNo(req.FromWalletAddress)

	quoteReq := bridgers.QuoteRequest{
		FromTokenAddress: UsdtPolygonTokenAddress,
		ToTokenAddress:   UsdtBscTokenAddress,
		FromTokenAmount:  fmt.Sprintf("%.0f000000", req.FromTokenAmount),
		FromTokenChain:   "POLYGON",
		ToTokenChain:     "BSC",
		UserAddr:         req.FromWalletAddress,
		EquipmentNo:      equipmentNo,
		SourceFlag:       "bridgebot",
	}

	quoteResp, err := bridgers.RequestQuote(ctx, quoteReq)
	if err != nil {
		log.Errorf("failed to fetch quote: %w", err)
	}
	return quoteResp, nil
}

func (s *SwapService) ProcessSwap(ctx context.Context, req QuoteReq) (*bridgers.QuoteResponse, error) {
	equipmentNo := services.GenerateEquipmentNo(req.FromWalletAddress)

	calldataReq := bridgers.QuoteRequest{
		FromTokenAddress: UsdtPolygonTokenAddress,
		ToTokenAddress:   UsdtBscTokenAddress,
		FromTokenAmount:  fmt.Sprintf("%.0f000000", req.FromTokenAmount),
		FromTokenChain:   "POLYGON",
		ToTokenChain:     "BSC",
		UserAddr:         req.FromWalletAddress,
		EquipmentNo:      equipmentNo,
		SourceFlag:       "bridgebot",
	}

	quoteResp, err := bridgers.RequestQuote(ctx, calldataReq)
	if err != nil {
		log.Errorf("failed to fetch quote: %w", err)
	}
	return quoteResp, nil
}