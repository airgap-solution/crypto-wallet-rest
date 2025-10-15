package litecoin

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/lamengao/go-electrum/electrum"
)

var (
	ErrInvalidLitecoinAddress = errors.New("invalid Litecoin address format")
)

const (
	DefaultExternalCount = 10
	DefaultChangeCount   = 10
	MaxRetryAttempts     = 3
	RetryDelay           = 2 * time.Second
	ReconnectDelay       = 5 * time.Second
	BalanceTimeout       = 10 * time.Second
	ConnectionTimeout    = 5 * time.Second
	SatoshiPerLTC        = 1e8
)

type Adapter struct {
	mu              sync.RWMutex
	electrumClient  *electrum.Client
	addresses       map[string][]btcutil.Address
	electrumAddress string
	isTestnet       bool
}

func NewAdapter(addr string, isTestnet bool) *Adapter {
	a := &Adapter{
		addresses:       make(map[string][]btcutil.Address),
		electrumAddress: addr,
		isTestnet:       isTestnet,
	}
	a.connectWithRetry()
	return a
}

func (a *Adapter) GetBalance(xpub string) (float64, error) {
	addresses, ok := a.addresses[xpub]
	if !ok {
		external, change, err := deriveLitecoinAddresses(xpub, DefaultExternalCount, DefaultChangeCount, a.isTestnet)
		if err != nil {
			return 0, err
		}
		combined := make([]btcutil.Address, 0, len(external)+len(change))
		combined = append(combined, external...)
		combined = append(combined, change...)
		addresses = combined
		a.addresses[xpub] = addresses
	}

	var lastErr error
	for i := range MaxRetryAttempts {
		client := a.getClient()
		if client.IsShutdown() {
			a.connectWithRetry()
			continue
		}

		balance, err := getXpubBalance(client, addresses, a.isTestnet)
		if err == nil {
			return balance, nil
		}

		lastErr = err
		log.Printf("[litecoin] balance fetch failed (attempt %d): %v", i+1, err)

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			a.connectWithRetry()
		}

		time.Sleep(RetryDelay)
	}

	return 0, lastErr
}

func (a *Adapter) connectWithRetry() {
	for {
		client, err := electrum.NewClientTCP(context.Background(), a.electrumAddress)
		if err == nil {
			log.Printf("[litecoin] connected to Electrum %s", a.electrumAddress)
			a.mu.Lock()
			a.electrumClient = client
			a.mu.Unlock()
			return
		}
		log.Printf("[litecoin] electrum connection failed: %v, retrying in 5s...", err)
		time.Sleep(ReconnectDelay)
	}
}

func (a *Adapter) getClient() *electrum.Client {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.electrumClient
}
