package bridge_swap

import (
	"bridgebot/internal/utils/logger"
	// "fmt"
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

	// message := fmt.Sprintf("toTokenAmount: %s \n toToken: USDT(BSC) \n fromTokenAmount: %s \n  fromToken: USDT(POLYGON) \n Your Quote ID: 1 \n", toAmount, fromAmount  )
	// message := fmt.Sprintf(
	// 	"toTokenAmount: %s\n toToken: %s\n fromTokenAmount: %s\n fromToken: %s\n Your Quote ID: %s\n",
	// 	toAmount,
	// 	"USDT(BSC)",
	// 	fromAmount,
	// 	"USDT(POLYGON)",
	// 	"1",
	// )
	return c.JSON(http.StatusOK, map[string]string{
		"toTokenAmount": toAmount,
        "toToken":       "USDT(BSC)",
        "fromTokenAmount": fromAmount,
        "fromToken":       "USDT(POLYGON)",
        "quoteId":         "1",
	})
}

// package bridgers_swap

// import (
//     "net/http"
//     "bridgebot/internal/utils/logger"
//     "github.com/labstack/echo/v4"
// )

// func HandleSwap(c echo.Context) error {
//     var req GetQuoteRequest

//     if err := c.Bind(&req); err != nil {
//         return c.JSON(http.StatusBadRequest, map[string]string{
//             "error": "Invalid JSON format",
//         })
//     }

//     if err := c.Validate(req); err != nil {
//         return c.JSON(http.StatusBadRequest, map[string]string{
//             "error": err.Error(),
//         })
//     }

//     logger.Infof("Received swap request: %+v", req)

//     service := SwapService{}
//     /*bridgeDestination, err :=*/ service.ProcessSwap(c.Request().Context(), req)
//     // if err != nil {
//     //     logger.Errorf("Swap failed: %v", err)
//     //     return c.JSON(http.StatusInternalServerError, map[string]string{
//     //         "error": err.Error(),
//     //     })
//     // }

//     // return c.JSON(http.StatusOK, map[string]string{
//     //     "bridge_destination": bridgeDestination,
//     // })
//     return c.JSON(http.StatusOK, map[string]string{
//         "message": "Swap request processed successfully",
//     })
// }
