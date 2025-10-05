package service_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/airgap-solution/crypto-wallet-rest/internal/core/service"
	internalportsmocks "github.com/airgap-solution/crypto-wallet-rest/mocks/internalports"
	cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProvider := internalportsmocks.NewMockProvider(ctrl)
	svc := service.New(mockProvider)
	assert.NotNil(t, svc)
}

func TestBalanceGet(t *testing.T) {
	symbol := "BTC"
	address := "address"
	balance := 1.23456

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
				mockProvider.EXPECT().GetBalance(symbol, address).Return(balance, nil)
			},
			expectedResponseBody: cryptowalletrest.BalanceGet200Response{
				Crypto:  symbol,
				Address: address,
				Balance: fmt.Sprintf("%.f", balance),
			},
			expectedResponseCode: 200,
			expectedError:        nil,
		},
		{
			name: "provider returns error",
			setupMocks: func(mockProvider *internalportsmocks.MockProvider) {
				mockProvider.EXPECT().GetBalance(symbol, address).Return(0.0, errors.New("provider error"))
			},
			expectedResponseBody: service.Error{Message: "provider error"},
			expectedResponseCode: http.StatusNotImplemented,
			expectedError:        nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProvider := internalportsmocks.NewMockProvider(ctrl)
			tc.setupMocks(mockProvider)

			svc := service.New(mockProvider)
			ctx := context.Background()

			resp, err := svc.BalanceGet(ctx, symbol, address)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)
	svc := service.New(mockProvider)
	ctx := context.Background()

	resp, err := svc.TransactionsGet(ctx, "BTC", "address")

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestUnsignedTxGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)
	svc := service.New(mockProvider)
	ctx := context.Background()

	resp, err := svc.UnsignedTxGet(ctx, "BTC", "from", "to", "amount")

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestBroadcastPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := internalportsmocks.NewMockProvider(ctrl)
	svc := service.New(mockProvider)
	ctx := context.Background()

	req := cryptowalletrest.BroadcastPostRequest{}
	resp, err := svc.BroadcastPost(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
}
