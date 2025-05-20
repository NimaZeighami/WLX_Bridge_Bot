package bridge_swap

import (
	"bridgebot/internal/database/models"
	log "bridgebot/internal/utils/logger"
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
	"strconv"
	"strings"
)


const BridgeProviderName = "The Bridgers cross-chain bridge"


func (s *SwapServer) HandleQuote(c echo.Context) error {
	var req QuoteReq
	var pairs []models.NetworkTokenPair

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON format",
		})
	}

	if err := s.DB.Find(&pairs).Error; err != nil {
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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid token pair",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	log.Infof("Received swap request: %+v", req)

	quoteResponse, quoteId, err := s.ProcessQuote(c.Request().Context(), req)
	if err != nil {
		log.Errorf("Swap failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	fromAmountDec, err := strconv.Atoi(quoteResponse.Data.TxData.FromTokenAmount)
	if err != nil {
		log.Errorf("Failed to convert fromAmount to integer: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid fromAmount format",
		})
	}

	//TODO : check this with your token const (map struct)
	fromAmountInt := fromAmountDec / int(math.Pow(10, float64(quoteResponse.Data.TxData.FromTokenDecimal)))
	fromAmount := strconv.Itoa(fromAmountInt)

	return c.JSON(http.StatusOK, map[string]string{
		"toTokenAmount":   quoteResponse.Data.TxData.ToTokenAmount,
		"toToken":         req.ToToken,
		"toTokenChain":    req.ToTokenChain,
		"fromTokenAmount": fromAmount,
		"fromToken":       req.FromToken,
		"fromTokenChain":  req.FromTokenChain,
		"bridge":          BridgeProviderName,
		"quoteId":         strconv.FormatUint(uint64(quoteId), 10),
		"estimatedTime":   strconv.FormatUint(uint64(quoteResponse.Data.TxData.EstimatedTime), 10),
		// ? Note: Decimal values are intentionally omitted from the response to simplify the user experience.
		// ? We get amount without decimal and we ourself send amount with decimal to the bridgers API. and
		// ? Return the response to the user without decimal values.
		// "fromTokenDecimal": quoteResponse.Data.TxData.FromTokenDecimal,
		// "toTokenDecimal":   quoteResponse.Data.TxData.ToTokenDecimal,
	})
}

func (s *SwapServer) HandleSwap(c echo.Context) error {
	var req SwapReq
	var quote models.Quote

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

	return c.JSON(http.StatusOK, map[string]string{
		"message":  "Swap submitted successfully",
		"tx_hash":  txHash,
		"quote_id": req.QuoteId,
	})
}
