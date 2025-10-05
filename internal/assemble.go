package internal

import (
	"net/http"

	"github.com/airgap-solution/crypto-wallet-rest/internal/config"
	cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"
)

func Assemble(cfg config.Config, servicer cryptowalletrest.DefaultAPIServicer) *http.Server {
	ctrl := cryptowalletrest.NewDefaultAPIController(servicer)
	router := cryptowalletrest.NewRouter(ctrl)
	srv := &http.Server{Addr: cfg.ListenAddr, Handler: router}
	return srv
}
