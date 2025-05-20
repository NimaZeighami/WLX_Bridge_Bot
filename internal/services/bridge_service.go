package services

import (
	"bridgebot/configs"
	"bridgebot/internal/blockchain/polygon"
	"bridgebot/internal/client/http/bridgers"
	log "bridgebot/internal/utils/logger"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/crypto"
)

// GenerateEquipmentNo ensures a fixed 32-character string for equipmentNo field.
func GenerateEquipmentNo(userAddr string) string {
	if len(userAddr) >= 32 {
		return userAddr[:32]
	}
	return fmt.Sprintf("%032s", userAddr)
}

// BuildCalldataRequest constructs a CallDataRequest based on token info and bridge response data.
func BuildCalldataRequest(userAddr, receiverAddr string, fromTokenAddr, toTokenAddr string, amountOutMin string, amount *big.Int) bridgers.CallDataRequest {
	return bridgers.CallDataRequest{
		FromTokenAddress: fromTokenAddr,
		ToTokenAddress:   toTokenAddr,
		FromAddress:      userAddr,
		ToAddress:        receiverAddr,
		FromTokenChain:   "POLYGON",
		ToTokenChain:     "BSC",
		FromTokenAmount: fmt.Sprintf("%d", amount),
		AmountOutMin:    amountOutMin,
		FromCoinCode:    "USDT(POL)",
		ToCoinCode:      "USDT(BSC)",
		EquipmentNo:     GenerateEquipmentNo(userAddr),
		SourceType:      "H5",
		SourceFlag:      "bridgers",
		Slippage:        "0.2",
	}
}

// ExecuteBridgeTransaction handles the complete bridging process
func ExecuteBridgeTransaction(ctx context.Context, request bridgers.CallDataRequest) (string, error) {
	log.Info("Initiating bridge transaction...")

	client, err := polygon.NewPolygonClient()
	if err != nil {
		return "", fmt.Errorf("error initializing Polygon client: %w", err)
	}

	configs.LoadEnv("../../.env")
	privateKeyHex := configs.GetPrivateKeyHex()
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	log.Infof("Using wallet address: %s", fromAddress.Hex())

	txHash, err := bridgers.ExecuteBridgersSwapTransaction(ctx, client, request, privateKey)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

