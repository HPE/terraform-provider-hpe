# \LicenseAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetLicense**](LicenseAPI.md#GetLicense) | **Get** /api/license | Get license
[**InstallLicense**](LicenseAPI.md#InstallLicense) | **Post** /api/license | Install license key
[**TestLicense**](LicenseAPI.md#TestLicense) | **Post** /api/license/test | Test license key
[**UninstallLicense**](LicenseAPI.md#UninstallLicense) | **Delete** /api/license | Uninstall license key



## GetLicense

> GetLicense200Response GetLicense(ctx).Execute()

Get license



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
	resp, r, err := apiClient.LicenseAPI.GetLicense(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LicenseAPI.GetLicense``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLicense`: GetLicense200Response
	fmt.Fprintf(os.Stdout, "Response from `LicenseAPI.GetLicense`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetLicenseRequest struct via the builder pattern


### Return type

[**GetLicense200Response**](GetLicense200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InstallLicense

> InstallLicense200Response InstallLicense(ctx).InstallLicenseRequest(installLicenseRequest).Execute()

Install license key



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
	installLicenseRequest := *openapiclient.NewInstallLicenseRequest("License_example") // InstallLicenseRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LicenseAPI.InstallLicense(context.Background()).InstallLicenseRequest(installLicenseRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LicenseAPI.InstallLicense``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InstallLicense`: InstallLicense200Response
	fmt.Fprintf(os.Stdout, "Response from `LicenseAPI.InstallLicense`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInstallLicenseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **installLicenseRequest** | [**InstallLicenseRequest**](InstallLicenseRequest.md) |  | 

### Return type

[**InstallLicense200Response**](InstallLicense200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TestLicense

> TestLicense200Response TestLicense(ctx).TestLicenseRequest(testLicenseRequest).Execute()

Test license key



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
	testLicenseRequest := *openapiclient.NewTestLicenseRequest("License_example") // TestLicenseRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LicenseAPI.TestLicense(context.Background()).TestLicenseRequest(testLicenseRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LicenseAPI.TestLicense``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TestLicense`: TestLicense200Response
	fmt.Fprintf(os.Stdout, "Response from `LicenseAPI.TestLicense`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTestLicenseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **testLicenseRequest** | [**TestLicenseRequest**](TestLicenseRequest.md) |  | 

### Return type

[**TestLicense200Response**](TestLicense200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UninstallLicense

> UpdateApprovalsItem200Response UninstallLicense(ctx).KeyId(keyId).Execute()

Uninstall license key



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
	keyId := int64(789) // int64 | Key ID (only the first 8 characters are required to identify license to uninstall). (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LicenseAPI.UninstallLicense(context.Background()).KeyId(keyId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LicenseAPI.UninstallLicense``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UninstallLicense`: UpdateApprovalsItem200Response
	fmt.Fprintf(os.Stdout, "Response from `LicenseAPI.UninstallLicense`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUninstallLicenseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **keyId** | **int64** | Key ID (only the first 8 characters are required to identify license to uninstall). | 

### Return type

[**UpdateApprovalsItem200Response**](UpdateApprovalsItem200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

