package services

import (
	"bridgebot/configs"
	"bridgebot/internal/blockchain/polygon"
	log "bridgebot/internal/utils/logger"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// CheckPolygonApproval determines whether a token approval is required for a given owner and amount on the Polygon network.
func CheckPolygonApproval(ctx context.Context, owner string, TokenContractAddress string, requiredAmount *big.Int) bool {
	client, err := polygon.NewPolygonClient()
	if err != nil {
		log.Fatalf("Error initializing Polygon client: %v", err)
	}
	tokenAddress := common.HexToAddress(TokenContractAddress)
	spender := common.HexToAddress(owner)

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
	return isNeeded
}

// SubmitPolygonApproval submits an approval transaction on the Polygon network, allowing the specified spender to spend a given amount of tokens on behalf of the owner.
func SubmitPolygonApproval(ctx context.Context, owner string, TokenContractAddress, spenderAddress string, requiredAmount *big.Int) error {
	client, err := polygon.NewPolygonClient()
	if err != nil {
		log.Fatalf("Error initializing Polygon client: %v", err)
		return err
	}

	privateKeyHex := configs.GetPrivateKeyHex()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to convert private key: %v", err)
		return err
	}
	tokenAddress := common.HexToAddress(TokenContractAddress)
	spender := common.HexToAddress(spenderAddress)

	txHash, err := polygon.ApproveContract(client, tokenAddress, spender, requiredAmount, privateKey)
	if err != nil {
		log.Fatalf("Error approving contract: %v", err)
		return err
	}

	log.Infof("Approval successful! Transaction hash: %s", txHash)
	log.Infof("Check on PolygonScan: https://polygonscan.com/tx/%s", txHash)

	return nil
}
