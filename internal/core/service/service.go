package service

import (
	"context"
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

func (s Service) BalanceGet(
	_ context.Context, cryptoSymbol string, address string, fiatSymbol string,
) (cryptowalletrest.ImplResponse, error) {
	balanceResult, err := s.adapter.GetBalance(cryptoSymbol, address, fiatSymbol)
	if err != nil {
		return handleError(err)
	}

	return cryptowalletrest.Response(http.StatusOK, cryptowalletrest.BalanceGet200Response{
		CryptoSymbol:  balanceResult.CryptoSymbol,
		Address:       balanceResult.Address,
		CryptoBalance: balanceResult.CryptoBalance,
		FiatSymbol:    balanceResult.FiatSymbol,
		FiatValue:     balanceResult.FiatValue,
		ExchangeRate:  balanceResult.ExchangeRate,
		Timestamp:     balanceResult.Timestamp,
	}), nil
}

func (s Service) TransactionsGet(
	_ context.Context, _ string, _ string, _ int32, _ int32,
) (cryptowalletrest.ImplResponse, error) {
	return cryptowalletrest.Response(http.StatusNotImplemented, nil), nil
}

func (s Service) UnsignedTxGet(
	_ context.Context, _ string, _ string, _ string, _ string, _ float64,
) (cryptowalletrest.ImplResponse, error) {
	return cryptowalletrest.Response(http.StatusNotImplemented, nil), nil
}

func (s Service) BroadcastPost(
	_ context.Context, _ cryptowalletrest.BroadcastPostRequest,
) (cryptowalletrest.ImplResponse, error) {
	return cryptowalletrest.Response(http.StatusNotImplemented, nil), nil
}
