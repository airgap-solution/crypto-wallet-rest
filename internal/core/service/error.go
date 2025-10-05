package service

import cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"

func handleError(err error) (cryptowalletrest.ImplResponse, error) {
	code := 200

	return cryptowalletrest.Response(code, err), nil
}
