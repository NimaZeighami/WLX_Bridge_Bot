package orchestration

import (
	"bridgebot/internal/blockchain/polygon"
	log "bridgebot/internal/utils/logger"
	"context"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	ApprovalAmount = 13000000 // Amount to approve
)

func CheckPolygonApproval(ctx context.Context, owner string, tokenAddressStr string) {
	client, err := polygon.NewPolygonClient()
	if err != nil {
		log.Fatalf("Error initializing Polygon client: %v", err)
	}
	tokenAddress := common.HexToAddress(tokenAddressStr)
	spender := common.HexToAddress(owner)
	requiredAmount := big.NewInt(ApprovalAmount)

	log.Info("Checking if approval is needed...")
	isNeeded, err := polygon.IsApprovalNeeded(client, tokenAddress, common.HexToAddress(owner), spender, requiredAmount)
	if err != nil {
		log.Fatalf("Error checking approval status: %v", err)
	}

	if !isNeeded {
		log.Info("Approval is already granted!")
	} else {
		log.Info("Approval is needed!!")
	}
}

func SubmitPolygonApproval(ctx context.Context, owner string, tokenAddressStr, spenderAddress string) {
	client, err := polygon.NewPolygonClient()
	if err != nil {
		log.Fatalf("Error initializing Polygon client: %v", err)
	}

	privateKeyHex := os.Getenv("POLYGON_PRIVATE_KEY")
	if privateKeyHex == "" {
		privateKeyHex = "dde619e9c94141eb5c60cf3c52e812f95db0a593543767a59e6b12e133a40c6d"
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to convert private key: %v", err)
	}
	tokenAddress := common.HexToAddress(tokenAddressStr)
	spender := common.HexToAddress(spenderAddress)
	requiredAmount := big.NewInt(ApprovalAmount)

	txHash, err := polygon.ApproveContract(client, tokenAddress, spender, requiredAmount, privateKey)
	if err != nil {
		log.Fatalf("Error approving contract: %v", err)
	}

	log.Infof("Approval successful! Transaction hash: %s", txHash)
	log.Infof("Check on PolygonScan: https://polygonscan.com/tx/%s", txHash)

}
