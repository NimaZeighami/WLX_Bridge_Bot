package services

import (
	"bridgebot/configs"
	"bridgebot/internal/blockchain/polygon"
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/database"
	log "bridgebot/internal/utils/logger"
	"context"
	"crypto/ecdsa"
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

// GenerateEquipmentNo ensures a fixed 32-character string for equipmentNo field.
func GenerateEquipmentNo(userAddr string) string {
	if len(userAddr) >= 32 {
		return userAddr[:32]
	}
	return fmt.Sprintf("%032s", userAddr)
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


	spenderAddress := common.HexToAddress("0xb685760ebd368a891f27ae547391f4e2a289895b") // Bridge contract
	tokenAddress := common.HexToAddress(polygon.TokenAddress) // USDT contract
	amount := big.NewInt(BridgingAmount)

	needsApproval, err := polygon.IsApprovalNeeded(client, tokenAddress, fromAddress, spenderAddress, amount)
	if err != nil {
		log.Warnf("Failed to check approval status: %v", err)
	}

	if needsApproval {
		log.Info("Approval needed for token spending. Sending approval transaction...")
		
		approvalAmount := new(big.Int).Mul(amount, big.NewInt(100)) // 100x the amount for future transactions
		
		txHash, err := polygon.ApproveContract(client, tokenAddress, spenderAddress, approvalAmount, privateKey)
		if err != nil {
			return "", fmt.Errorf("failed to approve token spending: %w", err)
		}
		
		log.Infof("Approval transaction sent: %s", txHash)
		log.Info("Waiting for approval confirmation...")
	}

	txHash, err := bridgers.ExecuteBridgersSwapTransaction(ctx, client, request, privateKey)
	if err != nil {
		return "", err  // Handle error from ExecuteBridgersSwapTransaction
	}

	return txHash, nil
}

func ExecuteFullBridgeProcess(ctx context.Context, userAddr string, from, to database.TokenInfo) (string, error) {
	amountOutMin := calculateAmountOutMin(BridgingAmount, 0.002) // 0.2% slippage
	
	request := BuildCalldataRequest(userAddr, userAddr, from, to, amountOutMin)
	
	return ExecuteBridgeTransaction(ctx, request)
}

// calculateAmountOutMin calculates the minimum acceptable output amount with slippage
func calculateAmountOutMin(amount int64, slippagePercent float64) string {
	amountFloat := new(big.Float).SetInt64(amount)
	
	slippageFloat := new(big.Float).SetFloat64(1.0 - slippagePercent)
	
	minAmount := new(big.Float).Mul(amountFloat, slippageFloat)
	
	result := new(big.Int)
	minAmount.Int(result)
	
	return result.String()
}
