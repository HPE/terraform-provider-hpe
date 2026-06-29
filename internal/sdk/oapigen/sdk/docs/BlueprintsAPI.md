# \BlueprintsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddBlueprint**](BlueprintsAPI.md#AddBlueprint) | **Post** /api/blueprints | Create a Blueprint
[**DeleteBlueprint**](BlueprintsAPI.md#DeleteBlueprint) | **Delete** /api/blueprints/{id} | Delete a Blueprint
[**GetBlueprint**](BlueprintsAPI.md#GetBlueprint) | **Get** /api/blueprints/{id} | Get a Specific Blueprint
[**ListBlueprints**](BlueprintsAPI.md#ListBlueprints) | **Get** /api/blueprints | Get All Blueprints
[**UpdateBlueprint**](BlueprintsAPI.md#UpdateBlueprint) | **Put** /api/blueprints/{id} | Updating a Blueprint
[**UpdateBlueprintImage**](BlueprintsAPI.md#UpdateBlueprintImage) | **Post** /api/blueprints/{id}/image | Update Blueprint Image
[**UpdateBlueprintPermissions**](BlueprintsAPI.md#UpdateBlueprintPermissions) | **Put** /api/blueprints/{id}/update-permissions | Update Blueprint Permissions



## AddBlueprint

> AddBlueprint200Response AddBlueprint(ctx).AddBlueprintRequest(addBlueprintRequest).Execute()

Create a Blueprint



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
	addBlueprintRequest := openapiclient.addBlueprint_request{AddBlueprintRequestOneOf: openapiclient.NewAddBlueprintRequestOneOf("Name_example", "Type_example", *openapiclient.NewAddBlueprintRequestOneOfArm("ConfigType_example"))} // AddBlueprintRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BlueprintsAPI.AddBlueprint(context.Background()).AddBlueprintRequest(addBlueprintRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BlueprintsAPI.AddBlueprint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddBlueprint`: AddBlueprint200Response
	fmt.Fprintf(os.Stdout, "Response from `BlueprintsAPI.AddBlueprint`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddBlueprintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addBlueprintRequest** | [**AddBlueprintRequest**](AddBlueprintRequest.md) |  | 

### Return type

[**AddBlueprint200Response**](AddBlueprint200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteBlueprint

> DeleteBlueprint200Response DeleteBlueprint(ctx, id).Execute()

Delete a Blueprint



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
	resp, r, err := apiClient.BlueprintsAPI.DeleteBlueprint(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BlueprintsAPI.DeleteBlueprint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteBlueprint`: DeleteBlueprint200Response
	fmt.Fprintf(os.Stdout, "Response from `BlueprintsAPI.DeleteBlueprint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteBlueprintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteBlueprint200Response**](DeleteBlueprint200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetBlueprint

> GetBlueprint200Response GetBlueprint(ctx, id).Execute()

Get a Specific Blueprint



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
	resp, r, err := apiClient.BlueprintsAPI.GetBlueprint(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BlueprintsAPI.GetBlueprint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetBlueprint`: GetBlueprint200Response
	fmt.Fprintf(os.Stdout, "Response from `BlueprintsAPI.GetBlueprint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetBlueprintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetBlueprint200Response**](GetBlueprint200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListBlueprints

> ListBlueprints200Response ListBlueprints(ctx).Max(max).Offset(offset).Name(name).Phrase(phrase).Labels(labels).AllLabels(allLabels).Execute()

Get All Blueprints



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
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BlueprintsAPI.ListBlueprints(context.Background()).Max(max).Offset(offset).Name(name).Phrase(phrase).Labels(labels).AllLabels(allLabels).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BlueprintsAPI.ListBlueprints``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListBlueprints`: ListBlueprints200Response
	fmt.Fprintf(os.Stdout, "Response from `BlueprintsAPI.ListBlueprints`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListBlueprintsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 

### Return type

[**ListBlueprints200Response**](ListBlueprints200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateBlueprint

> UpdateBlueprint200Response UpdateBlueprint(ctx, id).UpdateBlueprintRequest(updateBlueprintRequest).Execute()

Updating a Blueprint



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
	updateBlueprintRequest := openapiclient.updateBlueprint_request{UpdateBlueprintRequestOneOf: openapiclient.NewUpdateBlueprintRequestOneOf("Name_example", "Type_example", *openapiclient.NewUpdateBlueprintRequestOneOfArm("ConfigType_example"))} // UpdateBlueprintRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BlueprintsAPI.UpdateBlueprint(context.Background(), id).UpdateBlueprintRequest(updateBlueprintRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BlueprintsAPI.UpdateBlueprint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateBlueprint`: UpdateBlueprint200Response
	fmt.Fprintf(os.Stdout, "Response from `BlueprintsAPI.UpdateBlueprint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateBlueprintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateBlueprintRequest** | [**UpdateBlueprintRequest**](UpdateBlueprintRequest.md) |  | 

### Return type

[**UpdateBlueprint200Response**](UpdateBlueprint200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateBlueprintImage

> UpdateBlueprint200Response UpdateBlueprintImage(ctx, id).TemplateImage(templateImage).Execute()

Update Blueprint Image



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
	templateImage := os.NewFile(1234, "some_file") // *os.File |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BlueprintsAPI.UpdateBlueprintImage(context.Background(), id).TemplateImage(templateImage).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BlueprintsAPI.UpdateBlueprintImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateBlueprintImage`: UpdateBlueprint200Response
	fmt.Fprintf(os.Stdout, "Response from `BlueprintsAPI.UpdateBlueprintImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateBlueprintImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **templateImage** | ***os.File** |  | 

### Return type

[**UpdateBlueprint200Response**](UpdateBlueprint200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateBlueprintPermissions

> UpdateBlueprint200Response UpdateBlueprintPermissions(ctx, id).UpdateBlueprintPermissionsRequest(updateBlueprintPermissionsRequest).Execute()

Update Blueprint Permissions



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
	updateBlueprintPermissionsRequest := *openapiclient.NewUpdateBlueprintPermissionsRequest() // UpdateBlueprintPermissionsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BlueprintsAPI.UpdateBlueprintPermissions(context.Background(), id).UpdateBlueprintPermissionsRequest(updateBlueprintPermissionsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BlueprintsAPI.UpdateBlueprintPermissions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateBlueprintPermissions`: UpdateBlueprint200Response
	fmt.Fprintf(os.Stdout, "Response from `BlueprintsAPI.UpdateBlueprintPermissions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateBlueprintPermissionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateBlueprintPermissionsRequest** | [**UpdateBlueprintPermissionsRequest**](UpdateBlueprintPermissionsRequest.md) |  | 

### Return type

[**UpdateBlueprint200Response**](UpdateBlueprint200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

