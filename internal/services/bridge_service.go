package services

import (
	"bridgebot/configs"
	"bridgebot/internal/blockchain/polygon"
	"bridgebot/internal/client/http/bridgers"
	// "bridgebot/internal/database/models"
	log "bridgebot/internal/utils/logger"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strconv"
	"github.com/ethereum/go-ethereum/common"
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
		// FromTokenAmount:  fmt.Sprintf("%d", BridgingAmount),
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

	//Todo: get spender address from environment variable
	spenderAddress := common.HexToAddress("0xb685760ebd368a891f27ae547391f4e2a289895b") // Bridge contract
	tokenAddress := common.HexToAddress(polygon.TokenAddress)    
	FromTokenAmount , err:=strconv.ParseInt(request.FromTokenAmount,10 , 64)
	if  err != nil {
		log.Errorf("Error parsing token Amount : %v", err)
	}

	amount := big.NewInt(FromTokenAmount)

	needsApproval, err := polygon.IsApprovalNeeded(client, tokenAddress, fromAddress, spenderAddress, amount)
	if err != nil {
		log.Warnf("Failed to check approval status: %v", err)
	}

	if needsApproval {
		log.Info("Approval needed for token spending. Sending approval transaction...")

		// approvalAmount := new(big.Int).Mul(amount, big.NewInt(100)) // 100x the amount for future transactions

		txHash, err := polygon.ApproveContract(client, tokenAddress, spenderAddress, amount, privateKey)
		if err != nil {
			return "", fmt.Errorf("failed to approve token spending: %w", err)
		}

		log.Infof("Approval transaction sent: %s", txHash)
		log.Info("Waiting for approval confirmation...")
	}
	// ! Uncomment this if we need use without API verion ↑↑↑

	txHash, err := bridgers.ExecuteBridgersSwapTransaction(ctx, client, request, privateKey)
	if err != nil {
		return "", err
	}

	return txHash, nil
}
