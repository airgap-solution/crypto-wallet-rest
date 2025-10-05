package main

import (
	"fmt"
	"log"
	"os"

	"github.com/airgap-solution/crypto-wallet-rest/internal"
	"github.com/airgap-solution/crypto-wallet-rest/internal/adapters/provider"
	"github.com/airgap-solution/crypto-wallet-rest/internal/config"
	"github.com/airgap-solution/crypto-wallet-rest/internal/core/service"
	"github.com/restartfu/gophig"
)

func main() {
	conf, err := loadConfig("./config.toml")
	if err != nil {
		log.Fatalln(err)
	}

	providerAdapter := provider.NewAdapter()
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
