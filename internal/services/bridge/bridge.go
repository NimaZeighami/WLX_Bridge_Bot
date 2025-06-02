package bridge

import (
	"context"
	"math/big"
)

type BridgeProvider interface {
	Quote(fromAmount, fromToken, fromChain, toToken, toChain, fromWalletAddress string, ctx context.Context) (amountOut string, err error)
	ApprovalNeeded(fromAddress, fromTokenAddress, bridgeProviderContractaddress string, requiredAmount *big.Int, ctx context.Context) (bool, error) 
	Approve(fromAddress, fromTokenAddress, bridgeProviderContractaddress string, requiredAmount *big.Int, ctx context.Context) error 
	// CallData(any) (any, error)
	// Sign(txData any) (signedTx string, err error)
	// BroadCast(signedTx string) (txHash string, err error)
}

// todo: make parameters of BridgeProvider methods structs instead of multiple strings and ints
// ?> QuoteRequiredParams , ...

func FetchQuoteAmount(provider BridgeProvider, fromAmount, fromToken, fromChain, toToken, toChain, fromWalletAddr string, ctx context.Context) (string, error) {
	amountOut, err := provider.Quote(fromAmount, fromToken, fromChain, toToken, toChain, fromWalletAddr,ctx)
	if err != nil {
		return "0", err
	}
	return amountOut, nil
}

func CheckTokenApproval(provider BridgeProvider,fromAddress, fromTokenAddress, bridgeProviderContractaddress string, requiredAmount *big.Int, ctx context.Context) (bool, error) {
	isApprovalNeeded, err := provider.ApprovalNeeded(fromAddress, fromTokenAddress, bridgeProviderContractaddress , requiredAmount , ctx)
	if err != nil {
		return false, err
	}
	return isApprovalNeeded, nil
}

func RequestTokenApproval(provider BridgeProvider, fromAddress, fromTokenAddress, bridgeProviderContractaddress string, requiredAmount *big.Int, ctx context.Context) error {
	err := provider.Approve(fromAddress, fromTokenAddress, bridgeProviderContractaddress , requiredAmount , ctx )  
	if err != nil {
		return err
	}
	return nil
}
// func GenerateCallData(provider BridgeProvider, txData any) (any, error) {
// 	data, err := provider.CallData(txData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }
// func GenerateSignedTx(provider BridgeProvider, txData any) (string, error) {
// 	signedTx, err := provider.Sign(txData)
// 	if err != nil {
// 		return "", err
// 	}
// 	return signedTx, nil
// }
// func SendTransaction(provider BridgeProvider, signedTx string) (string, error) {
// 	txHash, err := provider.BroadCast(signedTx)
// 	if err != nil {
// 		return "", err
// 	}
// 	return txHash, nil
// }
