# \IdentitySourcesAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddIdentitySources**](IdentitySourcesAPI.md#AddIdentitySources) | **Post** /api/user-sources | Creates an Identity Source
[**GetIdentitySources**](IdentitySourcesAPI.md#GetIdentitySources) | **Get** /api/user-sources/{id} | Retrieves a Specific Identity Source
[**ListIdentitySources**](IdentitySourcesAPI.md#ListIdentitySources) | **Get** /api/user-sources | Retrieves all Identity Sources
[**RemoveIdentitySources**](IdentitySourcesAPI.md#RemoveIdentitySources) | **Delete** /api/user-sources/{id} | Deletes an Identity Source
[**UpdateIdentitySourceSubdomains**](IdentitySourcesAPI.md#UpdateIdentitySourceSubdomains) | **Put** /api/user-sources/{id}/subdomain | Updates an Identity Source Subdomain
[**UpdateIdentitySources**](IdentitySourcesAPI.md#UpdateIdentitySources) | **Put** /api/user-sources/{id} | Updates an Identity Source



## AddIdentitySources

> AddIdentitySources200Response AddIdentitySources(ctx).AccountId(accountId).AddIdentitySourcesRequest(addIdentitySourcesRequest).Execute()

Creates an Identity Source



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
	accountId := int64(3) // int64 | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. (optional)
	addIdentitySourcesRequest := *openapiclient.NewAddIdentitySourcesRequest(*openapiclient.NewAddIdentitySourcesRequestUserSource("mydomain AD", "Type_example")) // AddIdentitySourcesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IdentitySourcesAPI.AddIdentitySources(context.Background()).AccountId(accountId).AddIdentitySourcesRequest(addIdentitySourcesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IdentitySourcesAPI.AddIdentitySources``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddIdentitySources`: AddIdentitySources200Response
	fmt.Fprintf(os.Stdout, "Response from `IdentitySourcesAPI.AddIdentitySources`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddIdentitySourcesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountId** | **int64** | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. | 
 **addIdentitySourcesRequest** | [**AddIdentitySourcesRequest**](AddIdentitySourcesRequest.md) |  | 

### Return type

[**AddIdentitySources200Response**](AddIdentitySources200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetIdentitySources

> GetIdentitySources200Response GetIdentitySources(ctx, id).Execute()

Retrieves a Specific Identity Source



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
	resp, r, err := apiClient.IdentitySourcesAPI.GetIdentitySources(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IdentitySourcesAPI.GetIdentitySources``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetIdentitySources`: GetIdentitySources200Response
	fmt.Fprintf(os.Stdout, "Response from `IdentitySourcesAPI.GetIdentitySources`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetIdentitySourcesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetIdentitySources200Response**](GetIdentitySources200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListIdentitySources

> ListIdentitySources200Response ListIdentitySources(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).AccountId(accountId).Execute()

Retrieves all Identity Sources



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
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	accountId := int64(3) // int64 | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IdentitySourcesAPI.ListIdentitySources(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).AccountId(accountId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IdentitySourcesAPI.ListIdentitySources``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListIdentitySources`: ListIdentitySources200Response
	fmt.Fprintf(os.Stdout, "Response from `IdentitySourcesAPI.ListIdentitySources`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListIdentitySourcesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **accountId** | **int64** | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. | 

### Return type

[**ListIdentitySources200Response**](ListIdentitySources200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveIdentitySources

> RemoveExecuteSchedules200Response RemoveIdentitySources(ctx, id).Execute()

Deletes an Identity Source



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
	resp, r, err := apiClient.IdentitySourcesAPI.RemoveIdentitySources(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IdentitySourcesAPI.RemoveIdentitySources``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveIdentitySources`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `IdentitySourcesAPI.RemoveIdentitySources`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveIdentitySourcesRequest struct via the builder pattern


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


## UpdateIdentitySourceSubdomains

> UpdateIdentitySourceSubdomains200Response UpdateIdentitySourceSubdomains(ctx, id).UpdateIdentitySourceSubdomainsRequest(updateIdentitySourceSubdomainsRequest).Execute()

Updates an Identity Source Subdomain



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
	updateIdentitySourceSubdomainsRequest := *openapiclient.NewUpdateIdentitySourceSubdomainsRequest("mynewdomain") // UpdateIdentitySourceSubdomainsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IdentitySourcesAPI.UpdateIdentitySourceSubdomains(context.Background(), id).UpdateIdentitySourceSubdomainsRequest(updateIdentitySourceSubdomainsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IdentitySourcesAPI.UpdateIdentitySourceSubdomains``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateIdentitySourceSubdomains`: UpdateIdentitySourceSubdomains200Response
	fmt.Fprintf(os.Stdout, "Response from `IdentitySourcesAPI.UpdateIdentitySourceSubdomains`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateIdentitySourceSubdomainsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateIdentitySourceSubdomainsRequest** | [**UpdateIdentitySourceSubdomainsRequest**](UpdateIdentitySourceSubdomainsRequest.md) |  | 

### Return type

[**UpdateIdentitySourceSubdomains200Response**](UpdateIdentitySourceSubdomains200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateIdentitySources

> UpdateIdentitySources200Response UpdateIdentitySources(ctx, id).UpdateIdentitySourcesRequest(updateIdentitySourcesRequest).Execute()

Updates an Identity Source



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
	updateIdentitySourcesRequest := *openapiclient.NewUpdateIdentitySourcesRequest(*openapiclient.NewUpdateIdentitySourcesRequestUserSource("mydomain AD", "Type_example")) // UpdateIdentitySourcesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IdentitySourcesAPI.UpdateIdentitySources(context.Background(), id).UpdateIdentitySourcesRequest(updateIdentitySourcesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IdentitySourcesAPI.UpdateIdentitySources``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateIdentitySources`: UpdateIdentitySources200Response
	fmt.Fprintf(os.Stdout, "Response from `IdentitySourcesAPI.UpdateIdentitySources`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateIdentitySourcesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateIdentitySourcesRequest** | [**UpdateIdentitySourcesRequest**](UpdateIdentitySourcesRequest.md) |  | 

### Return type

[**UpdateIdentitySources200Response**](UpdateIdentitySources200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

