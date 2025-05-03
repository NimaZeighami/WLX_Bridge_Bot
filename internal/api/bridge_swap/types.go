package bridge_swap

type GetQuoteRequest struct {
	FromToken         string  `json:"fromToken" validate:"required"`
	FromTokenChain  string  `json:"fromTokenChain" validate:"required"`
	ToToken           string  `json:"toToken" validate:"required"`
	ToTokenChain     string  `json:"toTokenChain" validate:"required"`
	FromTokenAmount   float64 `json:"fromTokenAmount" validate:"required,gt=0"`
	FromWalletAddress string  `json:"fromWalletAddress" validate:"required"`
	ToWalletAddress   string  `json:"toWalletAddress" validate:"required"`
}

type Response struct {
	Message string `json:"message"`
}
