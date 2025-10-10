package config

type Config struct {
	ListenAddr         string `toml:"listen_addr"`
	CMCRestAddr        string `toml:"cmc_rest_addr"`
	BitcoinRPC         string `toml:"bitcoin_rpc"`
	EthereumRPC        string `toml:"ethereum_rpc"`
	EthereumTestnetRPC string `toml:"ethereum_testnet_rpc"`
	SolanaRPC          string `toml:"solana_rpc"`
	SolanaTestnetRPC   string `toml:"solana_testnet_rpc"`
}

func DefaultConfig() Config {
	cfg := Config{
		ListenAddr:         ":8399",
		CMCRestAddr:        ":7392",
		BitcoinRPC:         "electrum.blockstream.info:50001",
		EthereumRPC:        "https://eth.llamarpc.com",
		EthereumTestnetRPC: "https://eth-sepolia.public.blastapi.io",
		SolanaRPC:          "https://api.mainnet-beta.solana.com",
		SolanaTestnetRPC:   "https://api.testnet.solana.com",
	}

	return cfg
}
