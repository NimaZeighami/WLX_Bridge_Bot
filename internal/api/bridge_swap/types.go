package bridge_swap

type GetQuoteRequest struct {
	FromToken         string  `json:"fromToken" validate:"required"`
	ToToken           string  `json:"toToken" validate:"required"`
	FromTokenAmount   float64 `json:"fromTokenAmount" validate:"required,gt=0"`
	FromWalletAddress string  `json:"fromWalletAddress" validate:"required"`
	ToWalletAddress   string  `json:"toWalletAddress" validate:"required"`
}

type Response struct {
	Message string `json:"message"`
}
