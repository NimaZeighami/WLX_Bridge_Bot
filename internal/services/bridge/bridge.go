// Bridger interface (GetQuote, Approve, Sign, Broadcast)

package bridge

import (
    "context"
)

type Bridge interface {
    GetQuote(ctx context.Context, req QuoteRequest) (QuoteResponse, error)
    Approve(ctx context.Context, req ApproveRequest) (ApproveResponse, error)
    Sign(ctx context.Context, req SignRequest) (SignedTransaction, error)
    Broadcast(ctx context.Context, tx SignedTransaction) (BroadcastResult, error)
}

// Example request/response structs
type QuoteRequest struct { /* fields */ }
type QuoteResponse struct { /* fields */ }

type ApproveRequest struct { /* fields */ }
type ApproveResponse struct { /* fields */ }

type SignRequest struct { /* fields */ }
type SignedTransaction struct { /* fields */ }

type BroadcastResult struct { /* fields */ }
