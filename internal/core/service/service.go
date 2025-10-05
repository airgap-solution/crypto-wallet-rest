package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/airgap-solution/crypto-wallet-rest/internal/ports"
	cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"
)

type Service struct {
	adapter ports.Provider
}

func New(adapter ports.Provider) Service {
	return Service{adapter: adapter}
}

func (s Service) BalanceGet(_ context.Context, symbol string, address string) (cryptowalletrest.ImplResponse, error) {
	balance, err := s.adapter.GetBalance(symbol, address)
	if err != nil {
		return handleError(err)
	}

	return cryptowalletrest.Response(http.StatusOK, cryptowalletrest.BalanceGet200Response{
		Crypto:  symbol,
		Address: address,
		Balance: fmt.Sprintf("%.f", balance),
	}), nil
}
func (s Service) TransactionsGet(context.Context, string, string) (cryptowalletrest.ImplResponse, error) {
	return cryptowalletrest.ImplResponse{}, nil
}
func (s Service) UnsignedTxGet(context.Context, string, string, string, string) (cryptowalletrest.ImplResponse, error) {
	return cryptowalletrest.ImplResponse{}, nil
}
func (s Service) BroadcastPost(context.Context, cryptowalletrest.BroadcastPostRequest) (
	cryptowalletrest.ImplResponse, error) {
	return cryptowalletrest.ImplResponse{}, nil
}
