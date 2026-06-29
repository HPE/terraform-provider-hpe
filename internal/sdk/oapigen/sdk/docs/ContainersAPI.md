# \ContainersAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CloneImageContainerAction**](ContainersAPI.md#CloneImageContainerAction) | **Put** /api/containers/{id}/clone-image | Clone Specific Container to Image
[**ContainersAttachFloatingIp**](ContainersAPI.md#ContainersAttachFloatingIp) | **Put** /api/containers/{id}/attach-floating-ip | Attach Floating IP to Container
[**ContainersDetachFloatingIp**](ContainersAPI.md#ContainersDetachFloatingIp) | **Put** /api/containers/{id}/detach-floating-ip | Detach Floating IP from Container
[**EjectContainerAction**](ContainersAPI.md#EjectContainerAction) | **Put** /api/containers/{id}/eject | Eject a Specific Container
[**ExecuteContainerAction**](ContainersAPI.md#ExecuteContainerAction) | **Put** /api/containers/{id}/action | Execute Container Action
[**GetContainer**](ContainersAPI.md#GetContainer) | **Get** /api/containers/{id} | Get a Specific Container
[**GetContainerActions**](ContainersAPI.md#GetContainerActions) | **Get** /api/containers/{id}/actions | List Container Actions
[**ImportContainerAction**](ContainersAPI.md#ImportContainerAction) | **Put** /api/containers/{id}/import | Import a Specific Container
[**RestartContainerAction**](ContainersAPI.md#RestartContainerAction) | **Put** /api/containers/{id}/restart | Restart a Specific Container
[**StartContainerAction**](ContainersAPI.md#StartContainerAction) | **Put** /api/containers/{id}/start | Start a Specific Container
[**StopContainerAction**](ContainersAPI.md#StopContainerAction) | **Put** /api/containers/{id}/stop | Stop a Specific Container
[**SuspendContainerAction**](ContainersAPI.md#SuspendContainerAction) | **Put** /api/containers/{id}/suspend | Suspend a Specific Container



## CloneImageContainerAction

> CloneImageContainerAction200Response CloneImageContainerAction(ctx, id).CloneImageContainerActionRequest(cloneImageContainerActionRequest).Execute()

Clone Specific Container to Image



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
	cloneImageContainerActionRequest := *openapiclient.NewCloneImageContainerActionRequest() // CloneImageContainerActionRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ContainersAPI.CloneImageContainerAction(context.Background(), id).CloneImageContainerActionRequest(cloneImageContainerActionRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.CloneImageContainerAction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CloneImageContainerAction`: CloneImageContainerAction200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.CloneImageContainerAction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiCloneImageContainerActionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **cloneImageContainerActionRequest** | [**CloneImageContainerActionRequest**](CloneImageContainerActionRequest.md) |  | 

### Return type

[**CloneImageContainerAction200Response**](CloneImageContainerAction200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ContainersAttachFloatingIp

> ContainersAttachFloatingIp200Response ContainersAttachFloatingIp(ctx, id).ContainersAttachFloatingIpRequest(containersAttachFloatingIpRequest).Execute()

Attach Floating IP to Container



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
	containersAttachFloatingIpRequest := *openapiclient.NewContainersAttachFloatingIpRequest() // ContainersAttachFloatingIpRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ContainersAPI.ContainersAttachFloatingIp(context.Background(), id).ContainersAttachFloatingIpRequest(containersAttachFloatingIpRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.ContainersAttachFloatingIp``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ContainersAttachFloatingIp`: ContainersAttachFloatingIp200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.ContainersAttachFloatingIp`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiContainersAttachFloatingIpRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **containersAttachFloatingIpRequest** | [**ContainersAttachFloatingIpRequest**](ContainersAttachFloatingIpRequest.md) |  | 

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


## ContainersDetachFloatingIp

> ContainersAttachFloatingIp200Response ContainersDetachFloatingIp(ctx, id).Execute()

Detach Floating IP from Container



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
	resp, r, err := apiClient.ContainersAPI.ContainersDetachFloatingIp(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.ContainersDetachFloatingIp``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ContainersDetachFloatingIp`: ContainersAttachFloatingIp200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.ContainersDetachFloatingIp`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiContainersDetachFloatingIpRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ContainersAttachFloatingIp200Response**](ContainersAttachFloatingIp200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EjectContainerAction

> EjectContainerAction200Response EjectContainerAction(ctx, id).Execute()

Eject a Specific Container



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
	resp, r, err := apiClient.ContainersAPI.EjectContainerAction(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.EjectContainerAction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `EjectContainerAction`: EjectContainerAction200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.EjectContainerAction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiEjectContainerActionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**EjectContainerAction200Response**](EjectContainerAction200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExecuteContainerAction

> ExecuteContainerAction200Response ExecuteContainerAction(ctx, id).Code(code).Execute()

Execute Container Action



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
	code := "apache-remove-node" // string | Action code to be executed on a specific container.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ContainersAPI.ExecuteContainerAction(context.Background(), id).Code(code).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.ExecuteContainerAction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecuteContainerAction`: ExecuteContainerAction200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.ExecuteContainerAction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiExecuteContainerActionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **code** | **string** | Action code to be executed on a specific container. | 

### Return type

[**ExecuteContainerAction200Response**](ExecuteContainerAction200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetContainer

> GetContainer200Response GetContainer(ctx, id).Execute()

Get a Specific Container



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
	resp, r, err := apiClient.ContainersAPI.GetContainer(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.GetContainer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetContainer`: GetContainer200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.GetContainer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetContainerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetContainer200Response**](GetContainer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetContainerActions

> GetContainerActions200Response GetContainerActions(ctx, id).Execute()

List Container Actions



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
	resp, r, err := apiClient.ContainersAPI.GetContainerActions(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.GetContainerActions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetContainerActions`: GetContainerActions200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.GetContainerActions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetContainerActionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetContainerActions200Response**](GetContainerActions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ImportContainerAction

> ImportContainerAction200Response ImportContainerAction(ctx, id).ImportContainerActionRequest(importContainerActionRequest).Execute()

Import a Specific Container



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
	importContainerActionRequest := *openapiclient.NewImportContainerActionRequest() // ImportContainerActionRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ContainersAPI.ImportContainerAction(context.Background(), id).ImportContainerActionRequest(importContainerActionRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.ImportContainerAction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ImportContainerAction`: ImportContainerAction200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.ImportContainerAction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiImportContainerActionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **importContainerActionRequest** | [**ImportContainerActionRequest**](ImportContainerActionRequest.md) |  | 

### Return type

[**ImportContainerAction200Response**](ImportContainerAction200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RestartContainerAction

> RestartContainerAction200Response RestartContainerAction(ctx, id).Execute()

Restart a Specific Container



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
	resp, r, err := apiClient.ContainersAPI.RestartContainerAction(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.RestartContainerAction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RestartContainerAction`: RestartContainerAction200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.RestartContainerAction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRestartContainerActionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**RestartContainerAction200Response**](RestartContainerAction200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StartContainerAction

> StartContainerAction200Response StartContainerAction(ctx, id).Execute()

Start a Specific Container



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
	resp, r, err := apiClient.ContainersAPI.StartContainerAction(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.StartContainerAction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `StartContainerAction`: StartContainerAction200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.StartContainerAction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiStartContainerActionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**StartContainerAction200Response**](StartContainerAction200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StopContainerAction

> StopContainerAction200Response StopContainerAction(ctx, id).Execute()

Stop a Specific Container



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
	resp, r, err := apiClient.ContainersAPI.StopContainerAction(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.StopContainerAction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `StopContainerAction`: StopContainerAction200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.StopContainerAction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiStopContainerActionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**StopContainerAction200Response**](StopContainerAction200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SuspendContainerAction

> SuspendContainerAction200Response SuspendContainerAction(ctx, id).Execute()

Suspend a Specific Container



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
	resp, r, err := apiClient.ContainersAPI.SuspendContainerAction(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ContainersAPI.SuspendContainerAction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SuspendContainerAction`: SuspendContainerAction200Response
	fmt.Fprintf(os.Stdout, "Response from `ContainersAPI.SuspendContainerAction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiSuspendContainerActionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**SuspendContainerAction200Response**](SuspendContainerAction200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

