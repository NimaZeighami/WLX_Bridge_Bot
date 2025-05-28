package thebridgers

import (
	"bridgebot/internal/client/http/bridgers"
	"bridgebot/internal/services"
	log "bridgebot/internal/utils/logger"
	"context"
	"fmt"
	"strconv"
)

type TheBridgers struct{}

func BuildQuoteRequest(userAddress, fromTokenAddr, toTokenAddr, fromTokenChain, toTokenChain string, bridgingAmount uint) bridgers.QuoteRequest {
	return bridgers.QuoteRequest{
		FromTokenAddress: fromTokenAddr,
		ToTokenAddress:   toTokenAddr,
		FromTokenAmount:  fmt.Sprintf("%d", bridgingAmount),
		FromTokenChain:   fromTokenChain,
		ToTokenChain:     toTokenChain,
		EquipmentNo:      services.GenerateEquipmentNo(userAddress),
		SourceFlag:       "WBB",
		SourceType:       "",
		UserAddr:         userAddress,
	}
}

func RequestQuote(ctx context.Context, req bridgers.QuoteRequest) (*bridgers.QuoteResponse,error) {
	resp, err := bridgers.FetchQuote(ctx, req)
	if err != nil {
		log.Errorf("Error fetching quote: %v", err)
		return nil , err
	}
	return resp , nil
}

func (b *TheBridgers) Quote(userAddr string, fromAmount int, fromTokenAddr, fromChain, toTokenAddr, toChain string, ctx context.Context) (int64, error) {
	quoteReq := BuildQuoteRequest(userAddr, fromTokenAddr, toTokenAddr, fromChain, toChain, uint(fromAmount))
	quoteRes, err := RequestQuote(ctx, quoteReq)
	if err != nil {
		log.Errorf("Error fetching the bridgers quote: %v", err)
		return 0 , err
	}
	toTokenAmount , err:=strconv.ParseInt(quoteRes.Data.TxData.ToTokenAmount,10 , 64)
	if  err != nil {
		log.Errorf("Error parsing token Amount : %v", err)
		return 0 , err
	}

	return toTokenAmount, nil
}

func (b *TheBridgers) ApprovalNeeded(fromAddress, fromTokenAddress string, requiredAmount int) (bool, error) {

	return true, nil
}

func (b *TheBridgers) Approve(fromAddress, fromTokenAddress string, requiredAmount int) error {
	// از services.SubmitPolygonApproval استفاده کن
	return nil
}

func (b *TheBridgers) CallData() (any, error) {
	// داده تراکنش آماده سازی کن
	return nil, nil
}

func (b *TheBridgers) Sign(txData any) (string, error) {
	// امضای تراکنش
	return "signed_tx_string", nil
}

func (b *TheBridgers) BroadCast(signedTx string) (string, error) {
	// ارسال تراکنش
	return "0x123456...", nil
}
