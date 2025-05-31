package thebridgers

import (
	// "bridgebot/configs"
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/services"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
	// "math/big"
	// "strconv"
)

type TheBridgers struct{}


func (b TheBridgers) Quote(fromAmount, fromToken, fromChain, toToken, toChain, FromWalletAddress string, ctx context.Context) (amountOut string, err error) {

	quoteReq := bridgers.QuoteRequest{
		FromTokenAmount:  fromAmount,
		FromTokenAddress: fromToken,
		FromTokenChain:   fromChain,
		ToTokenAddress:   toToken,
		ToTokenChain:     toChain,
		EquipmentNo:      services.GenerateEquipmentNo(FromWalletAddress),
		UserAddr: 	  FromWalletAddress,
		SourceFlag:       "WBB",
	}

	quoteResp, err := bridgers.FetchQuote(ctx, quoteReq)
	if err != nil {
		log.Errorf("Error fetching quote: %v", err)
		return "0", fmt.Errorf("failed to fetch quote: %w", err)
	}

	amountOut = quoteResp.Data.TxData.ToTokenAmount
	 
	return amountOut, nil
}

// func (b *TheBridgers) ApprovalNeeded(fromAddress, fromTokenAddress string, requiredAmount int, ctx context.Context) (bool, error) {
// 	IsApprovalNeeded, err := services.CheckPolygonApproval(ctx, fromAddress,configs.GetBridgersContractAddr("POLYGON"), fromTokenAddress, big.NewInt(int64(requiredAmount)))
// 	if err != nil {
// 		log.Errorf("Error Checking Approval: %v", err)

// 		return false, err
// 	}
// 	return IsApprovalNeeded, nil
// }

// func (b *TheBridgers) Approve(fromAddress, fromTokenAddress, contractAddress string, requiredAmount int, ctx context.Context) error {
// 	services.SubmitPolygonApproval(ctx,fromTokenAddress,configs.GetBridgersContractAddr("POlYGON"),big.NewInt(int64(requiredAmount)))
// 	return nil
// }

// func (b *TheBridgers) CallData() (any, error) {
// 	// داده تراکنش آماده سازی کن
// 	return nil, nil
// }

// func (b *TheBridgers) Sign(txData any) (string, error) {
// 	// امضای تراکنش
// 	return "signed_tx_string", nil
// }

// func (b *TheBridgers) BroadCast(signedTx string) (string, error) {
// 	// ارسال تراکنش
// 	return "0x123456...", nil
// }
