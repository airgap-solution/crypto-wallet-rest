# BalancesGetRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**requests** | [**Array&lt;BalancesGetRequestRequestsInner&gt;**](BalancesGetRequestRequestsInner.md) |  | [default to undefined]
**fiat_symbol** | **string** | Default fiat currency symbol for all requests if not specified individually | [optional] [default to 'USD']

## Example

```typescript
import { BalancesGetRequest } from '@airgap-solution/crypto-wallet-rest';

const instance: BalancesGetRequest = {
    requests,
    fiat_symbol,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
