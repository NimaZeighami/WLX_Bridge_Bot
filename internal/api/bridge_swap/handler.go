package bridge_swap

import (
	"bridgebot/internal/database/models"
	log "bridgebot/internal/utils/logger"
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var universalWalletRegex = regexp.MustCompile(`^[a-zA-Z0-9]{26,64}$`)

func isValidUniversalWalletAddress(addr string) bool {
	addr = strings.TrimSpace(addr)
	return universalWalletRegex.MatchString(addr)
}


func (s *SwapServer) HandleQuote(c echo.Context) error {
	var req QuoteReq
	var pairs []models.NetworkTokenPair
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON format",
		})
	}
	
	if err := s.DB.Find(&pairs).Error; err != nil {
		log.Errorf("Error fetching token pairs: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "failed to fetch token pairs",
		})
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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid token pair",
		})
	}

	if !isValidUniversalWalletAddress(req.FromWalletAddress) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid fromWalletAddress: must be alphanumeric and 26–64 characters",
		})
	}

	if !isValidUniversalWalletAddress(req.ToWalletAddress) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid toWalletAddress: must be alphanumeric and 26–64 characters",
		})
	}

	log.Infof("Received swap request: %+v", req)

	amountIn ,toAmount, quoteId, err := s.ProcessQuote(c.Request().Context(), req)
	if err != nil {
		log.Errorf("Swap failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	fromAmountDec, err := strconv.ParseFloat(amountIn, 64)
	if err != nil {
		log.Errorf("Failed to convert fromAmount to integer: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid fromAmount format",
		})
	}

	fromAmountFloat := fromAmountDec / math.Pow(10, float64(ChainDecimal(req.FromTokenChain)))
	fromAmount := strconv.FormatFloat(fromAmountFloat, 'f', -1, 64)

	return c.JSON(http.StatusOK, map[string]string{
		"toTokenAmount":   toAmount,
		"toToken":         req.ToToken,
		"toTokenChain":    req.ToTokenChain,
		"fromTokenAmount": fromAmount,
		"fromToken":       req.FromToken,
		"fromTokenChain":  req.FromTokenChain,
		"bridge":          "The Bridgers1",
		"quoteId":         strconv.FormatUint(uint64(quoteId), 10),
		"estimatedTime":   strconv.FormatUint(uint64(10), 10),
		// ? Note: Decimal values are intentionally omitted from the response and request to simplify the user experience.
	})
}

func (s *SwapServer) HandleSwap(c echo.Context) error {
	var quote models.Quote
	var req SwapReq
	if err := c.Bind(&req); err != nil || req.QuoteId == "0" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Valid quoteId is required in the JSON req",
		})
	}

	log.Infof("Processing swap for quote ID: %d", req.QuoteId)

	quoteIdUint64, err := strconv.ParseUint(req.QuoteId, 10, 64)
	if err != nil {
		log.Errorf("Swap failed for quote ID %d: %v", req.QuoteId, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if err := s.DB.First(&quote, uint(quoteIdUint64)).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Quote not found",
		})
	}
	if quote.State != "started" {
		log.Warnf("Quote ID %d is in state '%s', not allowed for processing", quoteIdUint64, quote.State)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "The quote has already been used or is no longer in a 'started' state. Please create a new quote for this swap, or check the current state using the status endpoint.",
		})
	}

	txHash, err := s.ProcessSwap(c.Request().Context(), uint(quoteIdUint64))
	if err != nil {
		log.Errorf("Swap failed for quote ID %d: %v", req.QuoteId, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// TODO: ADD an state updater for quote in quotes table
	// I have it in in ProcessSwap
	return c.JSON(http.StatusOK, map[string]string{
		"message":  "Swap submitted successfully",
		"tx_hash":  txHash,
		"quote_id": req.QuoteId,
	})
	// TODO: ADD implement other bridgers API for tracking Transaction Status
}

// HandleSwapStatus handles GET /v1/swaps/:id to return the current status of a swap quote.
func (s *SwapServer) HandleSwapStatus(c echo.Context) error {
	quoteOrderId := c.Param("orderid")
	if quoteOrderId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing order ID in URL",
		})
	}

	var quote models.Quote
	if err := s.DB.Where("order_id = ?", quoteOrderId).First(&quote).Error; err != nil {
		log.Errorf("Quote not found for orderId: %s", quoteOrderId)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Quote not found for provided orderId",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"quoteId":      quote.ID,
		"state":        quote.State,
		"fromChain":    quote.FromChain,
		"toChain":      quote.ToChain,
		"fromToken":    quote.FromToken,
		"toToken":      quote.ToToken,
		"fromAddress":  quote.FromAddress,
		"toAddress":    quote.ToAddress,
		"txHash":       quote.TxHash,
		"orderId":      quote.OrderId,
		"fromAmount":   quote.FromAmount,
		"toAmountMin":  quote.ToAmountMin,
		"updatedAt":    quote.UpdatedAt,
	})
}

