package bridge_swap

import (
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/database/models"
	"bridgebot/internal/services"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

const (
	UsdtPolygonTokenAddress = "0xc2132d05d31c914a87c6611c10748aeb04b58e8f" // USDT(POLYGON)
	UsdtBscTokenAddress     = "0x55d398326f99059ff775485246999027b3197955" // USDT(BSC)
)

type SwapServer struct {
	DB *gorm.DB
}

// TODO:  Make ProcessQuote and ProcessSwap methods more generic to handle different bridge providers

type BridgeProvider interface{}

func (s *SwapServer) ProcessQuote(ctx context.Context, req QuoteReq) (*bridgers.QuoteResponse, uint, error) {
	equipmentNo := services.GenerateEquipmentNo(req.FromWalletAddress)

	log.Infof("DB: %#v", s.DB)
	var pairs []models.NetworkTokenPair
	if err := s.DB.Find(&pairs).Error; err != nil {
		log.Errorf("Error fetching token pairs: %v", err)
		return nil, 0, fmt.Errorf("failed to fetch token pairs")
	}

	isPairValid := false
	for _, p := range pairs {
		if strings.EqualFold(p.FromTokenSymbol, req.FromToken) &&
			strings.EqualFold(p.FromNetworkSymbol, req.FromTokenChain) &&
			strings.EqualFold(p.ToTokenSymbol, req.ToToken) &&
			strings.EqualFold(p.ToNetworkSymbol, req.ToTokenChain) {
			isPairValid = true
			break
		}
	}
	if !isPairValid {
		log.Errorf("Invalid token pair: %s-%s to %s-%s", req.FromToken, req.FromTokenChain, req.ToToken, req.ToTokenChain)
		return nil, 0, fmt.Errorf("invalid token pair")
	}

	//TODO: make fromTokenAmount to be dynamic based on decimals of the token

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
		return nil, 0, fmt.Errorf("failed to fetch quote")
	}
	quote := models.Quote{
		FromTokenAddress: quoteReq.FromTokenAddress,
		ToTokenAddress:   quoteReq.ToTokenAddress,
		FromChain:        quoteReq.FromTokenChain,
		ToChain:          quoteReq.ToTokenChain,
		FromAddress:      quoteReq.UserAddr,
		ToAddress:        req.ToWalletAddress,
		FromAmount:       quoteReq.FromTokenAmount,
		ToAmountMin:      quoteResp.Data.TxData.ToTokenAmount,
		TxHash:           "",
		State:            "pending", // initial state , other states can be submitted, confirmed, failed, expired and success.
	}

	if err := s.DB.Create(&quote).Error; err != nil {
		log.Errorf("failed to insert quote: %v", err)
		return nil, 0, fmt.Errorf("failed to store quote")
	}
	return quoteResp, quote.ID, nil
}

func (s *SwapServer) ProcessSwap(ctx context.Context, req QuoteReq) (*bridgers.QuoteResponse, error) {
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
