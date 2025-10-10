package config

type Config struct {
	ListenAddr  string `toml:"listen_addr"`
	CMCRestAddr string `toml:"cmc_rest_addr"`
}

func DefaultConfig() Config {
	cfg := Config{
		ListenAddr:  ":8399",
		CMCRestAddr: ":7392",
	}

	return cfg
}
