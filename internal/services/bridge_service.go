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
	// "github.com/ethereum/go-ethereum/ethclient"
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

	// Initialize Polygon client
	client, err := polygon.NewPolygonClient()
	if err != nil {
		return "", fmt.Errorf("error initializing Polygon client: %w", err)
	}

	// Load private key
	configs.LoadEnv("../../.env")
	privateKeyHex := configs.GetPrivateKeyHex()
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get wallet address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	
	log.Infof("Using wallet address: %s", fromAddress.Hex())

	// Check if approval is needed for the bridge contract
	// Example: check if we need to approve USDT spending
	spenderAddress := common.HexToAddress("0xb685760ebd368a891f27ae547391f4e2a289895b") // Bridge contract
	tokenAddress := common.HexToAddress(polygon.TokenAddress) // USDT contract
	amount := big.NewInt(BridgingAmount)

	needsApproval, err := polygon.IsApprovalNeeded(client, tokenAddress, fromAddress, spenderAddress, amount)
	if err != nil {
		log.Warnf("Failed to check approval status: %v", err)
		// Continue anyway, the transaction will fail if approval is needed
	}

	// If approval is needed, send approval transaction first
	if needsApproval {
		log.Info("Approval needed for token spending. Sending approval transaction...")
		
		// Set a large approval amount to avoid frequent approvals
		approvalAmount := new(big.Int).Mul(amount, big.NewInt(100)) // 100x the amount for future transactions
		
		txHash, err := polygon.ApproveContract(client, tokenAddress, spenderAddress, approvalAmount, privateKey)
		if err != nil {
			return "", fmt.Errorf("failed to approve token spending: %w", err)
		}
		
		log.Infof("Approval transaction sent: %s", txHash)
		log.Info("Waiting for approval confirmation...")
		
		// In a production system, you would wait for confirmation here
		// For simplicity, we're continuing immediately, but you might want to add a delay or check for confirmation
	}

	// Execute the bridge transaction
	txHash, err := bridgers.ExecuteBridgersSwapTransaction(ctx, client, request, privateKey)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// ExecuteFullBridgeProcess coordinates the entire bridging process
func ExecuteFullBridgeProcess(ctx context.Context, userAddr string, from, to database.TokenInfo) (string, error) {
	// Calculate minimum output amount with slippage (e.g., 0.2% slippage)
	amountOutMin := calculateAmountOutMin(BridgingAmount, 0.002) // 0.2% slippage
	
	// Build the request
	request := BuildCalldataRequest(userAddr, userAddr, from, to, amountOutMin)
	
	// Execute the transaction
	return ExecuteBridgeTransaction(ctx, request)
}

// calculateAmountOutMin calculates the minimum acceptable output amount with slippage
func calculateAmountOutMin(amount int64, slippagePercent float64) string {
	// Convert amount to big.Float for precise calculation
	amountFloat := new(big.Float).SetInt64(amount)
	
	// Calculate slippage amount
	slippageFloat := new(big.Float).SetFloat64(1.0 - slippagePercent)
	
	// Calculate minimum output amount
	minAmount := new(big.Float).Mul(amountFloat, slippageFloat)
	
	// Convert back to integer
	result := new(big.Int)
	minAmount.Int(result)
	
	return result.String()
}
// package services

// import (
// 	"bridgebot/configs"
// 	"bridgebot/internal/blockchain/polygon"
// 	"bridgebot/internal/client/http/bridgers"
// 	"bridgebot/internal/database"
// 	log "bridgebot/internal/utils/logger"
// 	"context"
// 	"crypto/ecdsa"
// 	"fmt"
// 	"math/big"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	// "github.com/ethereum/go-ethereum/ethclient"
// )



// // BuildCalldataRequest constructs a CallDataRequest based on token info and bridge response data.
// func BuildCalldataRequest(userAddr, receiverAddr string, from, to database.TokenInfo, amountOutMin string) bridgers.CallDataRequest {
// 	return bridgers.CallDataRequest{
// 		FromTokenAddress: from.TokenContractAddress,
// 		ToTokenAddress:   to.TokenContractAddress,
// 		FromAddress:      receiverAddr,
// 		ToAddress:        receiverAddr,
// 		FromTokenChain:   "POLYGON",
// 		ToTokenChain:     "BSC",
// 		FromTokenAmount:  fmt.Sprintf("%d", BridgingAmount),
// 		AmountOutMin:     amountOutMin,
// 		FromCoinCode:     "USDT(POL)",
// 		ToCoinCode:       "USDT(BSC)",
// 		EquipmentNo:      GenerateEquipmentNo(userAddr),
// 		SourceType:       "H5",
// 		SourceFlag:       "bridgers",
// 		Slippage:         "0.2",
// 	}
// }

// // GenerateEquipmentNo ensures a fixed 32-character string for equipmentNo field.
// func GenerateEquipmentNo(userAddr string) string {
// 	if len(userAddr) >= 32 {
// 		return userAddr[:32]
// 	}
// 	return fmt.Sprintf("%032s", userAddr)
// }

