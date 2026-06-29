# \MonitoringSettingsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetMonitoringSettings**](MonitoringSettingsAPI.md#GetMonitoringSettings) | **Get** /api/monitoring-settings | Get Monitoring Settings
[**UpdateMonitoringSettings**](MonitoringSettingsAPI.md#UpdateMonitoringSettings) | **Put** /api/monitoring-settings | Update Monitoring Settings



## GetMonitoringSettings

> GetMonitoringSettings200Response GetMonitoringSettings(ctx).Execute()

Get Monitoring Settings



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
	resp, r, err := apiClient.MonitoringSettingsAPI.GetMonitoringSettings(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MonitoringSettingsAPI.GetMonitoringSettings``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMonitoringSettings`: GetMonitoringSettings200Response
	fmt.Fprintf(os.Stdout, "Response from `MonitoringSettingsAPI.GetMonitoringSettings`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetMonitoringSettingsRequest struct via the builder pattern


### Return type

[**GetMonitoringSettings200Response**](GetMonitoringSettings200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMonitoringSettings

> UpdateMonitoringSettings200Response UpdateMonitoringSettings(ctx).UpdateMonitoringSettingsRequest(updateMonitoringSettingsRequest).Execute()

Update Monitoring Settings



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
	updateMonitoringSettingsRequest := *openapiclient.NewUpdateMonitoringSettingsRequest() // UpdateMonitoringSettingsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MonitoringSettingsAPI.UpdateMonitoringSettings(context.Background()).UpdateMonitoringSettingsRequest(updateMonitoringSettingsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MonitoringSettingsAPI.UpdateMonitoringSettings``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMonitoringSettings`: UpdateMonitoringSettings200Response
	fmt.Fprintf(os.Stdout, "Response from `MonitoringSettingsAPI.UpdateMonitoringSettings`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMonitoringSettingsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateMonitoringSettingsRequest** | [**UpdateMonitoringSettingsRequest**](UpdateMonitoringSettingsRequest.md) |  | 

### Return type

[**UpdateMonitoringSettings200Response**](UpdateMonitoringSettings200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

