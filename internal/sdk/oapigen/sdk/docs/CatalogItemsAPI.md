# \CatalogItemsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddCatalogItemType**](CatalogItemsAPI.md#AddCatalogItemType) | **Post** /api/catalog-item-types | Create a Catalog Item Type
[**GetCatalogItemType**](CatalogItemsAPI.md#GetCatalogItemType) | **Get** /api/catalog-item-types/{id} | Get a Specific Catalog Item Type
[**ListCatalogItemTypes**](CatalogItemsAPI.md#ListCatalogItemTypes) | **Get** /api/catalog-item-types | Get All Catalog Item Types
[**RemoveCatalogItemType**](CatalogItemsAPI.md#RemoveCatalogItemType) | **Delete** /api/catalog-item-types/{id} | Delete a Catalog Item Type
[**UpdateCatalogItemType**](CatalogItemsAPI.md#UpdateCatalogItemType) | **Put** /api/catalog-item-types/{id} | Update a Catalog Item Type
[**UpdateCatalogItemTypeLogo**](CatalogItemsAPI.md#UpdateCatalogItemTypeLogo) | **Put** /api/catalog-item-types/{id}/update-logo | Update Logo For Catalog Item Type



## AddCatalogItemType

> AddCatalogItemType200Response AddCatalogItemType(ctx).AddCatalogItemTypeRequest(addCatalogItemTypeRequest).Execute()

Create a Catalog Item Type



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
	addCatalogItemTypeRequest := *openapiclient.NewAddCatalogItemTypeRequest() // AddCatalogItemTypeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CatalogItemsAPI.AddCatalogItemType(context.Background()).AddCatalogItemTypeRequest(addCatalogItemTypeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CatalogItemsAPI.AddCatalogItemType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddCatalogItemType`: AddCatalogItemType200Response
	fmt.Fprintf(os.Stdout, "Response from `CatalogItemsAPI.AddCatalogItemType`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddCatalogItemTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addCatalogItemTypeRequest** | [**AddCatalogItemTypeRequest**](AddCatalogItemTypeRequest.md) |  | 

### Return type

[**AddCatalogItemType200Response**](AddCatalogItemType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json, multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCatalogItemType

> GetCatalogItemType200Response GetCatalogItemType(ctx, id).Execute()

Get a Specific Catalog Item Type



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
	resp, r, err := apiClient.CatalogItemsAPI.GetCatalogItemType(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CatalogItemsAPI.GetCatalogItemType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCatalogItemType`: GetCatalogItemType200Response
	fmt.Fprintf(os.Stdout, "Response from `CatalogItemsAPI.GetCatalogItemType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCatalogItemTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetCatalogItemType200Response**](GetCatalogItemType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCatalogItemTypes

> ListCatalogItemTypes200Response ListCatalogItemTypes(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Description(description).Enabled(enabled).Featured(featured).Labels(labels).AllLabels(allLabels).Code(code).Execute()

Get All Catalog Item Types



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
	description := "The desription of my cool object" // string | Filter by description, wildcard may be specified as %. eg. `example-%` (optional)
	enabled := false // bool | Filter by enabled (optional)
	featured := false // bool | Filter by featured (optional)
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)
	code := "azr" // string | If specified will return an exact match on code (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CatalogItemsAPI.ListCatalogItemTypes(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Description(description).Enabled(enabled).Featured(featured).Labels(labels).AllLabels(allLabels).Code(code).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CatalogItemsAPI.ListCatalogItemTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCatalogItemTypes`: ListCatalogItemTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `CatalogItemsAPI.ListCatalogItemTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListCatalogItemTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **description** | **string** | Filter by description, wildcard may be specified as %. eg. &#x60;example-%&#x60; | 
 **enabled** | **bool** | Filter by enabled | 
 **featured** | **bool** | Filter by featured | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 
 **code** | **string** | If specified will return an exact match on code | 

### Return type

[**ListCatalogItemTypes200Response**](ListCatalogItemTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveCatalogItemType

> RemoveExecuteSchedules200Response RemoveCatalogItemType(ctx, id).Execute()

Delete a Catalog Item Type



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
	resp, r, err := apiClient.CatalogItemsAPI.RemoveCatalogItemType(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CatalogItemsAPI.RemoveCatalogItemType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveCatalogItemType`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `CatalogItemsAPI.RemoveCatalogItemType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveCatalogItemTypeRequest struct via the builder pattern


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


## UpdateCatalogItemType

> UpdateCatalogItemType200Response UpdateCatalogItemType(ctx, id).UpdateCatalogItemTypeRequest(updateCatalogItemTypeRequest).Execute()

Update a Catalog Item Type



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
	updateCatalogItemTypeRequest := *openapiclient.NewUpdateCatalogItemTypeRequest() // UpdateCatalogItemTypeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CatalogItemsAPI.UpdateCatalogItemType(context.Background(), id).UpdateCatalogItemTypeRequest(updateCatalogItemTypeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CatalogItemsAPI.UpdateCatalogItemType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateCatalogItemType`: UpdateCatalogItemType200Response
	fmt.Fprintf(os.Stdout, "Response from `CatalogItemsAPI.UpdateCatalogItemType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateCatalogItemTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateCatalogItemTypeRequest** | [**UpdateCatalogItemTypeRequest**](UpdateCatalogItemTypeRequest.md) |  | 

### Return type

[**UpdateCatalogItemType200Response**](UpdateCatalogItemType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json, multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateCatalogItemTypeLogo

> UpdateCatalogItemTypeLogo200Response UpdateCatalogItemTypeLogo(ctx, id).CatalogItemTypeLogo(catalogItemTypeLogo).CatalogItemTypeDarkLogo(catalogItemTypeDarkLogo).Execute()

Update Logo For Catalog Item Type



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
	catalogItemTypeLogo := os.NewFile(1234, "some_file") // *os.File | Logo File png,jpg,svg (optional)
	catalogItemTypeDarkLogo := os.NewFile(1234, "some_file") // *os.File | Dark Logo File png,jpg,svg (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CatalogItemsAPI.UpdateCatalogItemTypeLogo(context.Background(), id).CatalogItemTypeLogo(catalogItemTypeLogo).CatalogItemTypeDarkLogo(catalogItemTypeDarkLogo).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CatalogItemsAPI.UpdateCatalogItemTypeLogo``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateCatalogItemTypeLogo`: UpdateCatalogItemTypeLogo200Response
	fmt.Fprintf(os.Stdout, "Response from `CatalogItemsAPI.UpdateCatalogItemTypeLogo`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateCatalogItemTypeLogoRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **catalogItemTypeLogo** | ***os.File** | Logo File png,jpg,svg | 
 **catalogItemTypeDarkLogo** | ***os.File** | Dark Logo File png,jpg,svg | 

### Return type

[**UpdateCatalogItemTypeLogo200Response**](UpdateCatalogItemTypeLogo200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

