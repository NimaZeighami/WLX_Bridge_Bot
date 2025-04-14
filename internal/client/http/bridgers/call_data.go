package bridgers

import (
	"bridgebot/internal/client/http"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
)

type CallDataRequest struct {
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
	Slippage         string `json:"slippage,omitempty"`
	UtmSource		string `json:"utmSource,omitempty"`
	OrderId		 string `json:"orderId,omitempty"`
	SessionUuid	 string `json:"sessionUuid,omitempty"`
	UserNo		 string `json:"userNo,omitempty"`
	
}

type CallDataResponse struct {
	ResCode int    `json:"resCode"`
	ResMsg  string `json:"resMsg"`
	Data    struct {
		TxData struct {
			TronRouterAddress string `json:"tronRouterAddrees"`
			FunctionName      string `json:"functionName"`
			Options           struct {
				FeeLimit  int64  `json:"feeLimit"`
				CallValue string `json:"callValue"`
			} `json:"options"`
			Parameter []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"parameter"`
			FromAddress string `json:"fromAddress"`
			To          string `json:"to"`
		} `json:"txData"`
	} `json:"data"`
}

// FetchBridgeCallData initiates the token swap process using the HTTP client
func FetchBridgeCallData(ctx context.Context, request CallDataRequest) (*CallDataResponse, error) {
	url := "https://api.bridgers.xyz/api/sswap/swap"

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	log.Infof("Sending CallData request: %+v", request)

	response, err := http.Post[CallDataResponse](ctx, url, headers, request)
	if err != nil {
		log.Errorf("CallData request failed: %v", err)
		return nil, fmt.Errorf("CallData request failed: %v", err)
	}

	if response.ResCode != 100 {
		log.Errorf("CallData request failed: %s", response.ResMsg)
		return nil, fmt.Errorf("CallData request failed: %s", response.ResMsg)
	}

	log.Infof("CallData request successful. Transaction Data: %+v", response.Data.TxData)
	return response, nil
}
