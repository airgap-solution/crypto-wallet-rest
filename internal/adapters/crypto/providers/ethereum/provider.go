package ethereum

import (
	"context"
	"errors"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	// ErrInvalidEthereumAddress indicates an invalid Ethereum address format.
	ErrInvalidEthereumAddress = errors.New("invalid Ethereum address format")
)

const (
	// Retry constants.
	MaxRetryAttempts = 3
	RetryDelay       = 2 * time.Second
	ReconnectDelay   = 5 * time.Second

	// Timeout constants.
	BalanceTimeout    = 10 * time.Second
	ConnectionTimeout = 5 * time.Second

	// Wei to ETH conversion factor.
	WeiPerEther = 1e18
)

type Adapter struct {
	mu        sync.RWMutex
	client    *ethclient.Client
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
	if !common.IsHexAddress(address) {
		return 0, ErrInvalidEthereumAddress
	}

	addr := common.HexToAddress(address)

	var lastErr error
	for i := range MaxRetryAttempts {
		client := a.getClient()
		if client == nil {
			a.connectWithRetry()
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), BalanceTimeout)
		balance, err := client.BalanceAt(ctx, addr, nil)
		cancel()

		if err == nil {
			// Convert Wei to ETH
			balanceFloat := new(big.Float)
			balanceFloat.SetString(balance.String())
			weiPerEth := new(big.Float)
			weiPerEth.SetFloat64(WeiPerEther)
			ethBalance := new(big.Float).Quo(balanceFloat, weiPerEth)

			result, _ := ethBalance.Float64()
			return result, nil
		}

		lastErr = err
		log.Printf("[ethereum] balance fetch failed (attempt %d): %v", i+1, err)

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			a.connectWithRetry()
		}

		time.Sleep(RetryDelay)
	}

	return 0, lastErr
}

func (a *Adapter) connectWithRetry() {
	for {
		client, err := ethclient.Dial(a.rpcURL)
		if err == nil {
			// Test the connection
			ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
			_, err = client.NetworkID(ctx)
			cancel()

			if err == nil {
				log.Printf("[ethereum] connected to RPC %s", a.rpcURL)
				a.mu.Lock()
				a.client = client
				a.connected = true
				a.mu.Unlock()
				return
			}
			client.Close()
		}

		log.Printf("[ethereum] connection failed: %v, retrying in %v...", err, ReconnectDelay)
		a.mu.Lock()
		a.connected = false
		a.mu.Unlock()
		time.Sleep(ReconnectDelay)
	}
}

func (a *Adapter) getClient() *ethclient.Client {
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
	if a.client != nil {
		a.client.Close()
		a.client = nil
		a.connected = false
	}
}
