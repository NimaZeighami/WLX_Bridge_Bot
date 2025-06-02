// File Name: approval_service.go
// This file contains the SubmitPolygonApproval function, which is responsible
// for submitting an approval transaction on the Polygon blockchain. It interacts
// with a specified token contract to approve a spender address for a required
// amount, enabling further operations such as token transfers or bridging.

//todo: move these functions to blockchain/polygon package

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

	txHash, err := polygon.ApproveContract(client, tokenAddress, bridgeProviderContractAddr, requiredAmount, privateKey)
	if err != nil {
		log.Fatalf("Error approving contract: %v", err)
		return err
	}

	log.Infof("Approval successful! Transaction hash: %s", txHash)
	log.Infof("Check on PolygonScan: https://polygonscan.com/tx/%s", txHash)

	return nil
}

//todo: pass polygon client as a dependency to the functions instead of creating a new client each time.
//todo: create a generic approval service that can handle multiple networks
//todo: create a network interface that defines the methods for checking approval, submitting approval, signing transactions, and sending them.
// network -> interface  , methods : checkApproval, submitApproval, sign , send 
// each network -> struct
func CheckApproval(ctx context.Context, fromWalletAddress, bridgeProviderContractAddr, TokenContractAddress string, requiredAmount *big.Int, network string ) (bool, error) {
	// This function is a placeholder for future implementations.
	// Currently, it does not perform any operations and returns false.
	log.Info("CheckApproval function is not implemented yet.")
	return false, nil
}
