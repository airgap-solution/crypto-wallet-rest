package internal

import (
	"net/http"
	"time"

	"github.com/airgap-solution/crypto-wallet-rest/internal/config"
	cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"
)

var readTimeout = time.Second * 10

func Assemble(cfg config.Config, servicer cryptowalletrest.DefaultAPIServicer) *http.Server {
	ctrl := cryptowalletrest.NewDefaultAPIController(servicer)
	router := cryptowalletrest.NewRouter(ctrl)
	srv := &http.Server{Addr: cfg.ListenAddr, Handler: router, ReadTimeout: readTimeout}
	return srv
}
