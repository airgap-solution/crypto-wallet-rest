package service

import (
	"net/http"

	cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"
)

type Error struct {
	Message string `json:"message"`
}

func handleError(err error) (cryptowalletrest.ImplResponse, error) {
	code := http.StatusNotImplemented
	return cryptowalletrest.Response(code, Error{Message: err.Error()}), nil
}
