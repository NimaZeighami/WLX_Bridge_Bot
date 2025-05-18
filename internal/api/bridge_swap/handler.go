package bridge_swap

import (
	log "bridgebot/internal/utils/logger"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

var universalWalletRegex = regexp.MustCompile(`^[a-zA-Z0-9]{26,64}$`)

func isValidUniversalWalletAddress(addr string) bool {
	addr = strings.TrimSpace(addr)
	return universalWalletRegex.MatchString(addr)
}

// helper function to check if a value exists in a slice
func isAllowed(value string, allowed []string) bool {
	for _, v := range allowed {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}

func (s *SwapServer) HandleQuote(c echo.Context) error {
	var req QuoteReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON format",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	allowedTokens := []string{"USDT"}
	if !isAllowed(req.FromToken, allowedTokens) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Unsupported fromToken: only USDT is allowed",
		})
	}
	if !isAllowed(req.ToToken, allowedTokens) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Unsupported toToken: only USDT is allowed",
		})
	}

	allowedChains := []string{"BSC", "POLYGON", "TRX"}
	if !isAllowed(req.FromTokenChain, allowedChains) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid fromTokenChain: must be one of BSC, Polygon, or TRX(not tron based on bridgers !)",
		})
	}
	if !isAllowed(req.ToTokenChain, allowedChains) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid toTokenChain: must be one of BSC, Polygon, or TRX(not tron based on bridgers !)",
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

	fromAmountInt := fromAmountDec / int(math.Pow(10, float64(quoteResponse.Data.TxData.FromTokenDecimal)))
	fromAmount := strconv.Itoa(fromAmountInt)

	return c.JSON(http.StatusOK, map[string]string{
		"toTokenAmount":   quoteResponse.Data.TxData.ToTokenAmount,
		"toToken":         req.ToToken,
		"toTokenChain":    req.ToTokenChain,
		"fromTokenAmount": fromAmount,
		"fromToken":       req.FromToken,
		"fromTokenChain":  req.FromTokenChain,
		"bridge":          "The Bridgers1",
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
	// TODO: ADD other messages like : "Swap failed" 
	// TODO: ADD implement other bridgers API for tracking Transaction Status
	// TODO: ADD other status updater function based on New API Response ...  expired, confirmed (mined), success (mind and funds recieved) 
}
