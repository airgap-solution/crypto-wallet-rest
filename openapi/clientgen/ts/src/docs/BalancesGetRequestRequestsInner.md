# BalancesGetRequestRequestsInner


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**crypto_symbol** | **string** | The cryptocurrency symbol (BTC, ETH, etc.) | [default to undefined]
**address** | **string** | The cryptocurrency address or xpub | [default to undefined]
**fiat_symbol** | **string** | The fiat currency symbol for conversion (USD, EUR, CAD, etc.) | [optional] [default to 'USD']

## Example

```typescript
import { BalancesGetRequestRequestsInner } from '@airgap-solution/crypto-wallet-rest';

const instance: BalancesGetRequestRequestsInner = {
    crypto_symbol,
    address,
    fiat_symbol,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
