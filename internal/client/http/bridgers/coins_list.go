package bridgers

import (
	"bridgebot/internal/client/http"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
)

// TokenResponse represents the response structure for token retrieval
type TokenResponse struct {
	ResCode int    `json:"resCode"`
	ResMsg  string `json:"resMsg"`
	Data    struct {
		Tokens []Token `json:"tokens"`
	} `json:"data"`
}

// Token represents a cryptocurrency token's details
type Token struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	Decimals      int    `json:"decimals"`
	LogoURI       string `json:"logoURI"`
	Chain         string `json:"chain"`
	IsCrossEnable int    `json:"isCrossEnable"`
	WithdrawGas   int    `json:"withdrawGas"`
	ChainId       string `json:"chainId"`
}

// RequestBody represents the request body for fetching tokens
type RequestBody struct {
	Chain string `json:"chain"`
}

// FetchTokens sends a POST request to retrieve available tokens for a given chain
func FetchTokens(ctx context.Context, requestData RequestBody) ([]map[string]string, error) {
	url := "https://api.bridgers.xyz/api/exchangeRecord/getToken"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	log.Infof("Fetching tokens for chain: %s", requestData.Chain)

	response, err := http.Post[TokenResponse](ctx, url, headers, requestData)
	if err != nil {
		log.Errorf("Error fetching tokens: %v", err)
		return nil, fmt.Errorf("failed to fetch tokens: %w", err)
	}

	if response.ResCode != 100 {
		log.Errorf("Error response from API: %s", response.ResMsg)
		return nil, fmt.Errorf("API error: %s", response.ResMsg)
	}

	tokenMaps := make([]map[string]string, len(response.Data.Tokens))
	for i, token := range response.Data.Tokens {
		tokenMaps[i] = map[string]string{
			"symbol":  token.Symbol,
			"address": token.Address,
		}
	}

	log.Infof("Successfully fetched %d tokens for chain: %s", len(tokenMaps), requestData.Chain)
	return tokenMaps, nil
}

// FetchAndMapTokens retrieves tokens for a given chain and stores them in the provided map.
func FetchAndMapTokens(ctx context.Context, chain string, tokenMap map[string]string) error {
	// log.Infof("Fetching tokens for chain: %s", chain)

	requestData := RequestBody{Chain: chain}
	tokens, err := FetchTokens(ctx, requestData)
	if err != nil {
		return fmt.Errorf("error fetching tokens for %s: %w", chain, err)
	}

	// log.Infof("Total Tokens on %s: %d", chain, len(tokens))
	for _, token := range tokens {
		tokenMap[token["symbol"]] = token["address"]
	}

	return nil
}
