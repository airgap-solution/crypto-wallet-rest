package ports

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
}

// Provider interface for getting balance with fiat conversion.
type Provider interface {
	GetBalance(symbol, address, fiatSymbol string) (*BalanceResult, error)
}

// CryptoProvider interface for individual cryptocurrency providers.
type CryptoProvider interface {
	GetBalance(address string) (float64, error)
}
