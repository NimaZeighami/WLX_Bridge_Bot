// ExecuteBridgeTransaction handles calldata construction and broadcasting the transaction.
// Package orchestration provides functionality for managing and executing
// complex workflows and operations related to blockchain transactions.
//
// The `bridge_executor.go` file contains the implementation for executing
// bridge transactions on the Polygon network. It handles fetching transaction
// calldata, encoding it, signing the transaction with a private key, and
// broadcasting it to the network. Additionally, it ensures proper gas price
// estimation and provides transaction status tracking via PolygonScan.

package services

import (
	"bridgebot/configs"
	"bridgebot/internal/blockchain/polygon"
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/database"
	log "bridgebot/internal/utils/logger"
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// BuildCalldataRequest constructs a CallDataRequest based on token info and bridge response data.
func BuildCalldataRequest(userAddr, receiverAddr string, from, to database.TokenInfo, amountOutMin string) bridgers.CallDataRequest {
	return bridgers.CallDataRequest{
		FromTokenAddress: from.TokenContractAddress,
		ToTokenAddress:   to.TokenContractAddress,
		FromAddress:      receiverAddr,
		ToAddress:        receiverAddr,
		FromTokenChain:   "POLYGON",
		ToTokenChain:     "BSC",
		FromTokenAmount:  fmt.Sprintf("%d", BridgingAmount),
		AmountOutMin:     amountOutMin,
		FromCoinCode:     "USDT(POL)",
		ToCoinCode:       "USDT(BSC)",
		EquipmentNo:      GenerateEquipmentNo(userAddr),
		SourceType:       "H5",
		SourceFlag:       "bridgers",
		Slippage:         "0.2",
	}
}

func ExecuteBridgeTransaction(ctx context.Context, request bridgers.CallDataRequest) {
	log.Info("Fetching bridge transaction calldata...")
	callData, err := bridgers.FetchBridgeCallData(ctx, request)
	if err != nil {
		log.Fatalf("Error fetching bridge calldata: %v", err)
	}

	calldataJSON, err := json.Marshal(callData.Data.TxData.Parameter)
	if err != nil {
		log.Fatalf("Failed to encode calldata: %v", err)
	}

	log.Infof("Transaction Destination: %s", callData.Data.TxData.To)
	configs.LoadEnv("../../.env")
	privateKeyHex := configs.GetPrivateKeyHex()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	client, err := polygon.NewPolygonClient()
	if err != nil {
		log.Fatalf("Error initializing Polygon client: %v", err)
	}

	// Estimate higher gas to avoid underpriced replacement
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Warnf("Failed to suggest gas price, defaulting: %v", err)
		gasPrice = big.NewInt(3e9) // default fallback 3 Gwei
	} else {
		gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(11))
		gasPrice = gasPrice.Div(gasPrice, big.NewInt(10)) // +10%
	}

	txHash, err := polygon.BroadcastTransactionWithCalldataWithGas(
		ctx,
		client,
		common.HexToAddress(callData.Data.TxData.To),
		calldataJSON,
		privateKey,
		gasPrice,
	)
	if err != nil {
		log.Fatalf("Error executing bridge transaction: %v", err)
	}

	log.Infof("Bridge transaction submitted! Tx Hash: %s", txHash)
	log.Infof("Check on PolygonScan: https://polygonscan.com/tx/%s", txHash)
}

// GenerateEquipmentNo ensures a fixed 32-character string for equipmentNo field.
func GenerateEquipmentNo(userAddr string) string {
	if len(userAddr) >= 32 {
		return userAddr[:32]
	}
	return fmt.Sprintf("%032s", userAddr)
}
