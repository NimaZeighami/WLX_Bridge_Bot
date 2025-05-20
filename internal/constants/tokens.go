package constants

// Token represents a token on a blockchain
type Token struct {
	Name        string
	Decimal     int
	ContractAddr string
	Symbol      string
	CoinCode    string
	ChainID     int
	ChainName   string
	ChainSymbol string
}

// Chain represents a blockchain network
type Chain struct {
	ID              int
	Name            string
	Symbol          string
	SupportedTokens []Token
}

// Chains initialized as a slice of Chain with embedded token data including chain info
var Chains = []Chain{
	{
		ID:     1,
		Name:   "Ethereum",
		Symbol: "ETH",
		SupportedTokens: []Token{
			{
				Name:        "Tether ERC20",
				Decimal:     6,
				ContractAddr:"0xdac17f958d2ee523a2206206994597c13d831ec7",
				Symbol:      "USDT(ETH)",
				CoinCode:    "USDT(ETH)",
				ChainID:     1,
				ChainName:   "Ethereum",
				ChainSymbol: "ETH",
			},
		},
	},
	{
		ID:     56,
		Name:   "Binance Smart Chain",
		Symbol: "BSC",
		SupportedTokens: []Token{
			{
				Name:        "Tether BSC",
				Decimal:     18,
				ContractAddr:"0x55d398326f99059ff775485246999027b3197955",
				Symbol:      "USDT(BSC)",
				CoinCode:    "USDT(BSC)",
				ChainID:     56,
				ChainName:   "Binance Smart Chain",
				ChainSymbol: "BSC",
			},
		},
	},
	{
		ID:     137,
		Name:   "Polygon",
		Symbol: "POLYGON",
		SupportedTokens: []Token{
			{
				Name:        "Tether Polygon",
				Decimal:     6,
				ContractAddr:"0xc2132d05d31c914a87c6611c10748aeb04b58e8f",
				Symbol:      "USDT(POL)",
				CoinCode:    "USDT(POLYGON)",
				ChainID:     137,
				ChainName:   "Polygon",
				ChainSymbol: "MATIC",
			},
		},
	},
	{
		ID:     728126428,
		Name:   "Tron",
		Symbol: "TRX",
		SupportedTokens: []Token{
			{
				Name:        "Tether TRC20",
				Decimal:     6,
				ContractAddr:"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
				Symbol:      "USDT(Tron)",
				CoinCode:    "USDT(TRX)",
				ChainID:     728126428,
				ChainName:   "Tron",
				ChainSymbol: "TRX",
			},
		},
	},

}




// TODO: convert to map enum (bsc) sturct  (contract addr , decimal, chainID, supportedTokensArr etc)
	