package bridgers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	bridgers "bridgebot/internal/client/http/bridgers"
	"github.com/stretchr/testify/assert"
)

func TestFetchTokens_Success(t *testing.T) {
	mockResponse := `{
		"resCode": 100,
		"resMsg": "success",
		"data": {
			"tokens": [
				{
					"symbol": "USDT",
					"name": "Tether",
					"address": "0x1234",
					"decimals": 6,
					"logoURI": "",
					"chain": "POLYGON",
					"isCrossEnable": 1,
					"withdrawGas": 10,
					"chainId": "137"
				}
			]
		}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	tokens, err := bridgers.FetchTokens(context.Background(), bridgers.RequestBody{Chain: "POLYGON"})
	assert.NoError(t, err)
	assert.Len(t, tokens, 1)
	assert.Equal(t, "USDT", tokens[0]["symbol"])
	assert.Equal(t, "0x1234", tokens[0]["address"])
}

func TestFetchTokens_ErrorResponse(t *testing.T) {
	mockResponse := `{"resCode": 400, "resMsg": "something went wrong"}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	_, err := bridgers.FetchTokens(context.Background(), bridgers.RequestBody{Chain: "TRON"})
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "API error"))
}
