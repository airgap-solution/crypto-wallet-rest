package config

type Config struct {
	ListenAddr  string `toml:"listen_addr"`
	CMCRestAddr string `toml:"cmc_rest_addr"`

	TLSEnabled bool `toml:"tls_enabled"`
	TLSConfig  struct {
		CertificatePath string `toml:"certificate_path"`
		PrivateKeyPath  string `toml:"private_key_path"`
	} `toml:"tls_config"`
}

func DefaultConfig() Config {
	cfg := Config{
		ListenAddr:  ":8399",
		CMCRestAddr: ":7392",
		TLSEnabled:  false,
	}

	cfg.TLSConfig.CertificatePath = "/etc/letsencrypt/live/restartfu.com/fullchain.pem"
	cfg.TLSConfig.PrivateKeyPath = "/etc/letsencrypt/live/restartfu.com/privkey.pem"
	return cfg
}
