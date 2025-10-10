package solana

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var (
	// ErrInvalidSolanaAddress indicates an invalid Solana address format.
	ErrInvalidSolanaAddress = errors.New("invalid Solana address format")
)

const (
	// Retry constants.
	MaxRetryAttempts = 3
	RetryDelay       = 2 * time.Second
	ReconnectDelay   = 5 * time.Second

	// Timeout constants.
	BalanceTimeout    = 10 * time.Second
	ConnectionTimeout = 5 * time.Second

	// Lamports to SOL conversion factor.
	LamportsPerSol = 1e9
)

type Adapter struct {
	mu        sync.RWMutex
	client    *rpc.Client
	rpcURL    string
	connected bool
}

func NewAdapter(rpcURL string) *Adapter {
	a := &Adapter{
		rpcURL: rpcURL,
	}
	a.connectWithRetry()
	return a
}

func (a *Adapter) GetBalance(address string) (float64, error) {
	pubkey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return 0, ErrInvalidSolanaAddress
	}

	var lastErr error
	for i := range MaxRetryAttempts {
		client := a.getClient()
		if client == nil {
			a.connectWithRetry()
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), BalanceTimeout)
		balance, err := client.GetBalance(ctx, pubkey, rpc.CommitmentFinalized)
		cancel()

		if err == nil {
			// Convert lamports to SOL
			solBalance := float64(balance.Value) / LamportsPerSol
			return solBalance, nil
		}

		lastErr = err
		log.Printf("[solana] balance fetch failed (attempt %d): %v", i+1, err)

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			a.connectWithRetry()
		}

		time.Sleep(RetryDelay)
	}

	return 0, lastErr
}

func (a *Adapter) connectWithRetry() {
	for {
		client := rpc.New(a.rpcURL)

		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
		_, err := client.GetVersion(ctx)
		cancel()

		if err == nil {
			log.Printf("[solana] connected to RPC %s", a.rpcURL)
			a.mu.Lock()
			a.client = client
			a.connected = true
			a.mu.Unlock()
			return
		}

		log.Printf("[solana] connection failed: %v, retrying in %v...", err, ReconnectDelay)
		a.mu.Lock()
		a.connected = false
		a.mu.Unlock()
		time.Sleep(ReconnectDelay)
	}
}

func (a *Adapter) getClient() *rpc.Client {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if !a.connected {
		return nil
	}
	return a.client
}

func (a *Adapter) Close() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.client = nil
	a.connected = false
}
