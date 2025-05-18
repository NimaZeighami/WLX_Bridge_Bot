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

const (
	BridgingAmount = 3_000_000 // Amount to approve, get qoute from bridgers and sing & broadcast transaction
)

func CheckPolygonApproval(ctx context.Context, owner string, tokenAddressStr string, /*requiredAmount *big.Int */) bool {
	client, err := polygon.NewPolygonClient()
	if err != nil {
		log.Fatalf("Error initializing Polygon client: %v", err)
	}
	tokenAddress := common.HexToAddress(tokenAddressStr)
	spender := common.HexToAddress(owner)
	requiredAmount := big.NewInt(BridgingAmount)

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

func SubmitPolygonApproval(ctx context.Context, owner string, tokenAddressStr, spenderAddress string, /*requiredAmount *big.Int */ ) error {
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
	tokenAddress := common.HexToAddress(tokenAddressStr)
	spender := common.HexToAddress(spenderAddress)
	requiredAmount := big.NewInt(BridgingAmount)

	txHash, err := polygon.ApproveContract(client, tokenAddress, spender, requiredAmount, privateKey)
	if err != nil {
		log.Fatalf("Error approving contract: %v", err)
		return err
	}

	log.Infof("Approval successful! Transaction hash: %s", txHash)
	log.Infof("Check on PolygonScan: https://polygonscan.com/tx/%s", txHash)

	return nil
}
