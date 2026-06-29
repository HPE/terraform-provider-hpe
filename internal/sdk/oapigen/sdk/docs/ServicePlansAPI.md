# \ServicePlansAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActivateServicePlans**](ServicePlansAPI.md#ActivateServicePlans) | **Put** /api/service-plans/{id}/activate | Activates a Service Plan
[**AddServicePlans**](ServicePlansAPI.md#AddServicePlans) | **Post** /api/service-plans | Creates a Service Plan
[**DeactivateServicePlans**](ServicePlansAPI.md#DeactivateServicePlans) | **Put** /api/service-plans/{id}/deactivate | Deactivates a Service Plan
[**GetServicePlans**](ServicePlansAPI.md#GetServicePlans) | **Get** /api/service-plans/{id} | Retrieves a Specific Service Plan
[**ListServicePlans**](ServicePlansAPI.md#ListServicePlans) | **Get** /api/service-plans | Retrieves all Service Plans
[**RemoveServicePlans**](ServicePlansAPI.md#RemoveServicePlans) | **Delete** /api/service-plans/{id} | Deletes a Service Plan
[**UpdateServicePlans**](ServicePlansAPI.md#UpdateServicePlans) | **Put** /api/service-plans/{id} | Updates a Service Plan



## ActivateServicePlans

> RemoveExecuteSchedules200Response ActivateServicePlans(ctx, id).Execute()

Activates a Service Plan



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
	resp, r, err := apiClient.ServicePlansAPI.ActivateServicePlans(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicePlansAPI.ActivateServicePlans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ActivateServicePlans`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `ServicePlansAPI.ActivateServicePlans`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiActivateServicePlansRequest struct via the builder pattern


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


## AddServicePlans

> AddServicePlans200Response AddServicePlans(ctx).AddServicePlansRequest(addServicePlansRequest).Execute()

Creates a Service Plan



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
	addServicePlansRequest := *openapiclient.NewAddServicePlansRequest(*openapiclient.NewAddServicePlansRequestServicePlan("Name_example", "Code_example", int64(123), int64(123), *openapiclient.NewAddServicePlansRequestServicePlanProvisionType(int64(123)))) // AddServicePlansRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServicePlansAPI.AddServicePlans(context.Background()).AddServicePlansRequest(addServicePlansRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicePlansAPI.AddServicePlans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddServicePlans`: AddServicePlans200Response
	fmt.Fprintf(os.Stdout, "Response from `ServicePlansAPI.AddServicePlans`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddServicePlansRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addServicePlansRequest** | [**AddServicePlansRequest**](AddServicePlansRequest.md) |  | 

### Return type

[**AddServicePlans200Response**](AddServicePlans200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeactivateServicePlans

> RemoveExecuteSchedules200Response DeactivateServicePlans(ctx, id).Execute()

Deactivates a Service Plan



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
	resp, r, err := apiClient.ServicePlansAPI.DeactivateServicePlans(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicePlansAPI.DeactivateServicePlans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeactivateServicePlans`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `ServicePlansAPI.DeactivateServicePlans`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeactivateServicePlansRequest struct via the builder pattern


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


## GetServicePlans

> GetServicePlans200Response GetServicePlans(ctx, id).Execute()

Retrieves a Specific Service Plan



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
	resp, r, err := apiClient.ServicePlansAPI.GetServicePlans(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicePlansAPI.GetServicePlans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetServicePlans`: GetServicePlans200Response
	fmt.Fprintf(os.Stdout, "Response from `ServicePlansAPI.GetServicePlans`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetServicePlansRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetServicePlans200Response**](GetServicePlans200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListServicePlans

> ListServicePlans200Response ListServicePlans(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).IncludeZones(includeZones).ProvisionTypeId(provisionTypeId).IncludeInactive(includeInactive).Execute()

Retrieves all Service Plans



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
	includeZones := true // bool | Indicates including zones in the service plan response payload (optional) (default to false)
	provisionTypeId := int64(22) // int64 | Provision type filter, restricts query to only load service plans of specified provision type (optional)
	includeInactive := true // bool | If true, include inactive prices in the results (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServicePlansAPI.ListServicePlans(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).IncludeZones(includeZones).ProvisionTypeId(provisionTypeId).IncludeInactive(includeInactive).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicePlansAPI.ListServicePlans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListServicePlans`: ListServicePlans200Response
	fmt.Fprintf(os.Stdout, "Response from `ServicePlansAPI.ListServicePlans`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListServicePlansRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **includeZones** | **bool** | Indicates including zones in the service plan response payload | [default to false]
 **provisionTypeId** | **int64** | Provision type filter, restricts query to only load service plans of specified provision type | 
 **includeInactive** | **bool** | If true, include inactive prices in the results | 

### Return type

[**ListServicePlans200Response**](ListServicePlans200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveServicePlans

> RemoveExecuteSchedules200Response RemoveServicePlans(ctx, id).Execute()

Deletes a Service Plan



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
	resp, r, err := apiClient.ServicePlansAPI.RemoveServicePlans(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicePlansAPI.RemoveServicePlans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveServicePlans`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `ServicePlansAPI.RemoveServicePlans`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveServicePlansRequest struct via the builder pattern


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


## UpdateServicePlans

> UpdateServicePlans200Response UpdateServicePlans(ctx, id).UpdateServicePlansRequest(updateServicePlansRequest).Execute()

Updates a Service Plan



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
	updateServicePlansRequest := *openapiclient.NewUpdateServicePlansRequest(*openapiclient.NewUpdateServicePlansRequestServicePlan()) // UpdateServicePlansRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServicePlansAPI.UpdateServicePlans(context.Background(), id).UpdateServicePlansRequest(updateServicePlansRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicePlansAPI.UpdateServicePlans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateServicePlans`: UpdateServicePlans200Response
	fmt.Fprintf(os.Stdout, "Response from `ServicePlansAPI.UpdateServicePlans`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateServicePlansRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateServicePlansRequest** | [**UpdateServicePlansRequest**](UpdateServicePlansRequest.md) |  | 

### Return type

[**UpdateServicePlans200Response**](UpdateServicePlans200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

