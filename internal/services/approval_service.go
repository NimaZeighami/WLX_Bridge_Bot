// File Name: approval_service.go
// This file contains the SubmitPolygonApproval function, which is responsible
// for submitting an approval transaction on the Polygon blockchain. It interacts
// with a specified token contract to approve a spender address for a required
// amount, enabling further operations such as token transfers or bridging.

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

// const (
// 	BridgingAmount = 3_000_000 // Amount to approve, get qoute from bridgers and sing & broadcast transaction
// )

func CheckPolygonApproval(ctx context.Context, fromWalletAddress, bridgeProviderContractAddr, TokenContractAddress string, requiredAmount *big.Int) (bool, error) {
	client, err := polygon.NewPolygonClient()
	if err != nil {
		log.Fatalf("Error initializing Polygon client: %v", err)
		return false, err
	}
	tokenAddress := common.HexToAddress(TokenContractAddress)

	log.Info("Checking if approval is needed...")
	isNeeded, err := polygon.IsApprovalNeeded(client, tokenAddress, common.HexToAddress(fromWalletAddress), common.HexToAddress(bridgeProviderContractAddr), requiredAmount)
	if err != nil {
		log.Fatalf("Error checking approval status: %v", err)
		return false, err
	}

	if !isNeeded {
		log.Info("Approval is already granted!")
	} else {
		log.Info("Approval is needed!!")
	}
	return isNeeded, nil
}

func SubmitPolygonApproval(ctx context.Context, TokenContractAddress, bridgeProviderContractAddress string, requiredAmount *big.Int) error {
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
	bridgeProviderContractAddr := common.HexToAddress(bridgeProviderContractAddress)
	// requiredAmount := big.NewInt(BridgingAmount)

	txHash, err := polygon.ApproveContract(client, tokenAddress, bridgeProviderContractAddr, requiredAmount, privateKey)
	if err != nil {
		log.Fatalf("Error approving contract: %v", err)
		return err
	}

	log.Infof("Approval successful! Transaction hash: %s", txHash)
	log.Infof("Check on PolygonScan: https://polygonscan.com/tx/%s", txHash)

	return nil
}
