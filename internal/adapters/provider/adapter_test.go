package provider_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	cmcrest "github.com/airgap-solution/cmc-rest/openapi/clientgen/go"
	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/provider"
	"github.com/airgap-solution/crypto-wallet-rest/internal/ports"
	internaladaptersprovidermocks "github.com/airgap-solution/crypto-wallet-rest/mocks/internaladaptersprovider"
	internalportsmocks "github.com/airgap-solution/crypto-wallet-rest/mocks/internalports"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewAdapter(t *testing.T) {
	t.Parallel()
	adapter := provider.NewAdapter(nil, nil)
	require.NotNil(t, adapter)
}

func TestGetBalance(t *testing.T) {
	t.Parallel()

	rate := 45000.50
	cryptoBalance := 0.00123456
	fiatValue := rate * cryptoBalance
	addr := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	cryptoSymbol := "BTC"
	fiatSymbol := "USD"
	req := cmcrest.ApiV1RateCurrencyFiatGetRequest{}

	testCases := []struct {
		name                     string
		expectedError            error
		expectedBalanceResult    *ports.BalanceResult
		setupMockCryptoProviders func(map[string]ports.CryptoProvider)
		setupMocks               func(
			*internaladaptersprovidermocks.MockCMCRestClient,
			*internalportsmocks.MockCryptoProvider,
		)
	}{
		{
			name:          "successful request returns no error",
			expectedError: nil,
			setupMocks: func(
				mockCMC *internaladaptersprovidermocks.MockCMCRestClient,
				mockBTC *internalportsmocks.MockCryptoProvider,
			) {
				mockCMC.EXPECT().V1RateCurrencyFiatGet(
					gomock.Any(), cryptoSymbol, fiatSymbol).Return(req)
				mockCMC.EXPECT().V1RateCurrencyFiatGetExecute(req).Return(
					&cmcrest.GetRateResponse{Rate: lo.ToPtr(rate)}, nil, nil)
				mockBTC.EXPECT().GetBalance(addr).Return(cryptoBalance, nil)
			},
			expectedBalanceResult: &ports.BalanceResult{
				CryptoSymbol:  "BTC",
				Address:       addr,
				CryptoBalance: cryptoBalance,
				FiatSymbol:    "USD",
				FiatValue:     fiatValue,
				ExchangeRate:  rate,
			},
		},
		{
			name:          "cmc client request fails",
			expectedError: assert.AnError,
			setupMocks: func(
				mockCMC *internaladaptersprovidermocks.MockCMCRestClient,
				_ *internalportsmocks.MockCryptoProvider,
			) {
				mockCMC.EXPECT().V1RateCurrencyFiatGet(
					gomock.Any(), cryptoSymbol, fiatSymbol).Return(req)
				mockCMC.EXPECT().V1RateCurrencyFiatGetExecute(req).Return(
					nil, nil, assert.AnError)
			},
		},
		{
			name:          "no provider found for symbol",
			expectedError: provider.ErrProviderNotFoundForSymbol,
			setupMocks: func(
				mockCMC *internaladaptersprovidermocks.MockCMCRestClient,
				_ *internalportsmocks.MockCryptoProvider,
			) {
				mockCMC.EXPECT().V1RateCurrencyFiatGet(
					gomock.Any(), cryptoSymbol, fiatSymbol).Return(req)
				mockCMC.EXPECT().V1RateCurrencyFiatGetExecute(req).Return(
					&cmcrest.GetRateResponse{Rate: lo.ToPtr(rate)}, nil, nil)
			},
			setupMockCryptoProviders: func(providers map[string]ports.CryptoProvider) {
				delete(providers, cryptoSymbol)
			},
		},
		{
			name:          "btc provider request fails",
			expectedError: assert.AnError,
			setupMocks: func(
				mockCMC *internaladaptersprovidermocks.MockCMCRestClient,
				mockBTC *internalportsmocks.MockCryptoProvider,
			) {
				mockCMC.EXPECT().V1RateCurrencyFiatGet(
					gomock.Any(), cryptoSymbol, fiatSymbol).Return(req)
				mockCMC.EXPECT().V1RateCurrencyFiatGetExecute(req).Return(
					&cmcrest.GetRateResponse{Rate: lo.ToPtr(rate)}, nil, nil)
				mockBTC.EXPECT().GetBalance(addr).Return(0.0, assert.AnError)
			},
		},
		{
			name:          "successful request with http response body",
			expectedError: nil,
			setupMocks: func(
				mockCMC *internaladaptersprovidermocks.MockCMCRestClient,
				mockBTC *internalportsmocks.MockCryptoProvider,
			) {
				mockCMC.EXPECT().V1RateCurrencyFiatGet(
					gomock.Any(), cryptoSymbol, fiatSymbol).Return(req)
				httpResp := &http.Response{
					Body: io.NopCloser(strings.NewReader("test response")),
				}
				mockCMC.EXPECT().V1RateCurrencyFiatGetExecute(req).Return(
					&cmcrest.GetRateResponse{Rate: lo.ToPtr(rate)}, httpResp, nil)
				mockBTC.EXPECT().GetBalance(addr).Return(cryptoBalance, nil)
			},
			expectedBalanceResult: &ports.BalanceResult{
				CryptoSymbol:  "BTC",
				Address:       addr,
				CryptoBalance: cryptoBalance,
				FiatSymbol:    "USD",
				FiatValue:     fiatValue,
				ExchangeRate:  rate,
			},
		},
		{
			name:          "successful request with nil http response",
			expectedError: nil,
			setupMocks: func(
				mockCMC *internaladaptersprovidermocks.MockCMCRestClient,
				mockBTC *internalportsmocks.MockCryptoProvider,
			) {
				mockCMC.EXPECT().V1RateCurrencyFiatGet(
					gomock.Any(), cryptoSymbol, fiatSymbol).Return(req)
				mockCMC.EXPECT().V1RateCurrencyFiatGetExecute(req).Return(
					&cmcrest.GetRateResponse{Rate: lo.ToPtr(rate)}, nil, nil)
				mockBTC.EXPECT().GetBalance(addr).Return(cryptoBalance, nil)
			},
			expectedBalanceResult: &ports.BalanceResult{
				CryptoSymbol:  "BTC",
				Address:       addr,
				CryptoBalance: cryptoBalance,
				FiatSymbol:    "USD",
				FiatValue:     fiatValue,
				ExchangeRate:  rate,
			},
		},
		{
			name:          "empty fiat symbol defaults to USD",
			expectedError: nil,
			setupMocks: func(
				mockCMC *internaladaptersprovidermocks.MockCMCRestClient,
				mockBTC *internalportsmocks.MockCryptoProvider,
			) {
				mockCMC.EXPECT().V1RateCurrencyFiatGet(
					gomock.Any(), cryptoSymbol, "USD").Return(req)
				mockCMC.EXPECT().V1RateCurrencyFiatGetExecute(req).Return(
					&cmcrest.GetRateResponse{Rate: lo.ToPtr(rate)}, nil, nil)
				mockBTC.EXPECT().GetBalance(addr).Return(cryptoBalance, nil)
			},
			expectedBalanceResult: &ports.BalanceResult{
				CryptoSymbol:  "BTC",
				Address:       addr,
				CryptoBalance: cryptoBalance,
				FiatSymbol:    "USD",
				FiatValue:     fiatValue,
				ExchangeRate:  rate,
			},
		},
		{
			name:          "different fiat currency (EUR)",
			expectedError: nil,
			setupMocks: func(
				mockCMC *internaladaptersprovidermocks.MockCMCRestClient,
				mockBTC *internalportsmocks.MockCryptoProvider,
			) {
				eurRate := 42000.75
				mockCMC.EXPECT().V1RateCurrencyFiatGet(
					gomock.Any(), cryptoSymbol, "EUR").Return(req)
				mockCMC.EXPECT().V1RateCurrencyFiatGetExecute(req).Return(
					&cmcrest.GetRateResponse{Rate: lo.ToPtr(eurRate)}, nil, nil)
				mockBTC.EXPECT().GetBalance(addr).Return(cryptoBalance, nil)
			},
			expectedBalanceResult: &ports.BalanceResult{
				CryptoSymbol:  "BTC",
				Address:       addr,
				CryptoBalance: cryptoBalance,
				FiatSymbol:    "EUR",
				FiatValue:     42000.75 * cryptoBalance,
				ExchangeRate:  42000.75,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockCMCRestClient := internaladaptersprovidermocks.NewMockCMCRestClient(ctrl)
			mockBTCProvider := internalportsmocks.NewMockCryptoProvider(ctrl)
			tc.setupMocks(mockCMCRestClient, mockBTCProvider)

			cryptoProviders := map[string]ports.CryptoProvider{
				cryptoSymbol: mockBTCProvider,
			}
			if tc.setupMockCryptoProviders != nil {
				tc.setupMockCryptoProviders(cryptoProviders)
			}

			adapter := provider.NewAdapter(mockCMCRestClient, cryptoProviders)

			// Use appropriate fiat symbol for test case
			testFiatSymbol := fiatSymbol
			if tc.name == "empty fiat symbol defaults to USD" {
				testFiatSymbol = ""
			} else if tc.name == "different fiat currency (EUR)" {
				testFiatSymbol = "EUR"
			}

			balanceResult, err := adapter.GetBalance(cryptoSymbol, addr, testFiatSymbol)

			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
				require.Nil(t, balanceResult)
			} else {
				require.NoError(t, err)
				require.NotNil(t, balanceResult)

				// Verify the structure and values
				assert.Equal(t, tc.expectedBalanceResult.CryptoSymbol, balanceResult.CryptoSymbol)
				assert.Equal(t, tc.expectedBalanceResult.Address, balanceResult.Address)
				assert.InEpsilon(t, tc.expectedBalanceResult.CryptoBalance, balanceResult.CryptoBalance, 1e-9)
				assert.Equal(t, tc.expectedBalanceResult.FiatSymbol, balanceResult.FiatSymbol)
				assert.InEpsilon(t, tc.expectedBalanceResult.FiatValue, balanceResult.FiatValue, 1e-9)
				assert.InEpsilon(t, tc.expectedBalanceResult.ExchangeRate, balanceResult.ExchangeRate, 1e-9)

				// Verify timestamp is recent (within last 5 seconds)
				assert.WithinDuration(t, time.Now(), balanceResult.Timestamp, 5*time.Second)
			}
		})
	}
}
