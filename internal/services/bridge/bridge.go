package bridge

import (
	"context"
	"math/big"
)

type BridgeProvider interface {
	Quote(fromAmount, fromToken, fromChain, toToken, toChain, fromWalletAddress string, ctx context.Context) (amountOut string, err error)
	ApprovalNeeded(fromAddress, fromTokenAddress, bridgeProviderContractaddress string, requiredAmount *big.Int, ctx context.Context) (bool, error)
	Approve(fromAddress, fromTokenAddress, bridgeProviderContractaddress string, requiredAmount *big.Int, ctx context.Context) error
	BroadCast(fromAddress, toAddress, fromToken, toToken, ToAmountMin string, fromAmount *big.Int, ctx context.Context) (txhash string, err error)
}

// todo: make parameters of BridgeProvider methods structs instead of multiple strings and ints
// ?> QuoteRequiredParams , ...

func FetchQuoteAmount(provider BridgeProvider, fromAmount, fromToken, fromChain, toToken, toChain, fromWalletAddr string, ctx context.Context) (string, error) {
	amountOut, err := provider.Quote(fromAmount, fromToken, fromChain, toToken, toChain, fromWalletAddr, ctx)
	if err != nil {
		return "0", err
	}
	return amountOut, nil
}

func CheckTokenApproval(provider BridgeProvider, fromAddress, fromTokenAddress, bridgeProviderContractaddress string, requiredAmount *big.Int, ctx context.Context) (bool, error) {
	isApprovalNeeded, err := provider.ApprovalNeeded(fromAddress, fromTokenAddress, bridgeProviderContractaddress, requiredAmount, ctx)
	if err != nil {
		return false, err
	}
	return isApprovalNeeded, nil
}

func RequestTokenApproval(provider BridgeProvider, fromAddress, fromTokenAddress, bridgeProviderContractaddress string, requiredAmount *big.Int, ctx context.Context) error {
	err := provider.Approve(fromAddress, fromTokenAddress, bridgeProviderContractaddress, requiredAmount, ctx)
	if err != nil {
		return err
	}
	return nil
}


func FinalizeTransaction(provider BridgeProvider,fromAddress, toAddress, fromToken, toToken, ToAmountMin string, fromAmount *big.Int, ctx context.Context) (string, error) {
	txHash, err := provider.BroadCast(fromAddress, toAddress, fromToken, toToken, ToAmountMin , fromAmount , ctx )
	if err != nil {
		return "", err
	}
	return txHash, nil
}
