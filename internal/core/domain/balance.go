package domain

import "time"

// BalanceResult represents the complete balance information including fiat conversion.
type BalanceResult struct {
	CryptoSymbol  string    `json:"cryptoSymbol"`
	Address       string    `json:"address"`
	CryptoBalance float64   `json:"cryptoBalance"`
	FiatSymbol    string    `json:"fiatSymbol"`
	FiatValue     float64   `json:"fiatValue"`
	ExchangeRate  float64   `json:"exchangeRate"`
	Timestamp     time.Time `json:"timestamp"`
	Change24h     float64   `json:"change24h"`
	Error         *string   `json:"error,omitempty"`
}

// BalanceRequest represents a single balance request in a batch.
type BalanceRequest struct {
	CryptoSymbol string `json:"cryptoSymbol"`
	Address      string `json:"address"`
	FiatSymbol   string `json:"fiatSymbol"`
}
