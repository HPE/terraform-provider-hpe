# \PoliciesAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddPolicies**](PoliciesAPI.md#AddPolicies) | **Post** /api/policies | Creates a Policy
[**AddPoliciesCloud**](PoliciesAPI.md#AddPoliciesCloud) | **Post** /api/zones/{cloudId}/policies | Creates a Policy for a Cloud
[**AddPoliciesGroup**](PoliciesAPI.md#AddPoliciesGroup) | **Post** /api/groups/{groupId}/policies | Creates a Policy for a Group
[**GetPolicies**](PoliciesAPI.md#GetPolicies) | **Get** /api/policies/{id} | Retrieves a Specific Policy
[**GetPoliciesCloud**](PoliciesAPI.md#GetPoliciesCloud) | **Get** /api/zones/{cloudId}/policies/{id} | Retrieves a Specific Policy for a Cloud
[**GetPoliciesGroup**](PoliciesAPI.md#GetPoliciesGroup) | **Get** /api/groups/{groupId}/policies/{id} | Retrieves a Specific Policy for a Group
[**ListPolicies**](PoliciesAPI.md#ListPolicies) | **Get** /api/policies | Retrieves all Policies
[**ListPoliciesCloud**](PoliciesAPI.md#ListPoliciesCloud) | **Get** /api/zones/{cloudId}/policies | Retrieves Policies for a Cloud
[**ListPoliciesGroup**](PoliciesAPI.md#ListPoliciesGroup) | **Get** /api/groups/{groupId}/policies | Retrieves Policies for a Group
[**ListPolicyTypes**](PoliciesAPI.md#ListPolicyTypes) | **Get** /api/policy-types | Retrieves all Policy Types
[**RemovePolicies**](PoliciesAPI.md#RemovePolicies) | **Delete** /api/policies/{id} | Deletes a Policy
[**RemovePoliciesCloud**](PoliciesAPI.md#RemovePoliciesCloud) | **Delete** /api/zones/{cloudId}/policies/{id} | Deletes a Policy for a Cloud
[**RemovePoliciesGroup**](PoliciesAPI.md#RemovePoliciesGroup) | **Delete** /api/groups/{groupId}/policies/{id} | Deletes a Policy for a Group
[**UpdatePolicies**](PoliciesAPI.md#UpdatePolicies) | **Put** /api/policies/{id} | Updates a Policy
[**UpdatePoliciesCloud**](PoliciesAPI.md#UpdatePoliciesCloud) | **Put** /api/zones/{cloudId}/policies/{id} | Updates a Policy for a Cloud
[**UpdatePoliciesGroup**](PoliciesAPI.md#UpdatePoliciesGroup) | **Put** /api/groups/{groupId}/policies/{id} | Updates a Policy for a Group



## AddPolicies

> AddPolicies200Response AddPolicies(ctx).AddPoliciesRequest(addPoliciesRequest).Execute()

Creates a Policy



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
	addPoliciesRequest := *openapiclient.NewAddPoliciesRequest(*openapiclient.NewAddPoliciesRequestPolicy("Sample Policy", *openapiclient.NewAddPoliciesRequestPolicyPolicyType(), *openapiclient.NewAddPoliciesRequestPolicyConfig(*openapiclient.NewMaxMemoryPolicyTypeConfiguration1MaxMemory()))) // AddPoliciesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.AddPolicies(context.Background()).AddPoliciesRequest(addPoliciesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.AddPolicies``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddPolicies`: AddPolicies200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.AddPolicies`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddPoliciesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addPoliciesRequest** | [**AddPoliciesRequest**](AddPoliciesRequest.md) |  | 

### Return type

[**AddPolicies200Response**](AddPolicies200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddPoliciesCloud

> AddPoliciesCloud200Response AddPoliciesCloud(ctx, cloudId).AddPoliciesCloudRequest(addPoliciesCloudRequest).Execute()

Creates a Policy for a Cloud



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
	cloudId := int64(7) // int64 | The ID of the cloud
	addPoliciesCloudRequest := *openapiclient.NewAddPoliciesCloudRequest(*openapiclient.NewAddPoliciesCloudRequestPolicy("Sample Policy", *openapiclient.NewAddPoliciesCloudRequestPolicyPolicyType())) // AddPoliciesCloudRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.AddPoliciesCloud(context.Background(), cloudId).AddPoliciesCloudRequest(addPoliciesCloudRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.AddPoliciesCloud``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddPoliciesCloud`: AddPoliciesCloud200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.AddPoliciesCloud`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**cloudId** | **int64** | The ID of the cloud | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddPoliciesCloudRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addPoliciesCloudRequest** | [**AddPoliciesCloudRequest**](AddPoliciesCloudRequest.md) |  | 

### Return type

[**AddPoliciesCloud200Response**](AddPoliciesCloud200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddPoliciesGroup

> AddPoliciesGroup200Response AddPoliciesGroup(ctx, groupId).AddPoliciesGroupRequest(addPoliciesGroupRequest).Execute()

Creates a Policy for a Group



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
	groupId := int64(7) // int64 | The ID of the group
	addPoliciesGroupRequest := *openapiclient.NewAddPoliciesGroupRequest(*openapiclient.NewAddPoliciesGroupRequestPolicy("Sample Policy", *openapiclient.NewAddPoliciesGroupRequestPolicyPolicyType())) // AddPoliciesGroupRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.AddPoliciesGroup(context.Background(), groupId).AddPoliciesGroupRequest(addPoliciesGroupRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.AddPoliciesGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddPoliciesGroup`: AddPoliciesGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.AddPoliciesGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**groupId** | **int64** | The ID of the group | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddPoliciesGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addPoliciesGroupRequest** | [**AddPoliciesGroupRequest**](AddPoliciesGroupRequest.md) |  | 

### Return type

[**AddPoliciesGroup200Response**](AddPoliciesGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPolicies

> GetPolicies200Response GetPolicies(ctx, id).Execute()

Retrieves a Specific Policy



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
	resp, r, err := apiClient.PoliciesAPI.GetPolicies(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.GetPolicies``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPolicies`: GetPolicies200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.GetPolicies`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPoliciesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetPolicies200Response**](GetPolicies200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPoliciesCloud

> GetPoliciesCloud200Response GetPoliciesCloud(ctx, cloudId, id).Execute()

Retrieves a Specific Policy for a Cloud



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
	cloudId := int64(7) // int64 | The ID of the cloud
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.GetPoliciesCloud(context.Background(), cloudId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.GetPoliciesCloud``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPoliciesCloud`: GetPoliciesCloud200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.GetPoliciesCloud`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**cloudId** | **int64** | The ID of the cloud | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPoliciesCloudRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetPoliciesCloud200Response**](GetPoliciesCloud200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPoliciesGroup

> GetPoliciesGroup200Response GetPoliciesGroup(ctx, groupId, id).Execute()

Retrieves a Specific Policy for a Group



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
	groupId := int64(7) // int64 | The ID of the group
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.GetPoliciesGroup(context.Background(), groupId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.GetPoliciesGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPoliciesGroup`: GetPoliciesGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.GetPoliciesGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**groupId** | **int64** | The ID of the group | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPoliciesGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetPoliciesGroup200Response**](GetPoliciesGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPolicies

> ListPolicies200Response ListPolicies(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Retrieves all Policies



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
	resp, r, err := apiClient.PoliciesAPI.ListPolicies(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.ListPolicies``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPolicies`: ListPolicies200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.ListPolicies`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListPoliciesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListPolicies200Response**](ListPolicies200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPoliciesCloud

> ListPoliciesCloud200Response ListPoliciesCloud(ctx, cloudId).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Retrieves Policies for a Cloud



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
	cloudId := int64(7) // int64 | The ID of the cloud
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.ListPoliciesCloud(context.Background(), cloudId).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.ListPoliciesCloud``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPoliciesCloud`: ListPoliciesCloud200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.ListPoliciesCloud`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**cloudId** | **int64** | The ID of the cloud | 

### Other Parameters

Other parameters are passed through a pointer to a apiListPoliciesCloudRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListPoliciesCloud200Response**](ListPoliciesCloud200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPoliciesGroup

> ListPoliciesGroup200Response ListPoliciesGroup(ctx, groupId).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Retrieves Policies for a Group



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
	groupId := int64(7) // int64 | The ID of the group
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.ListPoliciesGroup(context.Background(), groupId).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.ListPoliciesGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPoliciesGroup`: ListPoliciesGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.ListPoliciesGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**groupId** | **int64** | The ID of the group | 

### Other Parameters

Other parameters are passed through a pointer to a apiListPoliciesGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListPoliciesGroup200Response**](ListPoliciesGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPolicyTypes

> ListPolicyTypes200Response ListPolicyTypes(ctx).Execute()

Retrieves all Policy Types



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
	resp, r, err := apiClient.PoliciesAPI.ListPolicyTypes(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.ListPolicyTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPolicyTypes`: ListPolicyTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.ListPolicyTypes`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListPolicyTypesRequest struct via the builder pattern


### Return type

[**ListPolicyTypes200Response**](ListPolicyTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemovePolicies

> RemoveExecuteSchedules200Response RemovePolicies(ctx, id).Execute()

Deletes a Policy



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
	resp, r, err := apiClient.PoliciesAPI.RemovePolicies(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.RemovePolicies``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemovePolicies`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.RemovePolicies`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemovePoliciesRequest struct via the builder pattern


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


## RemovePoliciesCloud

> RemoveExecuteSchedules200Response RemovePoliciesCloud(ctx, cloudId, id).Execute()

Deletes a Policy for a Cloud



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
	cloudId := int64(7) // int64 | The ID of the cloud
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.RemovePoliciesCloud(context.Background(), cloudId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.RemovePoliciesCloud``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemovePoliciesCloud`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.RemovePoliciesCloud`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**cloudId** | **int64** | The ID of the cloud | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemovePoliciesCloudRequest struct via the builder pattern


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


## RemovePoliciesGroup

> RemoveExecuteSchedules200Response RemovePoliciesGroup(ctx, groupId, id).Execute()

Deletes a Policy for a Group



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
	groupId := int64(7) // int64 | The ID of the group
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.RemovePoliciesGroup(context.Background(), groupId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.RemovePoliciesGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemovePoliciesGroup`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.RemovePoliciesGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**groupId** | **int64** | The ID of the group | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemovePoliciesGroupRequest struct via the builder pattern


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


## UpdatePolicies

> UpdatePolicies200Response UpdatePolicies(ctx, id).UpdatePoliciesRequest(updatePoliciesRequest).Execute()

Updates a Policy



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
	updatePoliciesRequest := *openapiclient.NewUpdatePoliciesRequest(*openapiclient.NewUpdatePoliciesRequestPolicy()) // UpdatePoliciesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.UpdatePolicies(context.Background(), id).UpdatePoliciesRequest(updatePoliciesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.UpdatePolicies``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdatePolicies`: UpdatePolicies200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.UpdatePolicies`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdatePoliciesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updatePoliciesRequest** | [**UpdatePoliciesRequest**](UpdatePoliciesRequest.md) |  | 

### Return type

[**UpdatePolicies200Response**](UpdatePolicies200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdatePoliciesCloud

> UpdatePoliciesCloud200Response UpdatePoliciesCloud(ctx, cloudId, id).UpdatePoliciesCloudRequest(updatePoliciesCloudRequest).Execute()

Updates a Policy for a Cloud



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
	cloudId := int64(7) // int64 | The ID of the cloud
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updatePoliciesCloudRequest := *openapiclient.NewUpdatePoliciesCloudRequest(*openapiclient.NewUpdatePoliciesCloudRequestPolicy()) // UpdatePoliciesCloudRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.UpdatePoliciesCloud(context.Background(), cloudId, id).UpdatePoliciesCloudRequest(updatePoliciesCloudRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.UpdatePoliciesCloud``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdatePoliciesCloud`: UpdatePoliciesCloud200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.UpdatePoliciesCloud`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**cloudId** | **int64** | The ID of the cloud | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdatePoliciesCloudRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updatePoliciesCloudRequest** | [**UpdatePoliciesCloudRequest**](UpdatePoliciesCloudRequest.md) |  | 

### Return type

[**UpdatePoliciesCloud200Response**](UpdatePoliciesCloud200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdatePoliciesGroup

> UpdatePoliciesGroup200Response UpdatePoliciesGroup(ctx, groupId, id).UpdatePoliciesGroupRequest(updatePoliciesGroupRequest).Execute()

Updates a Policy for a Group



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
	groupId := int64(7) // int64 | The ID of the group
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updatePoliciesGroupRequest := *openapiclient.NewUpdatePoliciesGroupRequest(*openapiclient.NewUpdatePoliciesGroupRequestPolicy()) // UpdatePoliciesGroupRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoliciesAPI.UpdatePoliciesGroup(context.Background(), groupId, id).UpdatePoliciesGroupRequest(updatePoliciesGroupRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoliciesAPI.UpdatePoliciesGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdatePoliciesGroup`: UpdatePoliciesGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `PoliciesAPI.UpdatePoliciesGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**groupId** | **int64** | The ID of the group | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdatePoliciesGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updatePoliciesGroupRequest** | [**UpdatePoliciesGroupRequest**](UpdatePoliciesGroupRequest.md) |  | 

### Return type

[**UpdatePoliciesGroup200Response**](UpdatePoliciesGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

