package service

import (
	"context"
	"net/http"
	"time"

	"github.com/airgap-solution/crypto-wallet-rest/internal/core/domain"
	"github.com/airgap-solution/crypto-wallet-rest/internal/ports"
	cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"
)

type Service struct {
	adapter ports.Provider
}

func New(adapter ports.Provider) Service {
	return Service{adapter: adapter}
}

func (s Service) BalancesGet(
	_ context.Context, request cryptowalletrest.BalancesGetRequest,
) (cryptowalletrest.ImplResponse, error) {
	// Convert OpenAPI request to internal format
	balanceRequests := make([]domain.BalanceRequest, len(request.Requests))
	for i, req := range request.Requests {
		fiatSymbol := req.FiatSymbol
		if fiatSymbol == "" {
			fiatSymbol = "USD"
		}
		balanceRequests[i] = domain.BalanceRequest{
			CryptoSymbol: req.CryptoSymbol,
			Address:      req.Address,
			FiatSymbol:   fiatSymbol,
		}
	}

	// Get all balances using batch method
	results, err := s.adapter.GetBatchBalances(balanceRequests)

	if err != nil {
		return handleError(err)
	}

	// Convert results to OpenAPI format
	balances := make([]cryptowalletrest.BalancesGet200ResponseResultsInner, len(results))
	for i, result := range results {
		balance := cryptowalletrest.BalancesGet200ResponseResultsInner{
			CryptoSymbol:  result.CryptoSymbol,
			Address:       result.Address,
			CryptoBalance: result.CryptoBalance,
			FiatSymbol:    result.FiatSymbol,
			FiatValue:     result.FiatValue,
			ExchangeRate:  result.ExchangeRate,
			Change24h:     result.Change24h,
			Timestamp:     result.Timestamp,
		}
		if result.Error != nil {
			balance.Error = *result.Error
		}
		balances[i] = balance
	}

	return cryptowalletrest.Response(http.StatusOK, cryptowalletrest.BalancesGet200Response{
		Results:   balances,
		Timestamp: time.Now(),
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
