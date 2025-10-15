package ports

import "github.com/airgap-solution/crypto-wallet-rest/internal/core/domain"

// Provider interface for getting balance with fiat conversion.
type Provider interface {
	GetBalance(symbol, address, fiatSymbol string) (*domain.BalanceResult, error)
	GetBalances(requests []domain.BalanceRequest) ([]*domain.BalanceResult, error)
	GetBatchBalances(requests []domain.BalanceRequest) ([]*domain.BalanceResult, error)
}

// CryptoProvider interface for individual cryptocurrency providers.
type CryptoProvider interface {
	GetBalance(address string) (float64, error)
}
