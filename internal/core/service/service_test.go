package service_test

import (
	"testing"
	"time"

	"github.com/airgap-solution/crypto-wallet-rest/internal/core/service"
	"github.com/airgap-solution/crypto-wallet-rest/internal/ports"
	internalportsmocks "github.com/airgap-solution/crypto-wallet-rest/mocks/internalports"
	cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProvider := internalportsmocks.NewMockProvider(ctrl)
	svc := service.New(mockProvider)
	assert.NotNil(t, svc)
}

func TestBalanceGet(t *testing.T) {
	t.Parallel()
	cryptoSymbol := "BTC"
	address := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	fiatSymbol := "USD"
	cryptoBalance := 0.00123456
	fiatValue := 55.67
	exchangeRate := 45000.50
	timestamp := time.Now()

	for _, tc := range []struct {
		name                 string
		setupMocks           func(*internalportsmocks.MockProvider)
		expectedError        error
		expectedResponseBody any
		expectedResponseCode int
	}{
		{
			name: "provider returns no error",
			setupMocks: func(mockProvider *internalportsmocks.MockProvider) {
				balanceResult := &ports.BalanceResult{
					CryptoSymbol:  cryptoSymbol,
					Address:       address,
					CryptoBalance: cryptoBalance,
					FiatSymbol:    fiatSymbol,
					FiatValue:     fiatValue,
					ExchangeRate:  exchangeRate,
					Timestamp:     timestamp,
				}
				mockProvider.EXPECT().GetBalance(cryptoSymbol, address, fiatSymbol).Return(balanceResult, nil)
			},
			expectedResponseBody: cryptowalletrest.BalanceGet200Response{
				CryptoSymbol:  cryptoSymbol,
				Address:       address,
				CryptoBalance: cryptoBalance,
				FiatSymbol:    fiatSymbol,
				FiatValue:     fiatValue,
				ExchangeRate:  exchangeRate,
				Timestamp:     timestamp,
			},
			expectedResponseCode: 200,
			expectedError:        nil,
		},
		{
			name: "provider returns error",
			setupMocks: func(mockProvider *internalportsmocks.MockProvider) {
				mockProvider.EXPECT().GetBalance(cryptoSymbol, address, fiatSymbol).Return(nil, assert.AnError)
			},
			expectedResponseBody: service.Error{Message: assert.AnError.Error()},
			expectedResponseCode: 501, // StatusNotImplemented
			expectedError:        nil,
		},
		{
			name: "EUR fiat currency",
			setupMocks: func(mockProvider *internalportsmocks.MockProvider) {
				balanceResult := &ports.BalanceResult{
					CryptoSymbol:  cryptoSymbol,
					Address:       address,
					CryptoBalance: cryptoBalance,
					FiatSymbol:    "EUR",
					FiatValue:     52.34,
					ExchangeRate:  42400.25,
					Timestamp:     timestamp,
				}
				mockProvider.EXPECT().GetBalance(cryptoSymbol, address, "EUR").Return(balanceResult, nil)
			},
			expectedResponseBody: cryptowalletrest.BalanceGet200Response{
				CryptoSymbol:  cryptoSymbol,
				Address:       address,
				CryptoBalance: cryptoBalance,
				FiatSymbol:    "EUR",
				FiatValue:     52.34,
				ExchangeRate:  42400.25,
				Timestamp:     timestamp,
			},
			expectedResponseCode: 200,
			expectedError:        nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProvider := internalportsmocks.NewMockProvider(ctrl)
			tc.setupMocks(mockProvider)

			svc := service.New(mockProvider)

			testFiatSymbol := fiatSymbol
			if tc.name == "EUR fiat currency" {
				testFiatSymbol = "EUR"
			}

			resp, err := svc.BalanceGet(t.Context(), cryptoSymbol, address, testFiatSymbol)

			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedResponseCode, resp.Code)
			require.Equal(t, tc.expectedResponseBody, resp.Body)
		})
	}
}

func TestTransactionsGet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)
	svc := service.New(mockProvider)

	resp, err := svc.TransactionsGet(t.Context(), "BTC", "address", 50, 0)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 501, resp.Code) // StatusNotImplemented
}

func TestUnsignedTxGet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)
	svc := service.New(mockProvider)

	resp, err := svc.UnsignedTxGet(t.Context(), "BTC", "from", "to", "0.001", 10.5)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 501, resp.Code) // StatusNotImplemented
}

func TestBroadcastPost(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)
	svc := service.New(mockProvider)

	req := cryptowalletrest.BroadcastPostRequest{
		CryptoSymbol: "BTC",
		SignedTx:     "0200000001f5d8ee39a430...",
	}
	resp, err := svc.BroadcastPost(t.Context(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 501, resp.Code) // StatusNotImplemented
}
