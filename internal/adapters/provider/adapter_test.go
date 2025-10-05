package provider_test

import (
	"testing"

	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/provider"
	"github.com/stretchr/testify/require"
)

func TestNewAdapter(t *testing.T) {
	adapter := provider.NewAdapter()
	require.NotNil(t, adapter)
}

func TestGetBalance(t *testing.T) {
	adapter := provider.NewAdapter()

	balance, err := adapter.GetBalance("BTC", "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")

	require.NoError(t, err)
	require.Equal(t, 0.0, balance)
}
