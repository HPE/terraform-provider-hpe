# \KeyPairsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddKeyPairs**](KeyPairsAPI.md#AddKeyPairs) | **Post** /api/key-pairs | Creates a Key Pair
[**GenerateKeyPairs**](KeyPairsAPI.md#GenerateKeyPairs) | **Post** /api/key-pairs/generate | Generates a Key Pair
[**GetKeyPairs**](KeyPairsAPI.md#GetKeyPairs) | **Get** /api/key-pairs/{id} | Retrieves a Specific Key Pair
[**RemoveKeyPairs**](KeyPairsAPI.md#RemoveKeyPairs) | **Delete** /api/key-pairs/{id} | Deletes a Key Pair



## AddKeyPairs

> AddKeyPairs200Response AddKeyPairs(ctx).AddKeyPairsRequest(addKeyPairsRequest).Execute()

Creates a Key Pair



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	addKeyPairsRequest := *openapiclient.NewAddKeyPairsRequest(*openapiclient.NewAddKeyPairsRequestKeyPair("Name_example", "PublicKey_example")) // AddKeyPairsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KeyPairsAPI.AddKeyPairs(context.Background()).AddKeyPairsRequest(addKeyPairsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KeyPairsAPI.AddKeyPairs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddKeyPairs`: AddKeyPairs200Response
	fmt.Fprintf(os.Stdout, "Response from `KeyPairsAPI.AddKeyPairs`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddKeyPairsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addKeyPairsRequest** | [**AddKeyPairsRequest**](AddKeyPairsRequest.md) |  | 

### Return type

[**AddKeyPairs200Response**](AddKeyPairs200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GenerateKeyPairs

> GenerateKeyPairs200Response GenerateKeyPairs(ctx).GenerateKeyPairsRequest(generateKeyPairsRequest).Execute()

Generates a Key Pair



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	generateKeyPairsRequest := *openapiclient.NewGenerateKeyPairsRequest(*openapiclient.NewGenerateKeyPairsRequestKeyPair("Name_example")) // GenerateKeyPairsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KeyPairsAPI.GenerateKeyPairs(context.Background()).GenerateKeyPairsRequest(generateKeyPairsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KeyPairsAPI.GenerateKeyPairs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GenerateKeyPairs`: GenerateKeyPairs200Response
	fmt.Fprintf(os.Stdout, "Response from `KeyPairsAPI.GenerateKeyPairs`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGenerateKeyPairsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **generateKeyPairsRequest** | [**GenerateKeyPairsRequest**](GenerateKeyPairsRequest.md) |  | 

### Return type

[**GenerateKeyPairs200Response**](GenerateKeyPairs200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetKeyPairs

> GetKeyPairs200Response GetKeyPairs(ctx, id).Execute()

Retrieves a Specific Key Pair



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KeyPairsAPI.GetKeyPairs(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KeyPairsAPI.GetKeyPairs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetKeyPairs`: GetKeyPairs200Response
	fmt.Fprintf(os.Stdout, "Response from `KeyPairsAPI.GetKeyPairs`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetKeyPairsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetKeyPairs200Response**](GetKeyPairs200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveKeyPairs

> RemoveExecuteSchedules200Response RemoveKeyPairs(ctx, id).Execute()

Deletes a Key Pair



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KeyPairsAPI.RemoveKeyPairs(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KeyPairsAPI.RemoveKeyPairs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveKeyPairs`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `KeyPairsAPI.RemoveKeyPairs`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveKeyPairsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**RemoveExecuteSchedules200Response**](RemoveExecuteSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

