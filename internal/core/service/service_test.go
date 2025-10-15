package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/airgap-solution/crypto-wallet-rest/internal/core/domain"
	"github.com/airgap-solution/crypto-wallet-rest/internal/core/service"
	internalportsmocks "github.com/airgap-solution/crypto-wallet-rest/mocks/internalports"
	cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	errProviderGeneric = errors.New("provider error")
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
				balanceResult := &domain.BalanceResult{
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
		},
		{
			name: "provider returns error",
			setupMocks: func(mockProvider *internalportsmocks.MockProvider) {
				mockProvider.EXPECT().GetBalance(cryptoSymbol, address, fiatSymbol).Return(nil, errProviderGeneric)
			},
			expectedResponseBody: service.Error{Message: "provider error"},
			expectedResponseCode: 501,
		},
		{
			name: "EUR fiat currency",
			setupMocks: func(mockProvider *internalportsmocks.MockProvider) {
				balanceResult := &domain.BalanceResult{
					CryptoSymbol:  cryptoSymbol,
					Address:       address,
					CryptoBalance: cryptoBalance,
					FiatSymbol:    "EUR",
					FiatValue:     fiatValue,
					ExchangeRate:  exchangeRate,
					Timestamp:     timestamp,
				}
				mockProvider.EXPECT().GetBalance(cryptoSymbol, address, "EUR").Return(balanceResult, nil)
			},
			expectedResponseBody: cryptowalletrest.BalanceGet200Response{
				CryptoSymbol:  cryptoSymbol,
				Address:       address,
				CryptoBalance: cryptoBalance,
				FiatSymbol:    "EUR",
				FiatValue:     fiatValue,
				ExchangeRate:  exchangeRate,
				Timestamp:     timestamp,
			},
			expectedResponseCode: 200,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProvider := internalportsmocks.NewMockProvider(ctrl)
			tc.setupMocks(mockProvider)

			svc := service.New(mockProvider)

			fiat := fiatSymbol
			if tc.name == "EUR fiat currency" {
				fiat = "EUR"
			}

			response, err := svc.BalanceGet(t.Context(), cryptoSymbol, address, fiat)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponseCode, response.Code)
			assert.Equal(t, tc.expectedResponseBody, response.Body)
		})
	}
}

func TestService_BalancesPost_BasicFunctionality(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)

	btcResult := &domain.BalanceResult{
		CryptoSymbol:  "BTC",
		Address:       "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		CryptoBalance: 0.001,
		FiatSymbol:    "USD",
		FiatValue:     50.0,
		ExchangeRate:  50000.0,
		Change24h:     1.0,
		Timestamp:     time.Now(),
		Error:         nil,
	}

	ethResult := &domain.BalanceResult{
		CryptoSymbol:  "ETH",
		Address:       "0x742d35Cc6634C0532925a3b8D3A7F13f",
		CryptoBalance: 1.5,
		FiatSymbol:    "EUR",
		FiatValue:     3000.0,
		ExchangeRate:  2000.0,
		Change24h:     -30.0,
		Timestamp:     time.Now(),
		Error:         nil,
	}

	expectedRequests := []domain.BalanceRequest{
		{
			CryptoSymbol: "BTC",
			Address:      "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
			FiatSymbol:   "USD",
		},
		{
			CryptoSymbol: "ETH",
			Address:      "0x742d35Cc6634C0532925a3b8D3A7F13f",
			FiatSymbol:   "EUR",
		},
	}

	mockProvider.EXPECT().GetBatchBalances(expectedRequests).Return([]*domain.BalanceResult{btcResult, ethResult}, nil)

	svc := service.New(mockProvider)

	request := cryptowalletrest.BalancesPostRequest{
		Requests: []cryptowalletrest.BalancesPostRequestRequestsInner{
			{
				CryptoSymbol: "BTC",
				Address:      "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				FiatSymbol:   "USD",
			},
			{
				CryptoSymbol: "ETH",
				Address:      "0x742d35Cc6634C0532925a3b8D3A7F13f",
				FiatSymbol:   "EUR",
			},
		},
	}

	response, err := svc.BalancesPost(t.Context(), request)

	require.NoError(t, err)
	assert.Equal(t, 200, response.Code)

	responseBody, ok := response.Body.(cryptowalletrest.BalancesPost200Response)
	require.True(t, ok)
	require.Len(t, responseBody.Results, 2)

	btcBalance := responseBody.Results[0]
	assert.Equal(t, "BTC", btcBalance.CryptoSymbol)
	assert.Equal(t, "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", btcBalance.Address)
	assert.InEpsilon(t, 0.001, btcBalance.CryptoBalance, 0.001)
	assert.Equal(t, "USD", btcBalance.FiatSymbol)
	assert.InEpsilon(t, 50.0, btcBalance.FiatValue, 0.001)
	assert.Empty(t, btcBalance.Error)

	ethBalance := responseBody.Results[1]
	assert.Equal(t, "ETH", ethBalance.CryptoSymbol)
	assert.Equal(t, "0x742d35Cc6634C0532925a3b8D3A7F13f", ethBalance.Address)
	assert.InEpsilon(t, 1.5, ethBalance.CryptoBalance, 0.001)
	assert.Equal(t, "EUR", ethBalance.FiatSymbol)
	assert.InEpsilon(t, 3000.0, ethBalance.FiatValue, 0.001)
	assert.Empty(t, ethBalance.Error)
}

