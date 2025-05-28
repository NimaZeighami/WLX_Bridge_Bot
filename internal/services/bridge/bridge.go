package bridge
type BridgeProvider interface {
	Quote(fromAmount int, fromToken, fromChain, toToken, toChain string) (amountOut int, err error)
	ApprovalNeeded(fromAddress ,FromTokenAddress string, requiredAmount int) (bool, error)
	Approve(fromAddress ,FromTokenAddress string, requiredAmount int) error
	CallData(any) (any , error)
	Sign(txData any)(signedTx string, err error)
	BroadCast(signedTx string) (txHash string, err error)
}

//todo: make request and reponse struct 
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
