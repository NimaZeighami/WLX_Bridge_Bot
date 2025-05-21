package models

import (
	"time"

	"gorm.io/gorm"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type BridgeConfiguration struct {
	ID                           uint   `gorm:"primaryKey"`
	Network                      string `gorm:"not null"`
	ChainID                      int
	Token                        string `gorm:"not null"`
	TokenContractAddress         string `gorm:"not null"`
	TokenDecimals                int    `gorm:"not null"`
	BridgersSmartContractAddress string `gorm:"not null"`
	IsEnabled                    bool   `gorm:"default:true"`
	CreatedAt                    string `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt                    string `gorm:"default:CURRENT_TIMESTAMP"`
}

type TokenInfo struct {
	ChainID                      int
	TokenContractAddress         string
	TokenDecimals                int
	BridgersSmartContractAddress string
	IsEnabled                    bool
}

var DB *gorm.DB

type NetworkTokenPair struct {
	ID                uint   `gorm:"primaryKey"`
	FromNetworkSymbol string `gorm:"column:from_network_symbol"`
	FromTokenSymbol   string `gorm:"column:from_token_symbol"`
	ToNetworkSymbol   string `gorm:"column:to_network_symbol"`
	ToTokenSymbol     string `gorm:"column:to_token_symbol"`
	IsEnabled         bool   `gorm:"column:is_enabled"`
}

type Quote struct {
	ID               uint      `gorm:"primaryKey"`
	OrderId          string    `gorm:"column:order_id"`
	FromTokenAddress string    `gorm:"column:from_token_address"`
	ToTokenAddress   string    `gorm:"column:to_token_address"`
	FromChain        string    `gorm:"column:from_chain"`
	ToChain          string    `gorm:"column:to_chain"`
	FromToken        string    `gorm:"column:from_token"`
	ToToken          string    `gorm:"column:to_token"`
	FromCoinCode     string    `gorm:"column:from_coin_code"`
	ToCoinCode       string    `gorm:"column:to_coin_code"`
	FromAddress      string    `gorm:"column:from_address"`
	ToAddress        string    `gorm:"column:to_address"`
	FromAmount       string    `gorm:"column:from_amount"`
	ToAmountMin      string    `gorm:"column:to_amount_min"`
	TxHash           string    `gorm:"column:tx_hash"`
	State            string    `gorm:"column:state"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}
