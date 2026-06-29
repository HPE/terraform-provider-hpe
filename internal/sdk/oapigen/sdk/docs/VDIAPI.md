# \VDIAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddVDIApps**](VDIAPI.md#AddVDIApps) | **Post** /api/vdi-apps | Creates a VDI App
[**AddVDIGateways**](VDIAPI.md#AddVDIGateways) | **Post** /api/vdi-gateways | Creates a VDI Gateway
[**AddVDIPools**](VDIAPI.md#AddVDIPools) | **Post** /api/vdi-pools | Creates a VDI Pool
[**AddVdiAllocation**](VDIAPI.md#AddVdiAllocation) | **Post** /api/vdi/{id}/allocate | Allocate Virtual Desktop
[**GetVDIAllocations**](VDIAPI.md#GetVDIAllocations) | **Get** /api/vdi-allocations/{id} | Retrieves a Specific VDI Allocation
[**GetVDIApps**](VDIAPI.md#GetVDIApps) | **Get** /api/vdi-apps/{id} | Retrieves a Specific VDI App
[**GetVDIGateways**](VDIAPI.md#GetVDIGateways) | **Get** /api/vdi-gateways/{id} | Retrieves a Specific VDI Gateway
[**GetVDIPools**](VDIAPI.md#GetVDIPools) | **Get** /api/vdi-pools/{id} | Retrieves a Specific VDI Pool
[**GetVdi**](VDIAPI.md#GetVdi) | **Get** /api/vdi/{id} | Get a Specific Virtual Desktop
[**ListVDIAllocations**](VDIAPI.md#ListVDIAllocations) | **Get** /api/vdi-allocations | Retrieves all VDI Allocations
[**ListVDIApps**](VDIAPI.md#ListVDIApps) | **Get** /api/vdi-apps | Retrieves all VDI Apps
[**ListVDIGateways**](VDIAPI.md#ListVDIGateways) | **Get** /api/vdi-gateways | Retrieves all VDI Gateways
[**ListVDIPools**](VDIAPI.md#ListVDIPools) | **Get** /api/vdi-pools | Retrieves all VDI Pools
[**ListVdi**](VDIAPI.md#ListVdi) | **Get** /api/vdi | List Virtual Desktops
[**RemoveVDIApps**](VDIAPI.md#RemoveVDIApps) | **Delete** /api/vdi-apps/{id} | Deletes a VDI App
[**RemoveVDIGateways**](VDIAPI.md#RemoveVDIGateways) | **Delete** /api/vdi-gateways/{id} | Deletes a VDI Gateway
[**RemoveVDIPools**](VDIAPI.md#RemoveVDIPools) | **Delete** /api/vdi-pools/{id} | Deletes a VDI Pool
[**UpdateVDIApps**](VDIAPI.md#UpdateVDIApps) | **Put** /api/vdi-apps/{id} | Updates a VDI App Configuration or Icon
[**UpdateVDIGateways**](VDIAPI.md#UpdateVDIGateways) | **Put** /api/vdi-gateways/{id} | Updates a VDI Gateway Configuration
[**UpdateVDIPools**](VDIAPI.md#UpdateVDIPools) | **Put** /api/vdi-pools/{id} | Updates a VDI Pool Configuration or Icon



## AddVDIApps

> AddVDIApps200Response AddVDIApps(ctx).AddVDIAppsRequest(addVDIAppsRequest).Execute()

Creates a VDI App



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
	addVDIAppsRequest := *openapiclient.NewAddVDIAppsRequest(openapiclient.addVDIApps_request_vdiApp{AddVDIAppsRequestVdiAppOneOf: openapiclient.NewAddVDIAppsRequestVdiAppOneOf("Name_example")}) // AddVDIAppsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VDIAPI.AddVDIApps(context.Background()).AddVDIAppsRequest(addVDIAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.AddVDIApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddVDIApps`: AddVDIApps200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.AddVDIApps`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddVDIAppsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addVDIAppsRequest** | [**AddVDIAppsRequest**](AddVDIAppsRequest.md) |  | 

### Return type

[**AddVDIApps200Response**](AddVDIApps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json, multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddVDIGateways

> AddVDIGateways200Response AddVDIGateways(ctx).AddVDIGatewaysRequest(addVDIGatewaysRequest).Execute()

Creates a VDI Gateway



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
	addVDIGatewaysRequest := *openapiclient.NewAddVDIGatewaysRequest(openapiclient.addVDIGateways_request_vdiGateway{AddVDIGatewaysRequestVdiGatewayOneOf: openapiclient.NewAddVDIGatewaysRequestVdiGatewayOneOf("Name_example", "https://fqdn.com")}) // AddVDIGatewaysRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VDIAPI.AddVDIGateways(context.Background()).AddVDIGatewaysRequest(addVDIGatewaysRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.AddVDIGateways``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddVDIGateways`: AddVDIGateways200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.AddVDIGateways`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddVDIGatewaysRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addVDIGatewaysRequest** | [**AddVDIGatewaysRequest**](AddVDIGatewaysRequest.md) |  | 

### Return type

[**AddVDIGateways200Response**](AddVDIGateways200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddVDIPools

> AddVDIPools200Response AddVDIPools(ctx).AddVDIPoolsRequest(addVDIPoolsRequest).Execute()

Creates a VDI Pool



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
	addVDIPoolsRequest := *openapiclient.NewAddVDIPoolsRequest(openapiclient.addVDIPools_request_vdiPool{AddVDIPoolsRequestVdiPoolOneOf: openapiclient.NewAddVDIPoolsRequestVdiPoolOneOf("Name_example", float32(50), "InstanceConfig_example")}) // AddVDIPoolsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VDIAPI.AddVDIPools(context.Background()).AddVDIPoolsRequest(addVDIPoolsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.AddVDIPools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddVDIPools`: AddVDIPools200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.AddVDIPools`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddVDIPoolsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addVDIPoolsRequest** | [**AddVDIPoolsRequest**](AddVDIPoolsRequest.md) |  | 

### Return type

[**AddVDIPools200Response**](AddVDIPools200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json, multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddVdiAllocation

> AddVdiAllocation200Response AddVdiAllocation(ctx, id).Execute()

Allocate Virtual Desktop



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
	resp, r, err := apiClient.VDIAPI.AddVdiAllocation(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.AddVdiAllocation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddVdiAllocation`: AddVdiAllocation200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.AddVdiAllocation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddVdiAllocationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**AddVdiAllocation200Response**](AddVdiAllocation200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVDIAllocations

> GetVDIAllocations200Response GetVDIAllocations(ctx, id).Execute()

Retrieves a Specific VDI Allocation



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
	resp, r, err := apiClient.VDIAPI.GetVDIAllocations(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.GetVDIAllocations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVDIAllocations`: GetVDIAllocations200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.GetVDIAllocations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVDIAllocationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetVDIAllocations200Response**](GetVDIAllocations200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVDIApps

> GetVDIApps200Response GetVDIApps(ctx, id).Execute()

Retrieves a Specific VDI App



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
	resp, r, err := apiClient.VDIAPI.GetVDIApps(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.GetVDIApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVDIApps`: GetVDIApps200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.GetVDIApps`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVDIAppsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetVDIApps200Response**](GetVDIApps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVDIGateways

> GetVDIGateways200Response GetVDIGateways(ctx, id).Execute()

Retrieves a Specific VDI Gateway



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
	resp, r, err := apiClient.VDIAPI.GetVDIGateways(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.GetVDIGateways``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVDIGateways`: GetVDIGateways200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.GetVDIGateways`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVDIGatewaysRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetVDIGateways200Response**](GetVDIGateways200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVDIPools

> GetVDIPools200Response GetVDIPools(ctx, id).Execute()

Retrieves a Specific VDI Pool



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
	resp, r, err := apiClient.VDIAPI.GetVDIPools(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.GetVDIPools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVDIPools`: GetVDIPools200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.GetVDIPools`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVDIPoolsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetVDIPools200Response**](GetVDIPools200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVdi

> GetVdi200Response GetVdi(ctx, id).Execute()

Get a Specific Virtual Desktop



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
	resp, r, err := apiClient.VDIAPI.GetVdi(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.GetVdi``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVdi`: GetVdi200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.GetVdi`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVdiRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetVdi200Response**](GetVdi200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVDIAllocations

> ListVDIAllocations200Response ListVDIAllocations(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Id(id).Status(status).PoolId(poolId).UserId(userId).Execute()

Retrieves all VDI Allocations



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
	id := "preparing" // string | Filter by allocation ID (optional)
	status := "preparing" // string | Filter by allocation status (optional)
	poolId := int64(1) // int64 | Filter by `VDI Pool` ID (optional)
	userId := int64(6) // int64 | Filter by User ID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VDIAPI.ListVDIAllocations(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Id(id).Status(status).PoolId(poolId).UserId(userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.ListVDIAllocations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListVDIAllocations`: ListVDIAllocations200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.ListVDIAllocations`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListVDIAllocationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **id** | **string** | Filter by allocation ID | 
 **status** | **string** | Filter by allocation status | 
 **poolId** | **int64** | Filter by &#x60;VDI Pool&#x60; ID | 
 **userId** | **int64** | Filter by User ID | 

### Return type

[**ListVDIAllocations200Response**](ListVDIAllocations200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVDIApps

> ListVDIApps200Response ListVDIApps(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Description(description).Execute()

Retrieves all VDI Apps



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VDIAPI.ListVDIApps(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Description(description).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.ListVDIApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListVDIApps`: ListVDIApps200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.ListVDIApps`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListVDIAppsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **description** | **string** | Filter by description, wildcard may be specified as %. eg. &#x60;example-%&#x60; | 

### Return type

[**ListVDIApps200Response**](ListVDIApps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVDIGateways

> ListVDIGateways200Response ListVDIGateways(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Description(description).Execute()

Retrieves all VDI Gateways



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VDIAPI.ListVDIGateways(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Description(description).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.ListVDIGateways``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListVDIGateways`: ListVDIGateways200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.ListVDIGateways`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListVDIGatewaysRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **description** | **string** | Filter by description, wildcard may be specified as %. eg. &#x60;example-%&#x60; | 

### Return type

[**ListVDIGateways200Response**](ListVDIGateways200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVDIPools

> ListVDIPools200Response ListVDIPools(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Description(description).Enabled(enabled).Execute()

Retrieves all VDI Pools



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VDIAPI.ListVDIPools(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Description(description).Enabled(enabled).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.ListVDIPools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListVDIPools`: ListVDIPools200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.ListVDIPools`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListVDIPoolsRequest struct via the builder pattern


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

### Return type

[**ListVDIPools200Response**](ListVDIPools200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVdi

> ListVdi200Response ListVdi(ctx).Phrase(phrase).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Description(description).Execute()

List Virtual Desktops



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
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	name := "example" // string | Filter by name (optional)
	description := "The desription of my cool object" // string | Filter by description, wildcard may be specified as %. eg. `example-%` (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VDIAPI.ListVdi(context.Background()).Phrase(phrase).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Description(description).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.ListVdi``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListVdi`: ListVdi200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.ListVdi`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListVdiRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **name** | **string** | Filter by name | 
 **description** | **string** | Filter by description, wildcard may be specified as %. eg. &#x60;example-%&#x60; | 

### Return type

[**ListVdi200Response**](ListVdi200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveVDIApps

> RemoveExecuteSchedules200Response RemoveVDIApps(ctx, id).Execute()

Deletes a VDI App



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
	resp, r, err := apiClient.VDIAPI.RemoveVDIApps(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.RemoveVDIApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveVDIApps`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.RemoveVDIApps`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveVDIAppsRequest struct via the builder pattern


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


## RemoveVDIGateways

> RemoveExecuteSchedules200Response RemoveVDIGateways(ctx, id).Execute()

Deletes a VDI Gateway



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
	resp, r, err := apiClient.VDIAPI.RemoveVDIGateways(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.RemoveVDIGateways``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveVDIGateways`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.RemoveVDIGateways`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveVDIGatewaysRequest struct via the builder pattern


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


## RemoveVDIPools

> RemoveExecuteSchedules200Response RemoveVDIPools(ctx, id).Execute()

Deletes a VDI Pool



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
	resp, r, err := apiClient.VDIAPI.RemoveVDIPools(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.RemoveVDIPools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveVDIPools`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.RemoveVDIPools`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveVDIPoolsRequest struct via the builder pattern


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


## UpdateVDIApps

> UpdateVDIApps200Response UpdateVDIApps(ctx, id).UpdateVDIAppsRequest(updateVDIAppsRequest).Execute()

Updates a VDI App Configuration or Icon



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
	updateVDIAppsRequest := *openapiclient.NewUpdateVDIAppsRequest(openapiclient.updateVDIApps_request_vdiApp{UpdateVDIAppsRequestVdiAppOneOf: openapiclient.NewUpdateVDIAppsRequestVdiAppOneOf()}) // UpdateVDIAppsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VDIAPI.UpdateVDIApps(context.Background(), id).UpdateVDIAppsRequest(updateVDIAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.UpdateVDIApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateVDIApps`: UpdateVDIApps200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.UpdateVDIApps`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateVDIAppsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateVDIAppsRequest** | [**UpdateVDIAppsRequest**](UpdateVDIAppsRequest.md) |  | 

### Return type

[**UpdateVDIApps200Response**](UpdateVDIApps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json, multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateVDIGateways

> UpdateVDIGateways200Response UpdateVDIGateways(ctx, id).UpdateVDIGatewaysRequest(updateVDIGatewaysRequest).Execute()

Updates a VDI Gateway Configuration



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
	updateVDIGatewaysRequest := *openapiclient.NewUpdateVDIGatewaysRequest(openapiclient.updateVDIGateways_request_vdiGateway{UpdateVDIGatewaysRequestVdiGatewayOneOf: openapiclient.NewUpdateVDIGatewaysRequestVdiGatewayOneOf()}) // UpdateVDIGatewaysRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VDIAPI.UpdateVDIGateways(context.Background(), id).UpdateVDIGatewaysRequest(updateVDIGatewaysRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.UpdateVDIGateways``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateVDIGateways`: UpdateVDIGateways200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.UpdateVDIGateways`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateVDIGatewaysRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateVDIGatewaysRequest** | [**UpdateVDIGatewaysRequest**](UpdateVDIGatewaysRequest.md) |  | 

### Return type

[**UpdateVDIGateways200Response**](UpdateVDIGateways200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateVDIPools

> UpdateVDIPools200Response UpdateVDIPools(ctx, id).UpdateVDIPoolsRequest(updateVDIPoolsRequest).Execute()

Updates a VDI Pool Configuration or Icon



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
	updateVDIPoolsRequest := *openapiclient.NewUpdateVDIPoolsRequest(*openapiclient.NewUpdateVDIPoolsRequestVdiPool()) // UpdateVDIPoolsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VDIAPI.UpdateVDIPools(context.Background(), id).UpdateVDIPoolsRequest(updateVDIPoolsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VDIAPI.UpdateVDIPools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateVDIPools`: UpdateVDIPools200Response
	fmt.Fprintf(os.Stdout, "Response from `VDIAPI.UpdateVDIPools`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateVDIPoolsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateVDIPoolsRequest** | [**UpdateVDIPoolsRequest**](UpdateVDIPoolsRequest.md) |  | 

### Return type

[**UpdateVDIPools200Response**](UpdateVDIPools200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json, multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

