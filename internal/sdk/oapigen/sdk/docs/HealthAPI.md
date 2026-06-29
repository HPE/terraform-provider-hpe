# \HealthAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AcknowledgeHealthAlarm**](HealthAPI.md#AcknowledgeHealthAlarm) | **Put** /api/health/alarms/{id}/acknowledge | Acknowledge a Health Alarm
[**AcknowledgeHealthAlarms**](HealthAPI.md#AcknowledgeHealthAlarms) | **Put** /api/health/alarms/acknowledge | Acknowledge Many Health Alarms
[**ExportHealthLogs**](HealthAPI.md#ExportHealthLogs) | **Get** /api/health/logs/export | Export Appliance Health Logs
[**GetHealthAlarms**](HealthAPI.md#GetHealthAlarms) | **Get** /api/health/alarms/{id} | Retrieves a Specific Appliance Health Alarm
[**ListHealth**](HealthAPI.md#ListHealth) | **Get** /api/health | Retrieves Appliance Health
[**ListHealthAlarms**](HealthAPI.md#ListHealthAlarms) | **Get** /api/health/alarms | Retrieves Appliance Health Alarms
[**ListHealthLogs**](HealthAPI.md#ListHealthLogs) | **Get** /api/health/logs | Retrieves Appliance Health Logs



## AcknowledgeHealthAlarm

> RemoveExecuteSchedules200Response AcknowledgeHealthAlarm(ctx, id).AcknowledgeHealthAlarmRequest(acknowledgeHealthAlarmRequest).Execute()

Acknowledge a Health Alarm



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
	acknowledgeHealthAlarmRequest := *openapiclient.NewAcknowledgeHealthAlarmRequest(*openapiclient.NewAcknowledgeHealthAlarmRequestAlarm(false)) // AcknowledgeHealthAlarmRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.HealthAPI.AcknowledgeHealthAlarm(context.Background(), id).AcknowledgeHealthAlarmRequest(acknowledgeHealthAlarmRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HealthAPI.AcknowledgeHealthAlarm``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AcknowledgeHealthAlarm`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `HealthAPI.AcknowledgeHealthAlarm`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiAcknowledgeHealthAlarmRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **acknowledgeHealthAlarmRequest** | [**AcknowledgeHealthAlarmRequest**](AcknowledgeHealthAlarmRequest.md) |  | 

### Return type

[**RemoveExecuteSchedules200Response**](RemoveExecuteSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AcknowledgeHealthAlarms

> RemoveExecuteSchedules200Response AcknowledgeHealthAlarms(ctx).AcknowledgeHealthAlarmsRequest(acknowledgeHealthAlarmsRequest).Execute()

Acknowledge Many Health Alarms



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
	acknowledgeHealthAlarmsRequest := *openapiclient.NewAcknowledgeHealthAlarmsRequest(*openapiclient.NewAcknowledgeHealthAlarmsRequestAlarm(false)) // AcknowledgeHealthAlarmsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.HealthAPI.AcknowledgeHealthAlarms(context.Background()).AcknowledgeHealthAlarmsRequest(acknowledgeHealthAlarmsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HealthAPI.AcknowledgeHealthAlarms``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AcknowledgeHealthAlarms`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `HealthAPI.AcknowledgeHealthAlarms`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAcknowledgeHealthAlarmsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **acknowledgeHealthAlarmsRequest** | [**AcknowledgeHealthAlarmsRequest**](AcknowledgeHealthAlarmsRequest.md) |  | 

### Return type

[**RemoveExecuteSchedules200Response**](RemoveExecuteSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExportHealthLogs

> *os.File ExportHealthLogs(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Acknowledged(acknowledged).StartDate(startDate).EndDate(endDate).Reverse(reverse).Execute()

Export Appliance Health Logs



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
	acknowledged := false // bool | True or False flag for Acknowledged items (optional)
	startDate := "2019-01-01" // string | Filter by startDate greater than or equal to a specified date (optional)
	endDate := "2020-01-01" // string | Filter by endDate less than or equal to a specified date (optional)
	reverse := true // bool | Reverse order of records. This `true` by default when sort and direction are not passed, but `false` by default if either is passed. This means that by default the newest log entries are the bottom of the file.  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.HealthAPI.ExportHealthLogs(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Acknowledged(acknowledged).StartDate(startDate).EndDate(endDate).Reverse(reverse).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HealthAPI.ExportHealthLogs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExportHealthLogs`: *os.File
	fmt.Fprintf(os.Stdout, "Response from `HealthAPI.ExportHealthLogs`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiExportHealthLogsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **acknowledged** | **bool** | True or False flag for Acknowledged items | 
 **startDate** | **string** | Filter by startDate greater than or equal to a specified date | 
 **endDate** | **string** | Filter by endDate less than or equal to a specified date | 
 **reverse** | **bool** | Reverse order of records. This &#x60;true&#x60; by default when sort and direction are not passed, but &#x60;false&#x60; by default if either is passed. This means that by default the newest log entries are the bottom of the file.  | 

### Return type

[***os.File**](*os.File.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/octet-stream, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetHealthAlarms

> GetHealthAlarms200Response GetHealthAlarms(ctx, id).Execute()

Retrieves a Specific Appliance Health Alarm



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
	resp, r, err := apiClient.HealthAPI.GetHealthAlarms(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HealthAPI.GetHealthAlarms``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetHealthAlarms`: GetHealthAlarms200Response
	fmt.Fprintf(os.Stdout, "Response from `HealthAPI.GetHealthAlarms`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetHealthAlarmsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetHealthAlarms200Response**](GetHealthAlarms200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListHealth

> ListHealth200Response ListHealth(ctx).Execute()

Retrieves Appliance Health



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
	resp, r, err := apiClient.HealthAPI.ListHealth(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HealthAPI.ListHealth``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListHealth`: ListHealth200Response
	fmt.Fprintf(os.Stdout, "Response from `HealthAPI.ListHealth`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListHealthRequest struct via the builder pattern


### Return type

[**ListHealth200Response**](ListHealth200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListHealthAlarms

> ListHealthAlarms200Response ListHealthAlarms(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Acknowledged(acknowledged).Execute()

Retrieves Appliance Health Alarms



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
	acknowledged := false // bool | True or False flag for Acknowledged items (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.HealthAPI.ListHealthAlarms(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Acknowledged(acknowledged).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HealthAPI.ListHealthAlarms``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListHealthAlarms`: ListHealthAlarms200Response
	fmt.Fprintf(os.Stdout, "Response from `HealthAPI.ListHealthAlarms`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListHealthAlarmsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **acknowledged** | **bool** | True or False flag for Acknowledged items | 

### Return type

[**ListHealthAlarms200Response**](ListHealthAlarms200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListHealthLogs

> ListHealthLogs200Response ListHealthLogs(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Acknowledged(acknowledged).StartDate(startDate).EndDate(endDate).Execute()

Retrieves Appliance Health Logs



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
	acknowledged := false // bool | True or False flag for Acknowledged items (optional)
	startDate := "2019-01-01" // string | Filter by startDate greater than or equal to a specified date (optional)
	endDate := "2020-01-01" // string | Filter by endDate less than or equal to a specified date (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.HealthAPI.ListHealthLogs(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Acknowledged(acknowledged).StartDate(startDate).EndDate(endDate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `HealthAPI.ListHealthLogs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListHealthLogs`: ListHealthLogs200Response
	fmt.Fprintf(os.Stdout, "Response from `HealthAPI.ListHealthLogs`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListHealthLogsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **acknowledged** | **bool** | True or False flag for Acknowledged items | 
 **startDate** | **string** | Filter by startDate greater than or equal to a specified date | 
 **endDate** | **string** | Filter by endDate less than or equal to a specified date | 

### Return type

[**ListHealthLogs200Response**](ListHealthLogs200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

