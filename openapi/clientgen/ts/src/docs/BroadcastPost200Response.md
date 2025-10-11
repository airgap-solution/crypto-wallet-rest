# BroadcastPost200Response


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**crypto_symbol** | **string** |  | [default to undefined]
**transaction_id** | **string** | Transaction hash/ID | [default to undefined]
**status** | **string** |  | [default to undefined]
**message** | **string** | Human-readable status message | [default to undefined]
**network_fee** | **string** | Actual network fee paid | [optional] [default to undefined]
**timestamp** | **string** |  | [default to undefined]

## Example

```typescript
import { BroadcastPost200Response } from '@airgap-solution/crypto-wallet-rest';

const instance: BroadcastPost200Response = {
    crypto_symbol,
    transaction_id,
    status,
    message,
    network_fee,
    timestamp,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
