# \ClusterPackagesAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddClusterPackage**](ClusterPackagesAPI.md#AddClusterPackage) | **Post** /api/library/cluster-packages | Create a Cluster Package
[**DeleteClusterPackage**](ClusterPackagesAPI.md#DeleteClusterPackage) | **Delete** /api/library/cluster-packages/{id} | Delete a Cluster Package
[**GetClusterPackage**](ClusterPackagesAPI.md#GetClusterPackage) | **Get** /api/library/cluster-packages/{id} | Get a Specific Cluster Package
[**ListClusterPackages**](ClusterPackagesAPI.md#ListClusterPackages) | **Get** /api/library/cluster-packages | Get All Cluster Packages
[**UpdateClusterPackage**](ClusterPackagesAPI.md#UpdateClusterPackage) | **Put** /api/library/cluster-packages/{id} | Update a Cluster Package



## AddClusterPackage

> AddClusterPackage200Response AddClusterPackage(ctx).AddClusterPackageRequest(addClusterPackageRequest).Execute()

Create a Cluster Package



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
	addClusterPackageRequest := *openapiclient.NewAddClusterPackageRequest() // AddClusterPackageRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClusterPackagesAPI.AddClusterPackage(context.Background()).AddClusterPackageRequest(addClusterPackageRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClusterPackagesAPI.AddClusterPackage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddClusterPackage`: AddClusterPackage200Response
	fmt.Fprintf(os.Stdout, "Response from `ClusterPackagesAPI.AddClusterPackage`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddClusterPackageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addClusterPackageRequest** | [**AddClusterPackageRequest**](AddClusterPackageRequest.md) |  | 

### Return type

[**AddClusterPackage200Response**](AddClusterPackage200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterPackage

> DeleteClusterPackage200Response DeleteClusterPackage(ctx, id).Execute()

Delete a Cluster Package



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
	resp, r, err := apiClient.ClusterPackagesAPI.DeleteClusterPackage(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClusterPackagesAPI.DeleteClusterPackage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteClusterPackage`: DeleteClusterPackage200Response
	fmt.Fprintf(os.Stdout, "Response from `ClusterPackagesAPI.DeleteClusterPackage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterPackageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteClusterPackage200Response**](DeleteClusterPackage200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterPackage

> GetClusterPackage200Response GetClusterPackage(ctx, id).Execute()

Get a Specific Cluster Package



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
	resp, r, err := apiClient.ClusterPackagesAPI.GetClusterPackage(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClusterPackagesAPI.GetClusterPackage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterPackage`: GetClusterPackage200Response
	fmt.Fprintf(os.Stdout, "Response from `ClusterPackagesAPI.GetClusterPackage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterPackageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetClusterPackage200Response**](GetClusterPackage200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterPackages

> ListClusterPackages200Response ListClusterPackages(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Description(description).Phrase(phrase).Execute()

Get All Cluster Packages



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
	name := "example" // string | Filter by name (optional)
	description := "The desription of my cool object" // string | Filter by description, wildcard may be specified as %. eg. `example-%` (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClusterPackagesAPI.ListClusterPackages(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Description(description).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClusterPackagesAPI.ListClusterPackages``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterPackages`: ListClusterPackages200Response
	fmt.Fprintf(os.Stdout, "Response from `ClusterPackagesAPI.ListClusterPackages`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListClusterPackagesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **name** | **string** | Filter by name | 
 **description** | **string** | Filter by description, wildcard may be specified as %. eg. &#x60;example-%&#x60; | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListClusterPackages200Response**](ListClusterPackages200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateClusterPackage

> UpdateClusterPackage200Response UpdateClusterPackage(ctx, id).UpdateClusterPackageRequest(updateClusterPackageRequest).Execute()

Update a Cluster Package



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
	updateClusterPackageRequest := *openapiclient.NewUpdateClusterPackageRequest() // UpdateClusterPackageRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClusterPackagesAPI.UpdateClusterPackage(context.Background(), id).UpdateClusterPackageRequest(updateClusterPackageRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClusterPackagesAPI.UpdateClusterPackage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateClusterPackage`: UpdateClusterPackage200Response
	fmt.Fprintf(os.Stdout, "Response from `ClusterPackagesAPI.UpdateClusterPackage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateClusterPackageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateClusterPackageRequest** | [**UpdateClusterPackageRequest**](UpdateClusterPackageRequest.md) |  | 

### Return type

[**UpdateClusterPackage200Response**](UpdateClusterPackage200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

