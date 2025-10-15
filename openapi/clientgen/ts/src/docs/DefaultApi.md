# DefaultApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**balanceGet**](#balanceget) | **GET** /balance | Get balance for an address|
|[**balancesPost**](#balancespost) | **POST** /balances | Get balances for multiple addresses and cryptocurrencies|
|[**broadcastPost**](#broadcastpost) | **POST** /broadcast | Broadcast signed transaction|
|[**transactionsGet**](#transactionsget) | **GET** /transactions | Get transaction history for an address|
|[**unsignedTxGet**](#unsignedtxget) | **GET** /unsigned-tx | Generate an unsigned transaction|

# **balanceGet**
> BalanceGet200Response balanceGet()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from '@airgap-solution/crypto-wallet-rest';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let cryptoSymbol: string; // (default to undefined)
let address: string; // (default to undefined)
let fiatSymbol: string; // (optional) (default to 'USD')

const { status, data } = await apiInstance.balanceGet(
    cryptoSymbol,
    address,
    fiatSymbol
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **cryptoSymbol** | [**string**] |  | defaults to undefined|
| **address** | [**string**] |  | defaults to undefined|
| **fiatSymbol** | [**string**] |  | (optional) defaults to 'USD'|


### Return type

**BalanceGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Balance response with crypto, fiat values, and 24h change |  -  |
|**400** | Bad request (invalid parameters) |  -  |
|**404** | Cryptocurrency or address not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **balancesPost**
> BalancesPost200Response balancesPost(balancesPostRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    BalancesPostRequest
} from '@airgap-solution/crypto-wallet-rest';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let balancesPostRequest: BalancesPostRequest; //

const { status, data } = await apiInstance.balancesPost(
    balancesPostRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **balancesPostRequest** | **BalancesPostRequest**|  | |


### Return type

**BalancesPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Batch balance response with crypto, fiat values, and 24h changes |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **broadcastPost**
> BroadcastPost200Response broadcastPost(broadcastPostRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    BroadcastPostRequest
} from '@airgap-solution/crypto-wallet-rest';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let broadcastPostRequest: BroadcastPostRequest; //

const { status, data } = await apiInstance.broadcastPost(
    broadcastPostRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **broadcastPostRequest** | **BroadcastPostRequest**|  | |


### Return type

**BroadcastPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Transaction broadcast result |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **transactionsGet**
> TransactionsGet200Response transactionsGet()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from '@airgap-solution/crypto-wallet-rest';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let cryptoSymbol: string; // (default to undefined)
let address: string; // (default to undefined)
let limit: number; // (optional) (default to 50)
let offset: number; // (optional) (default to 0)

const { status, data } = await apiInstance.transactionsGet(
    cryptoSymbol,
    address,
    limit,
    offset
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **cryptoSymbol** | [**string**] |  | defaults to undefined|
| **address** | [**string**] |  | defaults to undefined|
| **limit** | [**number**] |  | (optional) defaults to 50|
| **offset** | [**number**] |  | (optional) defaults to 0|


### Return type

**TransactionsGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Transaction history |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **unsignedTxGet**
> UnsignedTxGet200Response unsignedTxGet()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from '@airgap-solution/crypto-wallet-rest';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let cryptoSymbol: string; // (default to undefined)
let fromAddress: string; // (default to undefined)
let toAddress: string; // (default to undefined)
let amount: string; // (default to undefined)
let feeRate: number; // (optional) (default to undefined)

const { status, data } = await apiInstance.unsignedTxGet(
    cryptoSymbol,
    fromAddress,
    toAddress,
    amount,
    feeRate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **cryptoSymbol** | [**string**] |  | defaults to undefined|
| **fromAddress** | [**string**] |  | defaults to undefined|
| **toAddress** | [**string**] |  | defaults to undefined|
| **amount** | [**string**] |  | defaults to undefined|
| **feeRate** | [**number**] |  | (optional) defaults to undefined|


### Return type

**UnsignedTxGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Unsigned transaction |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

