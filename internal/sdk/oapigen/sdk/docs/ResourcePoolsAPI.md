# \ResourcePoolsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateResourcePoolGroup**](ResourcePoolsAPI.md#CreateResourcePoolGroup) | **Post** /api/resource-pools/groups | Create a Resource Pool Group
[**DeleteResourcePoolGroup**](ResourcePoolsAPI.md#DeleteResourcePoolGroup) | **Delete** /api/resource-pools/groups/{id} | Delete a Resource Pool Group
[**GetResourcePoolGroups**](ResourcePoolsAPI.md#GetResourcePoolGroups) | **Get** /api/resource-pools/groups | Get all Resource Pool Groups
[**GetresourcePoolGroup**](ResourcePoolsAPI.md#GetresourcePoolGroup) | **Get** /api/resource-pools/groups/{id} | Get a Specific Resource Pool Group
[**UpdateResourcePoolGroup**](ResourcePoolsAPI.md#UpdateResourcePoolGroup) | **Put** /api/resource-pools/groups/{id} | Update a Resource Pool Group



## CreateResourcePoolGroup

> CreateResourcePoolGroup200Response CreateResourcePoolGroup(ctx).CreateResourcePoolGroupRequest(createResourcePoolGroupRequest).Execute()

Create a Resource Pool Group



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
	createResourcePoolGroupRequest := *openapiclient.NewCreateResourcePoolGroupRequest() // CreateResourcePoolGroupRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ResourcePoolsAPI.CreateResourcePoolGroup(context.Background()).CreateResourcePoolGroupRequest(createResourcePoolGroupRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ResourcePoolsAPI.CreateResourcePoolGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateResourcePoolGroup`: CreateResourcePoolGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `ResourcePoolsAPI.CreateResourcePoolGroup`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateResourcePoolGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createResourcePoolGroupRequest** | [**CreateResourcePoolGroupRequest**](CreateResourcePoolGroupRequest.md) |  | 

### Return type

[**CreateResourcePoolGroup200Response**](CreateResourcePoolGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteResourcePoolGroup

> DeleteResourcePoolGroup200Response DeleteResourcePoolGroup(ctx, id).Execute()

Delete a Resource Pool Group



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
	resp, r, err := apiClient.ResourcePoolsAPI.DeleteResourcePoolGroup(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ResourcePoolsAPI.DeleteResourcePoolGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteResourcePoolGroup`: DeleteResourcePoolGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `ResourcePoolsAPI.DeleteResourcePoolGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteResourcePoolGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteResourcePoolGroup200Response**](DeleteResourcePoolGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetResourcePoolGroups

> GetResourcePoolGroups200Response GetResourcePoolGroups(ctx).Execute()

Get all Resource Pool Groups



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ResourcePoolsAPI.GetResourcePoolGroups(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ResourcePoolsAPI.GetResourcePoolGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetResourcePoolGroups`: GetResourcePoolGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `ResourcePoolsAPI.GetResourcePoolGroups`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetResourcePoolGroupsRequest struct via the builder pattern


### Return type

[**GetResourcePoolGroups200Response**](GetResourcePoolGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetresourcePoolGroup

> GetresourcePoolGroup200Response GetresourcePoolGroup(ctx, id).Execute()

Get a Specific Resource Pool Group



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
	resp, r, err := apiClient.ResourcePoolsAPI.GetresourcePoolGroup(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ResourcePoolsAPI.GetresourcePoolGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetresourcePoolGroup`: GetresourcePoolGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `ResourcePoolsAPI.GetresourcePoolGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetresourcePoolGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetresourcePoolGroup200Response**](GetresourcePoolGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateResourcePoolGroup

> UpdateResourcePoolGroup200Response UpdateResourcePoolGroup(ctx, id).UpdateResourcePoolGroupRequest(updateResourcePoolGroupRequest).Execute()

Update a Resource Pool Group



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
	updateResourcePoolGroupRequest := *openapiclient.NewUpdateResourcePoolGroupRequest() // UpdateResourcePoolGroupRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ResourcePoolsAPI.UpdateResourcePoolGroup(context.Background(), id).UpdateResourcePoolGroupRequest(updateResourcePoolGroupRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ResourcePoolsAPI.UpdateResourcePoolGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateResourcePoolGroup`: UpdateResourcePoolGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `ResourcePoolsAPI.UpdateResourcePoolGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateResourcePoolGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateResourcePoolGroupRequest** | [**UpdateResourcePoolGroupRequest**](UpdateResourcePoolGroupRequest.md) |  | 

### Return type

[**UpdateResourcePoolGroup200Response**](UpdateResourcePoolGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

