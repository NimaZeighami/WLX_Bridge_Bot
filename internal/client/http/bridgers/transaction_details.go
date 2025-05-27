package bridgers

import (
	"bridgebot/internal/client/http"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
)

type OrderIdContainer struct {
	OrderID string `json:"orderId"`
}
type TXDetailsResponse struct {
	ResCode int    `json:"resCode"`
	ResMsg  string `json:"resMsg"`
	Data    struct {
		ID                 int     `json:"id"`
		OrderID            string  `json:"orderId"`
		FromTokenAddress   string  `json:"fromTokenAddress"`
		ToTokenAddress     string  `json:"toTokenAddress"`
		FromTokenAmount    string  `json:"fromTokenAmount"`
		ToTokenAmount      string  `json:"toTokenAmount"`
		FromAmount         string  `json:"fromAmount"`
		ToAmount           string  `json:"toAmount"`
		FromDecimals       string  `json:"fromDecimals"`
		ToDecimals         string  `json:"toDecimals"`
		FromAddress        string  `json:"fromAddress"`
		Slippage           string  `json:"slippage"`
		FromChain          string  `json:"fromChain"`
		ToChain            string  `json:"toChain"`
		Hash               string  `json:"hash"`
		DepositHashExplore string  `json:"depositHashExplore"`
		DexName            string  `json:"dexName"`
		Status             string  `json:"status"`
		CreateTime         string  `json:"createTime"`
		FinishTime         string  `json:"finishTime"`
		Source             string  `json:"source"`
		ToAddress          string  `json:"toAddress"`
		ToHash             string  `json:"toHash"`
		ReceiveHashExplore string  `json:"receiveHashExplore"`
		EquipmentNo        string  `json:"equipmentNo"`
		RefundCoinAmt      string  `json:"refundCoinAmt"`
		RefundHash         string  `json:"refundHash"`
		RefundHashExplore  string  `json:"refundHashExplore"`
		RefundReason       string  `json:"refundReason"`
		FromCoinCode       string  `json:"fromCoinCode"`
		ToCoinCode         string  `json:"toCoinCode"`
		EstimatedTime      int     `json:"estimatedTime"`
		FromGas            *string `json:"fromGas"` // Use *string to handle nulls
		ToGas              *string `json:"toGas"`
		PlatformSource     string  `json:"platformSource"`
		Fee                *string `json:"fee"`
		Confirms           string  `json:"confirms"`
	} `json:"data"`
}

// FetchTXDetails sends a order id  to the Bridgers API.
func FetchTXDetails(ctx context.Context, requestBody OrderIdContainer) (*TXDetailsResponse, error) {
	url := "https://api.bridgers.xyz/api/exchangeRecord/getTransDataById"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	log.Infof("Sending fetch transaction details: %+v", requestBody)

	response, err := http.Post[TXDetailsResponse](ctx, url, headers, requestBody)
	if err != nil {
		log.Errorf("Error fetching quote: %v", err)
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}

	if response.ResCode != 100 {
		log.Errorf("fetch transaction details failed: %s", response.ResMsg)
		return nil, fmt.Errorf("fetch transaction details failed: %s", response.ResMsg)
	}

	log.Infof("fetch transaction details successful. Data: %+v", response.Data)
	return response, nil
}
