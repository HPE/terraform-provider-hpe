# \PriceSetsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddPriceSets**](PriceSetsAPI.md#AddPriceSets) | **Post** /api/price-sets | Creates a Price Set
[**DeactivatePriceSets**](PriceSetsAPI.md#DeactivatePriceSets) | **Put** /api/price-sets/{id}/deactivate | Deactivates a Price Set
[**GetPriceSets**](PriceSetsAPI.md#GetPriceSets) | **Get** /api/price-sets/{id} | Retrieves a Specific Price Set
[**ListPriceSets**](PriceSetsAPI.md#ListPriceSets) | **Get** /api/price-sets | Retrieves all Price Sets
[**UpdatePriceSets**](PriceSetsAPI.md#UpdatePriceSets) | **Put** /api/price-sets/{id} | Updates a Price Set



## AddPriceSets

> AddPriceSets200Response AddPriceSets(ctx).AddPriceSetsRequest(addPriceSetsRequest).Execute()

Creates a Price Set



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
	addPriceSetsRequest := *openapiclient.NewAddPriceSetsRequest(*openapiclient.NewAddPriceSetsRequestPriceSet("testName", "priceSet1", "PriceUnit_example", "Type_example")) // AddPriceSetsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PriceSetsAPI.AddPriceSets(context.Background()).AddPriceSetsRequest(addPriceSetsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PriceSetsAPI.AddPriceSets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddPriceSets`: AddPriceSets200Response
	fmt.Fprintf(os.Stdout, "Response from `PriceSetsAPI.AddPriceSets`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddPriceSetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addPriceSetsRequest** | [**AddPriceSetsRequest**](AddPriceSetsRequest.md) |  | 

### Return type

[**AddPriceSets200Response**](AddPriceSets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeactivatePriceSets

> RemoveExecuteSchedules200Response DeactivatePriceSets(ctx, id).Execute()

Deactivates a Price Set



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
	resp, r, err := apiClient.PriceSetsAPI.DeactivatePriceSets(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PriceSetsAPI.DeactivatePriceSets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeactivatePriceSets`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `PriceSetsAPI.DeactivatePriceSets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeactivatePriceSetsRequest struct via the builder pattern


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


## GetPriceSets

> GetPriceSets200Response GetPriceSets(ctx, id).Execute()

Retrieves a Specific Price Set



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
	resp, r, err := apiClient.PriceSetsAPI.GetPriceSets(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PriceSetsAPI.GetPriceSets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPriceSets`: GetPriceSets200Response
	fmt.Fprintf(os.Stdout, "Response from `PriceSetsAPI.GetPriceSets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPriceSetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetPriceSets200Response**](GetPriceSets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPriceSets

> ListPriceSets200Response ListPriceSets(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).IncludeInactive(includeInactive).Type_(type_).Execute()

Retrieves all Price Sets



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
	includeInactive := true // bool | If true, include inactive prices in the results (optional)
	type_ := "type__example" // string | Filter by type code (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PriceSetsAPI.ListPriceSets(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).IncludeInactive(includeInactive).Type_(type_).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PriceSetsAPI.ListPriceSets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPriceSets`: ListPriceSets200Response
	fmt.Fprintf(os.Stdout, "Response from `PriceSetsAPI.ListPriceSets`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListPriceSetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **includeInactive** | **bool** | If true, include inactive prices in the results | 
 **type_** | **string** | Filter by type code | 

### Return type

[**ListPriceSets200Response**](ListPriceSets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdatePriceSets

> UpdatePriceSets200Response UpdatePriceSets(ctx, id).UpdatePriceSetsRequest(updatePriceSetsRequest).Execute()

Updates a Price Set



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
	updatePriceSetsRequest := *openapiclient.NewUpdatePriceSetsRequest(*openapiclient.NewUpdatePriceSetsRequestPriceSet()) // UpdatePriceSetsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PriceSetsAPI.UpdatePriceSets(context.Background(), id).UpdatePriceSetsRequest(updatePriceSetsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PriceSetsAPI.UpdatePriceSets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdatePriceSets`: UpdatePriceSets200Response
	fmt.Fprintf(os.Stdout, "Response from `PriceSetsAPI.UpdatePriceSets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdatePriceSetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updatePriceSetsRequest** | [**UpdatePriceSetsRequest**](UpdatePriceSetsRequest.md) |  | 

### Return type

[**UpdatePriceSets200Response**](UpdatePriceSets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

