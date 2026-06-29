# \HistoryAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CancelProcess**](HistoryAPI.md#CancelProcess) | **Post** /api/processes/{id}/cancel | Cancel a Specific Process
[**GetHistory**](HistoryAPI.md#GetHistory) | **Get** /api/processes/{id} | Retrieves a Specific Process
[**ListHistory**](HistoryAPI.md#ListHistory) | **Get** /api/processes | Retrieves Process History
[**RetryProcess**](HistoryAPI.md#RetryProcess) | **Post** /api/processes/{id}/retry | Retry a Specific Process



## CancelProcess

> ContainersAttachFloatingIp200Response CancelProcess(ctx, id).Execute()

Cancel a Specific Process



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
	id := int64(42) // int64 | Process ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.HistoryAPI.CancelProcess(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HistoryAPI.CancelProcess``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CancelProcess`: ContainersAttachFloatingIp200Response
	fmt.Fprintf(os.Stdout, "Response from `HistoryAPI.CancelProcess`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Process ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiCancelProcessRequest struct via the builder pattern


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


## GetHistory

> GetHistory200Response GetHistory(ctx, id).Execute()

Retrieves a Specific Process



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
	resp, r, err := apiClient.HistoryAPI.GetHistory(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HistoryAPI.GetHistory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetHistory`: GetHistory200Response
	fmt.Fprintf(os.Stdout, "Response from `HistoryAPI.GetHistory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetHistoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetHistory200Response**](GetHistory200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListHistory

> ListHistory200Response ListHistory(ctx).InstanceId(instanceId).ContainerId(containerId).ServerId(serverId).ZoneId(zoneId).AppId(appId).Phrase(phrase).Execute()

Retrieves Process History



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
	instanceId := int64(94) // int64 | The Instance ID for Filtering (optional)
	containerId := int64(135) // int64 | The Container ID for Filtering (optional)
	serverId := int64(97) // int64 | The Server ID for Filtering (optional)
	zoneId := int64(3) // int64 | The Cloud ID (Zone ID) for Filtering (optional)
	appId := int64(5) // int64 | The App ID for Filtering (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on message, displayName, output, event.message, event.output or event.error (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.HistoryAPI.ListHistory(context.Background()).InstanceId(instanceId).ContainerId(containerId).ServerId(serverId).ZoneId(zoneId).AppId(appId).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HistoryAPI.ListHistory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListHistory`: ListHistory200Response
	fmt.Fprintf(os.Stdout, "Response from `HistoryAPI.ListHistory`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListHistoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **instanceId** | **int64** | The Instance ID for Filtering | 
 **containerId** | **int64** | The Container ID for Filtering | 
 **serverId** | **int64** | The Server ID for Filtering | 
 **zoneId** | **int64** | The Cloud ID (Zone ID) for Filtering | 
 **appId** | **int64** | The App ID for Filtering | 
 **phrase** | **string** | Search phrase for partial matches on message, displayName, output, event.message, event.output or event.error | 

### Return type

[**ListHistory200Response**](ListHistory200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RetryProcess

> ContainersAttachFloatingIp200Response RetryProcess(ctx, id).Execute()

Retry a Specific Process



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
	id := int64(42) // int64 | Process ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.HistoryAPI.RetryProcess(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HistoryAPI.RetryProcess``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RetryProcess`: ContainersAttachFloatingIp200Response
	fmt.Fprintf(os.Stdout, "Response from `HistoryAPI.RetryProcess`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Process ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiRetryProcessRequest struct via the builder pattern


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

