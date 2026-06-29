# \GuidanceAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ExecuteGuidances**](GuidanceAPI.md#ExecuteGuidances) | **Put** /api/guidance/{id}/execute | Executes a Specific Guidance Recommendation
[**GetGuidanceStats**](GuidanceAPI.md#GetGuidanceStats) | **Get** /api/guidance/stats | Retrieves Guidance Stats
[**GetGuidanceTypes**](GuidanceAPI.md#GetGuidanceTypes) | **Get** /api/guidance/types | Retrieves Guidance Types
[**GetGuidances**](GuidanceAPI.md#GetGuidances) | **Get** /api/guidance/{id} | Retrieves a Specific Guidance Recommendation
[**IgnoreGuidances**](GuidanceAPI.md#IgnoreGuidances) | **Put** /api/guidance/{id}/ignore | Ignores a Specific Guidance Recommendation
[**ListGuidances**](GuidanceAPI.md#ListGuidances) | **Get** /api/guidance | Retrieves all Guidance Recommendations



## ExecuteGuidances

> ExecuteGuidances200Response ExecuteGuidances(ctx, id).Execute()

Executes a Specific Guidance Recommendation



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
	resp, r, err := apiClient.GuidanceAPI.ExecuteGuidances(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GuidanceAPI.ExecuteGuidances``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecuteGuidances`: ExecuteGuidances200Response
	fmt.Fprintf(os.Stdout, "Response from `GuidanceAPI.ExecuteGuidances`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiExecuteGuidancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ExecuteGuidances200Response**](ExecuteGuidances200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetGuidanceStats

> GetGuidanceStats200Response GetGuidanceStats(ctx).Execute()

Retrieves Guidance Stats



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
	resp, r, err := apiClient.GuidanceAPI.GetGuidanceStats(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GuidanceAPI.GetGuidanceStats``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetGuidanceStats`: GetGuidanceStats200Response
	fmt.Fprintf(os.Stdout, "Response from `GuidanceAPI.GetGuidanceStats`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetGuidanceStatsRequest struct via the builder pattern


### Return type

[**GetGuidanceStats200Response**](GetGuidanceStats200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetGuidanceTypes

> GetGuidanceTypes200Response GetGuidanceTypes(ctx).Execute()

Retrieves Guidance Types



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
	resp, r, err := apiClient.GuidanceAPI.GetGuidanceTypes(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GuidanceAPI.GetGuidanceTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetGuidanceTypes`: GetGuidanceTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `GuidanceAPI.GetGuidanceTypes`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetGuidanceTypesRequest struct via the builder pattern


### Return type

[**GetGuidanceTypes200Response**](GetGuidanceTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetGuidances

> GetGuidances200Response GetGuidances(ctx, id).Execute()

Retrieves a Specific Guidance Recommendation



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
	resp, r, err := apiClient.GuidanceAPI.GetGuidances(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GuidanceAPI.GetGuidances``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetGuidances`: GetGuidances200Response
	fmt.Fprintf(os.Stdout, "Response from `GuidanceAPI.GetGuidances`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetGuidancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetGuidances200Response**](GetGuidances200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## IgnoreGuidances

> IgnoreGuidances200Response IgnoreGuidances(ctx, id).Execute()

Ignores a Specific Guidance Recommendation



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
	resp, r, err := apiClient.GuidanceAPI.IgnoreGuidances(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GuidanceAPI.IgnoreGuidances``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `IgnoreGuidances`: IgnoreGuidances200Response
	fmt.Fprintf(os.Stdout, "Response from `GuidanceAPI.IgnoreGuidances`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiIgnoreGuidancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**IgnoreGuidances200Response**](IgnoreGuidances200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListGuidances

> ListGuidances200Response ListGuidances(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Severity(severity).Type_(type_).State(state).Execute()

Retrieves all Guidance Recommendations



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
	severity := "severity_example" // string | Filter by severity (optional)
	type_ := "type__example" // string | Filter by Guidance Type.  See `Retrieves Guidance Types` for most up to date list of types. (optional)
	state := "state_example" // string | Filter by state, restricts query to only load discoveries by state (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.GuidanceAPI.ListGuidances(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Severity(severity).Type_(type_).State(state).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GuidanceAPI.ListGuidances``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListGuidances`: ListGuidances200Response
	fmt.Fprintf(os.Stdout, "Response from `GuidanceAPI.ListGuidances`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListGuidancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **severity** | **string** | Filter by severity | 
 **type_** | **string** | Filter by Guidance Type.  See &#x60;Retrieves Guidance Types&#x60; for most up to date list of types. | 
 **state** | **string** | Filter by state, restricts query to only load discoveries by state | 

### Return type

[**ListGuidances200Response**](ListGuidances200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

