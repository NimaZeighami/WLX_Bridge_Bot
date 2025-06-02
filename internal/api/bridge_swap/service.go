package bridge_swap

import (
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/constants"
	"bridgebot/internal/database/models"
	"bridgebot/internal/services"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
	"gorm.io/gorm"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
	"bridgebot/configs"
	"bridgebot/internal/services/bridge"
	// "bridgebot/internal/services/bridge/thebridgers"
)

type SwapServer struct {
	DB *gorm.DB
}

// TODO:  Make ProcessQuote and ProcessSwap methods more generic to handle different bridge providers

//	TODO: Implement BridgeProvider Interface for each of them
//
// * ADD sturct and these 2 method for  The Bridgers and other bridge providers for implementing BridgeProvider Interface
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
	if coinCode, ok := USDTCoinCode[chainSymbol]; ok {
		return coinCode
	}
	return ""
}

// TokenSymbol returns the symbol of token for the given chain symbol.
func TokenSymbol(chainSymbol string) string {
	if TokenSymbol, ok := USDTSymbol[chainSymbol]; ok {
		return TokenSymbol
	}
	return ""
}

func (s *SwapServer) ProcessQuote(ctx context.Context, req QuoteReq) (amountIn, amoutOut string, quoteID uint, err error) {
	var bridger bridge.BridgeProvider

	fromAmount := float64(req.FromTokenAmount) * math.Pow(10, float64(ChainDecimal(req.FromTokenChain)))
	fromAmountStr := strconv.FormatFloat(fromAmount, 'f', -1, 64)

	bridger = bridge.SelectBestBridger()

	toAmount, err := bridge.FetchQuoteAmount(bridger, fromAmountStr, USDTContractAdress(req.FromTokenChain), req.FromTokenChain, USDTContractAdress(req.ToTokenChain), req.ToTokenChain, req.FromWalletAddress, ctx)
	if err != nil {
		log.Errorf("failed to fetch quote amount: %v", err)
		return fromAmountStr, "0", 0, fmt.Errorf("failed to fetch quote amount: %w", err)
	}

	quote := models.Quote{
		FromTokenAddress: USDTContractAdress(req.FromTokenChain),
		ToTokenAddress:   USDTContractAdress(req.ToTokenChain),
		FromChain:        req.FromTokenChain,
		ToChain:          req.ToTokenChain,
		FromToken:        TokenSymbol(req.FromTokenChain),
		ToToken:          TokenSymbol(req.ToTokenChain),
		FromCoinCode:     TokenCoinCode(req.FromTokenChain),
		ToCoinCode:       TokenCoinCode(req.ToTokenChain),
		FromAddress:      req.FromWalletAddress,
		ToAddress:        req.ToWalletAddress,
		FromAmount:       fromAmountStr,
		ToAmountMin:      toAmount,
		TxHash:           "",
		State:            "started",
	}

	if err := s.DB.Create(&quote).Error; err != nil {
		log.Errorf("failed to insert quote: %v", err)
		return fromAmountStr, toAmount, 0, fmt.Errorf("failed to store quote")
	}

	return fromAmountStr, toAmount, quote.ID, nil
}

