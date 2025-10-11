# UnsignedTxGet200Response


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**crypto_symbol** | **string** |  | [default to undefined]
**from_address** | **string** |  | [default to undefined]
**to_address** | **string** |  | [default to undefined]
**amount** | **string** |  | [default to undefined]
**fee_amount** | **string** | Calculated transaction fee | [default to undefined]
**unsigned_tx** | **string** | Base64 or hex encoded unsigned transaction or PSBT | [default to undefined]
**tx_size_bytes** | **number** | Estimated transaction size in bytes | [optional] [default to undefined]

## Example

```typescript
import { UnsignedTxGet200Response } from '@airgap-solution/crypto-wallet-rest';

const instance: UnsignedTxGet200Response = {
    crypto_symbol,
    from_address,
    to_address,
    amount,
    fee_amount,
    unsigned_tx,
    tx_size_bytes,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
