# \ChecksAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddCheckApps**](ChecksAPI.md#AddCheckApps) | **Post** /api/monitoring/apps | Create a New Check App
[**AddCheckGroups**](ChecksAPI.md#AddCheckGroups) | **Post** /api/monitoring/groups | Create a New Check Group
[**AddChecks**](ChecksAPI.md#AddChecks) | **Post** /api/monitoring/checks | Create a New Check
[**DeleteCheckApps**](ChecksAPI.md#DeleteCheckApps) | **Delete** /api/monitoring/apps/{id} | Delete a Specific Check App
[**DeleteCheckGroups**](ChecksAPI.md#DeleteCheckGroups) | **Delete** /api/monitoring/groups/{id} | Delete a Specific Check Group
[**DeleteChecks**](ChecksAPI.md#DeleteChecks) | **Delete** /api/monitoring/checks/{id} | Delete a Specific Check
[**GetCheckApps**](ChecksAPI.md#GetCheckApps) | **Get** /api/monitoring/apps/{id} | Get a Specific Check App
[**GetCheckGroups**](ChecksAPI.md#GetCheckGroups) | **Get** /api/monitoring/groups/{id} | Get a Specific Check Group
[**GetCheckTypes**](ChecksAPI.md#GetCheckTypes) | **Get** /api/monitoring/check-types/{id} | Get a Specific Check Type
[**GetChecks**](ChecksAPI.md#GetChecks) | **Get** /api/monitoring/checks/{id} | Get a Specific Check
[**ListCheckApps**](ChecksAPI.md#ListCheckApps) | **Get** /api/monitoring/apps | List All Check Apps
[**ListCheckGroups**](ChecksAPI.md#ListCheckGroups) | **Get** /api/monitoring/groups | List All Check Groups
[**ListCheckTypes**](ChecksAPI.md#ListCheckTypes) | **Get** /api/monitoring/check-types | List All Check Types
[**ListChecks**](ChecksAPI.md#ListChecks) | **Get** /api/monitoring/checks | List All Checks
[**UpdateCheckApps**](ChecksAPI.md#UpdateCheckApps) | **Put** /api/monitoring/apps/{id} | Update Check App
[**UpdateCheckGroups**](ChecksAPI.md#UpdateCheckGroups) | **Put** /api/monitoring/groups/{id} | Update Check Group
[**UpdateChecks**](ChecksAPI.md#UpdateChecks) | **Put** /api/monitoring/checks/{id} | Updates a Check
[**UpdateMuteAllCheckApps**](ChecksAPI.md#UpdateMuteAllCheckApps) | **Put** /api/monitoring/apps/mute-all | Mute All Check Apps
[**UpdateMuteAllCheckGroups**](ChecksAPI.md#UpdateMuteAllCheckGroups) | **Put** /api/monitoring/groups/mute-all | Mute All Check Groups
[**UpdateMuteAllChecks**](ChecksAPI.md#UpdateMuteAllChecks) | **Put** /api/monitoring/checks/mute-all | Mute All Checks
[**UpdateMuteCheckApps**](ChecksAPI.md#UpdateMuteCheckApps) | **Put** /api/monitoring/apps/{id}/mute | Mute Check App
[**UpdateMuteCheckGroups**](ChecksAPI.md#UpdateMuteCheckGroups) | **Put** /api/monitoring/groups/{id}/mute | Mute Check Group
[**UpdateMuteChecks**](ChecksAPI.md#UpdateMuteChecks) | **Put** /api/monitoring/checks/{id}/mute | Mute Check



## AddCheckApps

> AddCheckApps200Response AddCheckApps(ctx).AddCheckAppsRequest(addCheckAppsRequest).Execute()

Create a New Check App



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
	addCheckAppsRequest := *openapiclient.NewAddCheckAppsRequest(*openapiclient.NewAddCheckAppsRequestMonitorApp("My Check App")) // AddCheckAppsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.AddCheckApps(context.Background()).AddCheckAppsRequest(addCheckAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.AddCheckApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddCheckApps`: AddCheckApps200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.AddCheckApps`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddCheckAppsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addCheckAppsRequest** | [**AddCheckAppsRequest**](AddCheckAppsRequest.md) |  | 

### Return type

[**AddCheckApps200Response**](AddCheckApps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddCheckGroups

> AddCheckGroups200Response AddCheckGroups(ctx).AddCheckGroupsRequest(addCheckGroupsRequest).Execute()

Create a New Check Group



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
	addCheckGroupsRequest := *openapiclient.NewAddCheckGroupsRequest(*openapiclient.NewAddCheckGroupsRequestCheckGroup("My Check Group")) // AddCheckGroupsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.AddCheckGroups(context.Background()).AddCheckGroupsRequest(addCheckGroupsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.AddCheckGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddCheckGroups`: AddCheckGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.AddCheckGroups`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddCheckGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addCheckGroupsRequest** | [**AddCheckGroupsRequest**](AddCheckGroupsRequest.md) |  | 

### Return type

[**AddCheckGroups200Response**](AddCheckGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddChecks

> AddChecks200Response AddChecks(ctx).AddChecksRequest(addChecksRequest).Execute()

Create a New Check



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
	addChecksRequest := *openapiclient.NewAddChecksRequest(openapiclient.addChecks_request_check{CheckSocket: openapiclient.NewCheckSocket("My Check")}) // AddChecksRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.AddChecks(context.Background()).AddChecksRequest(addChecksRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.AddChecks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddChecks`: AddChecks200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.AddChecks`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddChecksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addChecksRequest** | [**AddChecksRequest**](AddChecksRequest.md) |  | 

### Return type

[**AddChecks200Response**](AddChecks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteCheckApps

> DeleteCheckApps200Response DeleteCheckApps(ctx, id).Execute()

Delete a Specific Check App



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
	resp, r, err := apiClient.ChecksAPI.DeleteCheckApps(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.DeleteCheckApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteCheckApps`: DeleteCheckApps200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.DeleteCheckApps`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCheckAppsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteCheckApps200Response**](DeleteCheckApps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteCheckGroups

> DeleteCheckGroups200Response DeleteCheckGroups(ctx, id).Execute()

Delete a Specific Check Group



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
	resp, r, err := apiClient.ChecksAPI.DeleteCheckGroups(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.DeleteCheckGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteCheckGroups`: DeleteCheckGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.DeleteCheckGroups`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCheckGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteCheckGroups200Response**](DeleteCheckGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteChecks

> DeleteChecks200Response DeleteChecks(ctx, id).Execute()

Delete a Specific Check



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
	resp, r, err := apiClient.ChecksAPI.DeleteChecks(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.DeleteChecks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteChecks`: DeleteChecks200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.DeleteChecks`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteChecksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteChecks200Response**](DeleteChecks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCheckApps

> GetCheckApps200Response GetCheckApps(ctx, id).Execute()

Get a Specific Check App



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
	resp, r, err := apiClient.ChecksAPI.GetCheckApps(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.GetCheckApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCheckApps`: GetCheckApps200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.GetCheckApps`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCheckAppsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetCheckApps200Response**](GetCheckApps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCheckGroups

> GetCheckGroups200Response GetCheckGroups(ctx, id).Execute()

Get a Specific Check Group



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
	resp, r, err := apiClient.ChecksAPI.GetCheckGroups(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.GetCheckGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCheckGroups`: GetCheckGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.GetCheckGroups`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCheckGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetCheckGroups200Response**](GetCheckGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCheckTypes

> GetCheckTypes200Response GetCheckTypes(ctx, id).Execute()

Get a Specific Check Type



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
	resp, r, err := apiClient.ChecksAPI.GetCheckTypes(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.GetCheckTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCheckTypes`: GetCheckTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.GetCheckTypes`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCheckTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetCheckTypes200Response**](GetCheckTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetChecks

> GetChecks200Response GetChecks(ctx, id).Execute()

Get a Specific Check



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
	resp, r, err := apiClient.ChecksAPI.GetChecks(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.GetChecks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetChecks`: GetChecks200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.GetChecks`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetChecksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetChecks200Response**](GetChecks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCheckApps

> ListCheckApps200Response ListCheckApps(ctx).Max(max).Offset(offset).Sort(sort).Name(name).Phrase(phrase).Status(status).LastUpdated(lastUpdated).Deleted(deleted).Execute()

List All Check Apps



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	status := "running" // string | The instance status for filtering. (optional)
	lastUpdated := time.Now() // time.Time | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) (optional)
	deleted := true // bool | If true, only deleted resources or instances in pending removal status are returned. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.ListCheckApps(context.Background()).Max(max).Offset(offset).Sort(sort).Name(name).Phrase(phrase).Status(status).LastUpdated(lastUpdated).Deleted(deleted).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.ListCheckApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCheckApps`: ListCheckApps200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.ListCheckApps`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListCheckAppsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **status** | **string** | The instance status for filtering. | 
 **lastUpdated** | **time.Time** | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) | 
 **deleted** | **bool** | If true, only deleted resources or instances in pending removal status are returned. | 

### Return type

[**ListCheckApps200Response**](ListCheckApps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCheckGroups

> ListCheckGroups200Response ListCheckGroups(ctx).Max(max).Offset(offset).Sort(sort).Name(name).Phrase(phrase).Status(status).LastUpdated(lastUpdated).Deleted(deleted).Execute()

List All Check Groups



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	status := "running" // string | The instance status for filtering. (optional)
	lastUpdated := time.Now() // time.Time | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) (optional)
	deleted := true // bool | If true, only deleted resources or instances in pending removal status are returned. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.ListCheckGroups(context.Background()).Max(max).Offset(offset).Sort(sort).Name(name).Phrase(phrase).Status(status).LastUpdated(lastUpdated).Deleted(deleted).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.ListCheckGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCheckGroups`: ListCheckGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.ListCheckGroups`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListCheckGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **status** | **string** | The instance status for filtering. | 
 **lastUpdated** | **time.Time** | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) | 
 **deleted** | **bool** | If true, only deleted resources or instances in pending removal status are returned. | 

### Return type

[**ListCheckGroups200Response**](ListCheckGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCheckTypes

> ListCheckTypes200Response ListCheckTypes(ctx).Max(max).Offset(offset).Sort(sort).Name(name).Phrase(phrase).Execute()

List All Check Types



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
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.ListCheckTypes(context.Background()).Max(max).Offset(offset).Sort(sort).Name(name).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.ListCheckTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCheckTypes`: ListCheckTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.ListCheckTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListCheckTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListCheckTypes200Response**](ListCheckTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListChecks

> ListChecks200Response ListChecks(ctx).Max(max).Offset(offset).Sort(sort).Name(name).Phrase(phrase).Status(status).LastUpdated(lastUpdated).Deleted(deleted).Execute()

List All Checks



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	status := "healthy" // string | Filter by health status. (optional)
	lastUpdated := time.Now() // time.Time | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) (optional)
	deleted := true // bool | If true, only deleted resources or instances in pending removal status are returned. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.ListChecks(context.Background()).Max(max).Offset(offset).Sort(sort).Name(name).Phrase(phrase).Status(status).LastUpdated(lastUpdated).Deleted(deleted).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.ListChecks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListChecks`: ListChecks200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.ListChecks`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListChecksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **status** | **string** | Filter by health status. | 
 **lastUpdated** | **time.Time** | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) | 
 **deleted** | **bool** | If true, only deleted resources or instances in pending removal status are returned. | 

### Return type

[**ListChecks200Response**](ListChecks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateCheckApps

> UpdateCheckApps200Response UpdateCheckApps(ctx, id).UpdateCheckAppsRequest(updateCheckAppsRequest).Execute()

Update Check App



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
	updateCheckAppsRequest := *openapiclient.NewUpdateCheckAppsRequest(*openapiclient.NewUpdateCheckAppsRequestMonitorApp()) // UpdateCheckAppsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.UpdateCheckApps(context.Background(), id).UpdateCheckAppsRequest(updateCheckAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.UpdateCheckApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateCheckApps`: UpdateCheckApps200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.UpdateCheckApps`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateCheckAppsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateCheckAppsRequest** | [**UpdateCheckAppsRequest**](UpdateCheckAppsRequest.md) |  | 

### Return type

[**UpdateCheckApps200Response**](UpdateCheckApps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateCheckGroups

> UpdateCheckGroups200Response UpdateCheckGroups(ctx, id).UpdateCheckGroupsRequest(updateCheckGroupsRequest).Execute()

Update Check Group



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
	updateCheckGroupsRequest := *openapiclient.NewUpdateCheckGroupsRequest(*openapiclient.NewUpdateCheckGroupsRequestCheckGroup()) // UpdateCheckGroupsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.UpdateCheckGroups(context.Background(), id).UpdateCheckGroupsRequest(updateCheckGroupsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.UpdateCheckGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateCheckGroups`: UpdateCheckGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.UpdateCheckGroups`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateCheckGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateCheckGroupsRequest** | [**UpdateCheckGroupsRequest**](UpdateCheckGroupsRequest.md) |  | 

### Return type

[**UpdateCheckGroups200Response**](UpdateCheckGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateChecks

> UpdateChecks200Response UpdateChecks(ctx, id).UpdateChecksRequest(updateChecksRequest).Execute()

Updates a Check



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
	updateChecksRequest := *openapiclient.NewUpdateChecksRequest(openapiclient.updateChecks_request_check{CheckSocket1: openapiclient.NewCheckSocket1("My Check")}) // UpdateChecksRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.UpdateChecks(context.Background(), id).UpdateChecksRequest(updateChecksRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.UpdateChecks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateChecks`: UpdateChecks200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.UpdateChecks`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateChecksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateChecksRequest** | [**UpdateChecksRequest**](UpdateChecksRequest.md) |  | 

### Return type

[**UpdateChecks200Response**](UpdateChecks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMuteAllCheckApps

> UpdateMuteAllCheckApps200Response UpdateMuteAllCheckApps(ctx).UpdateMuteAllCheckAppsRequest(updateMuteAllCheckAppsRequest).Execute()

Mute All Check Apps



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
	updateMuteAllCheckAppsRequest := *openapiclient.NewUpdateMuteAllCheckAppsRequest(false) // UpdateMuteAllCheckAppsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.UpdateMuteAllCheckApps(context.Background()).UpdateMuteAllCheckAppsRequest(updateMuteAllCheckAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.UpdateMuteAllCheckApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMuteAllCheckApps`: UpdateMuteAllCheckApps200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.UpdateMuteAllCheckApps`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMuteAllCheckAppsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateMuteAllCheckAppsRequest** | [**UpdateMuteAllCheckAppsRequest**](UpdateMuteAllCheckAppsRequest.md) |  | 

### Return type

[**UpdateMuteAllCheckApps200Response**](UpdateMuteAllCheckApps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMuteAllCheckGroups

> UpdateMuteAllCheckGroups200Response UpdateMuteAllCheckGroups(ctx).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()

Mute All Check Groups



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
	updateMuteCheckAppsRequest := *openapiclient.NewUpdateMuteCheckAppsRequest(false) // UpdateMuteCheckAppsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.UpdateMuteAllCheckGroups(context.Background()).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.UpdateMuteAllCheckGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMuteAllCheckGroups`: UpdateMuteAllCheckGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.UpdateMuteAllCheckGroups`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMuteAllCheckGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateMuteCheckAppsRequest** | [**UpdateMuteCheckAppsRequest**](UpdateMuteCheckAppsRequest.md) |  | 

### Return type

[**UpdateMuteAllCheckGroups200Response**](UpdateMuteAllCheckGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMuteAllChecks

> UpdateMuteAllChecks200Response UpdateMuteAllChecks(ctx).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()

Mute All Checks



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
	updateMuteCheckAppsRequest := *openapiclient.NewUpdateMuteCheckAppsRequest(false) // UpdateMuteCheckAppsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.UpdateMuteAllChecks(context.Background()).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.UpdateMuteAllChecks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMuteAllChecks`: UpdateMuteAllChecks200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.UpdateMuteAllChecks`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMuteAllChecksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateMuteCheckAppsRequest** | [**UpdateMuteCheckAppsRequest**](UpdateMuteCheckAppsRequest.md) |  | 

### Return type

[**UpdateMuteAllChecks200Response**](UpdateMuteAllChecks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMuteCheckApps

> UpdateMuteCheckApps200Response UpdateMuteCheckApps(ctx, id).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()

Mute Check App



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
	updateMuteCheckAppsRequest := *openapiclient.NewUpdateMuteCheckAppsRequest(false) // UpdateMuteCheckAppsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.UpdateMuteCheckApps(context.Background(), id).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.UpdateMuteCheckApps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMuteCheckApps`: UpdateMuteCheckApps200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.UpdateMuteCheckApps`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMuteCheckAppsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateMuteCheckAppsRequest** | [**UpdateMuteCheckAppsRequest**](UpdateMuteCheckAppsRequest.md) |  | 

### Return type

[**UpdateMuteCheckApps200Response**](UpdateMuteCheckApps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMuteCheckGroups

> UpdateMuteCheckGroups200Response UpdateMuteCheckGroups(ctx, id).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()

Mute Check Group



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
	updateMuteCheckAppsRequest := *openapiclient.NewUpdateMuteCheckAppsRequest(false) // UpdateMuteCheckAppsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.UpdateMuteCheckGroups(context.Background(), id).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.UpdateMuteCheckGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMuteCheckGroups`: UpdateMuteCheckGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.UpdateMuteCheckGroups`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMuteCheckGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateMuteCheckAppsRequest** | [**UpdateMuteCheckAppsRequest**](UpdateMuteCheckAppsRequest.md) |  | 

### Return type

[**UpdateMuteCheckGroups200Response**](UpdateMuteCheckGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMuteChecks

> UpdateMuteChecks200Response UpdateMuteChecks(ctx, id).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()

Mute Check



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
	updateMuteCheckAppsRequest := *openapiclient.NewUpdateMuteCheckAppsRequest(false) // UpdateMuteCheckAppsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChecksAPI.UpdateMuteChecks(context.Background(), id).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChecksAPI.UpdateMuteChecks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMuteChecks`: UpdateMuteChecks200Response
	fmt.Fprintf(os.Stdout, "Response from `ChecksAPI.UpdateMuteChecks`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMuteChecksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateMuteCheckAppsRequest** | [**UpdateMuteCheckAppsRequest**](UpdateMuteCheckAppsRequest.md) |  | 

### Return type

[**UpdateMuteChecks200Response**](UpdateMuteChecks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

