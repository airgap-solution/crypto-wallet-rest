package internal_test

import (
	"testing"

	"github.com/airgap-solution/crypto-wallet-rest/internal"
	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/provider"
	"github.com/airgap-solution/crypto-wallet-rest/internal/config"
	"github.com/airgap-solution/crypto-wallet-rest/internal/core/service"
	"github.com/stretchr/testify/require"
)

func TestAssemble(t *testing.T) {
	cfg := config.Config{
		ListenAddr: ":8080",
	}

	adapter := provider.NewAdapter()
	svc := service.New(adapter)

	srv := internal.Assemble(cfg, svc)

	require.NotNil(t, srv)
	require.Equal(t, ":8080", srv.Addr)
	require.NotNil(t, srv.Handler)
}
