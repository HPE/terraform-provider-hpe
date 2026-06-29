# \SecurityGroupsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddSecurityGroupLocations**](SecurityGroupsAPI.md#AddSecurityGroupLocations) | **Post** /api/security-groups/{id}/locations | Creates a Security Group Location
[**AddSecurityGroupRules**](SecurityGroupsAPI.md#AddSecurityGroupRules) | **Post** /api/security-groups/{id}/rules | Creates a Security Group Rule
[**AddSecurityGroups**](SecurityGroupsAPI.md#AddSecurityGroups) | **Post** /api/security-groups | Creates a Security Group
[**GetSecurityGroupRules**](SecurityGroupsAPI.md#GetSecurityGroupRules) | **Get** /api/security-groups/{id}/rules/{sgId} | Retrieves a Specific Security Group Rule
[**GetSecurityGroups**](SecurityGroupsAPI.md#GetSecurityGroups) | **Get** /api/security-groups/{id} | Retrieves a Specific Security Group
[**ListSecurityGroupRules**](SecurityGroupsAPI.md#ListSecurityGroupRules) | **Get** /api/security-groups/{id}/rules | Retrieves all Security Group Rules
[**ListSecurityGroups**](SecurityGroupsAPI.md#ListSecurityGroups) | **Get** /api/security-groups | Retrieves all Security Groups
[**RemoveSecurityGroupLocations**](SecurityGroupsAPI.md#RemoveSecurityGroupLocations) | **Delete** /api/security-groups/{id}/locations/{locationId} | Deletes a Security Group Location
[**RemoveSecurityGroupRules**](SecurityGroupsAPI.md#RemoveSecurityGroupRules) | **Delete** /api/security-groups/{id}/rules/{sgId} | Deletes a Security Group Rule
[**RemoveSecurityGroups**](SecurityGroupsAPI.md#RemoveSecurityGroups) | **Delete** /api/security-groups/{id} | Deletes a Security Group
[**UpdateSecurityGroupRules**](SecurityGroupsAPI.md#UpdateSecurityGroupRules) | **Put** /api/security-groups/{id}/rules/{sgId} | Updates a Security Group Rule
[**UpdateSecurityGroups**](SecurityGroupsAPI.md#UpdateSecurityGroups) | **Put** /api/security-groups/{id} | Updating a Security Group



## AddSecurityGroupLocations

> AddSecurityGroupLocations200Response AddSecurityGroupLocations(ctx, id).AddSecurityGroupLocationsRequest(addSecurityGroupLocationsRequest).Execute()

Creates a Security Group Location



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
	addSecurityGroupLocationsRequest := *openapiclient.NewAddSecurityGroupLocationsRequest(*openapiclient.NewAddSecurityGroupLocationsRequestSecurityGroupLocation(int64(1), *openapiclient.NewAddSecurityGroupLocationsRequestSecurityGroupLocationCustomOptions())) // AddSecurityGroupLocationsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityGroupsAPI.AddSecurityGroupLocations(context.Background(), id).AddSecurityGroupLocationsRequest(addSecurityGroupLocationsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.AddSecurityGroupLocations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddSecurityGroupLocations`: AddSecurityGroupLocations200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.AddSecurityGroupLocations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddSecurityGroupLocationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addSecurityGroupLocationsRequest** | [**AddSecurityGroupLocationsRequest**](AddSecurityGroupLocationsRequest.md) |  | 

### Return type

[**AddSecurityGroupLocations200Response**](AddSecurityGroupLocations200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddSecurityGroupRules

> AddSecurityGroupRules200Response AddSecurityGroupRules(ctx, id).AddSecurityGroupRulesRequest(addSecurityGroupRulesRequest).Execute()

Creates a Security Group Rule



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
	addSecurityGroupRulesRequest := *openapiclient.NewAddSecurityGroupRulesRequest(*openapiclient.NewAddSecurityGroupRulesRequestRule("Protocol_example", "RuleType_example")) // AddSecurityGroupRulesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityGroupsAPI.AddSecurityGroupRules(context.Background(), id).AddSecurityGroupRulesRequest(addSecurityGroupRulesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.AddSecurityGroupRules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddSecurityGroupRules`: AddSecurityGroupRules200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.AddSecurityGroupRules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddSecurityGroupRulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addSecurityGroupRulesRequest** | [**AddSecurityGroupRulesRequest**](AddSecurityGroupRulesRequest.md) |  | 

### Return type

[**AddSecurityGroupRules200Response**](AddSecurityGroupRules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddSecurityGroups

> AddSecurityGroups200Response AddSecurityGroups(ctx).AddSecurityGroupsRequest(addSecurityGroupsRequest).Execute()

Creates a Security Group



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
	addSecurityGroupsRequest := *openapiclient.NewAddSecurityGroupsRequest(*openapiclient.NewAddSecurityGroupsRequestSecurityGroup("Name_example", int64(3))) // AddSecurityGroupsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityGroupsAPI.AddSecurityGroups(context.Background()).AddSecurityGroupsRequest(addSecurityGroupsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.AddSecurityGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddSecurityGroups`: AddSecurityGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.AddSecurityGroups`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddSecurityGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addSecurityGroupsRequest** | [**AddSecurityGroupsRequest**](AddSecurityGroupsRequest.md) |  | 

### Return type

[**AddSecurityGroups200Response**](AddSecurityGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSecurityGroupRules

> GetSecurityGroupRules200Response GetSecurityGroupRules(ctx, id, sgId).Execute()

Retrieves a Specific Security Group Rule



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
	sgId := float32(2352) // float32 | Morpheus ID of the security group rule being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityGroupsAPI.GetSecurityGroupRules(context.Background(), id, sgId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.GetSecurityGroupRules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSecurityGroupRules`: GetSecurityGroupRules200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.GetSecurityGroupRules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 
**sgId** | **float32** | Morpheus ID of the security group rule being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSecurityGroupRulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetSecurityGroupRules200Response**](GetSecurityGroupRules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSecurityGroups

> GetSecurityGroups200Response GetSecurityGroups(ctx, id).Execute()

Retrieves a Specific Security Group



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
	resp, r, err := apiClient.SecurityGroupsAPI.GetSecurityGroups(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.GetSecurityGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSecurityGroups`: GetSecurityGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.GetSecurityGroups`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSecurityGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetSecurityGroups200Response**](GetSecurityGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSecurityGroupRules

> ListSecurityGroupRules200Response ListSecurityGroupRules(ctx, id).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Retrieves all Security Group Rules



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityGroupsAPI.ListSecurityGroupRules(context.Background(), id).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.ListSecurityGroupRules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListSecurityGroupRules`: ListSecurityGroupRules200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.ListSecurityGroupRules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiListSecurityGroupRulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListSecurityGroupRules200Response**](ListSecurityGroupRules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSecurityGroups

> ListSecurityGroups200Response ListSecurityGroups(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Retrieves all Security Groups



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityGroupsAPI.ListSecurityGroups(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.ListSecurityGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListSecurityGroups`: ListSecurityGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.ListSecurityGroups`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListSecurityGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListSecurityGroups200Response**](ListSecurityGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveSecurityGroupLocations

> RemoveExecuteSchedules200Response RemoveSecurityGroupLocations(ctx, id, locationId).Execute()

Deletes a Security Group Location



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
	locationId := float32(2) // float32 | The ID of the location

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityGroupsAPI.RemoveSecurityGroupLocations(context.Background(), id, locationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.RemoveSecurityGroupLocations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveSecurityGroupLocations`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.RemoveSecurityGroupLocations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 
**locationId** | **float32** | The ID of the location | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveSecurityGroupLocationsRequest struct via the builder pattern


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


## RemoveSecurityGroupRules

> RemoveExecuteSchedules200Response RemoveSecurityGroupRules(ctx, id, sgId).Execute()

Deletes a Security Group Rule



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
	sgId := float32(2352) // float32 | Morpheus ID of the security group rule being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityGroupsAPI.RemoveSecurityGroupRules(context.Background(), id, sgId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.RemoveSecurityGroupRules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveSecurityGroupRules`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.RemoveSecurityGroupRules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 
**sgId** | **float32** | Morpheus ID of the security group rule being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveSecurityGroupRulesRequest struct via the builder pattern


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


## RemoveSecurityGroups

> RemoveExecuteSchedules200Response RemoveSecurityGroups(ctx, id).Execute()

Deletes a Security Group



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
	resp, r, err := apiClient.SecurityGroupsAPI.RemoveSecurityGroups(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.RemoveSecurityGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveSecurityGroups`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.RemoveSecurityGroups`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveSecurityGroupsRequest struct via the builder pattern


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


## UpdateSecurityGroupRules

> UpdateSecurityGroupRules200Response UpdateSecurityGroupRules(ctx, id, sgId).UpdateSecurityGroupRulesRequest(updateSecurityGroupRulesRequest).Execute()

Updates a Security Group Rule



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
	sgId := float32(2352) // float32 | Morpheus ID of the security group rule being referenced
	updateSecurityGroupRulesRequest := *openapiclient.NewUpdateSecurityGroupRulesRequest(*openapiclient.NewUpdateSecurityGroupRulesRequestRule("Protocol_example", "RuleType_example")) // UpdateSecurityGroupRulesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityGroupsAPI.UpdateSecurityGroupRules(context.Background(), id, sgId).UpdateSecurityGroupRulesRequest(updateSecurityGroupRulesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.UpdateSecurityGroupRules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateSecurityGroupRules`: UpdateSecurityGroupRules200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.UpdateSecurityGroupRules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 
**sgId** | **float32** | Morpheus ID of the security group rule being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateSecurityGroupRulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateSecurityGroupRulesRequest** | [**UpdateSecurityGroupRulesRequest**](UpdateSecurityGroupRulesRequest.md) |  | 

### Return type

[**UpdateSecurityGroupRules200Response**](UpdateSecurityGroupRules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateSecurityGroups

> UpdateSecurityGroups200Response UpdateSecurityGroups(ctx, id).UpdateSecurityGroupsRequest(updateSecurityGroupsRequest).Execute()

Updating a Security Group



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
	updateSecurityGroupsRequest := *openapiclient.NewUpdateSecurityGroupsRequest(*openapiclient.NewUpdateSecurityGroupsRequestSecurityGroup()) // UpdateSecurityGroupsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SecurityGroupsAPI.UpdateSecurityGroups(context.Background(), id).UpdateSecurityGroupsRequest(updateSecurityGroupsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SecurityGroupsAPI.UpdateSecurityGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateSecurityGroups`: UpdateSecurityGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `SecurityGroupsAPI.UpdateSecurityGroups`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateSecurityGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateSecurityGroupsRequest** | [**UpdateSecurityGroupsRequest**](UpdateSecurityGroupsRequest.md) |  | 

### Return type

[**UpdateSecurityGroups200Response**](UpdateSecurityGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

