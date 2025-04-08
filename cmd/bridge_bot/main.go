package main

import (
	"bridgebot/internal/blockchain/tron"
	"bridgebot/internal/client/http/bridgers"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// generateEquipmentNo ensures a fixed 32-character length for user addresses.
func generateEquipmentNo(userAddr string) string {
	if len(userAddr) >= 32 {
		return userAddr[:32]
	}
	return fmt.Sprintf("%032s", userAddr) // Pad with zeros if less than 32
}

// getUSDTAddress retrieves the USDT address from a given token map.
func getUSDTAddress(tokenMap map[string]string, symbol string) (string, error) {
	address, exists := tokenMap[symbol]
	if !exists {
		return "", fmt.Errorf("USDT token not found on %s", symbol)
	}
	return address, nil
}

func main() {
	log.Info("Starting Bridge Bot...")

	// Create a root context with cancellation support
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals to gracefully shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Warn("Received termination signal. Shutting down gracefully...")
		cancel()
		os.Exit(0)
	}()

	// Fetch tokens
	userAddr := "TMGcWzEDiECVCwAxoprCedtXSeuJthq4AA"
	chains := []string{"TRX", "BSC"}
	crossChainTokens := make(map[string]string)

	for _, chain := range chains {
		if err := bridgers.FetchAndMapTokens(ctx, chain, crossChainTokens); err != nil {
			log.Error(err.Error())
			return
		}
	}

	// Retrieve USDT token addresses
	usdtTRX, err := getUSDTAddress(crossChainTokens, "USDT(TRON)")
	if err != nil {
		log.Error(err.Error())
		return
	}

	usdtBSC, err := getUSDTAddress(crossChainTokens, "USDT(BSC)")
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Infof("USDT(TRX) Address: %s", usdtTRX)
	log.Infof("USDT(BSC) Address: %s", usdtBSC)

	// Quote request
	request := bridgers.QuoteRequest{
		FromTokenAddress: usdtTRX,
		ToTokenAddress:   usdtBSC,
		FromTokenAmount:  "100000000",
		FromTokenChain:   "TRX",
		ToTokenChain:     "BSC",
		UserAddr:         "",
		EquipmentNo:      generateEquipmentNo(userAddr),
		SourceFlag:       "WLXBridgeApp",
		SourceType:       "",
	}

	response, err := bridgers.RequestQuote(ctx, request)
	if err != nil {
		log.Errorf("Error fetching quote: %v", err)
		return
	}

	log.Info("Quote Response:")
	fmt.Println("############## ##### USDT(TRX) to USDT(BSC)")
	fmt.Println("############## ##### FromTokenAmount: ", response.Data.TxData.FromTokenAmount, "in TRX")
	fmt.Println("############## ##### ToTokenAmount: ", response.Data.TxData.ToTokenAmount, "in BSC")

	privateKey := "73643967cb61b7e712370b97e8fadde9fe866d1809eec00bc484dd9fe2a7b8f3"

	// Create a TRON client with timeout-based context
	log.Info("Creating a client and Connecting to TRON network...")
	trxCtx, trxCancel := context.WithTimeout(ctx, 10*time.Second)
	defer trxCancel()

	client, err := tron.NewTronClient()
	if err != nil {
		log.Errorf("Error connecting to TRON network: %v", err)
		return
	}

	// Check if approval is needed
	walletAddr := "TMGcWzEDiECVCwAxoprCedtXSeuJthq4AA"
	log.Info("Checking if approval is needed...")

	needsApproval, approvedAmount, err := tron.IsApprovalNeeded(trxCtx, client, walletAddr)
	if err != nil {
		log.Errorf("Error checking approval status: %v", err)
		return
	}

	log.Infof("Current approval amount: %s USDT (TRC-20 units)", approvedAmount.String())
	log.Infof("Current approval amount: %s USDT ", approvedAmount.Div(approvedAmount, big.NewInt(1e6)).String())

	if !needsApproval {
		log.Info("This wallet is already approved, no further approval is needed.")
	} else {
		// Execute approval with a timeout context
		approveCtx, approveCancel := context.WithTimeout(ctx, 20*time.Second)
		defer approveCancel()

		txHash, err := tron.ApproveContract(approveCtx, client, privateKey)
		if err != nil {
			log.Errorf("Error executing approval: %v", err)
			return
		}

		log.Infof("Approval transaction successfully submitted! Tx Hash: %s", txHash)
		log.Infof("Check the transaction on TronScan: https://nile.tronscan.org/#/transaction/%s", txHash)
	}

	receiverAddr := "0x5aA96F60C1aFf555c43552931a177728f32fcA27"

	CallDataRequest := bridgers.CallDataRequest{
		FromTokenAddress: usdtTRX,
		ToTokenAddress:   usdtBSC,
		FromAddress:      userAddr,
		ToAddress:        receiverAddr,
		FromTokenChain:   "TRX",
		ToTokenChain:     "BSC",
		FromTokenAmount:  "100000000", // 100 USDT (USDT has 6 decimals)
		AmountOutMin:     "99000000",  // Expected min amount after slippage
		FromCoinCode:     "USDT(TRX)",
		ToCoinCode:       "USDT(BSC)",
		EquipmentNo:      generateEquipmentNo(userAddr),
		SourceType:       "",
		SourceFlag:       "BridgeBot",
		Slippage:         "0.1",
	}
	log.Info("Prepare Bridge Transaction Data ...")
	callData, err := bridgers.FetchBridgeCallData(ctx, CallDataRequest)
	if err != nil {
		log.Errorf("Swap request failed: %v", err)
		panic(err)
	}

	log.Infof("Transaction Destination : %s", callData.Data.TxData.To)

	// Execute the transaction using calldata
	log.Info("Executing Bridge Transaction on TRON...")

	// Set transaction parameters
	contractAddress := callData.Data.TxData.To               // Smart contract address
	functionSignature := callData.Data.TxData.FunctionName   // Function to be called
	calldata := callData.Data.TxData.Parameter               // Encoded parameters for the contract call
	feeLimit := int64(callData.Data.TxData.Options.FeeLimit) // Maximum TRX fee allowed
	callValue := int64(0)                                    // No TRX is being sent in this contract call

	// Convert parameters to JSON string format
	calldataStr := "["
	for i, param := range calldata {
		calldataStr += fmt.Sprintf(`{"type":"%s","value":"%s"}`, param.Type, param.Value)
		if i < len(calldata)-1 {
			calldataStr += ","
		}
	}
	calldataStr += "]"

	// Broadcast transaction
	txHash, err := tron.BroadcastTransactionWithCalldata(ctx, client, contractAddress, functionSignature, calldataStr, privateKey, feeLimit, callValue)
	if err != nil {
		log.Errorf("Error executing bridge transaction: %v", err)
		return
	}

	log.Infof("Bridge transaction successfully submitted! Tx Hash: %s", txHash)
	log.Infof("Check the transaction on TronScan: https://nile.tronscan.org/#/transaction/%s", txHash)
}
