package ports

type Provider interface {
	GetBalance(symbol, address string) (float64, error)
}

type CryptoProvider interface {
	GetBalance(address string) (float64, error)
}