func (s *SwapServer) ProcessSwap(ctx context.Context, quoteID uint) (string, error) {
	var quote models.Quote

	log.Warnf("quote.FromAmount raw value: '%s'", quote.FromAmount)

	if err := s.DB.First(&quote, uint(quoteID)).Error; err != nil {
		return "", fmt.Errorf("quote not found: %w", err)
	}

	fromAmountInt, err := strconv.ParseInt(quote.FromAmount, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid from amount: %v", err)
	}
	fromAmount := big.NewInt(fromAmountInt)

	// ? If we want to revoke approval before submitting a new one, we can uncomment the following lines
	// revokeErr := services.SubmitPolygonApproval(ctx, quote.FromAddress, quote.FromTokenAddress, quote.ToTokenAddress, big.NewInt(0))
	// if revokeErr != nil {
	// 	s.DB.Model(&quote).Update("state", "failed")
	// 	return "", fmt.Errorf("approval failed: %v", revokeErr)
	// }

	//todo: make approvement generic for all chains
	// switch quote.FromChain {
	// case "ETH", "BSC", "POLYGON" : // Supported EVM Based Chains
	// case "TRX": // TRON Chain
	// default:
	// 	return "", fmt.Errorf("unsupported chain: %s", quote.FromChain)
	// }

	bridger := bridge.SelectBestBridger()
	// todo: like Polygon we should check chain and based on that have approval (Switch-Case)
	if strings.ToUpper(quote.FromChain) == "POLYGON" {
		isApprovalNeeded, _ := bridge.CheckTokenApproval(bridger, quote.FromAddress, quote.FromTokenAddress, configs.GetBridgersContractAddr("POLYGON"), fromAmount, ctx)
		if isApprovalNeeded {
			err := bridge.RequestTokenApproval(bridger, quote.FromAddress, quote.FromTokenAddress, configs.GetBridgersContractAddr("POLYGON"), fromAmount, ctx)
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


	txHash, broadcastErr := bridge.FinalizeTransaction(bridger, quote.FromAddress,
		quote.ToAddress,
		fromToken,
		toToken,
		quote.ToAmountMin,
		fromAmount, ctx)
	if broadcastErr != nil {
		s.DB.Model(&quote).Updates(map[string]interface{}{
			"state":   "Broadcast_failed",
			"tx_hash": "",
		})
		return "", fmt.Errorf("transaction failed: %v", broadcastErr)
	}

	s.DB.Model(&quote).Updates(map[string]interface{}{
		"state":   "broadcast",
		"tx_hash": txHash,
	})

	orderResp, err := bridgers.FetchOrderId(ctx, bridgers.GenerateOrderIdRequest{
		Hash:             txHash,
		FromTokenAddress: quote.FromTokenAddress,
		ToTokenAddress:   quote.ToTokenAddress,
		FromAddress:      quote.FromAddress,
		ToAddress:        quote.ToAddress,
		FromTokenChain:   quote.FromChain,
		ToTokenChain:     quote.ToChain,
		FromTokenAmount:  quote.FromAmount,
		AmountOutMin:     quote.ToAmountMin,
		FromCoinCode:     quote.FromCoinCode,
		ToCoinCode:       quote.ToCoinCode,
		EquipmentNo:      services.GenerateEquipmentNo(quote.FromAddress),
		SourceFlag:       "WBB",
	})
	if err != nil {
		return "", fmt.Errorf("failed to get order ID: %v", err)
	}
	orderId := orderResp.Data.OrderID
	s.DB.Model(&quote).Updates(map[string]interface{}{
		"order_id": orderId,
	})

	log.Infof("Order ID: %s, Quote ID: %v", orderId, quote.ID)
	PollQuoteStatus(s.DB, orderId, quote.ID)

	return txHash, nil
}



func handlePollingStep(ctx context.Context, db *gorm.DB, orderId string, quoteId uint, attempt int) (shouldStop bool) {
	resp, err := bridgers.FetchTXDetails(ctx, orderId)
	if err != nil {
		log.Errorf("Polling failed for quote %d (attempt %d): %v", quoteId, attempt, err)
		return false 
	}

	status := resp.Data.Status
	log.Infof("Polling attempt %d for quote %d: status=%s", attempt, quoteId, status)

	switch {
	case status == "receive_complete":
		if err := db.Model(&models.Quote{}).Where("id = ?", quoteId).Update("state", "completed").Error; err != nil {
			log.Errorf("DB update (completed) failed: %v", err)
		}
		return true
	case strings.HasPrefix(status, "error"), status == "expired":
		if err := db.Model(&models.Quote{}).Where("id = ?", quoteId).Update("state", "failed").Error; err != nil {
			log.Errorf("DB update (failed) failed: %v", err)
		}
		return true 
	}

	return false 
}


func PollQuoteStatus(db *gorm.DB, orderId string, quoteId uint) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Infof("Recovered in polling goroutine for quote %d: %v", quoteId, r)
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
		defer cancel()

		log.Infof("Polling (immediate) for quote %d", quoteId)
		if stop := handlePollingStep(ctx, db, orderId, quoteId, 1); stop {
			return
		}

		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for attempt := 2; attempt <= 24; attempt++ {
			select {
			case <-ctx.Done():
				log.Infof("Polling context done for quote %d: %v", quoteId, ctx.Err())
				return
			case <-ticker.C:
				log.Infof("Polling (attempt %d) for quote %d", attempt, quoteId)
				if stop := handlePollingStep(ctx, db, orderId, quoteId, attempt); stop {
					return
				}
			}
		}
	}()
}

