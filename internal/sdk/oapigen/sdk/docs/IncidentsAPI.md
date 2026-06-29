# \IncidentsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddIncident**](IncidentsAPI.md#AddIncident) | **Post** /api/monitoring/incidents | Create a New Incident
[**DeleteIncidents**](IncidentsAPI.md#DeleteIncidents) | **Delete** /api/monitoring/incidents/{id} | Close a Specific Incident
[**GetIncidents**](IncidentsAPI.md#GetIncidents) | **Get** /api/monitoring/incidents/{id} | Get a Specific Incident
[**ListIncidents**](IncidentsAPI.md#ListIncidents) | **Get** /api/monitoring/incidents | List All Incidents
[**UpdateIncidents**](IncidentsAPI.md#UpdateIncidents) | **Put** /api/monitoring/incidents/{id} | Update Incident
[**UpdateIncidentsReopen**](IncidentsAPI.md#UpdateIncidentsReopen) | **Get** /api/monitoring/incidents/{id}/reopen | Reopen a Specific Incident
[**UpdateMuteAllIncidents**](IncidentsAPI.md#UpdateMuteAllIncidents) | **Put** /api/monitoring/incidents/mute-all | Mute All Incidents
[**UpdateMuteIncidents**](IncidentsAPI.md#UpdateMuteIncidents) | **Put** /api/monitoring/incidents/{id}/mute | Mute Incident



## AddIncident

> AddIncident200Response AddIncident(ctx).AddIncidentRequest(addIncidentRequest).Execute()

Create a New Incident



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
	addIncidentRequest := *openapiclient.NewAddIncidentRequest(*openapiclient.NewAddIncidentRequestIncident("Incident Name")) // AddIncidentRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IncidentsAPI.AddIncident(context.Background()).AddIncidentRequest(addIncidentRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IncidentsAPI.AddIncident``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddIncident`: AddIncident200Response
	fmt.Fprintf(os.Stdout, "Response from `IncidentsAPI.AddIncident`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddIncidentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addIncidentRequest** | [**AddIncidentRequest**](AddIncidentRequest.md) |  | 

### Return type

[**AddIncident200Response**](AddIncident200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteIncidents

> DeleteIncidents200Response DeleteIncidents(ctx, id).Execute()

Close a Specific Incident



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
	resp, r, err := apiClient.IncidentsAPI.DeleteIncidents(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IncidentsAPI.DeleteIncidents``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteIncidents`: DeleteIncidents200Response
	fmt.Fprintf(os.Stdout, "Response from `IncidentsAPI.DeleteIncidents`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteIncidentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteIncidents200Response**](DeleteIncidents200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetIncidents

> GetIncidents200Response GetIncidents(ctx, id).Execute()

Get a Specific Incident



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
	resp, r, err := apiClient.IncidentsAPI.GetIncidents(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IncidentsAPI.GetIncidents``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetIncidents`: GetIncidents200Response
	fmt.Fprintf(os.Stdout, "Response from `IncidentsAPI.GetIncidents`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetIncidentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetIncidents200Response**](GetIncidents200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListIncidents

> ListIncidents200Response ListIncidents(ctx).Max(max).Offset(offset).Status(status).Severity(severity).Execute()

List All Incidents



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
	status := "running" // string | The instance status for filtering. (optional)
	severity := "severity_example" // string | Filter by severity (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IncidentsAPI.ListIncidents(context.Background()).Max(max).Offset(offset).Status(status).Severity(severity).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IncidentsAPI.ListIncidents``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListIncidents`: ListIncidents200Response
	fmt.Fprintf(os.Stdout, "Response from `IncidentsAPI.ListIncidents`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListIncidentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **status** | **string** | The instance status for filtering. | 
 **severity** | **string** | Filter by severity | 

### Return type

[**ListIncidents200Response**](ListIncidents200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateIncidents

> UpdateIncidents200Response UpdateIncidents(ctx, id).UpdateIncidentsRequest(updateIncidentsRequest).Execute()

Update Incident



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
	updateIncidentsRequest := *openapiclient.NewUpdateIncidentsRequest(*openapiclient.NewUpdateIncidentsRequestIncident()) // UpdateIncidentsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IncidentsAPI.UpdateIncidents(context.Background(), id).UpdateIncidentsRequest(updateIncidentsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IncidentsAPI.UpdateIncidents``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateIncidents`: UpdateIncidents200Response
	fmt.Fprintf(os.Stdout, "Response from `IncidentsAPI.UpdateIncidents`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateIncidentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateIncidentsRequest** | [**UpdateIncidentsRequest**](UpdateIncidentsRequest.md) |  | 

### Return type

[**UpdateIncidents200Response**](UpdateIncidents200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateIncidentsReopen

> UpdateIncidentsReopen200Response UpdateIncidentsReopen(ctx, id).Execute()

Reopen a Specific Incident



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
	resp, r, err := apiClient.IncidentsAPI.UpdateIncidentsReopen(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IncidentsAPI.UpdateIncidentsReopen``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateIncidentsReopen`: UpdateIncidentsReopen200Response
	fmt.Fprintf(os.Stdout, "Response from `IncidentsAPI.UpdateIncidentsReopen`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateIncidentsReopenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**UpdateIncidentsReopen200Response**](UpdateIncidentsReopen200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMuteAllIncidents

> UpdateMuteAllIncidents200Response UpdateMuteAllIncidents(ctx).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()

Mute All Incidents



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
	resp, r, err := apiClient.IncidentsAPI.UpdateMuteAllIncidents(context.Background()).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IncidentsAPI.UpdateMuteAllIncidents``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMuteAllIncidents`: UpdateMuteAllIncidents200Response
	fmt.Fprintf(os.Stdout, "Response from `IncidentsAPI.UpdateMuteAllIncidents`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMuteAllIncidentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateMuteCheckAppsRequest** | [**UpdateMuteCheckAppsRequest**](UpdateMuteCheckAppsRequest.md) |  | 

### Return type

[**UpdateMuteAllIncidents200Response**](UpdateMuteAllIncidents200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMuteIncidents

> UpdateMuteIncidents200Response UpdateMuteIncidents(ctx, id).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()

Mute Incident



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
	resp, r, err := apiClient.IncidentsAPI.UpdateMuteIncidents(context.Background(), id).UpdateMuteCheckAppsRequest(updateMuteCheckAppsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IncidentsAPI.UpdateMuteIncidents``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMuteIncidents`: UpdateMuteIncidents200Response
	fmt.Fprintf(os.Stdout, "Response from `IncidentsAPI.UpdateMuteIncidents`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMuteIncidentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateMuteCheckAppsRequest** | [**UpdateMuteCheckAppsRequest**](UpdateMuteCheckAppsRequest.md) |  | 

### Return type

[**UpdateMuteIncidents200Response**](UpdateMuteIncidents200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

