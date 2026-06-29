# \ProvisioningSettingsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ListProvisioningSettings**](ProvisioningSettingsAPI.md#ListProvisioningSettings) | **Get** /api/provisioning-settings | Get Provisioning Settings
[**UpdateProvisioningSettings**](ProvisioningSettingsAPI.md#UpdateProvisioningSettings) | **Put** /api/provisioning-settings | Update Provisioning Settings



## ListProvisioningSettings

> ListProvisioningSettings200Response ListProvisioningSettings(ctx).Execute()

Get Provisioning Settings



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
	resp, r, err := apiClient.ProvisioningSettingsAPI.ListProvisioningSettings(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProvisioningSettingsAPI.ListProvisioningSettings``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListProvisioningSettings`: ListProvisioningSettings200Response
	fmt.Fprintf(os.Stdout, "Response from `ProvisioningSettingsAPI.ListProvisioningSettings`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListProvisioningSettingsRequest struct via the builder pattern


### Return type

[**ListProvisioningSettings200Response**](ListProvisioningSettings200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateProvisioningSettings

> SetApplianceSettingsMaintenanceMode200Response UpdateProvisioningSettings(ctx).UpdateProvisioningSettingsRequest(updateProvisioningSettingsRequest).Execute()

Update Provisioning Settings



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
	updateProvisioningSettingsRequest := *openapiclient.NewUpdateProvisioningSettingsRequest() // UpdateProvisioningSettingsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProvisioningSettingsAPI.UpdateProvisioningSettings(context.Background()).UpdateProvisioningSettingsRequest(updateProvisioningSettingsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProvisioningSettingsAPI.UpdateProvisioningSettings``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateProvisioningSettings`: SetApplianceSettingsMaintenanceMode200Response
	fmt.Fprintf(os.Stdout, "Response from `ProvisioningSettingsAPI.UpdateProvisioningSettings`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateProvisioningSettingsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateProvisioningSettingsRequest** | [**UpdateProvisioningSettingsRequest**](UpdateProvisioningSettingsRequest.md) |  | 

### Return type

[**SetApplianceSettingsMaintenanceMode200Response**](SetApplianceSettingsMaintenanceMode200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

