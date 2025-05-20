


package bridgers

import (
	"bridgebot/internal/client/http"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
	"math/big"

	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
	UtmSource        string `json:"utmSource,omitempty"`
	OrderId          string `json:"orderId,omitempty"`
	SessionUuid      string `json:"sessionUuid,omitempty"`
	UserNo           string `json:"userNo,omitempty"`
}

type TronCallDataResponse struct {
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

type PolygonCallDataResponse struct {
	ResCode int    `json:"resCode"`
	ResMsg  string `json:"resMsg"`
	Data    struct {
		TxData struct {
			Data  string `json:"data"`
			To    string `json:"to"`
			Value string `json:"value"`
		} `json:"txData"`
	} `json:"data"`
}

// FetchBridgeCallData initiates the token swap process using the HTTP client
func FetchBridgeCallData(ctx context.Context, request CallDataRequest) (*PolygonCallDataResponse, error) {
	url := "https://api.bridgers.xyz/api/sswap/swap"

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	log.Infof("Sending CallData request: %+v", request)

	response, err := http.Post[PolygonCallDataResponse](ctx, url, headers, request)
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

// ExecuteBridgersSwapTransaction constructs calldata and broadcasts a swap transaction using the Bridgers protocol, returning the transaction hash or an error.
func ExecuteBridgersSwapTransaction(ctx context.Context, client *ethclient.Client, request CallDataRequest, privateKey *ecdsa.PrivateKey) (string, error) {
	log.Info("Fetching bridge transaction calldata...")

	callData, err := FetchBridgeCallData(ctx, request)
	if err != nil {
		return "", fmt.Errorf("error fetching bridge calldata: %v", err)
	}

	if callData.ResCode != 100 {
		return "", fmt.Errorf("bridgers API error: %s", callData.ResMsg)
	}

	log.Infof("Bridgers API response: %+v", callData)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		gasPrice = big.NewInt(100e9) 
		log.Warnf("Failed to get suggested gas price, using default: %v", err)
	}

	tipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		tipCap = big.NewInt(20e9)
		log.Warnf("Failed to get suggested gas tip cap, using default: %v", err)
	} else {
		log.Infof("Suggested gas tip cap: %s wei", tipCap.String())
	}

	// maxPriorityFeePerGas
	gasPrice = new(big.Int).Add(gasPrice, tipCap)

	toAddress := common.HexToAddress(callData.Data.TxData.To)

	value := big.NewInt(0)
	if callData.Data.TxData.Value != "" && callData.Data.TxData.Value != "0" && callData.Data.TxData.Value != "0x0" {
		value, ok = new(big.Int).SetString(callData.Data.TxData.Value, 0)
		if !ok {
			log.Warnf("Failed to parse transaction value '%s', using 0", callData.Data.TxData.Value)
			value = big.NewInt(0)
		}
	}

	data := common.FromHex(callData.Data.TxData.Data)

	gasLimit := uint64(100000) // Default gas limit
	estimatedGas, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		Gas:      0,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	})
	if err != nil {
		log.Warnf("Gas estimation failed, using default: %v", err)
	} else {
		gasLimit = estimatedGas + (estimatedGas / 5)
		log.Infof("Estimated gas: %d, using gas limit: %d", estimatedGas, gasLimit)
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get network ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	txHash := signedTx.Hash().Hex()
	log.Infof("Bridge transaction submitted! Tx Hash: %s", txHash)

	return txHash, nil
}
