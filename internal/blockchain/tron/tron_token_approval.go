package tron

import (
	log "bridgebot/internal/utils/logger"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

const (
	TokenAddress = "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf" // USDT Contract Address on Nile Testnet
	// TokenAddress        = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" // USDT Contract Address on Mainnet
	ContractAddress = "TPwezUWpEGmFBENNWJHwXHRG1D2NCEEt5s"
	ApprovalAmount  = 5                             // USDT (Dont Mention Decimal it get multiplied by 1e6)
	TronNode        = "grpc.nile.trongrid.io:50051" // We Can have Array of Nodes and if one of them fails we can switch to another
	// TronNode       = "grpc.trongrid.io:50051"

)

// NewTronClient initializes and starts a new TRON gRPC client.
func NewTronClient() (*client.GrpcClient, error) {
	c := client.NewGrpcClient(TronNode)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()), //add encrypt
	}
	// opts = append(opts, grpc.WithBlock())

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// opts = append(opts, grpc.WithBlock())

	// import "google.golang.org/grpc/keepalive"

	// keepAliveParams := keepalive.ClientParameters{
	//     Time:                10 * time.Second,
	//     Timeout:             5 * time.Second,
	//     PermitWithoutStream: true,
	// }
	// opts = append(opts, grpc.WithKeepaliveParams(keepAliveParams))

	if err := c.Start(opts...); err != nil {
		log.Errorf("Error in Conneting to tron node: %v", err)
		return nil, fmt.Errorf("failed to start TRON client: %v", err)
	}

	return c, nil
}

// IsApprovalNeeded checks if the wallet has sufficient allowance for the contract.
// It checks if the wallet has approved the contract to spend at least 1 USDT (1,000,000 in TRC-20 decimal format).
func IsApprovalNeeded(ctx context.Context, client *client.GrpcClient, walletAddr string) (bool, *big.Int, error) {

	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	params := fmt.Sprintf(`[{"address":"%s"},{"address":"%s"}]`, walletAddr, ContractAddress)

	// This does not require TRX gas because it only reads data (not modifying the blockchain).
	result, err := client.TriggerConstantContract(
		walletAddr,
		TokenAddress,
		"allowance(address,address)",
		params,
	)
	if err != nil {
		log.Errorf("TriggerConstantContract error: %v", err)
		return false, nil, err
	}

	if result == nil || len(result.ConstantResult) == 0 {
		log.Errorf("Empty constant result from contract call: %+v", result)
		return false, nil, fmt.Errorf("invalid contract query result")
	}

	allowance := new(big.Int).SetBytes(result.ConstantResult[0])
	required := big.NewInt(ApprovalAmount)

	// return allowance.Cmp(requiredAmount) < 0, nil

	// Cmp() compares two big.Int values:
	// Cmp() == 0 → Both numbers are equal.
	// Cmp() == 1 → Allowance is greater.
	// Cmp() == -1 → Allowance is smaller.
	if allowance.Cmp(required) >= 0 {
		log.Info("Already approved")

		return false, allowance, nil
	}

	log.Info("Approval needed")
	return true, allowance, nil
}

