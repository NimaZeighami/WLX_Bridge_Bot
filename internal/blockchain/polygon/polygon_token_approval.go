// package polygon

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"math/big"
// 	"strings"

// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/ethereum/go-ethereum/ethclient"
// )

// const (
// 	TokenAddress = "0xc2132D05D31c914a87C6611C10748AEb04B58e8F" // Polygon  USDT Contract Address on Mainnet
// 	ContractAddress = "TPwezUWpEGmFBENNWJHwXHRG1D2NCEEt5s"
// 	PolygonNode     = "https://polygon-rpc.com" 
// 	PrivateKey = "73de51d8df89c729e384b228a7f093e9d55c58278b6270e15513bfc97cbc0746"
// )

// func ApproveERC20(tokenAddressStr, spenderAddressStr string, amount *big.Int) (string, error) {
// 	// Load environment variables for RPC URL and private key
// 	rpcURL := PolygonNode
// 	privKeyHex := PrivateKey
// 	if rpcURL == "" || privKeyHex == "" {
// 		return "", fmt.Errorf("missing RPC URL or Private Key in environment")
// 	}

// 	// Connect to the Polygon RPC endpoint
// 	client, err := ethclient.Dial(rpcURL)
// 	if err != nil {
// 		log.Printf("Failed to connect to Polygon RPC: %v", err)
// 		return "", err
// 	}
// 	// It's good practice to close the client when done, but if this function is short-lived, it may close on function exit.
// 	// defer client.Close()

// 	// Parse the private key to an *ecdsa.PrivateKey
// 	privateKey, err := crypto.HexToECDSA(privKeyHex)
// 	if err != nil {
// 		log.Printf("Failed to parse private key: %v", err)
// 		return "", err
// 	}
// 	// Derive the from address (public key) for logging or nonce retrieval
// 	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
// 	log.Printf("Using address %s to approve", fromAddress.Hex())

// 	// Get current chain ID for EIP-155 (Polygon mainnet = 137)
// 	ctx := context.Background()
// 	chainID, err := client.NetworkID(ctx)
// 	if err != nil {
// 		log.Printf("Failed to get network chain ID: %v", err)
// 		return "", err
// 	}

// 	// Prepare the ABI for ERC20 approve(address,uint256)
// 	erc20ABI := `[{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`
// 	contractABI, err := abi.JSON(strings.NewReader(erc20ABI))
// 	if err != nil {
// 		log.Printf("Failed to parse ERC20 ABI: %v", err)
// 		return "", err
// 	}

// 	// Convert input addresses to common.Address
// 	tokenAddress := common.HexToAddress(tokenAddressStr)
// 	spenderAddress := common.HexToAddress(spenderAddressStr)

// 	// Pack the ABI-encoded function call data for approve(spender, amount)
// 	data, err := contractABI.Pack("approve", spenderAddress, amount)
// 	if err != nil {
// 		log.Printf("Failed to pack approve call data: %v", err)
// 		return "", err
// 	}

// 	// Get the account nonce (using pending state to include pending transactions)
// 	nonce, err := client.PendingNonceAt(ctx, fromAddress)
// 	if err != nil {
// 		log.Printf("Failed to fetch nonce for %s: %v", fromAddress.Hex(), err)
// 		return "", err
// 	}

// 	// Get suggested gas price (legacy)
// 	gasPrice, err := client.SuggestGasPrice(ctx)
// 	if err != nil {
// 		log.Printf("Failed to get gas price: %v", err)
// 		return "", err
// 	}

// 	// Estimate gas limit required for the transaction
// 	callMsg := ethereum.CallMsg{From: fromAddress, To: &tokenAddress, Value: big.NewInt(0), Data: data, GasPrice: gasPrice}
// 	gasLimit, err := client.EstimateGas(ctx, callMsg)
// 	if err != nil {
// 		// If estimation fails, use a reasonable default for ERC-20 approval
// 		gasLimit = 100000 // 100k gas as a fallback
// 		log.Printf("Gas estimation failed, using default gas limit %d: %v", gasLimit, err)
// 	}

// 	// Create the unsigned transaction (here we use a legacy transaction format)
// 	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), gasLimit, gasPrice, data)

// 	// Sign the transaction with the private key and chain ID (EIP-155 signer)
// 	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
// 	if err != nil {
// 		log.Printf("Failed to sign transaction: %v", err)
// 		return "", err
// 	}

// 	// Send the signed transaction to the network
// 	err = client.SendTransaction(ctx, signedTx)
// 	if err != nil {
// 		log.Printf("Failed to send transaction: %v", err)
// 		return "", err
// 	}

