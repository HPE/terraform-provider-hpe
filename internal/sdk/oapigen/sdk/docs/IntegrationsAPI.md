# \IntegrationsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddIntegrationSnowObjects**](IntegrationsAPI.md#AddIntegrationSnowObjects) | **Post** /api/integrations/{id}/objects | Creates an Exposed ServiceNow Catalog Item
[**AddIntegrations**](IntegrationsAPI.md#AddIntegrations) | **Post** /api/integrations | Creates an Integration
[**GetIntegrationInventory**](IntegrationsAPI.md#GetIntegrationInventory) | **Get** /api/integrations/{id}/inventory/{inventoryId} | Get a Specific Integration Inventory
[**GetIntegrationObjects**](IntegrationsAPI.md#GetIntegrationObjects) | **Get** /api/integrations/{id}/objects/{objectId} | Get a Specific ServiceNow Integration Object
[**GetIntegrationTypeOptionTypes**](IntegrationsAPI.md#GetIntegrationTypeOptionTypes) | **Get** /api/integration-types/{id}/option-types | Retrieves a Option Types for a Specific Integration Type
[**GetIntegrationTypes**](IntegrationsAPI.md#GetIntegrationTypes) | **Get** /api/integration-types/{id} | Retrieves a Specific Integration Type
[**GetIntegrations**](IntegrationsAPI.md#GetIntegrations) | **Get** /api/integrations/{id} | Retrieves a Specific Integration
[**ListIntegrationInventory**](IntegrationsAPI.md#ListIntegrationInventory) | **Get** /api/integrations/{id}/inventory | Get All Integration Inventory
[**ListIntegrationObjects**](IntegrationsAPI.md#ListIntegrationObjects) | **Get** /api/integrations/{id}/objects | Get ServiceNow Integration Objects
[**ListIntegrationTypes**](IntegrationsAPI.md#ListIntegrationTypes) | **Get** /api/integration-types | Retrieves all Integration Types
[**ListIntegrations**](IntegrationsAPI.md#ListIntegrations) | **Get** /api/integrations | Retrieves all Integrations
[**RefreshIntegrations**](IntegrationsAPI.md#RefreshIntegrations) | **Post** /api/integrations/{id}/refresh | Refresh an Integration
[**RemoveIntegrationObjects**](IntegrationsAPI.md#RemoveIntegrationObjects) | **Delete** /api/integrations/{id}/objects/{objectId} | Deletes a ServiceNow Integration object
[**RemoveIntegrations**](IntegrationsAPI.md#RemoveIntegrations) | **Delete** /api/integrations/{id} | Deletes an Integration
[**UpdateIntegrationInventory**](IntegrationsAPI.md#UpdateIntegrationInventory) | **Put** /api/integrations/{id}/inventory/{inventoryId} | Updating an Integration Inventory
[**UpdateIntegrations**](IntegrationsAPI.md#UpdateIntegrations) | **Put** /api/integrations/{id} | Updates an Integration



## AddIntegrationSnowObjects

> AddIntegrationSnowObjects200Response AddIntegrationSnowObjects(ctx, id).AddIntegrationSnowObjectsRequest(addIntegrationSnowObjectsRequest).Execute()

Creates an Exposed ServiceNow Catalog Item



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
	addIntegrationSnowObjectsRequest := *openapiclient.NewAddIntegrationSnowObjectsRequest() // AddIntegrationSnowObjectsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.AddIntegrationSnowObjects(context.Background(), id).AddIntegrationSnowObjectsRequest(addIntegrationSnowObjectsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.AddIntegrationSnowObjects``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddIntegrationSnowObjects`: AddIntegrationSnowObjects200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.AddIntegrationSnowObjects`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddIntegrationSnowObjectsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addIntegrationSnowObjectsRequest** | [**AddIntegrationSnowObjectsRequest**](AddIntegrationSnowObjectsRequest.md) |  | 

### Return type

[**AddIntegrationSnowObjects200Response**](AddIntegrationSnowObjects200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddIntegrations

> AddIntegrations200Response AddIntegrations(ctx).AddIntegrationsRequest(addIntegrationsRequest).Execute()

Creates an Integration



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
	addIntegrationsRequest := openapiclient.addIntegrations_request{AddIntegrationsRequestOneOf: openapiclient.NewAddIntegrationsRequestOneOf(*openapiclient.NewAddIntegrationsRequestOneOfIntegration("Sample Integration", "Type_example"))} // AddIntegrationsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.AddIntegrations(context.Background()).AddIntegrationsRequest(addIntegrationsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.AddIntegrations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddIntegrations`: AddIntegrations200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.AddIntegrations`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddIntegrationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addIntegrationsRequest** | [**AddIntegrationsRequest**](AddIntegrationsRequest.md) |  | 

### Return type

[**AddIntegrations200Response**](AddIntegrations200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetIntegrationInventory

> GetIntegrationInventory200Response GetIntegrationInventory(ctx, id, inventoryId).Execute()

Get a Specific Integration Inventory



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
	id := int64(9) // int64 | Morpheus ID of the Integration
	inventoryId := int64(1) // int64 | Morpheus ID of the Integration Inventory

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.GetIntegrationInventory(context.Background(), id, inventoryId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.GetIntegrationInventory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetIntegrationInventory`: GetIntegrationInventory200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.GetIntegrationInventory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Integration | 
**inventoryId** | **int64** | Morpheus ID of the Integration Inventory | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetIntegrationInventoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetIntegrationInventory200Response**](GetIntegrationInventory200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetIntegrationObjects

> GetIntegrationObjects200Response GetIntegrationObjects(ctx, id, objectId).Execute()

Get a Specific ServiceNow Integration Object



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
	objectId := int64(122) // int64 | Morpheus ID of the Object being created or referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.GetIntegrationObjects(context.Background(), id, objectId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.GetIntegrationObjects``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetIntegrationObjects`: GetIntegrationObjects200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.GetIntegrationObjects`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 
**objectId** | **int64** | Morpheus ID of the Object being created or referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetIntegrationObjectsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetIntegrationObjects200Response**](GetIntegrationObjects200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetIntegrationTypeOptionTypes

> GetIntegrationTypeOptionTypes200Response GetIntegrationTypeOptionTypes(ctx, id).Execute()

Retrieves a Option Types for a Specific Integration Type



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
	resp, r, err := apiClient.IntegrationsAPI.GetIntegrationTypeOptionTypes(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.GetIntegrationTypeOptionTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetIntegrationTypeOptionTypes`: GetIntegrationTypeOptionTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.GetIntegrationTypeOptionTypes`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetIntegrationTypeOptionTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetIntegrationTypeOptionTypes200Response**](GetIntegrationTypeOptionTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetIntegrationTypes

> GetIntegrationTypes200Response GetIntegrationTypes(ctx, id).Optiontypes(optiontypes).Execute()

Retrieves a Specific Integration Type



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
	optiontypes := true // bool | Pass `true` to include optionTypes in the response for each integration type (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.GetIntegrationTypes(context.Background(), id).Optiontypes(optiontypes).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.GetIntegrationTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetIntegrationTypes`: GetIntegrationTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.GetIntegrationTypes`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetIntegrationTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **optiontypes** | **bool** | Pass &#x60;true&#x60; to include optionTypes in the response for each integration type | [default to false]

### Return type

[**GetIntegrationTypes200Response**](GetIntegrationTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetIntegrations

> GetIntegrations200Response GetIntegrations(ctx, id).Execute()

Retrieves a Specific Integration



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
	resp, r, err := apiClient.IntegrationsAPI.GetIntegrations(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.GetIntegrations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetIntegrations`: GetIntegrations200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.GetIntegrations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetIntegrationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetIntegrations200Response**](GetIntegrations200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListIntegrationInventory

> ListIntegrationInventory200Response ListIntegrationInventory(ctx, id).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Get All Integration Inventory



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
	id := int64(9) // int64 | Morpheus ID of the Integration
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.ListIntegrationInventory(context.Background(), id).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.ListIntegrationInventory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListIntegrationInventory`: ListIntegrationInventory200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.ListIntegrationInventory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Integration | 

### Other Parameters

Other parameters are passed through a pointer to a apiListIntegrationInventoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListIntegrationInventory200Response**](ListIntegrationInventory200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListIntegrationObjects

> ListIntegrationObjects200Response ListIntegrationObjects(ctx, id).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Value(value).RefId(refId).Execute()

Get ServiceNow Integration Objects



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
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	value := "value_example" // string | The type of data being stored, string or object. The data type depends on the cypher mount being used. Most mounts use string as their data type, but secret uses object by default.  You can store a string instead by passing `type=string`. This means the data value returned by the API will be a string instead of an object.  (optional)
	refId := int64(3) // int64 | If specified will return an exact match on refId (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.ListIntegrationObjects(context.Background(), id).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Value(value).RefId(refId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.ListIntegrationObjects``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListIntegrationObjects`: ListIntegrationObjects200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.ListIntegrationObjects`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiListIntegrationObjectsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **value** | **string** | The type of data being stored, string or object. The data type depends on the cypher mount being used. Most mounts use string as their data type, but secret uses object by default.  You can store a string instead by passing &#x60;type&#x3D;string&#x60;. This means the data value returned by the API will be a string instead of an object.  | 
 **refId** | **int64** | If specified will return an exact match on refId | 

### Return type

[**ListIntegrationObjects200Response**](ListIntegrationObjects200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListIntegrationTypes

> ListIntegrationTypes200Response ListIntegrationTypes(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).Optiontypes(optiontypes).Description(description).Category(category).Creatable(creatable).Enabled(enabled).HasCMDB(hasCMDB).HasCMDBDiscovery(hasCMDBDiscovery).HasCM(hasCM).HasDNS(hasDNS).HasApprovals(hasApprovals).IsPlugin(isPlugin).Execute()

Retrieves all Integration Types



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
	code := "azr" // string | If specified will return an exact match on code (optional)
	optiontypes := true // bool | Pass `true` to include optionTypes in the response for each integration type (optional) (default to false)
	description := "The desription of my cool object" // string | Filter by description, wildcard may be specified as %. eg. `example-%` (optional)
	category := "category_example" // string | If specified will return an exact match on category (optional)
	creatable := true // bool | Filter by creatable (optional)
	enabled := false // bool | Filter by enabled (optional)
	hasCMDB := true // bool | Filter by hasCMDB (optional)
	hasCMDBDiscovery := true // bool | Filter by hasCMDBDiscovery (optional)
	hasCM := true // bool | Filter by hasCM (optional)
	hasDNS := true // bool | Filter by hasDNS (optional)
	hasApprovals := true // bool | Filter by hasApprovals (optional)
	isPlugin := true // bool | Filter by isPlugin (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.ListIntegrationTypes(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).Optiontypes(optiontypes).Description(description).Category(category).Creatable(creatable).Enabled(enabled).HasCMDB(hasCMDB).HasCMDBDiscovery(hasCMDBDiscovery).HasCM(hasCM).HasDNS(hasDNS).HasApprovals(hasApprovals).IsPlugin(isPlugin).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.ListIntegrationTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListIntegrationTypes`: ListIntegrationTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.ListIntegrationTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListIntegrationTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **optiontypes** | **bool** | Pass &#x60;true&#x60; to include optionTypes in the response for each integration type | [default to false]
 **description** | **string** | Filter by description, wildcard may be specified as %. eg. &#x60;example-%&#x60; | 
 **category** | **string** | If specified will return an exact match on category | 
 **creatable** | **bool** | Filter by creatable | 
 **enabled** | **bool** | Filter by enabled | 
 **hasCMDB** | **bool** | Filter by hasCMDB | 
 **hasCMDBDiscovery** | **bool** | Filter by hasCMDBDiscovery | 
 **hasCM** | **bool** | Filter by hasCM | 
 **hasDNS** | **bool** | Filter by hasDNS | 
 **hasApprovals** | **bool** | Filter by hasApprovals | 
 **isPlugin** | **bool** | Filter by isPlugin | 

### Return type

[**ListIntegrationTypes200Response**](ListIntegrationTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListIntegrations

> ListIntegrations200Response ListIntegrations(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Id(id).Url(url).Host(host).Port(port).Username(username).Version(version).WindowsVersion(windowsVersion).Status(status).Objects(objects).Execute()

Retrieves all Integrations



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
	id := int64(7) // int64 | Morpheus ID of the Object being created or referenced (optional)
	url := "https://example.com/testimage.ovf" // string | Download the file from a remote url. This can be used instead of uploading a local file. (optional)
	host := "host_example" // string | Filter by integration Host (optional)
	port := "port_example" // string | Filter by integration Port (optional)
	username := "administrator" // string | Username (optional)
	version := int64(5) // int64 | Filter by version number (userVersion) (optional)
	windowsVersion := "windowsVersion_example" // string | Filter by integration Windows Version (optional)
	status := "running" // string | The instance status for filtering. (optional)
	objects := true // bool | Include `objects=true` to return the Integration Objects for each integration.  Available in api version 5.2.8/5.3.2.  (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.ListIntegrations(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Id(id).Url(url).Host(host).Port(port).Username(username).Version(version).WindowsVersion(windowsVersion).Status(status).Objects(objects).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.ListIntegrations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListIntegrations`: ListIntegrations200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.ListIntegrations`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListIntegrationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **id** | **int64** | Morpheus ID of the Object being created or referenced | 
 **url** | **string** | Download the file from a remote url. This can be used instead of uploading a local file. | 
 **host** | **string** | Filter by integration Host | 
 **port** | **string** | Filter by integration Port | 
 **username** | **string** | Username | 
 **version** | **int64** | Filter by version number (userVersion) | 
 **windowsVersion** | **string** | Filter by integration Windows Version | 
 **status** | **string** | The instance status for filtering. | 
 **objects** | **bool** | Include &#x60;objects&#x3D;true&#x60; to return the Integration Objects for each integration.  Available in api version 5.2.8/5.3.2.  | [default to false]

### Return type

[**ListIntegrations200Response**](ListIntegrations200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RefreshIntegrations

> RefreshIntegrations200Response RefreshIntegrations(ctx, id).Execute()

Refresh an Integration



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
	resp, r, err := apiClient.IntegrationsAPI.RefreshIntegrations(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.RefreshIntegrations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RefreshIntegrations`: RefreshIntegrations200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.RefreshIntegrations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRefreshIntegrationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**RefreshIntegrations200Response**](RefreshIntegrations200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveIntegrationObjects

> RemoveExecuteSchedules200Response RemoveIntegrationObjects(ctx, id, objectId).Execute()

Deletes a ServiceNow Integration object



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
	objectId := int64(122) // int64 | Morpheus ID of the Object being created or referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.RemoveIntegrationObjects(context.Background(), id, objectId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.RemoveIntegrationObjects``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveIntegrationObjects`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.RemoveIntegrationObjects`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 
**objectId** | **int64** | Morpheus ID of the Object being created or referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveIntegrationObjectsRequest struct via the builder pattern


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


## RemoveIntegrations

> RemoveExecuteSchedules200Response RemoveIntegrations(ctx, id).Execute()

Deletes an Integration



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
	resp, r, err := apiClient.IntegrationsAPI.RemoveIntegrations(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.RemoveIntegrations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveIntegrations`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.RemoveIntegrations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveIntegrationsRequest struct via the builder pattern


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


## UpdateIntegrationInventory

> UpdateIntegrationInventory200Response UpdateIntegrationInventory(ctx, id, inventoryId).UpdateIntegrationInventoryRequest(updateIntegrationInventoryRequest).Execute()

Updating an Integration Inventory



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
	id := int64(9) // int64 | Morpheus ID of the Integration
	inventoryId := int64(1) // int64 | Morpheus ID of the Integration Inventory
	updateIntegrationInventoryRequest := *openapiclient.NewUpdateIntegrationInventoryRequest(*openapiclient.NewUpdateIntegrationInventoryRequestInventory()) // UpdateIntegrationInventoryRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.UpdateIntegrationInventory(context.Background(), id, inventoryId).UpdateIntegrationInventoryRequest(updateIntegrationInventoryRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.UpdateIntegrationInventory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateIntegrationInventory`: UpdateIntegrationInventory200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.UpdateIntegrationInventory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Integration | 
**inventoryId** | **int64** | Morpheus ID of the Integration Inventory | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateIntegrationInventoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateIntegrationInventoryRequest** | [**UpdateIntegrationInventoryRequest**](UpdateIntegrationInventoryRequest.md) |  | 

### Return type

[**UpdateIntegrationInventory200Response**](UpdateIntegrationInventory200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateIntegrations

> UpdateIntegrations200Response UpdateIntegrations(ctx, id).UpdateIntegrationsRequest(updateIntegrationsRequest).Execute()

Updates an Integration



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
	updateIntegrationsRequest := openapiclient.updateIntegrations_request{UpdateIntegrationsRequestOneOf: openapiclient.NewUpdateIntegrationsRequestOneOf(*openapiclient.NewUpdateIntegrationsRequestOneOfIntegration("Sample Integration", "Type_example"))} // UpdateIntegrationsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IntegrationsAPI.UpdateIntegrations(context.Background(), id).UpdateIntegrationsRequest(updateIntegrationsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IntegrationsAPI.UpdateIntegrations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateIntegrations`: UpdateIntegrations200Response
	fmt.Fprintf(os.Stdout, "Response from `IntegrationsAPI.UpdateIntegrations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateIntegrationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateIntegrationsRequest** | [**UpdateIntegrationsRequest**](UpdateIntegrationsRequest.md) |  | 

### Return type

[**UpdateIntegrations200Response**](UpdateIntegrations200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

