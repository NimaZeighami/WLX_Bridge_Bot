package polygon

import (
	"bridgebot/configs"
	log "bridgebot/internal/utils/logger"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const TokenAddress = "0xc2132D05D31c914a87C6611C10748AEb04B58e8F" // Polygon  USDT Contract Address on Mainnet

// extendedERC20ABI includes functions for approve, allowance, increaseAllowance, and decreaseAllowance.
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
	configs.LoadEnv("../../.env")

	client, err := ethclient.Dial(configs.GetRpcURL())
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
		// return "", fmt.Errorf("failed to suggest gas price: %w", err)
		gasPrice = big.NewInt(50e9)
		log.Warnf("Failed to get suggested gas price, using default: %v", err)
	}

	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120)) // Increase by 20%  // TODO: it might be omitted

	msg := ethereum.CallMsg{
		From:     from,
		To:       &to,
		Data:     data,
		Value:    big.NewInt(0),
		GasPrice: gasPrice,
	}

	// TODO: Check this
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		// If estimation fails, use a default gas limit
		gasLimit = uint64(21000)
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
