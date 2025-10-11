# Transaction


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**transaction_id** | **string** | Transaction hash/ID | [default to undefined]
**block_height** | **number** | Block height (null if unconfirmed) | [optional] [default to undefined]
**timestamp** | **string** | Transaction timestamp | [default to undefined]
**amount** | **string** | Transaction amount in crypto units | [default to undefined]
**direction** | **string** | Transaction direction relative to the queried address | [default to undefined]
**confirmations** | **number** | Number of confirmations | [default to undefined]
**fee_amount** | **string** | Transaction fee (for outgoing transactions) | [optional] [default to undefined]
**from_addresses** | **Array&lt;string&gt;** | Source addresses | [optional] [default to undefined]
**to_addresses** | **Array&lt;string&gt;** | Destination addresses | [optional] [default to undefined]

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