// ApproveContract sends an approval transaction to the token contract.
func ApproveContract(ctx context.Context, client *client.GrpcClient, privateKey string) (string, error) {
	amount := new(big.Int).Mul(big.NewInt(ApprovalAmount), big.NewInt(1e6))
	params := fmt.Sprintf(`[{"address":"%s"},{"uint256":"%s"}]`, ContractAddress, amount.String())

	// Convert private key to TRON Base58 address
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// Convert public key to TRON address
	pubKey := privKey.PublicKey
	ethAddress := crypto.PubkeyToAddress(pubKey).Hex() // Ethereum-style hex address
	tronHexAddress := "41" + ethAddress[2:]            // Convert to Tron Hex format
	tronHexBytes, err := common.Hex2Bytes(tronHexAddress)
	if err != nil {
		return "", fmt.Errorf("failed to convert hex to bytes: %v", err)
	}
	// Convert Hex Tron address to Base58
	fromAddress := common.EncodeCheck(tronHexBytes)

	feeLimit := int64(1000000)

	// Call the contract
	txExt, err := client.TriggerContract(
		fromAddress, TokenAddress, "approve(address,uint256)", params,
		feeLimit, 0, "", 0,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create approve transaction: %v", err)
	}

	tx := txExt.GetTransaction()
	if tx == nil {
		return "", fmt.Errorf("transaction is nil")
	}

	// Sign the transaction
	signedTx, err := SignTransaction(tx, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Broadcast the transaction
	result, err := client.Broadcast(signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %v", err)
	}
	// Convert transaction hash from bytes to hex string
	txHash := hex.EncodeToString(txExt.Txid)

	fmt.Printf("Approval transaction sent: %s\n", result.String())
	return txHash, nil
}

// RevokeApproval revokes the approval by setting the allowance to zero.
func RevokeApproval(ctx context.Context, client *client.GrpcClient, privateKey string) (string, error) {
	amount := big.NewInt(0)
	params := fmt.Sprintf(`[{"address":"%s"},{"uint256":"%s"}]`, ContractAddress, amount.String())

	// Convert private key to TRON Base58 address
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// Convert public key to TRON address
	pubKey := privKey.PublicKey
	ethAddress := crypto.PubkeyToAddress(pubKey).Hex() // Ethereum-style hex address
	tronHexAddress := "41" + ethAddress[2:]            // Convert to Tron Hex format
	tronHexBytes, err := common.Hex2Bytes(tronHexAddress)
	if err != nil {
		return "", fmt.Errorf("failed to convert hex to bytes: %v", err)
	}
	// Convert Hex Tron address to Base58
	fromAddress := common.EncodeCheck(tronHexBytes)

	feeLimit := int64(1_000_000)

	// Call the contract
	txExt, err := client.TriggerContract(
		fromAddress, TokenAddress, "approve(address,uint256)", params,
		feeLimit, 0, "", 0,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create approve transaction: %v", err)
	}

	tx := txExt.GetTransaction()
	if tx == nil {
		return "", fmt.Errorf("transaction is nil")
	}

	// Sign the transaction
	signedTx, err := SignTransaction(tx, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Broadcast the transaction
	result, err := client.Broadcast(signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %v", err)
	}
	// Convert transaction hash from bytes to hex string
	txHash := hex.EncodeToString(txExt.Txid)

	fmt.Printf("Approval transaction sent: %s\n", result.String())
	return txHash, nil
}

// SignTransaction signs a TRON transaction using the provided private key.
func SignTransaction(tx *core.Transaction, privateKey string) (*core.Transaction, error) {
	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction raw data: %v", err)
	}

	hash := sha256.Sum256(rawData)

	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	signature, err := crypto.Sign(hash[:], privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	tx.Signature = append(tx.Signature, signature)
	return tx, nil
}

// BroadcastTransaction signs and broadcasts a TRON transaction
func BroadcastTransaction(client *client.GrpcClient, tx *core.Transaction, privateKey string) (string, error) {
	// Sign the transaction
	signedTx, err := SignTransaction(tx, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Compute transaction ID (Txid) manually
	rawData, err := proto.Marshal(signedTx.GetRawData())
	if err != nil {
		return "", fmt.Errorf("failed to marshal transaction raw data: %v", err)
	}
	txHash := sha256.Sum256(rawData)
	txHashHex := hex.EncodeToString(txHash[:])

	// Broadcast the transaction
	result, err := client.Broadcast(signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %v", err)
	}

	fmt.Printf("Transaction broadcasted successfully: %s\n", result)
	return txHashHex, nil
}


// BroadcastTransactionWithCalldata executes a TRON smart contract transaction using calldata.
func BroadcastTransactionWithCalldata(ctx context.Context, client *client.GrpcClient, contractAddress, functionSignature, calldata, privateKey string, feeLimit, callValue int64) (string, error) {
	// Convert private key to TRON Base58 address
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// Convert public key to TRON address
	pubKey := privKey.PublicKey
	ethAddress := crypto.PubkeyToAddress(pubKey).Hex() // Ethereum-style hex address
	tronHexAddress := "41" + ethAddress[2:]            // Convert to Tron Hex format
	tronHexBytes, err := common.Hex2Bytes(tronHexAddress)
	if err != nil {
		return "", fmt.Errorf("failed to convert hex to bytes: %v", err)
	}
	// Convert Hex Tron address to Base58
	fromAddress := common.EncodeCheck(tronHexBytes)

	// Call the contract
	txExt, err := client.TriggerContract(
		fromAddress,
		contractAddress,
		functionSignature, // Function name
		calldata,          // Encoded parameters
		feeLimit,
		callValue,
		"",
		0,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %v", err)
	}

	tx := txExt.GetTransaction()
	if tx == nil {
		return "", fmt.Errorf("transaction is nil")
	}

	// Sign the transaction
	signedTx, err := SignTransaction(tx, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Compute transaction ID (Txid) manually
	rawData, err := proto.Marshal(signedTx.GetRawData())
	if err != nil {
		return "", fmt.Errorf("failed to marshal transaction raw data: %v", err)
	}
	txHash := sha256.Sum256(rawData)
	txHashHex := hex.EncodeToString(txHash[:])

	// Broadcast the transaction
	result, err := client.Broadcast(signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %v", err)
	}

	fmt.Printf("Transaction broadcasted successfully: %s\n",result)
	return txHashHex, nil
}