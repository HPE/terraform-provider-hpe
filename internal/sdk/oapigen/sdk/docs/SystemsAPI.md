# \SystemsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddSystem**](SystemsAPI.md#AddSystem) | **Post** /api/infrastructure/systems | Create a System
[**AddUninitializedSystem**](SystemsAPI.md#AddUninitializedSystem) | **Post** /api/infrastructure/systems/uninitialized | Create an Uninitialized System
[**GetSystem**](SystemsAPI.md#GetSystem) | **Get** /api/infrastructure/systems/{id} | Get a System
[**InitializeSystem**](SystemsAPI.md#InitializeSystem) | **Put** /api/infrastructure/systems/{id}/initialize | Initialize a System
[**ListSystemTypeLayouts**](SystemsAPI.md#ListSystemTypeLayouts) | **Get** /api/infrastructure/system-types/{typeId}/layouts | List System Type Layouts
[**ListSystemTypes**](SystemsAPI.md#ListSystemTypes) | **Get** /api/infrastructure/system-types | List System Types
[**ListSystems**](SystemsAPI.md#ListSystems) | **Get** /api/infrastructure/systems | List Systems
[**RemoveSystem**](SystemsAPI.md#RemoveSystem) | **Delete** /api/infrastructure/systems/{id} | Delete a System
[**UpdateSystem**](SystemsAPI.md#UpdateSystem) | **Put** /api/infrastructure/systems/{id} | Update a System
[**ValidateSystem**](SystemsAPI.md#ValidateSystem) | **Get** /api/infrastructure/systems/{id}/validate | Validate a System



## AddSystem

> AddSystem200Response AddSystem(ctx).AddSystemRequest(addSystemRequest).Execute()

Create a System



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
	addSystemRequest := *openapiclient.NewAddSystemRequest() // AddSystemRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SystemsAPI.AddSystem(context.Background()).AddSystemRequest(addSystemRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemsAPI.AddSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddSystem`: AddSystem200Response
	fmt.Fprintf(os.Stdout, "Response from `SystemsAPI.AddSystem`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddSystemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addSystemRequest** | [**AddSystemRequest**](AddSystemRequest.md) |  | 

### Return type

[**AddSystem200Response**](AddSystem200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddUninitializedSystem

> AddUninitializedSystem200Response AddUninitializedSystem(ctx).AddUninitializedSystemRequest(addUninitializedSystemRequest).Execute()

Create an Uninitialized System



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
	addUninitializedSystemRequest := *openapiclient.NewAddUninitializedSystemRequest() // AddUninitializedSystemRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SystemsAPI.AddUninitializedSystem(context.Background()).AddUninitializedSystemRequest(addUninitializedSystemRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemsAPI.AddUninitializedSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddUninitializedSystem`: AddUninitializedSystem200Response
	fmt.Fprintf(os.Stdout, "Response from `SystemsAPI.AddUninitializedSystem`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddUninitializedSystemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addUninitializedSystemRequest** | [**AddUninitializedSystemRequest**](AddUninitializedSystemRequest.md) |  | 

### Return type

[**AddUninitializedSystem200Response**](AddUninitializedSystem200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSystem

> GetSystem200Response GetSystem(ctx, id).Execute()

Get a System



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
	resp, r, err := apiClient.SystemsAPI.GetSystem(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemsAPI.GetSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSystem`: GetSystem200Response
	fmt.Fprintf(os.Stdout, "Response from `SystemsAPI.GetSystem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSystemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetSystem200Response**](GetSystem200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InitializeSystem

> ContainersAttachFloatingIp200Response InitializeSystem(ctx, id).InitializeSystemRequest(initializeSystemRequest).Execute()

Initialize a System



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
	initializeSystemRequest := *openapiclient.NewInitializeSystemRequest() // InitializeSystemRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SystemsAPI.InitializeSystem(context.Background(), id).InitializeSystemRequest(initializeSystemRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemsAPI.InitializeSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InitializeSystem`: ContainersAttachFloatingIp200Response
	fmt.Fprintf(os.Stdout, "Response from `SystemsAPI.InitializeSystem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiInitializeSystemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **initializeSystemRequest** | [**InitializeSystemRequest**](InitializeSystemRequest.md) |  | 

### Return type

[**ContainersAttachFloatingIp200Response**](ContainersAttachFloatingIp200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSystemTypeLayouts

> ListSystemTypeLayouts200Response ListSystemTypeLayouts(ctx, typeId).Execute()

List System Type Layouts



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
	typeId := int64(1) // int64 | ID of the system type whose layouts to retrieve.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SystemsAPI.ListSystemTypeLayouts(context.Background(), typeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemsAPI.ListSystemTypeLayouts``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListSystemTypeLayouts`: ListSystemTypeLayouts200Response
	fmt.Fprintf(os.Stdout, "Response from `SystemsAPI.ListSystemTypeLayouts`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**typeId** | **int64** | ID of the system type whose layouts to retrieve. | 

### Other Parameters

Other parameters are passed through a pointer to a apiListSystemTypeLayoutsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListSystemTypeLayouts200Response**](ListSystemTypeLayouts200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSystemTypes

> ListSystemTypes200Response ListSystemTypes(ctx).Execute()

List System Types



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
	resp, r, err := apiClient.SystemsAPI.ListSystemTypes(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemsAPI.ListSystemTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListSystemTypes`: ListSystemTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `SystemsAPI.ListSystemTypes`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListSystemTypesRequest struct via the builder pattern


### Return type

[**ListSystemTypes200Response**](ListSystemTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSystems

> ListSystems200Response ListSystems(ctx).Execute()

List Systems



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
	resp, r, err := apiClient.SystemsAPI.ListSystems(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemsAPI.ListSystems``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListSystems`: ListSystems200Response
	fmt.Fprintf(os.Stdout, "Response from `SystemsAPI.ListSystems`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListSystemsRequest struct via the builder pattern


### Return type

[**ListSystems200Response**](ListSystems200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveSystem

> RemoveExecuteSchedules200Response RemoveSystem(ctx, id).Execute()

Delete a System



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
	resp, r, err := apiClient.SystemsAPI.RemoveSystem(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemsAPI.RemoveSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveSystem`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `SystemsAPI.RemoveSystem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveSystemRequest struct via the builder pattern


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


## UpdateSystem

> ContainersAttachFloatingIp200Response UpdateSystem(ctx, id).UpdateSystemRequest(updateSystemRequest).Execute()

Update a System



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
	updateSystemRequest := *openapiclient.NewUpdateSystemRequest() // UpdateSystemRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SystemsAPI.UpdateSystem(context.Background(), id).UpdateSystemRequest(updateSystemRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemsAPI.UpdateSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateSystem`: ContainersAttachFloatingIp200Response
	fmt.Fprintf(os.Stdout, "Response from `SystemsAPI.UpdateSystem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateSystemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateSystemRequest** | [**UpdateSystemRequest**](UpdateSystemRequest.md) |  | 

### Return type

[**ContainersAttachFloatingIp200Response**](ContainersAttachFloatingIp200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ValidateSystem

> ValidateSystem200Response ValidateSystem(ctx, id).Execute()

Validate a System



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
	resp, r, err := apiClient.SystemsAPI.ValidateSystem(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemsAPI.ValidateSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ValidateSystem`: ValidateSystem200Response
	fmt.Fprintf(os.Stdout, "Response from `SystemsAPI.ValidateSystem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiValidateSystemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ValidateSystem200Response**](ValidateSystem200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

