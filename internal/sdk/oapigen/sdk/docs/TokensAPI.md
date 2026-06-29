# \TokensAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddToken**](TokensAPI.md#AddToken) | **Post** /api/tokens | Create Token
[**BulkDeleteTokens**](TokensAPI.md#BulkDeleteTokens) | **Delete** /api/tokens | Delete Tokens in Bulk
[**DeleteToken**](TokensAPI.md#DeleteToken) | **Delete** /api/tokens/{id} | Delete Token
[**GetToken**](TokensAPI.md#GetToken) | **Get** /api/tokens/{id} | Retrieves a Specific Token
[**ListTokens**](TokensAPI.md#ListTokens) | **Get** /api/tokens | Get All Tokens



## AddToken

> AddToken200Response AddToken(ctx).UserId(userId).AddTokenRequest(addTokenRequest).Execute()

Create Token



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
	userId := int64(789) // int64 | ID of User to create the token for. Operating on a user other than your own requires admin permissions to manage users. (optional)
	addTokenRequest := *openapiclient.NewAddTokenRequest(*openapiclient.NewAddTokenRequestToken("morph-api")) // AddTokenRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TokensAPI.AddToken(context.Background()).UserId(userId).AddTokenRequest(addTokenRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TokensAPI.AddToken``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddToken`: AddToken200Response
	fmt.Fprintf(os.Stdout, "Response from `TokensAPI.AddToken`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddTokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int64** | ID of User to create the token for. Operating on a user other than your own requires admin permissions to manage users. | 
 **addTokenRequest** | [**AddTokenRequest**](AddTokenRequest.md) |  | 

### Return type

[**AddToken200Response**](AddToken200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BulkDeleteTokens

> RemoveExecuteSchedules200Response BulkDeleteTokens(ctx).Id(id).Execute()

Delete Tokens in Bulk



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
	id := []int64{int64(123)} // []int64 | Token ID to delete, can specify multiple for bulk delete

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TokensAPI.BulkDeleteTokens(context.Background()).Id(id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TokensAPI.BulkDeleteTokens``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BulkDeleteTokens`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `TokensAPI.BulkDeleteTokens`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiBulkDeleteTokensRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **[]int64** | Token ID to delete, can specify multiple for bulk delete | 

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


## DeleteToken

> RemoveExecuteSchedules200Response DeleteToken(ctx, id).UserId(userId).Execute()

Delete Token



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
	userId := int64(789) // int64 | ID of User that owns the token being deleted. Operating on a user other than your own requires admin permissions to manage users. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TokensAPI.DeleteToken(context.Background(), id).UserId(userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TokensAPI.DeleteToken``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteToken`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `TokensAPI.DeleteToken`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteTokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **userId** | **int64** | ID of User that owns the token being deleted. Operating on a user other than your own requires admin permissions to manage users. | 

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


## GetToken

> GetToken200Response GetToken(ctx, id).UserId(userId).Execute()

Retrieves a Specific Token



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
	userId := int64(789) // int64 | ID of User that owns the token being fetched. Operating on a user other than your own requires admin permissions to manage users. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TokensAPI.GetToken(context.Background(), id).UserId(userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TokensAPI.GetToken``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetToken`: GetToken200Response
	fmt.Fprintf(os.Stdout, "Response from `TokensAPI.GetToken`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **userId** | **int64** | ID of User that owns the token being fetched. Operating on a user other than your own requires admin permissions to manage users. | 

### Return type

[**GetToken200Response**](GetToken200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListTokens

> ListTokens200Response ListTokens(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).ClientId(clientId).Name(name).Token(token).UserId(userId).Execute()

Get All Tokens



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
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name (optional)
	clientId := "clientId_example" // string | If specified will return an exact match on clientId (optional)
	name := "name_example" // string | If specified will return an exact match on name (optional)
	token := "token_example" // string | If specified will return an exact match on access token value (should be kept secret) (optional)
	userId := int64(789) // int64 | ID of User to fetch tokens for. Operating on a user other than your own requires admin permissions to manage users. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TokensAPI.ListTokens(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).ClientId(clientId).Name(name).Token(token).UserId(userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TokensAPI.ListTokens``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListTokens`: ListTokens200Response
	fmt.Fprintf(os.Stdout, "Response from `TokensAPI.ListTokens`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListTokensRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name | 
 **clientId** | **string** | If specified will return an exact match on clientId | 
 **name** | **string** | If specified will return an exact match on name | 
 **token** | **string** | If specified will return an exact match on access token value (should be kept secret) | 
 **userId** | **int64** | ID of User to fetch tokens for. Operating on a user other than your own requires admin permissions to manage users. | 

### Return type

[**ListTokens200Response**](ListTokens200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

