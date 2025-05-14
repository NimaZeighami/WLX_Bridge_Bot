// Bridger interface (GetQuote, Approve, Sign, Broadcast)

package bridge

import (
	"context"
)

type Bridger interface {
	GetQuote(ctx context.Context, req QuoteRequest) (QuoteResponse, error)
    // TODO: Uncomment and implement the following methods as needed
	// IsApproveNeeded(ctx context.Context, req ApproveRequest) (ApproveResponse, error)
	// Approve(ctx context.Context, req ApproveRequest) (ApproveResponse, error)
	// Sign(ctx context.Context, req SignRequest) (SignedTransaction, error)
	// Broadcast(ctx context.Context, tx SignedTransaction) (BroadcastResult, error)
}

type QuoteRequest struct {
	FromToken         string  `json:"fromToken"`
	FromTokenChain         string  `json:"fromTokenChain"`
	FromTokenAmount   float64 `json:"fromTokenAmount"`
	ToToken           string  `json:"toToken"`
	ToTokenChain           string  `json:"ToTokenChain"`
	FromWalletAddress string  `json:"fromWalletAddress"`
	ToWalletAddress   string  `json:"toWalletAddress"`
}

type QuoteResponse struct { 
    QuoteId int `json:"quoteId"`
    ToTokenAmount float64 `json:"toTokenAmount"`
}

type ApproveFields struct { 
    TokenAddress string `json:"tokenAddress"`
    TokenChain string `json:"tokenChain"`
    RequiredAmount string `json:"requiredAmount"`
}

type IsApproveNedded struct { 
    RequiredAmount float64 `json:"requiredAmount"`
    CurrentApprovedAmount float64 `json:"currentApprovedAmount"`
}








type SignRequest struct { /* fields */
}
type SignedTransaction struct { /* fields */
}

type BroadcastResult struct { /* fields */
}
