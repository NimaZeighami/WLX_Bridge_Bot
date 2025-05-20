package bridge_swap

import (
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/constants"
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
	"ETH":     constants.Chains[0].SupportedTokens[0].Decimal,
	"BSC":     constants.Chains[1].SupportedTokens[0].Decimal,
	"POLYGON": constants.Chains[2].SupportedTokens[0].Decimal,
	"TRX":     constants.Chains[3].SupportedTokens[0].Decimal,
}

var USDTContractAdresses = map[string]string{
	"ETH":     constants.Chains[0].SupportedTokens[0].ContractAddr,
	"BSC":     constants.Chains[1].SupportedTokens[0].ContractAddr,
	"POLYGON": constants.Chains[2].SupportedTokens[0].ContractAddr,
	"TRX":     constants.Chains[3].SupportedTokens[0].ContractAddr,
}

var USDTSymbol = map[string]string{
	"ETH":     constants.Chains[0].SupportedTokens[0].Symbol,
	"BSC":     constants.Chains[1].SupportedTokens[0].Symbol,
	"POLYGON": constants.Chains[2].SupportedTokens[0].Symbol,
	"TRX":     constants.Chains[3].SupportedTokens[0].Symbol,
}

var USDTCoinCode = map[string]string{
	"ETH":     constants.Chains[0].SupportedTokens[0].CoinCode,
	"BSC":     constants.Chains[1].SupportedTokens[0].CoinCode,
	"POLYGON": constants.Chains[2].SupportedTokens[0].CoinCode,
	"TRX":     constants.Chains[3].SupportedTokens[0].CoinCode,
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

// TokenCoinCode returns the coin code for the given chain symbol.
func TokenCoinCode(chainSymbol string) string {
	if contractAddress, ok := USDTCoinCode[chainSymbol]; ok {
		return contractAddress
	}
	return ""
}

// TokenSymbol returns the symbol of token for the given chain symbol.
func TokenSymbol(chainSymbol string) string {
	if contractAddress, ok := USDTSymbol[chainSymbol]; ok {
		return contractAddress
	}
	return ""
}

func (s *SwapServer) ProcessQuote(ctx context.Context, req QuoteReq) (*bridgers.QuoteResponse, uint, error) {
	equipmentNo := services.GenerateEquipmentNo(req.FromWalletAddress)

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

	fromSymbol := TokenSymbol(req.FromTokenChain)
	toSymbol := TokenSymbol(req.ToTokenChain)
	fromCode := TokenSymbol(req.FromTokenChain)
	toCode := TokenSymbol(req.ToTokenChain)

	quote := models.Quote{
		FromTokenAddress: quoteReq.FromTokenAddress,
		ToTokenAddress:   quoteReq.ToTokenAddress,
		FromChain:        quoteReq.FromTokenChain,
		ToChain:          quoteReq.ToTokenChain,
		FromAddress:      quoteReq.UserAddr,
		ToAddress:        req.ToWalletAddress,
		FromTokenSymbol:  fromSymbol,
		ToTokenSymbol:    toSymbol,
		FromCoinCode:     fromCode,
		ToCoinCode:       toCode,
		FromAmount:       quoteReq.FromTokenAmount,
		ToAmountMin:      quoteResp.Data.TxData.AmountOutMin,
		TxHash:           "",
		State:            "started",
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
		return "", fmt.Errorf("quote not found: %v", err)
	}
	
	fromAmountInt, err := strconv.ParseInt(quote.FromAmount, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid from amount: %v", err)
	}
	fromAmount := big.NewInt(fromAmountInt)

	// ! If you need revoke approval for some reason uncomment this part and re-run the code
	// revokeErr := services.SubmitPolygonApproval(ctx, quote.FromAddress, quote.FromTokenAddress, quote.ToTokenAddress, big.NewInt(0))
	// if revokeErr != nil {
	// 	s.DB.Model(&quote).Update("state", "failed")
	// 	return "", fmt.Errorf("approval failed: %v", revokeErr)
	// }

	// todo: like Polygon we should check chain and based on that have approval (Switch-Case)
	if strings.ToUpper(quote.FromChain) == "POLYGON" {
		isApprovalNeeded := services.CheckPolygonApproval(ctx, quote.FromAddress, quote.FromTokenAddress, fromAmount)
		if isApprovalNeeded {
			err := services.SubmitPolygonApproval(ctx, quote.FromAddress, quote.FromTokenAddress, quote.ToTokenAddress, fromAmount)
			if err != nil {
				s.DB.Model(&quote).Update("state", "approval_failed")
				return "", fmt.Errorf("approval failed: %v", err)
			}
		}
	}
	s.DB.Model(&quote).Updates(map[string]interface{}{
		"state": "approved",
	})

	fromToken := quote.FromTokenAddress
	toToken := quote.ToTokenAddress

	callReq := services.BuildCalldataRequest(
		quote.FromAddress,
		quote.ToAddress,
		fromToken,
		toToken,
		quote.ToAmountMin,
		fromAmount)

	txHash, err := services.ExecuteBridgeTransaction(ctx, callReq)
	if err != nil {
		s.DB.Model(&quote).Updates(map[string]interface{}{
			"state":   "Broadcast_failed",
			"tx_hash": "",
		})
		return "", fmt.Errorf("transaction failed: %v", err)
	}

	s.DB.Model(&quote).Updates(map[string]interface{}{
		"state": "broadcast",
	})

	return txHash, nil
}
