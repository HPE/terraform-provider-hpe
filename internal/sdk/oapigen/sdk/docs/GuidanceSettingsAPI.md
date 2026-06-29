# \GuidanceSettingsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetGuidanceSettings**](GuidanceSettingsAPI.md#GetGuidanceSettings) | **Get** /api/guidance-settings | Get Guidance Settings
[**UpdateGuidanceSettings**](GuidanceSettingsAPI.md#UpdateGuidanceSettings) | **Put** /api/guidance-settings | Update Guidance Settings



## GetGuidanceSettings

> GetGuidanceSettings200Response GetGuidanceSettings(ctx).Execute()

Get Guidance Settings



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
	resp, r, err := apiClient.GuidanceSettingsAPI.GetGuidanceSettings(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GuidanceSettingsAPI.GetGuidanceSettings``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetGuidanceSettings`: GetGuidanceSettings200Response
	fmt.Fprintf(os.Stdout, "Response from `GuidanceSettingsAPI.GetGuidanceSettings`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetGuidanceSettingsRequest struct via the builder pattern


### Return type

[**GetGuidanceSettings200Response**](GetGuidanceSettings200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateGuidanceSettings

> UpdateGuidanceSettings200Response UpdateGuidanceSettings(ctx).UpdateGuidanceSettingsRequest(updateGuidanceSettingsRequest).Execute()

Update Guidance Settings



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
	updateGuidanceSettingsRequest := *openapiclient.NewUpdateGuidanceSettingsRequest() // UpdateGuidanceSettingsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.GuidanceSettingsAPI.UpdateGuidanceSettings(context.Background()).UpdateGuidanceSettingsRequest(updateGuidanceSettingsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GuidanceSettingsAPI.UpdateGuidanceSettings``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateGuidanceSettings`: UpdateGuidanceSettings200Response
	fmt.Fprintf(os.Stdout, "Response from `GuidanceSettingsAPI.UpdateGuidanceSettings`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateGuidanceSettingsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateGuidanceSettingsRequest** | [**UpdateGuidanceSettingsRequest**](UpdateGuidanceSettingsRequest.md) |  | 

### Return type

[**UpdateGuidanceSettings200Response**](UpdateGuidanceSettings200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

