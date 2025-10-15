## @airgap-solution/crypto-wallet-rest@1.0.2

This generator creates TypeScript/JavaScript client that utilizes [axios](https://github.com/axios/axios). The generated Node module can be used in the following environments:

Environment
* Node.js
* Webpack
* Browserify

Language level
* ES5 - you must have a Promises/A+ library installed
* ES6

Module system
* CommonJS
* ES6 module system

It can be used in both TypeScript and JavaScript. In TypeScript, the definition will be automatically resolved via `package.json`. ([Reference](https://www.typescriptlang.org/docs/handbook/declaration-files/consumption.html))

### Building

To build and compile the typescript sources to javascript use:
```
npm install
npm run build
```

### Publishing

First build the package then run `npm publish`

### Consuming

navigate to the folder of your consuming project and run one of the following commands.

_published:_

```
npm install @airgap-solution/crypto-wallet-rest@1.0.2 --save
```

_unPublished (not recommended):_

```
npm install PATH_TO_GENERATED_PACKAGE --save
```

### Documentation for API Endpoints

All URIs are relative to *http://localhost*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*DefaultApi* | [**balancesGet**](docs/DefaultApi.md#balancesget) | **GET** /balances | Get balances for multiple addresses and cryptocurrencies
*DefaultApi* | [**broadcastPost**](docs/DefaultApi.md#broadcastpost) | **POST** /broadcast | Broadcast signed transaction
*DefaultApi* | [**transactionsGet**](docs/DefaultApi.md#transactionsget) | **GET** /transactions | Get transaction history for an address
*DefaultApi* | [**unsignedTxGet**](docs/DefaultApi.md#unsignedtxget) | **GET** /unsigned-tx | Generate an unsigned transaction


### Documentation For Models

 - [BalancesGet200Response](docs/BalancesGet200Response.md)
 - [BalancesGet200ResponseResultsInner](docs/BalancesGet200ResponseResultsInner.md)
 - [BalancesGetRequest](docs/BalancesGetRequest.md)
 - [BalancesGetRequestRequestsInner](docs/BalancesGetRequestRequestsInner.md)
 - [BroadcastPost200Response](docs/BroadcastPost200Response.md)
 - [BroadcastPostRequest](docs/BroadcastPostRequest.md)
 - [ErrorResponse](docs/ErrorResponse.md)
 - [Transaction](docs/Transaction.md)
 - [TransactionsGet200Response](docs/TransactionsGet200Response.md)
 - [UnsignedTxGet200Response](docs/UnsignedTxGet200Response.md)


<a id="documentation-for-authorization"></a>
## Documentation For Authorization

Endpoints do not require authorization.

