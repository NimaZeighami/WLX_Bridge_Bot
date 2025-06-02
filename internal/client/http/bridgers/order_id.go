package bridgers

import (
	"bridgebot/internal/client/http"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
)

type GenerateOrderIdRequest struct {
	Hash             string `json:"hash"`
	FromTokenAddress string `json:"fromTokenAddress"`
	ToTokenAddress   string `json:"toTokenAddress"`
	FromAddress      string `json:"fromAddress"`
	ToAddress        string `json:"toAddress"`
	FromTokenChain   string `json:"fromTokenChain"`
	ToTokenChain     string `json:"toTokenChain"`
	FromTokenAmount  string `json:"fromTokenAmount"`
	AmountOutMin     string `json:"amountOutMin"`
	FromCoinCode     string `json:"fromCoinCode"`
	ToCoinCode       string `json:"toCoinCode"`
	EquipmentNo      string `json:"equipmentNo"`
	SourceType       string `json:"sourceType,omitempty"`
	SourceFlag       string `json:"sourceFlag"`
}

type OrderIdResponse struct {
	ResCode int              `json:"resCode"`
	ResMsg  string           `json:"resMsg"`
	Data    OrderIdContainer `json:"data"`
}

// FetchOrderId initiates the token swap process using the HTTP client
func FetchOrderId(ctx context.Context, request GenerateOrderIdRequest) (*OrderIdResponse, error) {
	url := "https://api.bridgers.xyz/api/exchangeRecord/updateDataAndStatus"

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	response, err := http.Post[OrderIdResponse](ctx, url, headers, request)
	if err != nil {
		log.Errorf("order id request  failed: %v", err)
		return nil, fmt.Errorf("order id request  failed: %v", err)
	}

	if response.ResCode != 100 {
		log.Errorf("order id request  failed: %s", response.ResMsg)
		return nil, fmt.Errorf("order id request  failed: %s", response.ResMsg)
	}

	log.Infof("order id request  successful. Transaction Data: %+v", response.Data.OrderID)
	return response, nil
}
