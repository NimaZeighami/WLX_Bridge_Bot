package bridge_swap

import (
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/contracts"
	"bridgebot/internal/database/models"
	"bridgebot/internal/services"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type SwapServer struct {
	DB *gorm.DB
}

// TODO:  Make ProcessQuote and ProcessSwap methods more generic to handle different bridge providers

type BridgeProvider interface {
	Quote()
	Bridge()
}

//  TODO: Implement BridgeProvider Interface for each of them
//* ADD sturct and these 2 method for  The Bridgers and other bridge providers for implementing BridgeProvider Interface

var ChainsDecimal = map[string]int{
	"BSC":     contracts.USDT_BSC_Decimal,
	"ETH":     contracts.USDT_BSC_Decimal,
	"POLYGON": contracts.USDT_POLYGON_Decimal,
	"TRX":     contracts.USDT_TRC20_Decimal,
}

var USDTContractAdresses = map[string]string{
	"BSC":     contracts.USDT_BSC_Addr,
	"ETH":     contracts.USDT_BSC_Addr,
	"POLYGON": contracts.USDT_POLYGON_Addr,
	"TRX":     contracts.USDT_TRC20_Addr,
}

// ChainDecimal returns the decimal precision for the given chain symbol.
func ChainDecimal(chainSymbol string) int {
	if decimal, ok := ChainsDecimal[chainSymbol]; ok {
		return decimal
	}
	return 0
}

// USDTContractAdress returns the contract address for the given chain symbol.
func USDTContractAdress(chainSymbol string) string {
	if contractAddress, ok := USDTContractAdresses[chainSymbol]; ok {
		return contractAddress
	}
	return ""
}

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

	// Handling TokenDecimal and return in string type
	fromAmount := float64(req.FromTokenAmount) * math.Pow(10, float64(ChainDecimal(req.FromTokenChain)))
	fromAmountStr := strconv.FormatFloat(fromAmount, 'f', -1, 64)

	log.Warnf("From Amount: %v", fromAmountStr)

	fromTokenAddr := USDTContractAdress(req.FromTokenChain)
	toTokenAddr := USDTContractAdress(req.ToTokenChain)

	quoteReq := bridgers.QuoteRequest{
		FromTokenAddress: fromTokenAddr,
		ToTokenAddress:   toTokenAddr,
		FromTokenAmount:  fromAmountStr,
		FromTokenChain:   req.FromTokenChain,
		ToTokenChain:     req.ToTokenChain,
		UserAddr:         req.FromWalletAddress,
		EquipmentNo:      equipmentNo,
		SourceFlag:       "WBB",
	}

	quoteResp, err := bridgers.RequestQuote(ctx, quoteReq)
	if err != nil {
		log.Errorf("failed to fetch quote: %w", err)
		return nil, 0, fmt.Errorf("failed to fetch quote")
	}
	// TODO: Tokens table can be omit because it is additional
	quote := models.Quote{
		FromTokenAddress: quoteReq.FromTokenAddress,
		ToTokenAddress:   quoteReq.ToTokenAddress,
		FromChain:        quoteReq.FromTokenChain,
		ToChain:          quoteReq.ToTokenChain,
		FromAddress:      quoteReq.UserAddr,
		ToAddress:        req.ToWalletAddress,
		FromAmount:       quoteReq.FromTokenAmount,
		ToAmountMin:      quoteResp.Data.TxData.AmountOutMin,
		TxHash:           "",
		State:            "pending", // initial state , other states can be submitted, confirmed, failed, expired and success.
	}

	if err := s.DB.Create(&quote).Error; err != nil {
		log.Errorf("failed to insert quote: %v", err)
		return nil, 0, fmt.Errorf("failed to store quote")
	}
	return quoteResp, quote.ID, nil
}

func (s *SwapServer) ProcessSwap(ctx context.Context, quoteID uint) (string, error) {
	var quote models.Quote
	if err := s.DB.First(&quote, quoteID).Error; err != nil {
		return "", fmt.Errorf("quote not found")
	}
	// ! Uncomment
	fromAmountInt, err := strconv.ParseInt(quote.FromAmount, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid from amount: %v", err)
	}
	fromAmount := big.NewInt(fromAmountInt)

	// revokeErr := services.SubmitPolygonApproval(ctx, quote.FromAddress, quote.FromTokenAddress, quote.ToTokenAddress, big.NewInt(0))
	// if revokeErr != nil {
	// 	s.DB.Model(&quote).Update("state", "failed")
	// 	return "", fmt.Errorf("approval failed: %v", revokeErr)
	// }

	// TODO: like Polygon we should check chain and based on that have approval (Switch-Case)
	if strings.ToUpper(quote.FromChain) == "POLYGON" {
		isApprovalNeeded := services.CheckPolygonApproval(ctx, quote.FromAddress, quote.FromTokenAddress, fromAmount)
		if isApprovalNeeded {
			err := services.SubmitPolygonApproval(ctx, quote.FromAddress, quote.FromTokenAddress, quote.ToTokenAddress, fromAmount)
			if err != nil {
				s.DB.Model(&quote).Update("state", "failed")
				return "", fmt.Errorf("approval failed: %v", err)
			}
		}
	}

	fromToken := quote.FromTokenAddress


	toToken := quote.ToTokenAddress


	callReq := services.BuildCalldataRequest(
		quote.FromAddress,
		quote.ToAddress,
		fromToken,
		toToken,
		quote.ToAmountMin, // TODO: toAmountMin has wrong value in database it is equal to fromAmount and it should fixed 
		fromAmount)

	txHash, err := services.ExecuteBridgeTransaction(ctx, callReq)
	if err != nil {
		s.DB.Model(&quote).Updates(map[string]interface{}{
			"state":   "failed",
			"tx_hash": "",
		})
		return "", fmt.Errorf("transaction failed: %v", err)
	}

	s.DB.Model(&quote).Updates(map[string]interface{}{
		"state":   "submitted",
		"tx_hash": txHash,
	})

	return txHash, nil
}
