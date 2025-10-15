# BalancesGet200ResponseResultsInner


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**crypto_symbol** | **string** |  | [default to undefined]
**address** | **string** |  | [default to undefined]
**crypto_balance** | **number** |  | [default to undefined]
**fiat_symbol** | **string** |  | [default to undefined]
**fiat_value** | **number** |  | [default to undefined]
**exchange_rate** | **number** |  | [default to undefined]
**change24h** | **number** | Absolute change in fiat value over the last 24 hours | [default to undefined]
**timestamp** | **string** |  | [default to undefined]
**error** | **string** | Error message if this specific balance fetch failed | [optional] [default to undefined]

## Example

```typescript
import { BalancesGet200ResponseResultsInner } from '@airgap-solution/crypto-wallet-rest';

const instance: BalancesGet200ResponseResultsInner = {
    crypto_symbol,
    address,
    crypto_balance,
    fiat_symbol,
    fiat_value,
    exchange_rate,
    change24h,
    timestamp,
    error,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
