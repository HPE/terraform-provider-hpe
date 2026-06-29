# \SecurityScansAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetSecurityScans**](SecurityScansAPI.md#GetSecurityScans) | **Get** /api/security-scans/{id} | Retrieves a Specific Security Scan
[**ListSecurityScans**](SecurityScansAPI.md#ListSecurityScans) | **Get** /api/security-scans | Retrieves all Security Scans



## GetSecurityScans

> GetSecurityScans200Response GetSecurityScans(ctx, id).Results(results).Execute()

Retrieves a Specific Security Scan



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
	results := true // bool | Include the `results` object in the response under the security scan. This is a potentially very large object containing the raw results of the scan. (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityScansAPI.GetSecurityScans(context.Background(), id).Results(results).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityScansAPI.GetSecurityScans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSecurityScans`: GetSecurityScans200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityScansAPI.GetSecurityScans`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSecurityScansRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **results** | **bool** | Include the &#x60;results&#x60; object in the response under the security scan. This is a potentially very large object containing the raw results of the scan. | [default to false]

### Return type

[**GetSecurityScans200Response**](GetSecurityScans200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSecurityScans

> ListSecurityScans200Response ListSecurityScans(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).SecurityPackageId(securityPackageId).ServerId(serverId).Results(results).Execute()

Retrieves all Security Scans



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
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "scanDate")
	direction := "desc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "desc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description of security package (optional)
	securityPackageId := int64(789) // int64 | Filter results by security package id(s). This parameter can be passed multiple times to match more than one id. (optional)
	serverId := int64(97) // int64 | The Server ID for Filtering (optional)
	results := true // bool | Include the `results` object in the response under each security scan. This is a potentially very large object containing the raw results of the scan. (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityScansAPI.ListSecurityScans(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).SecurityPackageId(securityPackageId).ServerId(serverId).Results(results).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityScansAPI.ListSecurityScans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListSecurityScans`: ListSecurityScans200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityScansAPI.ListSecurityScans`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListSecurityScansRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;scanDate&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;desc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description of security package | 
 **securityPackageId** | **int64** | Filter results by security package id(s). This parameter can be passed multiple times to match more than one id. | 
 **serverId** | **int64** | The Server ID for Filtering | 
 **results** | **bool** | Include the &#x60;results&#x60; object in the response under each security scan. This is a potentially very large object containing the raw results of the scan. | [default to false]

### Return type

[**ListSecurityScans200Response**](ListSecurityScans200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

