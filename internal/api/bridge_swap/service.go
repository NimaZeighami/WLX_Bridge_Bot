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

func (s *SwapService) ProcessSwap(ctx context.Context, req GetQuoteRequest) (*bridgers.QuoteResponse, error) {
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
	log.Infof("Quote Response: %s USDT(BSC) for %s USDT(POLYGON)", "heeeeeeeeeeeeeeeeeeeeeeee",
		quoteResp.Data.TxData.ToTokenAmount,
		quoteResp.Data.TxData.FromTokenAmount,
	)

	// callDataReq := bridgers.CallDataRequest{
	// 	FromTokenAddress: req.FromToken,
	// 	ToTokenAddress:   req.ToToken,
	// 	FromAddress:      req.FromWalletAddress,
	// 	ToAddress:        req.ToWalletAddress,
	// 	FromTokenChain:   "POLYGON",
	// 	ToTokenChain:     "BSC",
	// 	FromTokenAmount:  fmt.Sprintf("%.0f", req.FromTokenAmount),
	// 	AmountOutMin:     quoteResp.Data.TxData.ToTokenAmount,
	// 	FromCoinCode:     "USDT(POL)", // ⚡ TODO: read dynamically if needed
	// 	ToCoinCode:       "USDT(BSC)",
	// 	EquipmentNo:      equipmentNo,
	// 	SourceType:       "H5",
	// 	SourceFlag:       "bridgebot",
	// }

	// callDataResp, err := bridgers.FetchBridgeCallData(ctx, callDataReq)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to fetch call data: %w", err)
	// }

	return quoteResp, nil
}
