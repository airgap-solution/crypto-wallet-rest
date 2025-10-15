package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	cmcrest "github.com/airgap-solution/cmc-rest/openapi/clientgen/go"
	"github.com/airgap-solution/crypto-wallet-rest/internal/core/domain"
	"github.com/airgap-solution/crypto-wallet-rest/internal/ports"
)

var ErrProviderNotFoundForSymbol = errors.New("provider not found for symbol")

const (
	RateCacheTTL    = 5 * time.Second
	BalanceCacheTTL = 30 * time.Second
)

type CachedRateResult struct {
	Rate      float64
	Change24h float64
}

type CMCRestClient interface {
	V1RateCurrencyFiatGet(ctx context.Context, from, to string) cmcrest.ApiV1RateCurrencyFiatGetRequest
	V1RateCurrencyFiatGetExecute(
		r cmcrest.ApiV1RateCurrencyFiatGetRequest,
	) (*cmcrest.GetRateResponse, *http.Response, error)
}

type Adapter struct {
	cmcRest         CMCRestClient
	cryptoProviders map[string]ports.CryptoProvider
	rateCache       *Cache[*CachedRateResult]
	balanceCache    *Cache[float64]
}

func NewAdapter(cmcRest CMCRestClient, cryptoProviders map[string]ports.CryptoProvider) *Adapter {
	return &Adapter{
		cmcRest:         cmcRest,
		cryptoProviders: cryptoProviders,
		rateCache:       NewCache[*CachedRateResult](),
		balanceCache:    NewCache[float64](),
	}
}

func (a *Adapter) GetBalance(symbol, addr, fiatSymbol string) (*domain.BalanceResult, error) {
	if fiatSymbol == "" {
		fiatSymbol = "USD"
	}

	prov, ok := a.cryptoProviders[strings.ToUpper(symbol)]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFoundForSymbol, symbol)
	}

	cryptoBalance, err := a.getCachedOrFetchBalance(prov, symbol, addr)
	if err != nil {
		return nil, err
	}

	rate, change24h, err := a.getCachedOrFetchRate(symbol, fiatSymbol)
	if err != nil {
		return nil, err
	}

	return a.buildBalanceResult(symbol, addr, fiatSymbol, cryptoBalance, rate, change24h), nil
}

func (a *Adapter) getCachedOrFetchBalance(prov ports.CryptoProvider, symbol, addr string) (float64, error) {
	balanceKey := fmt.Sprintf("balance:%s:%s", strings.ToUpper(symbol), addr)

	if cachedBalance, found := a.balanceCache.Get(balanceKey); found {
		return cachedBalance, nil
	}

	balance, err := prov.GetBalance(addr)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance from provider: %w", err)
	}

	a.balanceCache.Set(balanceKey, balance, BalanceCacheTTL)
	return balance, nil
}

func (a *Adapter) getCachedOrFetchRate(symbol, fiatSymbol string) (float64, float64, error) {
	rateSymbol := strings.TrimSuffix(symbol, "_TESTNET")
	rateKey := fmt.Sprintf("rate:%s:%s", strings.ToUpper(rateSymbol), strings.ToUpper(fiatSymbol))

	if cachedRate, found := a.rateCache.Get(rateKey); found {
		return cachedRate.Rate, cachedRate.Change24h, nil
	}

	req := a.cmcRest.V1RateCurrencyFiatGet(context.Background(), rateSymbol, fiatSymbol)
	resp, httpResp, err := a.cmcRest.V1RateCurrencyFiatGetExecute(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get rate from CMC: %w", err)
	}
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	rate := resp.GetRate()
	var change24h float64
	if resp.Change24h != nil {
		change24h = *resp.Change24h
	}

	rateResult := &CachedRateResult{
		Rate:      rate,
		Change24h: change24h,
	}
	a.rateCache.Set(rateKey, rateResult, RateCacheTTL)

	return rate, change24h, nil
}

func (a *Adapter) buildBalanceResult(
	symbol, addr, fiatSymbol string, cryptoBalance, rate, change24h float64,
) *domain.BalanceResult {
	return &domain.BalanceResult{
		CryptoSymbol:  strings.ToUpper(symbol),
		Address:       addr,
		CryptoBalance: cryptoBalance,
		FiatSymbol:    strings.ToUpper(fiatSymbol),
		FiatValue:     cryptoBalance * rate,
		ExchangeRate:  rate,
		Timestamp:     time.Now(),
		Change24h:     cryptoBalance * change24h,
	}
}

func (a *Adapter) GetBalances(requests []domain.BalanceRequest) ([]*domain.BalanceResult, error) {
	return a.GetBatchBalances(requests)
}

func (a *Adapter) GetBatchBalances(requests []domain.BalanceRequest) ([]*domain.BalanceResult, error) {
	results := make([]*domain.BalanceResult, len(requests))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, req := range requests {
		wg.Add(1)
		go func(index int, request domain.BalanceRequest) {
			defer wg.Done()

			result, err := a.GetBalance(request.CryptoSymbol, request.Address, request.FiatSymbol)

			mu.Lock()
			if err != nil {
				errorMsg := err.Error()
				results[index] = &domain.BalanceResult{
					CryptoSymbol:  strings.ToUpper(request.CryptoSymbol),
					Address:       request.Address,
					CryptoBalance: 0,
					FiatSymbol:    strings.ToUpper(request.FiatSymbol),
					FiatValue:     0,
					ExchangeRate:  0,
					Timestamp:     time.Now(),
					Change24h:     0,
					Error:         &errorMsg,
				}
			} else {
				results[index] = result
			}
			mu.Unlock()
		}(i, req)
	}

	wg.Wait()
	return results, nil
}
