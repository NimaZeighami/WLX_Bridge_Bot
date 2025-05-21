package bridgers

import (
	"bridgebot/internal/client/http"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
)

type QuoteRequest struct {
	FromTokenAddress string `json:"fromTokenAddress"`
	ToTokenAddress   string `json:"toTokenAddress"`
	FromTokenAmount  string `json:"fromTokenAmount"`
	FromTokenChain   string `json:"fromTokenChain"`
	ToTokenChain     string `json:"toTokenChain"`
	UserAddr         string `json:"userAddr"`
	EquipmentNo      string `json:"equipmentNo"`
	SourceFlag       string `json:"sourceFlag"`
	SourceType       string `json:"sourceType"`
}

type QuoteResponse struct {
	ResCode int    `json:"resCode"`
	ResMsg  string `json:"resMsg"`
	Data    struct {
		TxData struct {
			FromTokenAmount  string  `json:"fromTokenAmount"`
			FromTokenDecimal int     `json:"fromTokenDecimal"`
			ToTokenAmount    string  `json:"toTokenAmount"`
			ToTokenDecimal   int     `json:"toTokenDecimal"`
			Dex              string  `json:"dex"`
			Path             []any   `json:"path"`
			Fee              float64 `json:"fee"`
			FeeToken         string  `json:"feeToken"`
			ContractAddress  string  `json:"contractAddress"`
			LogoURL          string  `json:"logoUrl"`
			ChainFee         string  `json:"chainFee"`
			DepositMin       string  `json:"depositMin"`
			DepositMax       string  `json:"depositMax"`
			AmountOutMin     string  `json:"amountOutMin"`
			EstimatedTime    int     `json:"estimatedTime"`
			InstantRate      string  `json:"instantRate"`
		} `json:"txData"`
	} `json:"data"`
}

// FetchQuote sends a quote request to the Bridgers API.
func FetchQuote(ctx context.Context, requestBody QuoteRequest) (*QuoteResponse, error) {
	url := "https://api.bridgers.xyz/api/sswap/quote"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	log.Infof("Sending quote request: %+v", requestBody)

	response, err := http.Post[QuoteResponse](ctx, url, headers, requestBody)
	if err != nil {
		log.Errorf("Error fetching quote: %v", err)
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}

	if response.ResCode != 100 {
		log.Errorf("Quote request failed: %s", response.ResMsg)
		return nil, fmt.Errorf("quote request failed: %s", response.ResMsg)
	}

	log.Infof("Quote request successful. Data: %+v", response.Data.TxData)
	return response, nil
}
