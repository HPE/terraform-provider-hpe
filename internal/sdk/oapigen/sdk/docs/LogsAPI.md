# \LogsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ListLogs**](LogsAPI.md#ListLogs) | **Get** /api/logs | Retrieves Logs



## ListLogs

> ListLogs200Response ListLogs(ctx).Max(max).Offset(offset).Sort(sort).Order(order).Query(query).Message(message).SourceType(sourceType).TypeCode(typeCode).ObjectId(objectId).Token(token).Level(level).StartMs(startMs).EndMs(endMs).StartDateTime(startDateTime).EndDateTime(endDateTime).Containers(containers).Servers(servers).ClusterId(clusterId).Execute()

Retrieves Logs



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
	order := "order_example" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	query := "query_example" // string | Alias for phrase (optional)
	message := "message_example" // string | Filter by message (optional)
	sourceType := "sourceType_example" // string | Filter by source type (optional)
	typeCode := "typeCode_example" // string | Filter by code type (optional)
	objectId := int64(15) // int64 | Filter by objectId (optional)
	token := "token_example" // string | Filter by token (optional)
	level := "ERROR" // string | Filter by log level. Multiple values can be passed pipe delimited. (optional)
	startMs := int64(1657309873) // int64 | Date filter in milliseconds (unix epoch), restricts query to only load logs updated more recently than the time specified. (optional)
	endMs := int64(1657309873) // int64 | Date filter in milliseconds (unix epoch), restricts query to only load logs updated before the time specified. (optional)
	startDateTime := time.Now() // time.Time | Start Date timestamp (ISO 8601) (optional)
	endDateTime := time.Now() // time.Time | End Date timestamp (ISO 8601) (optional)
	containers := int64(135) // int64 | The Container ID(s) for filtering. Accepts multiple values. (optional)
	servers := int64(135) // int64 | The Server ID(s) for filtering. Accepts multiple values. (optional)
	clusterId := int64(135) // int64 | The Cluster ID(s) for filtering. Accepts multiple values. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LogsAPI.ListLogs(context.Background()).Max(max).Offset(offset).Sort(sort).Order(order).Query(query).Message(message).SourceType(sourceType).TypeCode(typeCode).ObjectId(objectId).Token(token).Level(level).StartMs(startMs).EndMs(endMs).StartDateTime(startDateTime).EndDateTime(endDateTime).Containers(containers).Servers(servers).ClusterId(clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LogsAPI.ListLogs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLogs`: ListLogs200Response
	fmt.Fprintf(os.Stdout, "Response from `LogsAPI.ListLogs`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListLogsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **order** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **query** | **string** | Alias for phrase | 
 **message** | **string** | Filter by message | 
 **sourceType** | **string** | Filter by source type | 
 **typeCode** | **string** | Filter by code type | 
 **objectId** | **int64** | Filter by objectId | 
 **token** | **string** | Filter by token | 
 **level** | **string** | Filter by log level. Multiple values can be passed pipe delimited. | 
 **startMs** | **int64** | Date filter in milliseconds (unix epoch), restricts query to only load logs updated more recently than the time specified. | 
 **endMs** | **int64** | Date filter in milliseconds (unix epoch), restricts query to only load logs updated before the time specified. | 
 **startDateTime** | **time.Time** | Start Date timestamp (ISO 8601) | 
 **endDateTime** | **time.Time** | End Date timestamp (ISO 8601) | 
 **containers** | **int64** | The Container ID(s) for filtering. Accepts multiple values. | 
 **servers** | **int64** | The Server ID(s) for filtering. Accepts multiple values. | 
 **clusterId** | **int64** | The Cluster ID(s) for filtering. Accepts multiple values. | 

### Return type

[**ListLogs200Response**](ListLogs200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

