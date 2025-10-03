package servicer

import (
	"context"

	cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"
)

type Adapter struct{}

func (a Adapter) BalanceGet(context.Context, string, string) (cryptowalletrest.ImplResponse, error) {
	return cryptowalletrest.ImplResponse{}, nil
}
func (a Adapter) TransactionsGet(context.Context, string, string) (cryptowalletrest.ImplResponse, error) {
	return cryptowalletrest.ImplResponse{}, nil
}
func (a Adapter) UnsignedTxGet(context.Context, string, string, string, string) (cryptowalletrest.ImplResponse, error) {
	return cryptowalletrest.ImplResponse{}, nil
}
func (a Adapter) BroadcastPost(context.Context, cryptowalletrest.BroadcastPostRequest) (cryptowalletrest.ImplResponse, error) {
	return cryptowalletrest.ImplResponse{}, nil
}
