package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	cmcrest "github.com/airgap-solution/cmc-rest/openapi/clientgen/go"
	"github.com/airgap-solution/crypto-wallet-rest/internal/ports"
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

func (a *Adapter) GetBalance(symbol, addr string) (float64, error) {
	req := a.cmcRest.V1RateCurrencyFiatGet(context.Background(), symbol, "CAD")
	resp, httpResp, err := a.cmcRest.V1RateCurrencyFiatGetExecute(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get rate from CMC: %w", err)
	}
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	rate := resp.GetRate()

	prov, ok := a.cryptoProviders[strings.ToUpper(symbol)]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrProviderNotFoundForSymbol, symbol)
	}

	balance, err := prov.GetBalance(addr)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance from provider: %w", err)
	}

	return balance * rate, nil
}
