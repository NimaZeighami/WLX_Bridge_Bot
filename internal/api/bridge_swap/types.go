package bridge_swap

type QuoteReq struct {
	FromToken         string  `json:"fromToken" validate:"required"`
	FromTokenChain    string  `json:"fromTokenChain" validate:"required"`
	ToToken           string  `json:"toToken" validate:"required"`
	ToTokenChain      string  `json:"toTokenChain" validate:"required"`
	FromTokenAmount   float64 `json:"fromTokenAmount" validate:"required,gt=0"`
	FromWalletAddress string  `json:"fromWalletAddress" validate:"required"`
	ToWalletAddress   string  `json:"toWalletAddress" validate:"required"`
}

type QuoteRes struct {
	Message string `json:"message"`
}

type SwapReq struct {
	QuoteId  string  `json:"quoteId" validate:"required"`
}

type SwapRes struct {
	Message string `json:"message"`
}