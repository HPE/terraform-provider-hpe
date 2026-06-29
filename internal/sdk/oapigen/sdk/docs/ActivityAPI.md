# \ActivityAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ListActivity**](ActivityAPI.md#ListActivity) | **Get** /api/activity | Retrieves Activity



## ListActivity

> ListActivity200Response ListActivity(ctx).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).Name(name).UserId(userId).TenantId(tenantId).Timeframe(timeframe).Start(start).End(end).Execute()

Retrieves Activity



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
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	userId := int64(6) // int64 | Filter by User ID (optional)
	tenantId := float32(3) // float32 | Filter by Tenant ID. Only available to the master account. (optional)
	timeframe := "timeframe_example" // string | Filter by a timeframe (ex - today, yesterday, week, month, 3months) (optional) (default to "month")
	start := time.Now() // time.Time | Filter by activity on or after a date(time). Default is 1 month prior (optional)
	end := time.Now() // time.Time | Filter by activity on or before a date(time). Default is current date (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ActivityAPI.ListActivity(context.Background()).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).Name(name).UserId(userId).TenantId(tenantId).Timeframe(timeframe).Start(start).End(end).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ActivityAPI.ListActivity``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListActivity`: ListActivity200Response
	fmt.Fprintf(os.Stdout, "Response from `ActivityAPI.ListActivity`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListActivityRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **order** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **userId** | **int64** | Filter by User ID | 
 **tenantId** | **float32** | Filter by Tenant ID. Only available to the master account. | 
 **timeframe** | **string** | Filter by a timeframe (ex - today, yesterday, week, month, 3months) | [default to &quot;month&quot;]
 **start** | **time.Time** | Filter by activity on or after a date(time). Default is 1 month prior | 
 **end** | **time.Time** | Filter by activity on or before a date(time). Default is current date | 

### Return type

[**ListActivity200Response**](ListActivity200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

