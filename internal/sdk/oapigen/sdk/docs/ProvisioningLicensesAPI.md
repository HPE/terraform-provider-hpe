# \ProvisioningLicensesAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddProvisioningLicense**](ProvisioningLicensesAPI.md#AddProvisioningLicense) | **Post** /api/provisioning-licenses | Create a License
[**GetProvisioningLicense**](ProvisioningLicensesAPI.md#GetProvisioningLicense) | **Get** /api/provisioning-licenses/{id} | Get a Specific License
[**GetProvisioningLicenseReservations**](ProvisioningLicensesAPI.md#GetProvisioningLicenseReservations) | **Get** /api/provisioning-licenses/{id}/reservations | Get Reservations for Specific License
[**ListProvisioningLicenses**](ProvisioningLicensesAPI.md#ListProvisioningLicenses) | **Get** /api/provisioning-licenses | Get All Licenses
[**RemoveProvisioningLicense**](ProvisioningLicensesAPI.md#RemoveProvisioningLicense) | **Delete** /api/provisioning-licenses/{id} | Delete a License
[**UpdateProvisioningLicense**](ProvisioningLicensesAPI.md#UpdateProvisioningLicense) | **Put** /api/provisioning-licenses/{id} | Update a License



## AddProvisioningLicense

> AddProvisioningLicense200Response AddProvisioningLicense(ctx).AddProvisioningLicenseRequest(addProvisioningLicenseRequest).Execute()

Create a License



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
	addProvisioningLicenseRequest := *openapiclient.NewAddProvisioningLicenseRequest() // AddProvisioningLicenseRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProvisioningLicensesAPI.AddProvisioningLicense(context.Background()).AddProvisioningLicenseRequest(addProvisioningLicenseRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProvisioningLicensesAPI.AddProvisioningLicense``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddProvisioningLicense`: AddProvisioningLicense200Response
	fmt.Fprintf(os.Stdout, "Response from `ProvisioningLicensesAPI.AddProvisioningLicense`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddProvisioningLicenseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addProvisioningLicenseRequest** | [**AddProvisioningLicenseRequest**](AddProvisioningLicenseRequest.md) |  | 

### Return type

[**AddProvisioningLicense200Response**](AddProvisioningLicense200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetProvisioningLicense

> GetProvisioningLicense200Response GetProvisioningLicense(ctx, id).Execute()

Get a Specific License



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
	resp, r, err := apiClient.ProvisioningLicensesAPI.GetProvisioningLicense(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProvisioningLicensesAPI.GetProvisioningLicense``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetProvisioningLicense`: GetProvisioningLicense200Response
	fmt.Fprintf(os.Stdout, "Response from `ProvisioningLicensesAPI.GetProvisioningLicense`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetProvisioningLicenseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetProvisioningLicense200Response**](GetProvisioningLicense200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetProvisioningLicenseReservations

> GetProvisioningLicenseReservations200Response GetProvisioningLicenseReservations(ctx, id).Execute()

Get Reservations for Specific License



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
	resp, r, err := apiClient.ProvisioningLicensesAPI.GetProvisioningLicenseReservations(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProvisioningLicensesAPI.GetProvisioningLicenseReservations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetProvisioningLicenseReservations`: GetProvisioningLicenseReservations200Response
	fmt.Fprintf(os.Stdout, "Response from `ProvisioningLicensesAPI.GetProvisioningLicenseReservations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetProvisioningLicenseReservationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetProvisioningLicenseReservations200Response**](GetProvisioningLicenseReservations200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListProvisioningLicenses

> ListProvisioningLicenses200Response ListProvisioningLicenses(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).LicenseType(licenseType).LicenseVersion(licenseVersion).OrgName(orgName).FullName(fullName).Execute()

Get All Licenses



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
	licenseType := "win" // string | If specified will return an exact match on licenseType code (optional)
	licenseVersion := "2012 R2 Standard" // string | If specified will return an exact match on licenseVersion (optional)
	orgName := "Acme Motors" // string | If specified will return an exact match on orgName (optional)
	fullName := "Bugs Bunny" // string | If specified will return an exact match on fullName (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProvisioningLicensesAPI.ListProvisioningLicenses(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).LicenseType(licenseType).LicenseVersion(licenseVersion).OrgName(orgName).FullName(fullName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProvisioningLicensesAPI.ListProvisioningLicenses``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListProvisioningLicenses`: ListProvisioningLicenses200Response
	fmt.Fprintf(os.Stdout, "Response from `ProvisioningLicensesAPI.ListProvisioningLicenses`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListProvisioningLicensesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **licenseType** | **string** | If specified will return an exact match on licenseType code | 
 **licenseVersion** | **string** | If specified will return an exact match on licenseVersion | 
 **orgName** | **string** | If specified will return an exact match on orgName | 
 **fullName** | **string** | If specified will return an exact match on fullName | 

### Return type

[**ListProvisioningLicenses200Response**](ListProvisioningLicenses200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveProvisioningLicense

> RemoveExecuteSchedules200Response RemoveProvisioningLicense(ctx, id).Execute()

Delete a License



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
	resp, r, err := apiClient.ProvisioningLicensesAPI.RemoveProvisioningLicense(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProvisioningLicensesAPI.RemoveProvisioningLicense``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveProvisioningLicense`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `ProvisioningLicensesAPI.RemoveProvisioningLicense`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveProvisioningLicenseRequest struct via the builder pattern


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


## UpdateProvisioningLicense

> UpdateProvisioningLicense200Response UpdateProvisioningLicense(ctx, id).UpdateProvisioningLicenseRequest(updateProvisioningLicenseRequest).Execute()

Update a License



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
	updateProvisioningLicenseRequest := *openapiclient.NewUpdateProvisioningLicenseRequest() // UpdateProvisioningLicenseRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProvisioningLicensesAPI.UpdateProvisioningLicense(context.Background(), id).UpdateProvisioningLicenseRequest(updateProvisioningLicenseRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProvisioningLicensesAPI.UpdateProvisioningLicense``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateProvisioningLicense`: UpdateProvisioningLicense200Response
	fmt.Fprintf(os.Stdout, "Response from `ProvisioningLicensesAPI.UpdateProvisioningLicense`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateProvisioningLicenseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateProvisioningLicenseRequest** | [**UpdateProvisioningLicenseRequest**](UpdateProvisioningLicenseRequest.md) |  | 

### Return type

[**UpdateProvisioningLicense200Response**](UpdateProvisioningLicense200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

