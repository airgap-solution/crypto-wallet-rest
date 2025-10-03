package main

import (
	"log"
	"os"

	"github.com/airgap-solution/crypto-wallet-rest/internal"
	"github.com/airgap-solution/crypto-wallet-rest/internal/config"
	"github.com/restartfu/gophig"
)

func main() {
	conf, err := loadConfig("./config")
	if err != nil {
		log.Fatalln(err)
	}

	if err := internal.ListenAndServe(conf, nil); err != nil {
		log.Fatalln(err)
	}
}

func loadConfig(configPath string) (config.Config, error) {
	defaultConfig := config.DefaultConfig()
	g := gophig.NewGophig[config.Config](configPath, gophig.TOMLMarshaler{}, 0777)
	conf, err := g.LoadConf()
	if err != nil {
		if os.IsNotExist(err) {
			err = g.SaveConf(defaultConfig)
			return defaultConfig, err
		}
		return config.Config{}, err
	}
	return conf, nil
}
