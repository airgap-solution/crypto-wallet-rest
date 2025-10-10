package bitcoin

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/lamengao/go-electrum/electrum"
)

const (
	// Address derivation constants.
	DefaultExternalCount = 10
	DefaultChangeCount   = 10

	// Retry constants.
	MaxRetryAttempts = 3
	RetryDelay       = 2 * time.Second
	ReconnectDelay   = 5 * time.Second
)

type Adapter struct {
	mu              sync.RWMutex
	electrumClient  *electrum.Client
	addresses       map[string][]btcutil.Address
	electrumAddress string
}

func NewAdapter(addr string) *Adapter {
	a := &Adapter{
		addresses:       make(map[string][]btcutil.Address),
		electrumAddress: addr,
	}
	a.connectWithRetry()
	return a
}

func (a *Adapter) GetBalance(xpub string) (float64, error) {
	addresses, ok := a.addresses[xpub]
	if !ok {
		external, change, err := deriveTaprootAddresses(xpub, DefaultExternalCount, DefaultChangeCount)
		if err != nil {
			return 0, err
		}
		// Fix: assign the result of append to the same slice
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

		bal, err := getXpubBalance(client, addresses)
		if err == nil {
			return bal, nil
		}

		lastErr = err
		log.Printf("[bitcoin] balance fetch failed (attempt %d): %v", i+1, err)

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
			log.Printf("[bitcoin] connected to Electrum %s", a.electrumAddress)
			a.mu.Lock()
			a.electrumClient = client
			a.mu.Unlock()
			return
		}
		log.Printf("[bitcoin] electrum connection failed: %v, retrying in 5s...", err)
		time.Sleep(ReconnectDelay)
	}
}

func (a *Adapter) getClient() *electrum.Client {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.electrumClient
}
