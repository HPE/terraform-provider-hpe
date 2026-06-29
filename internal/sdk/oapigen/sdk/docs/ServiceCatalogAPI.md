# \ServiceCatalogAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddCatalogCart**](ServiceCatalogAPI.md#AddCatalogCart) | **Post** /api/catalog/checkout | Checkout Catalog Cart
[**AddCatalogCartItem**](ServiceCatalogAPI.md#AddCatalogCartItem) | **Put** /api/catalog/cart/items | Add Catalog Item to Cart
[**AddCatalogOrder**](ServiceCatalogAPI.md#AddCatalogOrder) | **Post** /api/catalog/orders | Place Catalog Order
[**DeleteCatalogCart**](ServiceCatalogAPI.md#DeleteCatalogCart) | **Delete** /api/catalog/cart | Clear Catalog Cart
[**DeleteCatalogCartItem**](ServiceCatalogAPI.md#DeleteCatalogCartItem) | **Delete** /api/catalog/cart/items/{id} | Remove a Catalog Item From Cart
[**DeleteCatalogItem**](ServiceCatalogAPI.md#DeleteCatalogItem) | **Delete** /api/catalog/items/{id} | Delete a Catalog Inventory Item
[**GetCatalogItem**](ServiceCatalogAPI.md#GetCatalogItem) | **Get** /api/catalog/items/{id} | Get a Specific Catalog Inventory Item
[**GetCatalogType**](ServiceCatalogAPI.md#GetCatalogType) | **Get** /api/catalog/types/{id} | Get a Specific Catalog Type
[**ListCatalogCart**](ServiceCatalogAPI.md#ListCatalogCart) | **Get** /api/catalog/cart | Get Catalog Cart
[**ListCatalogItems**](ServiceCatalogAPI.md#ListCatalogItems) | **Get** /api/catalog/items | List Catalog Inventory Items
[**ListCatalogTypes**](ServiceCatalogAPI.md#ListCatalogTypes) | **Get** /api/catalog/types | List Catalog Types



## AddCatalogCart

> AddCatalogCart200Response AddCatalogCart(ctx).Execute()

Checkout Catalog Cart



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
	resp, r, err := apiClient.ServiceCatalogAPI.AddCatalogCart(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServiceCatalogAPI.AddCatalogCart``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddCatalogCart`: AddCatalogCart200Response
	fmt.Fprintf(os.Stdout, "Response from `ServiceCatalogAPI.AddCatalogCart`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiAddCatalogCartRequest struct via the builder pattern


### Return type

[**AddCatalogCart200Response**](AddCatalogCart200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddCatalogCartItem

> AddCatalogCartItem200Response AddCatalogCartItem(ctx).Validate(validate).AddCatalogCartItemRequest(addCatalogCartItemRequest).Execute()

Add Catalog Item to Cart



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
	validate := true // bool | Validate Only. Use true to validate the configuration without actually placing the order or adding items to your inventory (optional) (default to false)
	addCatalogCartItemRequest := *openapiclient.NewAddCatalogCartItemRequest() // AddCatalogCartItemRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServiceCatalogAPI.AddCatalogCartItem(context.Background()).Validate(validate).AddCatalogCartItemRequest(addCatalogCartItemRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServiceCatalogAPI.AddCatalogCartItem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddCatalogCartItem`: AddCatalogCartItem200Response
	fmt.Fprintf(os.Stdout, "Response from `ServiceCatalogAPI.AddCatalogCartItem`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddCatalogCartItemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **validate** | **bool** | Validate Only. Use true to validate the configuration without actually placing the order or adding items to your inventory | [default to false]
 **addCatalogCartItemRequest** | [**AddCatalogCartItemRequest**](AddCatalogCartItemRequest.md) |  | 

### Return type

[**AddCatalogCartItem200Response**](AddCatalogCartItem200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddCatalogOrder

> AddCatalogOrder200Response AddCatalogOrder(ctx).Validate(validate).AddCatalogOrderRequest(addCatalogOrderRequest).Execute()

Place Catalog Order



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
	validate := true // bool | Validate Only. Use true to validate the configuration without actually placing the order or adding items to your inventory (optional) (default to false)
	addCatalogOrderRequest := *openapiclient.NewAddCatalogOrderRequest() // AddCatalogOrderRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServiceCatalogAPI.AddCatalogOrder(context.Background()).Validate(validate).AddCatalogOrderRequest(addCatalogOrderRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServiceCatalogAPI.AddCatalogOrder``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddCatalogOrder`: AddCatalogOrder200Response
	fmt.Fprintf(os.Stdout, "Response from `ServiceCatalogAPI.AddCatalogOrder`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddCatalogOrderRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **validate** | **bool** | Validate Only. Use true to validate the configuration without actually placing the order or adding items to your inventory | [default to false]
 **addCatalogOrderRequest** | [**AddCatalogOrderRequest**](AddCatalogOrderRequest.md) |  | 

### Return type

[**AddCatalogOrder200Response**](AddCatalogOrder200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteCatalogCart

> RemoveExecuteSchedules200Response DeleteCatalogCart(ctx).Execute()

Clear Catalog Cart



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
	resp, r, err := apiClient.ServiceCatalogAPI.DeleteCatalogCart(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServiceCatalogAPI.DeleteCatalogCart``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteCatalogCart`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `ServiceCatalogAPI.DeleteCatalogCart`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCatalogCartRequest struct via the builder pattern


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


## DeleteCatalogCartItem

> RemoveExecuteSchedules200Response DeleteCatalogCartItem(ctx, id).Execute()

Remove a Catalog Item From Cart



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
	resp, r, err := apiClient.ServiceCatalogAPI.DeleteCatalogCartItem(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServiceCatalogAPI.DeleteCatalogCartItem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteCatalogCartItem`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `ServiceCatalogAPI.DeleteCatalogCartItem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCatalogCartItemRequest struct via the builder pattern


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


## DeleteCatalogItem

> DeleteCatalogItem200Response DeleteCatalogItem(ctx, id).PreserveVolumes(preserveVolumes).KeepBackups(keepBackups).ReleaseFloatingIps(releaseFloatingIps).ReleaseEIPs(releaseEIPs).RemoveInstances(removeInstances).Force(force).Execute()

Delete a Catalog Inventory Item



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
	preserveVolumes := "on" // string | Preserve Volumes (optional) (default to "off")
	keepBackups := "on" // string | Preserve copy of backups (optional) (default to "off")
	releaseFloatingIps := "off" // string | Release Floating IPs (optional) (default to "on")
	releaseEIPs := "off" // string | Alias for releaseFloatingIps (optional) (default to "on")
	removeInstances := "off" // string | Remove Instances. Only applies to type `blueprint` (Apps) (optional) (default to "on")
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServiceCatalogAPI.DeleteCatalogItem(context.Background(), id).PreserveVolumes(preserveVolumes).KeepBackups(keepBackups).ReleaseFloatingIps(releaseFloatingIps).ReleaseEIPs(releaseEIPs).RemoveInstances(removeInstances).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServiceCatalogAPI.DeleteCatalogItem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteCatalogItem`: DeleteCatalogItem200Response
	fmt.Fprintf(os.Stdout, "Response from `ServiceCatalogAPI.DeleteCatalogItem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCatalogItemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **preserveVolumes** | **string** | Preserve Volumes | [default to &quot;off&quot;]
 **keepBackups** | **string** | Preserve copy of backups | [default to &quot;off&quot;]
 **releaseFloatingIps** | **string** | Release Floating IPs | [default to &quot;on&quot;]
 **releaseEIPs** | **string** | Alias for releaseFloatingIps | [default to &quot;on&quot;]
 **removeInstances** | **string** | Remove Instances. Only applies to type &#x60;blueprint&#x60; (Apps) | [default to &quot;on&quot;]
 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**DeleteCatalogItem200Response**](DeleteCatalogItem200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCatalogItem

> GetCatalogItem200Response GetCatalogItem(ctx, id).Execute()

Get a Specific Catalog Inventory Item



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
	resp, r, err := apiClient.ServiceCatalogAPI.GetCatalogItem(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServiceCatalogAPI.GetCatalogItem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCatalogItem`: GetCatalogItem200Response
	fmt.Fprintf(os.Stdout, "Response from `ServiceCatalogAPI.GetCatalogItem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCatalogItemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetCatalogItem200Response**](GetCatalogItem200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCatalogType

> GetCatalogType200Response GetCatalogType(ctx, id).Execute()

Get a Specific Catalog Type



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
	resp, r, err := apiClient.ServiceCatalogAPI.GetCatalogType(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServiceCatalogAPI.GetCatalogType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCatalogType`: GetCatalogType200Response
	fmt.Fprintf(os.Stdout, "Response from `ServiceCatalogAPI.GetCatalogType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCatalogTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetCatalogType200Response**](GetCatalogType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCatalogCart

> ListCatalogCart200Response ListCatalogCart(ctx).Execute()

Get Catalog Cart



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
	resp, r, err := apiClient.ServiceCatalogAPI.ListCatalogCart(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServiceCatalogAPI.ListCatalogCart``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCatalogCart`: ListCatalogCart200Response
	fmt.Fprintf(os.Stdout, "Response from `ServiceCatalogAPI.ListCatalogCart`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListCatalogCartRequest struct via the builder pattern


### Return type

[**ListCatalogCart200Response**](ListCatalogCart200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCatalogItems

> ListCatalogItems200Response ListCatalogItems(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

List Catalog Inventory Items



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServiceCatalogAPI.ListCatalogItems(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServiceCatalogAPI.ListCatalogItems``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCatalogItems`: ListCatalogItems200Response
	fmt.Fprintf(os.Stdout, "Response from `ServiceCatalogAPI.ListCatalogItems`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListCatalogItemsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListCatalogItems200Response**](ListCatalogItems200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCatalogTypes

> ListCatalogTypes200Response ListCatalogTypes(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Featured(featured).Execute()

List Catalog Types



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
	featured := false // bool | Filter by featured (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServiceCatalogAPI.ListCatalogTypes(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Featured(featured).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServiceCatalogAPI.ListCatalogTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCatalogTypes`: ListCatalogTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `ServiceCatalogAPI.ListCatalogTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListCatalogTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **featured** | **bool** | Filter by featured | 

### Return type

[**ListCatalogTypes200Response**](ListCatalogTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