func TestService_BalancesPost_ErrorHandling(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)

	btcResult := &domain.BalanceResult{
		CryptoSymbol:  "BTC",
		Address:       "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		CryptoBalance: 0.001,
		FiatSymbol:    "USD",
		FiatValue:     50.0,
		ExchangeRate:  50000.0,
		Change24h:     1.0,
		Timestamp:     time.Now(),
		Error:         nil,
	}

	errorMsg := "provider not found for symbol"
	ethResult := &domain.BalanceResult{
		CryptoSymbol:  "ETH",
		Address:       "0x742d35Cc6634C0532925a3b8D3A7F13f",
		CryptoBalance: 0,
		FiatSymbol:    "USD",
		FiatValue:     0,
		ExchangeRate:  0,
		Change24h:     0,
		Timestamp:     time.Now(),
		Error:         &errorMsg,
	}

	expectedRequests := []domain.BalanceRequest{
		{
			CryptoSymbol: "BTC",
			Address:      "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
			FiatSymbol:   "USD",
		},
		{
			CryptoSymbol: "ETH",
			Address:      "0x742d35Cc6634C0532925a3b8D3A7F13f",
			FiatSymbol:   "USD",
		},
	}

	mockProvider.EXPECT().GetBatchBalances(expectedRequests).Return([]*domain.BalanceResult{btcResult, ethResult}, nil)

	svc := service.New(mockProvider)

	request := cryptowalletrest.BalancesPostRequest{
		Requests: []cryptowalletrest.BalancesPostRequestRequestsInner{
			{
				CryptoSymbol: "BTC",
				Address:      "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				FiatSymbol:   "USD",
			},
			{
				CryptoSymbol: "ETH",
				Address:      "0x742d35Cc6634C0532925a3b8D3A7F13f",
				FiatSymbol:   "USD",
			},
		},
	}

	response, err := svc.BalancesPost(t.Context(), request)

	require.NoError(t, err)
	assert.Equal(t, 200, response.Code)

	responseBody, ok := response.Body.(cryptowalletrest.BalancesPost200Response)
	require.True(t, ok)
	require.Len(t, responseBody.Results, 2)

	btcBalance := responseBody.Results[0]
	assert.Equal(t, "BTC", btcBalance.CryptoSymbol)
	assert.Empty(t, btcBalance.Error)

	ethBalance := responseBody.Results[1]
	assert.Equal(t, "ETH", ethBalance.CryptoSymbol)
	assert.Equal(t, "provider not found for symbol", ethBalance.Error)
}

