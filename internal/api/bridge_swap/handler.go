package bridge_swap

import (
	"bridgebot/internal/utils/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleSwap(c echo.Context) error {
	var req GetQuoteRequest

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

	logger.Infof("Received swap request: %+v", req)

	service := SwapService{}

	toAmount, fromAmount, err := service.ProcessSwap(c.Request().Context(), req)
	if err != nil {
		logger.Errorf("Swap failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"toTokenAmount": toAmount,
        "toToken":       "USDT(BSC)",
        "fromTokenAmount": fromAmount,
        "fromToken":       "USDT(POLYGON)",
        "quoteId":         "1",
	})
}