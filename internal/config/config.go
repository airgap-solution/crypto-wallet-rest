package config

type CryptoConfig struct {
	Bitcoin  BitcoinConfig  `toml:"bitcoin"`
	Ethereum EthereumConfig `toml:"ethereum"`
	Solana   SolanaConfig   `toml:"solana"`
}

type BitcoinConfig struct {
	RPC string `toml:"rpc"`
}

type EthereumConfig struct {
	MainnetRPC string `toml:"mainnet_rpc"`
	TestnetRPC string `toml:"testnet_rpc"`
}

type SolanaConfig struct {
	MainnetRPC string `toml:"mainnet_rpc"`
	TestnetRPC string `toml:"testnet_rpc"`
}

type Config struct {
	ListenAddr  string       `toml:"listen_addr"`
	CMCRestAddr string       `toml:"cmc_rest_addr"`
	Crypto      CryptoConfig `toml:"crypto"`
}

func DefaultConfig() Config {
	cfg := Config{
		ListenAddr:  ":8399",
		CMCRestAddr: ":7392",
		Crypto: CryptoConfig{
			Bitcoin: BitcoinConfig{
				RPC: "electrum.blockstream.info:50001",
			},
			Ethereum: EthereumConfig{
				MainnetRPC: "https://eth.llamarpc.com",
				TestnetRPC: "https://eth-sepolia.public.blastapi.io",
			},
			Solana: SolanaConfig{
				MainnetRPC: "https://api.mainnet-beta.solana.com",
				TestnetRPC: "https://api.testnet.solana.com",
			},
		},
	}

	return cfg
}
