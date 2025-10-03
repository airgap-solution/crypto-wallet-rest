package internal

import (
	"net/http"

	"github.com/airgap-solution/crypto-wallet-rest/internal/config"
	cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"
)

func ListenAndServe(cfg config.Config, servicer cryptowalletrest.DefaultAPIServicer) error {
	ctrl := cryptowalletrest.NewDefaultAPIController(servicer)
	router := cryptowalletrest.NewRouter(ctrl)

	return http.ListenAndServe(cfg.ListenAddr, router)
}
