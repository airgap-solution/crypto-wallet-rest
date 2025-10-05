package service_test

import (
	"fmt"
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
		setupMocks           func(mocksProvider internalportsmocks.MockProvider)
		expectedError        error
		expectedResponseBody any
		expectedResponseCode int
	}{
		{
			name: "provider returns no error",
			setupMocks: func(mocksProvider internalportsmocks.MockProvider) {
				mocksProvider.EXPECT().GetBalance("symbol", "address").Return(balance, nil)
			},
			expectedResponseBody: cryptowalletrest.BalanceGet200Response{
				Crypto:  symbol,
				Address: address,
				Balance: fmt.Sprintf("%.f", balance),
			},
			expectedError: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedError != nil {
				require.ErrorIs(t, nil, tc.expectedError)

			}
		})
	}
}
