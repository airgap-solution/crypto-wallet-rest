# TransactionsGet200Response


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**crypto_symbol** | **string** |  | [default to undefined]
**address** | **string** |  | [default to undefined]
**transactions** | [**Array&lt;Transaction&gt;**](Transaction.md) |  | [default to undefined]
**total_count** | **number** |  | [default to undefined]
**has_more** | **boolean** |  | [default to undefined]

## Example

```typescript
import { TransactionsGet200Response } from '@airgap-solution/crypto-wallet-rest';

const instance: TransactionsGet200Response = {
    crypto_symbol,
    address,
    transactions,
    total_count,
    has_more,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
