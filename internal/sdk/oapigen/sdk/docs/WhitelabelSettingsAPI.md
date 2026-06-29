# \WhitelabelSettingsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetWhitelabelImage**](WhitelabelSettingsAPI.md#GetWhitelabelImage) | **Get** /api/whitelabel-settings/images/{imageType} | Download Image
[**ListWhitelabelSettings**](WhitelabelSettingsAPI.md#ListWhitelabelSettings) | **Get** /api/whitelabel-settings | Get Whitelabel Settings
[**RemoveWhitelabelImage**](WhitelabelSettingsAPI.md#RemoveWhitelabelImage) | **Delete** /api/whitelabel-settings/images/{imageType} | Reset Image
[**UpdateWhitelabelImages**](WhitelabelSettingsAPI.md#UpdateWhitelabelImages) | **Post** /api/whitelabel-settings/images | Update Images
[**UpdateWhitelabelSettings**](WhitelabelSettingsAPI.md#UpdateWhitelabelSettings) | **Put** /api/whitelabel-settings | Update Whitelabel Settings



## GetWhitelabelImage

> *os.File GetWhitelabelImage(ctx, imageType).Execute()

Download Image



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
	imageType := "headerLogo" // string | Valid image types

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WhitelabelSettingsAPI.GetWhitelabelImage(context.Background(), imageType).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WhitelabelSettingsAPI.GetWhitelabelImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetWhitelabelImage`: *os.File
	fmt.Fprintf(os.Stdout, "Response from `WhitelabelSettingsAPI.GetWhitelabelImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**imageType** | **string** | Valid image types | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetWhitelabelImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[***os.File**](*os.File.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: image/ico, image/jpeg, image/png, image/svg+xml, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListWhitelabelSettings

> ListWhitelabelSettings200Response ListWhitelabelSettings(ctx).Execute()

Get Whitelabel Settings



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
	resp, r, err := apiClient.WhitelabelSettingsAPI.ListWhitelabelSettings(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WhitelabelSettingsAPI.ListWhitelabelSettings``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListWhitelabelSettings`: ListWhitelabelSettings200Response
	fmt.Fprintf(os.Stdout, "Response from `WhitelabelSettingsAPI.ListWhitelabelSettings`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListWhitelabelSettingsRequest struct via the builder pattern


### Return type

[**ListWhitelabelSettings200Response**](ListWhitelabelSettings200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveWhitelabelImage

> RemoveExecuteSchedules200Response RemoveWhitelabelImage(ctx, imageType).Execute()

Reset Image



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
	imageType := "headerLogo" // string | Valid image types

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WhitelabelSettingsAPI.RemoveWhitelabelImage(context.Background(), imageType).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WhitelabelSettingsAPI.RemoveWhitelabelImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveWhitelabelImage`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `WhitelabelSettingsAPI.RemoveWhitelabelImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**imageType** | **string** | Valid image types | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveWhitelabelImageRequest struct via the builder pattern


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


## UpdateWhitelabelImages

> SetApplianceSettingsMaintenanceMode200Response UpdateWhitelabelImages(ctx).HeaderLogoFile(headerLogoFile).ResetHeaderLogo(resetHeaderLogo).FooterLogoFile(footerLogoFile).ResetFooterLogo(resetFooterLogo).LoginLogoFile(loginLogoFile).ResetLoginLogo(resetLoginLogo).FaviconFile(faviconFile).ResetFaviconLogo(resetFaviconLogo).Execute()

Update Images



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
	headerLogoFile := os.NewFile(1234, "some_file") // *os.File | Header logo image file, valid image types `png|jpg|svg` (optional)
	resetHeaderLogo := true // bool | Resets header logo to default (optional)
	footerLogoFile := os.NewFile(1234, "some_file") // *os.File | Footer logo image file, valid image types `png|jpg|svg` (optional)
	resetFooterLogo := true // bool | Resets footer logo to default (optional)
	loginLogoFile := os.NewFile(1234, "some_file") // *os.File | Login logo image file, valid image types `png|jpg|svg` (optional)
	resetLoginLogo := true // bool | Resets login logo to default (optional)
	faviconFile := os.NewFile(1234, "some_file") // *os.File | Favicon image file, valid image type ico (optional)
	resetFaviconLogo := true // bool | Resets favicon logo to default (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WhitelabelSettingsAPI.UpdateWhitelabelImages(context.Background()).HeaderLogoFile(headerLogoFile).ResetHeaderLogo(resetHeaderLogo).FooterLogoFile(footerLogoFile).ResetFooterLogo(resetFooterLogo).LoginLogoFile(loginLogoFile).ResetLoginLogo(resetLoginLogo).FaviconFile(faviconFile).ResetFaviconLogo(resetFaviconLogo).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WhitelabelSettingsAPI.UpdateWhitelabelImages``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateWhitelabelImages`: SetApplianceSettingsMaintenanceMode200Response
	fmt.Fprintf(os.Stdout, "Response from `WhitelabelSettingsAPI.UpdateWhitelabelImages`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateWhitelabelImagesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **headerLogoFile** | ***os.File** | Header logo image file, valid image types &#x60;png|jpg|svg&#x60; | 
 **resetHeaderLogo** | **bool** | Resets header logo to default | 
 **footerLogoFile** | ***os.File** | Footer logo image file, valid image types &#x60;png|jpg|svg&#x60; | 
 **resetFooterLogo** | **bool** | Resets footer logo to default | 
 **loginLogoFile** | ***os.File** | Login logo image file, valid image types &#x60;png|jpg|svg&#x60; | 
 **resetLoginLogo** | **bool** | Resets login logo to default | 
 **faviconFile** | ***os.File** | Favicon image file, valid image type ico | 
 **resetFaviconLogo** | **bool** | Resets favicon logo to default | 

### Return type

[**SetApplianceSettingsMaintenanceMode200Response**](SetApplianceSettingsMaintenanceMode200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateWhitelabelSettings

> SetApplianceSettingsMaintenanceMode200Response UpdateWhitelabelSettings(ctx).UpdateWhitelabelSettingsRequest(updateWhitelabelSettingsRequest).Execute()

Update Whitelabel Settings



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
	updateWhitelabelSettingsRequest := *openapiclient.NewUpdateWhitelabelSettingsRequest() // UpdateWhitelabelSettingsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WhitelabelSettingsAPI.UpdateWhitelabelSettings(context.Background()).UpdateWhitelabelSettingsRequest(updateWhitelabelSettingsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WhitelabelSettingsAPI.UpdateWhitelabelSettings``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateWhitelabelSettings`: SetApplianceSettingsMaintenanceMode200Response
	fmt.Fprintf(os.Stdout, "Response from `WhitelabelSettingsAPI.UpdateWhitelabelSettings`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateWhitelabelSettingsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateWhitelabelSettingsRequest** | [**UpdateWhitelabelSettingsRequest**](UpdateWhitelabelSettingsRequest.md) |  | 

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