// 	txHash := signedTx.Hash().Hex()
// 	log.Printf("Approve transaction submitted: %s", txHash)
// 	return txHash, nil
// }

// Sample code for Polygon contract approval and transaction signing
// package polygon

// import ()

// func NewPolygonClient() () {
// 	// Initialize the Polygon client here

// }

// func IsApprovalNeeded() bool {
// 	// Implement the logic to determine if approval is needed
// 	return true
// }

// func ApproveContract(){

// }

// func RevokeApproval(){

// }

// func SignTransaction(){

// }

// func BroadcastTransaction(){

// }

// func BroadcastTransactionWithCalldata(){

// }

// (Recommended Additions)

// GetCurrentAllowance()

// Checks how much allowance a spender has.

// Why? Helps IsApprovalNeeded() make decisions.

// IncreaseAllowance() / DecreaseAllowance()

// Safer than approve() (avoids race conditions).

// BatchApprove()

// Approves multiple tokens in one transaction (gas efficiency).

// polygon_token_approval.go
// Sample code for Polygon contract approval and transaction signing.
// This file implements complete functionality for interacting with an ERC20 token on Polygon.
package polygon

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	log "bridgebot/internal/utils/logger"
	"math/big"
	"os"
	"strings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum"
)

// extendedERC20ABI includes functions for approve, allowance, increaseAllowance, and decreaseAllowance.

const (
	TokenAddress = "0xc2132D05D31c914a87C6611C10748AEb04B58e8F" // Polygon  USDT Contract Address on Mainnet
	ContractAddress = "TPwezUWpEGmFBENNWJHwXHRG1D2NCEEt5s"
	PolygonNode     = "https://polygon-rpc.com" 
	PrivateKey = "73de51d8df89c729e384b228a7f093e9d55c58278b6270e15513bfc97cbc0746"
)

const extendedERC20ABI = `
[
	{
		"constant": false,
		"inputs": [
			{"name": "spender", "type": "address"},
			{"name": "value", "type": "uint256"}
		],
		"name": "approve",
		"outputs": [{"name": "", "type": "bool"}],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{"name": "owner", "type": "address"},
			{"name": "spender", "type": "address"}
		],
		"name": "allowance",
		"outputs": [{"name": "", "type": "uint256"}],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{"name": "spender", "type": "address"},
			{"name": "addedValue", "type": "uint256"}
		],
		"name": "increaseAllowance",
		"outputs": [{"name": "", "type": "bool"}],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{"name": "spender", "type": "address"},
			{"name": "subtractedValue", "type": "uint256"}
		],
		"name": "decreaseAllowance",
		"outputs": [{"name": "", "type": "bool"}],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]
`

// erc20ParsedABI is the parsed ABI for ERC20 functions.
var erc20ParsedABI abi.ABI

func init() {
	var err error
	erc20ParsedABI, err = abi.JSON(strings.NewReader(extendedERC20ABI))
	if err != nil {
		log.Fatalf("Failed to parse ERC20 ABI: %v", err)
	}
}

// NewPolygonClient creates and returns an Ethereum client connected to Polygon using POLYGON_RPC_URL.
func NewPolygonClient() (*ethclient.Client, error) {
	rpcURL := os.Getenv("POLYGON_RPC_URL")
	if rpcURL == "" {
		return nil, fmt.Errorf("POLYGON_RPC_URL not set in environment")
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Polygon RPC: %w", err)
	}
	return client, nil
}

// signAndSendTx is a helper that builds, signs, and sends a transaction with the given calldata.
func signAndSendTx(ctx context.Context, client *ethclient.Client, from common.Address, to common.Address, data []byte, privateKey *ecdsa.PrivateKey) (string, error) {
	// Get the account nonce.
	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to suggest gas price: %w", err)
	}

	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120)) // Increase by 20%

	// Estimate gas limit.
	msg := ethereum.CallMsg{
		From:     from,
		To:       &to,
		Data:     data,
		Value:    big.NewInt(0),
		GasPrice: gasPrice,
	}

gasLimit, err := client.EstimateGas(ctx, msg)
if err != nil {
    // If estimation fails, use a default gas limit
    gasLimit = 100000
    log.Warnf("Gas estimation failed, using default gas limit: %v", err)
}

	// Build the transaction.
	tx := types.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, data)

	// Get chain ID.
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		chainID = big.NewInt(137) // Default to Polygon mainnet
	log.Infof("Failed to get chain ID, defaulting to 137: %v", err)
	}

	// Sign the transaction.
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send the transaction.
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}

