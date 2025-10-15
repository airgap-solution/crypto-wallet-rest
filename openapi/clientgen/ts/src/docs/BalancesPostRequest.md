# BalancesPostRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**requests** | [**Array&lt;BalancesPostRequestRequestsInner&gt;**](BalancesPostRequestRequestsInner.md) |  | [default to undefined]
**fiat_symbol** | **string** | Default fiat currency symbol for all requests if not specified individually | [optional] [default to 'USD']

## Example

```typescript
import { BalancesPostRequest } from '@airgap-solution/crypto-wallet-rest';

const instance: BalancesPostRequest = {
    requests,
    fiat_symbol,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
