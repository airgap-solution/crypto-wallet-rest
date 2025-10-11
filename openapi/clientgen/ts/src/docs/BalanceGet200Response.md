# BalanceGet200Response


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**crypto_symbol** | **string** | The cryptocurrency symbol | [default to undefined]
**address** | **string** | The queried address | [default to undefined]
**crypto_balance** | **number** | Balance in the native cryptocurrency units | [default to undefined]
**fiat_symbol** | **string** | The fiat currency symbol used for conversion | [default to undefined]
**fiat_value** | **number** | Equivalent value in the specified fiat currency | [default to undefined]
**exchange_rate** | **number** | Current exchange rate (crypto to fiat) | [default to undefined]
**timestamp** | **string** | Timestamp when the balance was retrieved | [default to undefined]

## Example

```typescript
import { BalanceGet200Response } from '@airgap-solution/crypto-wallet-rest';

const instance: BalanceGet200Response = {
    crypto_symbol,
    address,
    crypto_balance,
    fiat_symbol,
    fiat_value,
    exchange_rate,
    timestamp,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