// // ExecuteBridgeTransaction handles the complete bridging process
// func ExecuteBridgeTransaction(ctx context.Context, request bridgers.CallDataRequest) (string, error) {
// 	log.Info("Initiating bridge transaction...")

// 	// Initialize Polygon client
// 	client, err := polygon.NewPolygonClient()
// 	if err != nil {
// 		return "", fmt.Errorf("error initializing Polygon client: %w", err)
// 	}

// 	// Load private key
// 	configs.LoadEnv("../../.env")
// 	privateKeyHex := configs.GetPrivateKeyHex()
// 	privateKey, err := crypto.HexToECDSA(privateKeyHex)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to parse private key: %w", err)
// 	}

// 	// Get wallet address from private key
// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		return "", fmt.Errorf("error casting public key to ECDSA")
// 	}
// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	
// 	log.Infof("Using wallet address: %s", fromAddress.Hex())

// 	// Check if approval is needed for the bridge contract
// 	// Example: check if we need to approve USDT spending
// 	spenderAddress := common.HexToAddress("0xb685760ebd368a891f27ae547391f4e2a289895b") // Bridge contract
// 	tokenAddress := common.HexToAddress(polygon.TokenAddress) // USDT contract
// 	amount := big.NewInt(BridgingAmount)

// 	needsApproval, err := polygon.IsApprovalNeeded(client, tokenAddress, fromAddress, spenderAddress, amount)
// 	if err != nil {
// 		log.Warnf("Failed to check approval status: %v", err)
// 		// Continue anyway, the transaction will fail if approval is needed
// 	}

// 	// If approval is needed, send approval transaction first
// 	if needsApproval {
// 		log.Info("Approval needed for token spending. Sending approval transaction...")
		
// 		// Set a large approval amount to avoid frequent approvals
// 		approvalAmount := new(big.Int).Mul(amount, big.NewInt(100)) // 100x the amount for future transactions
		
// 		txHash, err := polygon.ApproveContract(client, tokenAddress, spenderAddress, approvalAmount, privateKey)
// 		if err != nil {
// 			return "", fmt.Errorf("failed to approve token spending: %w", err)
// 		}
		
// 		log.Infof("Approval transaction sent: %s", txHash)
// 		log.Info("Waiting for approval confirmation...")
		
// 		// In a production system, you would wait for confirmation here
// 		// For simplicity, we're continuing immediately, but you might want to add a delay or check for confirmation
// 	}

// 	// Execute the bridge transaction
// 	txHash, err := bridgers.ExecuteBridgeTransaction(ctx, client, request, privateKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	return txHash, nil
// }

// // ExecuteFullBridgeProcess coordinates the entire bridging process
// func ExecuteFullBridgeProcess(ctx context.Context, userAddr string, from, to database.TokenInfo) (string, error) {
// 	// Calculate minimum output amount with slippage (e.g., 0.2% slippage)
// 	amountOutMin := calculateAmountOutMin(BridgingAmount, 0.002) // 0.2% slippage
	
// 	// Build the request
// 	request := BuildCalldataRequest(userAddr, userAddr, from, to, amountOutMin)
	
// 	// Execute the transaction
// 	return ExecuteBridgeTransaction(ctx, request)
// }

// // calculateAmountOutMin calculates the minimum acceptable output amount with slippage
// func calculateAmountOutMin(amount int64, slippagePercent float64) string {
// 	// Convert amount to big.Float for precise calculation
// 	amountFloat := new(big.Float).SetInt64(amount)
	
// 	// Calculate slippage amount
// 	slippageFloat := new(big.Float).SetFloat64(1.0 - slippagePercent)
	
// 	// Calculate minimum output amount
// 	minAmount := new(big.Float).Mul(amountFloat, slippageFloat)
	
// 	// Convert back to integer
// 	result := new(big.Int)
// 	minAmount.Int(result)
	
// 	return result.String()
// }

// // ExecuteBridgeTransaction handles calldata construction and broadcasting the transaction.
// // Package orchestration provides functionality for managing and executing
// // complex workflows and operations related to blockchain transactions.
// //
// // The `bridge_executor.go` file contains the implementation for executing
// // bridge transactions on the Polygon network. It handles fetching transaction
// // calldata, encoding it, signing the transaction with a private key, and
// // broadcasting it to the network. Additionally, it ensures proper gas price
// // estimation and provides transaction status tracking via PolygonScan.

// package services

// import (
// 	"bridgebot/configs"
// 	"bridgebot/internal/blockchain/polygon"
// 	"bridgebot/internal/client/http/bridgers"
// 	"bridgebot/internal/database"
// 	log "bridgebot/internal/utils/logger"
// 	"context"
// 	"crypto/ecdsa"
// 	"fmt"
// 	"math/big"

// 	// "github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"golang.org/x/crypto/sha3"
// )

