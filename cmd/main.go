package main

import (
	"fmt"
	"log"
	"os"

	cmcrest "github.com/airgap-solution/cmc-rest/openapi/clientgen/go"
	"github.com/airgap-solution/crypto-wallet-rest/internal"
	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/crypto/providers/bitcoin"
	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/crypto/providers/ethereum"
	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/crypto/providers/kaspa"
	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/crypto/providers/litecoin"
	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/crypto/providers/solana"
	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/provider"
	"github.com/airgap-solution/crypto-wallet-rest/internal/config"
	"github.com/airgap-solution/crypto-wallet-rest/internal/core/service"
	"github.com/airgap-solution/crypto-wallet-rest/internal/ports"
	"github.com/restartfu/gophig"
)

func main() {
	conf, err := loadConfig("./config.toml")
	if err != nil {
		log.Fatalln(err)
	}

	cmcRestCfg := cmcrest.NewConfiguration()
	cmcRestCfg.Scheme = "https"
	cmcRestCfg.Host = conf.CMCRestAddr
	cmcRestClient := cmcrest.NewAPIClient(cmcRestCfg)

	providerAdapter := provider.NewAdapter(cmcRestClient.DefaultAPI, map[string]ports.CryptoProvider{
		"KAS":         kaspa.NewAdapter(conf.Crypto.Kaspa.MainnetRPC),
		"BTC":         bitcoin.NewAdapter(conf.Crypto.Bitcoin.MainnetRPC, false),
		"BTC_TESTNET": bitcoin.NewAdapter(conf.Crypto.Bitcoin.TestnetRPC, true),
		"LTC":         litecoin.NewAdapter(conf.Crypto.Litecoin.MainnetRPC, false),
		"LTC_TESTNET": litecoin.NewAdapter(conf.Crypto.Litecoin.TestnetRPC, true),
		"ETH":         ethereum.NewAdapter(conf.Crypto.Ethereum.MainnetRPC),
		"ETH_TESTNET": ethereum.NewAdapter(conf.Crypto.Ethereum.TestnetRPC),
		"SOL":         solana.NewAdapter(conf.Crypto.Solana.MainnetRPC),
		"SOL_TESTNET": solana.NewAdapter(conf.Crypto.Solana.TestnetRPC),
	})
	servicer := service.New(providerAdapter)

	srv := internal.Assemble(conf, servicer)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func loadConfig(configPath string) (config.Config, error) {
	defaultConfig := config.DefaultConfig()
	g := gophig.NewGophig[config.Config](configPath, gophig.TOMLMarshaler{}, os.ModePerm)
	conf, err := g.LoadConf()
	if err != nil {
		if os.IsNotExist(err) {
			err = g.SaveConf(defaultConfig)
			return defaultConfig, fmt.Errorf("could not save default config: %w", err)
		}
		return config.Config{}, fmt.Errorf("could not load config: %w", err)
	}
	return conf, nil
}
