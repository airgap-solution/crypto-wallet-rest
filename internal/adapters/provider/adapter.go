package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	cmcrest "github.com/airgap-solution/cmc-rest/openapi/clientgen/go"
	"github.com/airgap-solution/crypto-wallet-rest/internal/ports"
	"github.com/samber/lo"
)

var ErrProviderNotFoundForSymbol = errors.New("provider not found for symbol")

type CMCRestClient interface {
	V1RateCurrencyFiatGet(ctx context.Context, from, to string) cmcrest.ApiV1RateCurrencyFiatGetRequest
	V1RateCurrencyFiatGetExecute(
		r cmcrest.ApiV1RateCurrencyFiatGetRequest,
	) (*cmcrest.GetRateResponse, *http.Response, error)
}

type Adapter struct {
	cmcRest         CMCRestClient
	cryptoProviders map[string]ports.CryptoProvider
}

func NewAdapter(cmcRest CMCRestClient, cryptoProviders map[string]ports.CryptoProvider) *Adapter {
	return &Adapter{
		cmcRest:         cmcRest,
		cryptoProviders: cryptoProviders,
	}
}

func (a *Adapter) GetBalance(symbol, addr, fiatSymbol string) (*ports.BalanceResult, error) {
	if fiatSymbol == "" {
		fiatSymbol = "USD"
	}

	// Strip _TESTNET suffix for fiat rate fetching
	rateSymbol := strings.TrimSuffix(symbol, "_TESTNET")

	req := a.cmcRest.V1RateCurrencyFiatGet(context.Background(), rateSymbol, fiatSymbol)
	resp, httpResp, err := a.cmcRest.V1RateCurrencyFiatGetExecute(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate from CMC: %w", err)
	}
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	prov, ok := a.cryptoProviders[strings.ToUpper(symbol)]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFoundForSymbol, symbol)
	}

	cryptoBalance, err := prov.GetBalance(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance from provider: %w", err)
	}

	rate := resp.GetRate()
	fiatValue := cryptoBalance * rate
	return &ports.BalanceResult{
		CryptoSymbol:  strings.ToUpper(symbol),
		Address:       addr,
		CryptoBalance: cryptoBalance,
		FiatSymbol:    strings.ToUpper(fiatSymbol),
		FiatValue:     fiatValue,
		ExchangeRate:  rate,
		Timestamp:     time.Now(),
		Change24h:     cryptoBalance * lo.FromPtr(resp.Change24h),
	}, nil
}
