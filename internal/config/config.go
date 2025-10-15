package config

type CryptoConfig struct {
	Kaspa    RPCConfig `toml:"kaspa"`
	Bitcoin  RPCConfig `toml:"bitcoin"`
	Litecoin RPCConfig `toml:"litecoin"`
	Ethereum RPCConfig `toml:"ethereum"`
	Solana   RPCConfig `toml:"solana"`
}

type RPCConfig struct {
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
		CMCRestAddr: "192.168.2.71:8765",
		Crypto: CryptoConfig{
			Bitcoin: RPCConfig{
				MainnetRPC: "electrum.blockstream.info:50001",
				TestnetRPC: "electrum.blockstream.info:60001",
			},
			Litecoin: RPCConfig{
				MainnetRPC: "electrum-ltc.bysh.me:50001",
				TestnetRPC: "electrum-ltc.bysh.me:51001",
			},
			Ethereum: RPCConfig{
				MainnetRPC: "https://eth.llamarpc.com",
				TestnetRPC: "https://eth-sepolia.public.blastapi.io",
			},
			Solana: RPCConfig{
				MainnetRPC: "https://api.mainnet-beta.solana.com",
				TestnetRPC: "https://api.testnet.solana.com",
			},
		},
	}

	return cfg
}