// GetCurrentAllowance queries the token contract for the current allowance granted by owner to spender.
func GetCurrentAllowance(client *ethclient.Client, tokenAddress, owner, spender common.Address) (*big.Int, error) {
	ctx := context.Background()
	data, err := erc20ParsedABI.Pack("allowance", owner, spender)
	if err != nil {
		return nil, fmt.Errorf("failed to pack allowance data: %w", err)
	}
	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}
	result, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}
	results, err := erc20ParsedABI.Unpack("allowance", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack allowance result: %w", err)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no result from allowance call")
	}
	allowance, ok := results[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected type for allowance")
	}
	return allowance, nil
}

// IsApprovalNeeded checks whether the current allowance is less than the required amount.
func IsApprovalNeeded(client *ethclient.Client, tokenAddress, owner, spender common.Address, requiredAmount *big.Int) (bool, error) {
	currentAllowance, err := GetCurrentAllowance(client, tokenAddress, owner, spender)
	if err != nil {
		return false, err
	}
	return currentAllowance.Cmp(requiredAmount) < 0, nil
}

// ApproveContract sends an approval transaction to allow spender to spend a specified amount.
func ApproveContract(client *ethclient.Client, tokenAddress, spender common.Address, amount *big.Int, privateKey *ecdsa.PrivateKey) (string, error) {
	ctx := context.Background()
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	data, err := erc20ParsedABI.Pack("approve", spender, amount)
	if err != nil {
		return "", fmt.Errorf("failed to pack approve data: %w", err)
	}
	return signAndSendTx(ctx, client, fromAddress, tokenAddress, data, privateKey)
}

// RevokeApproval revokes any previously granted approval by setting the allowance to zero.
func RevokeApproval(client *ethclient.Client, tokenAddress, spender common.Address, privateKey *ecdsa.PrivateKey) (string, error) {
	// Revocation is typically done by approving 0.
	return ApproveContract(client, tokenAddress, spender, big.NewInt(0), privateKey)
}

// IncreaseAllowance increases the current allowance by the specified added value.
func IncreaseAllowance(client *ethclient.Client, tokenAddress, spender common.Address, addedValue *big.Int, privateKey *ecdsa.PrivateKey) (string, error) {
	ctx := context.Background()
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	data, err := erc20ParsedABI.Pack("increaseAllowance", spender, addedValue)
	if err != nil {
		return "", fmt.Errorf("failed to pack increaseAllowance data: %w", err)
	}
	return signAndSendTx(ctx, client, fromAddress, tokenAddress, data, privateKey)
}

// DecreaseAllowance decreases the current allowance by the specified subtracted value.
func DecreaseAllowance(client *ethclient.Client, tokenAddress, spender common.Address, subtractedValue *big.Int, privateKey *ecdsa.PrivateKey) (string, error) {
	ctx := context.Background()
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	data, err := erc20ParsedABI.Pack("decreaseAllowance", spender, subtractedValue)
	if err != nil {
		return "", fmt.Errorf("failed to pack decreaseAllowance data: %w", err)
	}
	return signAndSendTx(ctx, client, fromAddress, tokenAddress, data, privateKey)
}

// BatchApprovalItem represents a single approval instruction.
type BatchApprovalItem struct {
	Token   common.Address
	Spender common.Address
	Amount  *big.Int
}

// BatchApprove processes multiple approval transactions sequentially and returns their transaction hashes.
func BatchApprove(client *ethclient.Client, items []BatchApprovalItem, privateKey *ecdsa.PrivateKey) ([]string, error) {
	var txHashes []string
	for _, item := range items {
		txHash, err := ApproveContract(client, item.Token, item.Spender, item.Amount, privateKey)
		if err != nil {
			return txHashes, fmt.Errorf("failed to approve token %s for spender %s: %w", item.Token.Hex(), item.Spender.Hex(), err)
		}
		txHashes = append(txHashes, txHash)
	}
	return txHashes, nil
}

// SignTransaction signs a given unsigned transaction using the provided private key.
// This function does not broadcast the transaction.
func SignTransaction(client *ethclient.Client, tx *types.Transaction, privateKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	ctx := context.Background()
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		chainID = big.NewInt(137)
		log.Infof("Failed to get chain ID, defaulting to 137: %v", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}
	return signedTx, nil
}

// BroadcastTransaction sends a pre-signed transaction to the network.
func BroadcastTransaction(client *ethclient.Client, signedTx *types.Transaction) (string, error) {
	ctx := context.Background()
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}
	return signedTx.Hash().Hex(), nil
}

// BroadcastTransactionWithCalldata builds, signs, and broadcasts a transaction using custom calldata.
func BroadcastTransactionWithCalldata(client *ethclient.Client, to common.Address, calldata []byte, privateKey *ecdsa.PrivateKey) (string, error) {
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	return signAndSendTx(context.Background(), client, fromAddress, to, calldata, privateKey)
}
