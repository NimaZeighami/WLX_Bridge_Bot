package bridge_swap

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	log "bridgebot/internal/utils/logger"

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

func HandleQuote(c echo.Context) error {
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

	service := SwapService{}

	quoteResponse, err := service.ProcessQuote(c.Request().Context(), req)
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

	fromAmountInt := fromAmountDec / quoteResponse.Data.TxData.FromTokenDecimal
	fromAmount := strconv.Itoa(fromAmountInt)

	return c.JSON(http.StatusOK, map[string]string{
		"toTokenAmount":   quoteResponse.Data.TxData.ToTokenAmount,
		"toToken":         req.ToToken,
		"toTokenChain":    req.ToTokenChain,
		"fromTokenAmount": fromAmount,
		"fromToken":       req.FromToken,
		"fromTokenChain":  req.FromTokenChain,
		"bridge":          "The Bridgers1",
		"quoteId":         "1",
		"estimatedTime":   string(quoteResponse.Data.TxData.EstimatedTime),
		// ? Note: Decimal values are intentionally omitted from the response to simplify the user experience.
		// ? We get amount without decimal and we ourself send amount with decimal to the bridgers API. and
		// ? Return the response to the user without decimal values.
		// "fromTokenDecimal": quoteResponse.Data.TxData.FromTokenDecimal,
		// "toTokenDecimal":   quoteResponse.Data.TxData.ToTokenDecimal,
	})
}

func HandleSwap(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Swap endpoint is not implemented yet",
	})
}
