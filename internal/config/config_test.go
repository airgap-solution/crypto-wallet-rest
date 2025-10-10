package config_test

import (
	"testing"

	"github.com/airgap-solution/crypto-wallet-rest/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()
	cfg := config.DefaultConfig()

	assert.Equal(t, ":8399", cfg.ListenAddr)
	assert.Equal(t, ":7392", cfg.CMCRestAddr)
}
