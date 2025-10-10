package provider_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

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

	rate := 0.00000601
	balance := 5.12345678
	value := rate * float64(balance)
	addr := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	cryptoSymbol := "BTC"
	fiatSymbol := "CAD"
	req := cmcrest.ApiV1RateCurrencyFiatGetRequest{}

	testCases := []struct {
		name                     string
		expectedError            error
		expectedBalance          float64
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
				mockBTC.EXPECT().GetBalance(addr).Return(balance, nil)
			},
			expectedBalance: value,
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
				mockBTC.EXPECT().GetBalance(addr).Return(balance, nil)
			},
			expectedBalance: value,
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
				mockBTC.EXPECT().GetBalance(addr).Return(balance, nil)
			},
			expectedBalance: value,
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

			balance, err := adapter.GetBalance(cryptoSymbol, addr)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
				require.InDelta(t, 0.0, balance, 1e-9)
			} else {
				require.NoError(t, err)
				require.InEpsilon(t, tc.expectedBalance, balance, 1e-9)
			}
		})
	}
}
