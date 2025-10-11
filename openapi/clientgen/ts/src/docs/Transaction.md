# Transaction


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**transaction_id** | **string** |  | [default to undefined]
**block_height** | **number** |  | [optional] [default to undefined]
**timestamp** | **string** |  | [default to undefined]
**amount** | **string** |  | [default to undefined]
**direction** | **string** |  | [default to undefined]
**confirmations** | **number** |  | [default to undefined]
**fee_amount** | **string** |  | [optional] [default to undefined]
**from_addresses** | **Array&lt;string&gt;** |  | [optional] [default to undefined]
**to_addresses** | **Array&lt;string&gt;** |  | [optional] [default to undefined]

## Example

```typescript
import { Transaction } from '@airgap-solution/crypto-wallet-rest';

const instance: Transaction = {
    transaction_id,
    block_height,
    timestamp,
    amount,
    direction,
    confirmations,
    fee_amount,
    from_addresses,
    to_addresses,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
