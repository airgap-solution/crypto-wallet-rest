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
	assert.False(t, cfg.TLSEnabled)
	assert.Equal(t, "/etc/letsencrypt/live/restartfu.com/fullchain.pem", cfg.TLSConfig.CertificatePath)
	assert.Equal(t, "/etc/letsencrypt/live/restartfu.com/privkey.pem", cfg.TLSConfig.PrivateKeyPath)
}