// // BuildCalldataRequest constructs a CallDataRequest based on token info and bridge response data.
// func BuildCalldataRequest(userAddr, receiverAddr string, from, to database.TokenInfo, amountOutMin string) bridgers.CallDataRequest {
// 	return bridgers.CallDataRequest{
// 		FromTokenAddress: from.TokenContractAddress,
// 		ToTokenAddress:   to.TokenContractAddress,
// 		FromAddress:      receiverAddr,
// 		ToAddress:        receiverAddr,
// 		FromTokenChain:   "POLYGON",
// 		ToTokenChain:     "BSC",
// 		FromTokenAmount:  fmt.Sprintf("%d", BridgingAmount),
// 		AmountOutMin:     amountOutMin,
// 		FromCoinCode:     "USDT(POL)",
// 		ToCoinCode:       "USDT(BSC)",
// 		EquipmentNo:      GenerateEquipmentNo(userAddr),
// 		SourceType:       "H5",
// 		SourceFlag:       "bridgers",
// 		Slippage:         "0.2",
// 	}
// }


// //ExecuteBridgeTransaction handles the complete bridge swap process
// func ExecuteBridgeTransaction(ctx context.Context, request bridgers.CallDataRequest) (string, error) {
// 	log.Info("Fetching bridge transaction calldata...")

// 	client, err := polygon.NewPolygonClient()
// 	if err != nil {
// 		return "", fmt.Errorf("error initializing Polygon client: %v", err)
// 	}

// 	callData, err := bridgers.FetchBridgeCallData(ctx, request)
// 	if err != nil {
// 		return "", fmt.Errorf("error fetching bridge calldata: %v", err)
// 	}

// 	if callData.ResCode != 100 {
// 		return "", fmt.Errorf("bridgers API error: %s", callData.ResMsg)
// 	}

// 	log.Infof("Bridgers API response: %+v", callData)

// 	privateKey, err := crypto.HexToECDSA(configs.GetPrivateKeyHex())
// 	if err != nil {
// 		return "", fmt.Errorf("failed to parse private key: %v", err)
// 	}

// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		return "", fmt.Errorf("error casting public key to ECDSA")
// 	}

// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

// 	nonce, err := client.PendingNonceAt(ctx, fromAddress)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get nonce: %v", err)
// 	}

// 	// Get suggested gas price
// 	// gasPrice, err := client.SuggestGasPrice(ctx)
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("failed to get gas price: %v", err)
// 	// }
// 	gasPrice := big.NewInt(100e9)

// 	toAddress := common.HexToAddress(callData.Data.TxData.To)
// 	value := big.NewInt(0) // No ETH transfer

// 	functionSignature := []byte(callData.Data.TxData.FunctionName)
// 	hash := sha3.NewLegacyKeccak256()
// 	hash.Write(functionSignature)
// 	methodID := hash.Sum(nil)[:4]

// 	var data []byte
// 	data = append(data, methodID...)

// 	for _, param := range callData.Data.TxData.Parameter {
// 		switch param.Type {
// 		case "address":
// 			addr := common.HexToAddress(param.Value)
// 			paddedAddr := common.LeftPadBytes(addr.Bytes(), 32)
// 			data = append(data, paddedAddr...)
// 		case "string":
// 			strBytes := []byte(param.Value)
// 			length := big.NewInt(int64(len(strBytes)))
// 			paddedLength := common.LeftPadBytes(length.Bytes(), 32)
// 			data = append(data, paddedLength...)
// 			data = append(data, strBytes...)
// 			remainder := len(strBytes) % 32
// 			if remainder > 0 {
// 				padding := make([]byte, 32-remainder)
// 				data = append(data, padding...)
// 			}
// 		case "uint256":
// 			val := new(big.Int)
// 			val.SetString(param.Value, 0) // 0 means auto-detect base from prefix
// 			paddedVal := common.LeftPadBytes(val.Bytes(), 32)
// 			data = append(data, paddedVal...)
// 		default:
// 			return "", fmt.Errorf("unsupported parameter type: %s", param.Type)
// 		}
// 	}

// 	// gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
// 	// 	From:     fromAddress,
// 	// 	To:       &toAddress,
// 	// 	Data:     data,
// 	// 	Value:    value,
// 	// 	GasPrice: gasPrice,
// 	// })
// 	// if err != nil {
// 		// Use a default gas limit if estimation fails
// 		gasLimit := uint64(30000)
// 	// 	log.Warnf("Gas estimation failed, using default: %v", err)
// 	// }

// 	// Create the transaction
// 	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

// 	// Get chain ID
// 	chainID, err := client.NetworkID(ctx)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get network ID: %v", err)
// 	}

// 	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to sign transaction: %v", err)
// 	}

// 	err = client.SendTransaction(ctx, signedTx)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to send transaction: %v", err)
// 	}

// 	log.Warnf("Bridge transaction submitted! Tx Hash: %s", signedTx.Hash().Hex())
// 	return signedTx.Hash().Hex(), nil
// }

// // GenerateEquipmentNo ensures a fixed 32-character string for equipmentNo field.
// func GenerateEquipmentNo(userAddr string) string {
// 	if len(userAddr) >= 32 {
// 		return userAddr[:32]
// 	}
// 	return fmt.Sprintf("%032s", userAddr)
// }