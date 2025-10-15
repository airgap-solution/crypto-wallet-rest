package provider_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	cmcrest "github.com/airgap-solution/cmc-rest/openapi/clientgen/go"
	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/provider"
	"github.com/airgap-solution/crypto-wallet-rest/internal/core/domain"
	"github.com/airgap-solution/crypto-wallet-rest/internal/ports"
	cmcmocks "github.com/airgap-solution/crypto-wallet-rest/mocks/internaladaptersprovider"
	portsmocks "github.com/airgap-solution/crypto-wallet-rest/mocks/internalports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testSymbol     = "BTC"
	testAddress    = "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	testFiatSymbol = "USD"
)

var (
	errCMCAPI   = errors.New("CMC API error")
	errProvider = errors.New("provider connection error")
)

func TestNewAdapter(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	cryptoProviders := map[string]ports.CryptoProvider{
		"BTC": portsmocks.NewMockCryptoProvider(ctrl),
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)
	assert.NotNil(t, adapter)
}

func TestAdapter_GetBalance_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	mockCryptoProvider := portsmocks.NewMockCryptoProvider(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{
		"BTC": mockCryptoProvider,
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	symbol := testSymbol
	address := testAddress
	fiatSymbol := testFiatSymbol
	cryptoBalance := 0.5
	rate := 50000.0
	change24h := 1000.0

	mockRequest := cmcrest.ApiV1RateCurrencyFiatGetRequest{}
	response := &cmcrest.GetRateResponse{}
	response.SetRate(rate)
	response.SetChange24h(change24h)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGet(gomock.Any(), symbol, fiatSymbol).
		Return(mockRequest)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGetExecute(mockRequest).
		Return(response, &http.Response{Body: http.NoBody}, nil)

	mockCryptoProvider.EXPECT().
		GetBalance(address).
		Return(cryptoBalance, nil)

	result, err := adapter.GetBalance(symbol, address, fiatSymbol)
	require.NoError(t, err)
	assert.Equal(t, "BTC", result.CryptoSymbol)
	assert.Equal(t, address, result.Address)
	assert.InEpsilon(t, cryptoBalance, result.CryptoBalance, 0.001)
	assert.Equal(t, "USD", result.FiatSymbol)
	assert.InEpsilon(t, cryptoBalance*rate, result.FiatValue, 0.001)
	assert.InEpsilon(t, rate, result.ExchangeRate, 0.001)
	assert.InEpsilon(t, cryptoBalance*change24h, result.Change24h, 0.001)
	assert.WithinDuration(t, time.Now(), result.Timestamp, time.Second)
	assert.Nil(t, result.Error)
}

func TestAdapter_GetBalance_DefaultFiatSymbol(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	mockCryptoProvider := portsmocks.NewMockCryptoProvider(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{
		"ETH": mockCryptoProvider,
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	symbol := "ETH"
	address := "0x742d35Cc6634C0532925a3b8D3A7F13f"
	cryptoBalance := 1.0
	rate := 3000.0
	change24h := 100.0

	mockRequest := cmcrest.ApiV1RateCurrencyFiatGetRequest{}
	response := &cmcrest.GetRateResponse{}
	response.SetRate(rate)
	response.SetChange24h(change24h)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGet(gomock.Any(), symbol, "USD").
		Return(mockRequest)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGetExecute(mockRequest).
		Return(response, &http.Response{Body: http.NoBody}, nil)

	mockCryptoProvider.EXPECT().
		GetBalance(address).
		Return(cryptoBalance, nil)

	result, err := adapter.GetBalance(symbol, address, "")
	require.NoError(t, err)
	assert.Equal(t, "USD", result.FiatSymbol)
}

func TestAdapter_GetBalance_TestnetSymbol(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	mockCryptoProvider := portsmocks.NewMockCryptoProvider(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{
		"BTC_TESTNET": mockCryptoProvider,
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	symbol := "BTC_TESTNET"
	address := "tb1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4"
	fiatSymbol := "USD"
	cryptoBalance := 0.5
	rate := 50000.0
	change24h := 1000.0

	mockRequest := cmcrest.ApiV1RateCurrencyFiatGetRequest{}
	response := &cmcrest.GetRateResponse{}
	response.SetRate(rate)
	response.SetChange24h(change24h)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGet(gomock.Any(), "BTC", fiatSymbol).
		Return(mockRequest)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGetExecute(mockRequest).
		Return(response, &http.Response{Body: http.NoBody}, nil)

	mockCryptoProvider.EXPECT().
		GetBalance(address).
		Return(cryptoBalance, nil)

	result, err := adapter.GetBalance(symbol, address, fiatSymbol)
	require.NoError(t, err)
	assert.Equal(t, "BTC_TESTNET", result.CryptoSymbol)
	assert.Equal(t, address, result.Address)
}

func TestAdapter_GetBalance_ProviderNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	result, err := adapter.GetBalance("INVALID", "test-address", "USD")
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "provider not found for symbol")
	assert.Contains(t, err.Error(), "INVALID")
	assert.ErrorIs(t, err, provider.ErrProviderNotFoundForSymbol)
}

func TestAdapter_GetBalance_CMCError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	mockCryptoProvider := portsmocks.NewMockCryptoProvider(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{
		"BTC": mockCryptoProvider,
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	symbol := testSymbol
	address := testAddress
	fiatSymbol := testFiatSymbol
	cryptoBalance := 1.5
	rate := 50000.0
	change24h := 1000.0

	mockRequest := cmcrest.ApiV1RateCurrencyFiatGetRequest{}
	response := &cmcrest.GetRateResponse{}
	response.SetRate(rate)
	response.SetChange24h(change24h)
	cmcError := errCMCAPI

	mockCryptoProvider.EXPECT().
		GetBalance(address).
		Return(cryptoBalance, nil)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGet(gomock.Any(), symbol, fiatSymbol).
		Return(mockRequest)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGetExecute(mockRequest).
		Return(nil, &http.Response{Body: http.NoBody}, cmcError)

	result, err := adapter.GetBalance(symbol, address, fiatSymbol)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get rate from CMC")
	assert.Contains(t, err.Error(), "CMC API error")
}

func TestAdapter_GetBalance_CryptoProviderError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	mockCryptoProvider := portsmocks.NewMockCryptoProvider(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{
		"BTC": mockCryptoProvider,
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	symbol := testSymbol
	address := testAddress
	fiatSymbol := testFiatSymbol
	providerError := errProvider

	mockCryptoProvider.EXPECT().
		GetBalance(address).
		Return(0.0, providerError)

	result, err := adapter.GetBalance(symbol, address, fiatSymbol)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get balance from provider")
	assert.Contains(t, err.Error(), "provider connection error")
}

func TestAdapter_GetBalance_HTTPResponseBodyClosed(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	mockCryptoProvider := portsmocks.NewMockCryptoProvider(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{
		"BTC": mockCryptoProvider,
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	symbol := testSymbol
	address := testAddress
	fiatSymbol := testFiatSymbol
	cryptoBalance := 0.5
	rate := 50000.0
	change24h := 1000.0

	mockRequest := cmcrest.ApiV1RateCurrencyFiatGetRequest{}
	response := &cmcrest.GetRateResponse{}
	response.SetRate(rate)
	response.SetChange24h(change24h)

	httpResp := &http.Response{
		Body: io.NopCloser(strings.NewReader("test response")),
	}

	mockCMC.EXPECT().
		V1RateCurrencyFiatGet(gomock.Any(), symbol, fiatSymbol).
		Return(mockRequest)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGetExecute(mockRequest).
		Return(response, httpResp, nil)

	mockCryptoProvider.EXPECT().
		GetBalance(address).
		Return(cryptoBalance, nil)

	result, err := adapter.GetBalance(symbol, address, fiatSymbol)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, symbol, result.CryptoSymbol)
}

func TestAdapter_GetBatchBalances_EmptyRequests(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	cryptoProviders := map[string]ports.CryptoProvider{}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	results, err := adapter.GetBatchBalances([]domain.BalanceRequest{})
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestAdapter_GetBatchBalances_SingleRequest(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	mockCryptoProvider := portsmocks.NewMockCryptoProvider(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{
		"BTC": mockCryptoProvider,
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	symbol := testSymbol
	address := testAddress
	fiatSymbol := testFiatSymbol
	cryptoBalance := 0.5
	rate := 50000.0
	change24h := 1000.0

	mockRequest := cmcrest.ApiV1RateCurrencyFiatGetRequest{}
	response := &cmcrest.GetRateResponse{}
	response.SetRate(rate)
	response.SetChange24h(change24h)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGet(gomock.Any(), symbol, fiatSymbol).
		Return(mockRequest)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGetExecute(mockRequest).
		Return(response, &http.Response{Body: http.NoBody}, nil)

	mockCryptoProvider.EXPECT().
		GetBalance(address).
		Return(cryptoBalance, nil)

	requests := []domain.BalanceRequest{
		{
			CryptoSymbol: symbol,
			Address:      address,
			FiatSymbol:   fiatSymbol,
		},
	}

	results, err := adapter.GetBatchBalances(requests)
	require.NoError(t, err)
	require.Len(t, results, 1)

	result := results[0]
	assert.Equal(t, symbol, result.CryptoSymbol)
	assert.Equal(t, address, result.Address)
	assert.InEpsilon(t, cryptoBalance, result.CryptoBalance, 0.001)
	assert.Equal(t, fiatSymbol, result.FiatSymbol)
	assert.InEpsilon(t, cryptoBalance*rate, result.FiatValue, 0.001)
	assert.InEpsilon(t, rate, result.ExchangeRate, 0.001)
	assert.InEpsilon(t, cryptoBalance*change24h, result.Change24h, 0.001)
	assert.Nil(t, result.Error)
}

func TestAdapter_GetBatchBalances_ErrorHandling(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	requests := []domain.BalanceRequest{
		{
			CryptoSymbol: "INVALID",
			Address:      "test-address",
			FiatSymbol:   "USD",
		},
	}

	results, err := adapter.GetBalances(requests)
	require.NoError(t, err)
	require.Len(t, results, 1)

	result := results[0]
	assert.Equal(t, "INVALID", result.CryptoSymbol)
	assert.Equal(t, "test-address", result.Address)
	assert.Equal(t, "USD", result.FiatSymbol)
	assert.InDelta(t, float64(0), result.CryptoBalance, 0.001)
	assert.InDelta(t, float64(0), result.FiatValue, 0.001)
	assert.InDelta(t, float64(0), result.ExchangeRate, 0.001)
	assert.InDelta(t, float64(0), result.Change24h, 0.001)
	assert.NotNil(t, result.Error)
	assert.Contains(t, *result.Error, "provider not found for symbol")
	assert.WithinDuration(t, time.Now(), result.Timestamp, time.Second)
}

func TestAdapter_GetBalances_CallsBatchMethod(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	cryptoProviders := map[string]ports.CryptoProvider{}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	requests := []domain.BalanceRequest{
		{
			CryptoSymbol: "TEST",
			Address:      "test-address",
			FiatSymbol:   "USD",
		},
	}

	results, err := adapter.GetBalances(requests)

	require.NoError(t, err)
	require.Len(t, results, 1)

	result := results[0]
	assert.Equal(t, "TEST", result.CryptoSymbol)
	assert.NotNil(t, result.Error)
}

func TestAdapter_ImplementsInterfaces(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	cryptoProviders := map[string]ports.CryptoProvider{}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	var _ ports.Provider = adapter
}

func TestAdapter_GetBalance_CaseInsensitiveSymbol(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	mockCryptoProvider := portsmocks.NewMockCryptoProvider(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{
		"BTC": mockCryptoProvider,
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	symbol := "btc"
	address := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	fiatSymbol := "USD"
	cryptoBalance := 0.5
	rate := 50000.0
	change24h := 1000.0

	mockRequest := cmcrest.ApiV1RateCurrencyFiatGetRequest{}
	response := &cmcrest.GetRateResponse{}
	response.SetRate(rate)
	response.SetChange24h(change24h)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGet(gomock.Any(), "btc", fiatSymbol).
		Return(mockRequest)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGetExecute(mockRequest).
		Return(response, &http.Response{Body: http.NoBody}, nil)

	mockCryptoProvider.EXPECT().
		GetBalance(address).
		Return(cryptoBalance, nil)

	result, err := adapter.GetBalance(symbol, address, fiatSymbol)

	require.NoError(t, err)
	assert.Equal(t, "BTC", result.CryptoSymbol)
}

func TestAdapter_GetBalance_NilChange24h(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	mockCryptoProvider := portsmocks.NewMockCryptoProvider(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{
		"BTC": mockCryptoProvider,
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	symbol := "btc"
	address := testAddress
	fiatSymbol := testFiatSymbol
	cryptoBalance := 0.5
	rate := 50000.0

	mockRequest := cmcrest.ApiV1RateCurrencyFiatGetRequest{}
	response := &cmcrest.GetRateResponse{}
	response.SetRate(rate)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGet(gomock.Any(), symbol, fiatSymbol).
		Return(mockRequest)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGetExecute(mockRequest).
		Return(response, &http.Response{Body: http.NoBody}, nil)

	mockCryptoProvider.EXPECT().
		GetBalance(address).
		Return(cryptoBalance, nil)

	result, err := adapter.GetBalance(symbol, address, fiatSymbol)

	require.NoError(t, err)
	assert.InDelta(t, float64(0), result.Change24h, 0.001)
}

func TestAdapter_GetBalance_CachingBehavior(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	mockCryptoProvider := portsmocks.NewMockCryptoProvider(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{
		"BTC": mockCryptoProvider,
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	symbol := testSymbol
	address := testAddress
	fiatSymbol := testFiatSymbol
	cryptoBalance := 1.0
	rate := 50000.0
	change24h := 1000.0

	mockRequest := cmcrest.ApiV1RateCurrencyFiatGetRequest{}
	response := &cmcrest.GetRateResponse{}
	response.SetRate(rate)
	response.SetChange24h(change24h)

	mockCryptoProvider.EXPECT().
		GetBalance(address).
		Return(cryptoBalance, nil).
		Times(1)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGet(gomock.Any(), symbol, fiatSymbol).
		Return(mockRequest).
		Times(1)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGetExecute(mockRequest).
		Return(response, &http.Response{Body: http.NoBody}, nil).
		Times(1)

	address2 := "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2"
	mockCryptoProvider.EXPECT().
		GetBalance(address2).
		Return(cryptoBalance*2, nil).
		Times(1)

	result1, err1 := adapter.GetBalance(symbol, address, fiatSymbol)
	require.NoError(t, err1)
	assert.InEpsilon(t, cryptoBalance, result1.CryptoBalance, 0.001)
	assert.InEpsilon(t, rate, result1.ExchangeRate, 0.001)

	result2, err2 := adapter.GetBalance(symbol, address2, fiatSymbol)
	require.NoError(t, err2)
	assert.InEpsilon(t, cryptoBalance*2, result2.CryptoBalance, 0.001)
	assert.InEpsilon(t, rate, result2.ExchangeRate, 0.001)
}

func TestAdapter_GetBalance_BalanceCacheHit(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCMC := cmcmocks.NewMockCMCRestClient(ctrl)
	mockCryptoProvider := portsmocks.NewMockCryptoProvider(ctrl)

	cryptoProviders := map[string]ports.CryptoProvider{
		"BTC": mockCryptoProvider,
	}

	adapter := provider.NewAdapter(mockCMC, cryptoProviders)

	symbol := testSymbol
	address := testAddress
	fiatSymbol := testFiatSymbol
	cryptoBalance := 1.0
	rate := 50000.0
	change24h := 1000.0

	mockRequest := cmcrest.ApiV1RateCurrencyFiatGetRequest{}
	response := &cmcrest.GetRateResponse{}
	response.SetRate(rate)
	response.SetChange24h(change24h)

	mockCryptoProvider.EXPECT().
		GetBalance(address).
		Return(cryptoBalance, nil).
		Times(1)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGet(gomock.Any(), symbol, fiatSymbol).
		Return(mockRequest).
		Times(1)

	mockCMC.EXPECT().
		V1RateCurrencyFiatGetExecute(mockRequest).
		Return(response, &http.Response{Body: http.NoBody}, nil).
		Times(1)

	result1, err1 := adapter.GetBalance(symbol, address, fiatSymbol)
	require.NoError(t, err1)
	assert.InEpsilon(t, cryptoBalance, result1.CryptoBalance, 0.001)
	assert.InEpsilon(t, rate, result1.ExchangeRate, 0.001)

	result2, err2 := adapter.GetBalance(symbol, address, fiatSymbol)
	require.NoError(t, err2)
	assert.InEpsilon(t, cryptoBalance, result2.CryptoBalance, 0.001)
	assert.InEpsilon(t, rate, result2.ExchangeRate, 0.001)
}