func TestService_BalancesPost_DefaultFiatSymbol(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)

	btcResult := &domain.BalanceResult{
		CryptoSymbol:  "BTC",
		Address:       "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		CryptoBalance: 0.001,
		FiatSymbol:    "USD",
		FiatValue:     50.0,
		ExchangeRate:  50000.0,
		Change24h:     1.0,
		Timestamp:     time.Now(),
		Error:         nil,
	}

	// Request with no fiat symbol should default to USD
	expectedRequests := []domain.BalanceRequest{
		{
			CryptoSymbol: "BTC",
			Address:      "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
			FiatSymbol:   "USD",
		},
	}

	mockProvider.EXPECT().GetBatchBalances(expectedRequests).Return([]*domain.BalanceResult{btcResult}, nil)

	svc := service.New(mockProvider)

	request := cryptowalletrest.BalancesPostRequest{
		Requests: []cryptowalletrest.BalancesPostRequestRequestsInner{
			{
				CryptoSymbol: "BTC",
				Address:      "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				// No FiatSymbol specified
			},
		},
	}

	response, err := svc.BalancesPost(t.Context(), request)

	require.NoError(t, err)
	assert.Equal(t, 200, response.Code)

	responseBody, ok := response.Body.(cryptowalletrest.BalancesPost200Response)
	require.True(t, ok)
	require.Len(t, responseBody.Results, 1)

	btcBalance := responseBody.Results[0]
	assert.Equal(t, "BTC", btcBalance.CryptoSymbol)
	assert.Equal(t, "USD", btcBalance.FiatSymbol)
}

func TestService_BalancesPost_EmptyRequests(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)

	mockProvider.EXPECT().GetBatchBalances([]domain.BalanceRequest{}).Return([]*domain.BalanceResult{}, nil)

	svc := service.New(mockProvider)

	request := cryptowalletrest.BalancesPostRequest{
		Requests: []cryptowalletrest.BalancesPostRequestRequestsInner{},
	}

	response, err := svc.BalancesPost(t.Context(), request)

	require.NoError(t, err)
	assert.Equal(t, 200, response.Code)

	responseBody, ok := response.Body.(cryptowalletrest.BalancesPost200Response)
	require.True(t, ok)
	assert.Empty(t, responseBody.Results)
}

func TestService_BalancesPost_BatchOperationError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)

	expectedRequests := []domain.BalanceRequest{
		{
			CryptoSymbol: "BTC",
			Address:      "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
			FiatSymbol:   "USD",
		},
	}

	mockProvider.EXPECT().GetBatchBalances(expectedRequests).Return(nil, errProviderGeneric)

	svc := service.New(mockProvider)

	request := cryptowalletrest.BalancesPostRequest{
		Requests: []cryptowalletrest.BalancesPostRequestRequestsInner{
			{
				CryptoSymbol: "BTC",
				Address:      "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				FiatSymbol:   "USD",
			},
		},
	}

	response, err := svc.BalancesPost(t.Context(), request)

	require.NoError(t, err)
	assert.Equal(t, 501, response.Code)

	errorResponse, ok := response.Body.(service.Error)
	require.True(t, ok)
	assert.Equal(t, "provider error", errorResponse.Message)
}

func TestTransactionsGet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)
	svc := service.New(mockProvider)

	response, err := svc.TransactionsGet(t.Context(), "BTC", "address", 10, 0)

	require.NoError(t, err)
	assert.Equal(t, 501, response.Code)
}

func TestUnsignedTxGet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)
	svc := service.New(mockProvider)

	response, err := svc.UnsignedTxGet(t.Context(), "BTC", "from", "to", "USD", 1.0)

	require.NoError(t, err)
	assert.Equal(t, 501, response.Code)
}

func TestBroadcastPost(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)
	svc := service.New(mockProvider)

	response, err := svc.BroadcastPost(t.Context(), cryptowalletrest.BroadcastPostRequest{})

	require.NoError(t, err)
	assert.Equal(t, 501, response.Code)
}
