# \LibraryAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddFileTemplate**](LibraryAPI.md#AddFileTemplate) | **Post** /api/library/container-templates | Create a File Template
[**AddInstanceType**](LibraryAPI.md#AddInstanceType) | **Post** /api/library/instance-types | Create an Instance Type
[**AddLayout**](LibraryAPI.md#AddLayout) | **Post** /api/library/instance-types/{instanceTypeId}/layouts | Create a Layout
[**AddNodeType**](LibraryAPI.md#AddNodeType) | **Post** /api/library/container-types | Create a Node Type
[**AddOptionForm**](LibraryAPI.md#AddOptionForm) | **Post** /api/library/option-type-forms | Create an Option Form
[**AddOptionList**](LibraryAPI.md#AddOptionList) | **Post** /api/library/option-type-lists | Create an Option List
[**AddOptionType**](LibraryAPI.md#AddOptionType) | **Post** /api/library/option-types | Create an Input
[**AddOsTypeImage**](LibraryAPI.md#AddOsTypeImage) | **Post** /api/library/operating-systems/os-types/create-image | Create an OsTypeImage
[**AddOsTypes**](LibraryAPI.md#AddOsTypes) | **Post** /api/library/operating-systems/os-types | Create an OsType
[**AddScript**](LibraryAPI.md#AddScript) | **Post** /api/library/container-scripts | Create a Script
[**AddSpecTemplate**](LibraryAPI.md#AddSpecTemplate) | **Post** /api/library/spec-templates | Create a Spec Template
[**AddVirtualImage**](LibraryAPI.md#AddVirtualImage) | **Post** /api/virtual-images | Create a Virtual Image
[**AddVirtualImageFile**](LibraryAPI.md#AddVirtualImageFile) | **Post** /api/virtual-images/{virtualImageId}/upload | Upload Virtual Image File
[**ConvertImage**](LibraryAPI.md#ConvertImage) | **Post** /api/virtual-images/{virtualImageId}/convert | Convert Virtual Image File
[**DeleteFileTemplate**](LibraryAPI.md#DeleteFileTemplate) | **Delete** /api/library/container-templates/{id} | Delete a File Template
[**DeleteInstanceType**](LibraryAPI.md#DeleteInstanceType) | **Delete** /api/library/instance-types/{instanceTypeId} | Delete an Instance Type
[**DeleteLayout**](LibraryAPI.md#DeleteLayout) | **Delete** /api/library/layouts/{id} | Delete a Layout
[**DeleteNodeType**](LibraryAPI.md#DeleteNodeType) | **Delete** /api/library/container-types/{id} | Delete a Node Type
[**DeleteOptionForm**](LibraryAPI.md#DeleteOptionForm) | **Delete** /api/library/option-type-forms/{id} | Delete an Option Form
[**DeleteOptionList**](LibraryAPI.md#DeleteOptionList) | **Delete** /api/library/option-type-lists/{id} | Delete an Option List
[**DeleteOptionType**](LibraryAPI.md#DeleteOptionType) | **Delete** /api/library/option-types/{id} | Delete an Input
[**DeleteOsType**](LibraryAPI.md#DeleteOsType) | **Delete** /api/library/operating-systems/os-types/{id} | Delete an OsType
[**DeleteOsTypeImage**](LibraryAPI.md#DeleteOsTypeImage) | **Delete** /api/library/operating-systems/os-types/images/{id} | Delete an OsTypeImage
[**DeleteScript**](LibraryAPI.md#DeleteScript) | **Delete** /api/library/container-scripts/{id} | Delete a Script
[**DeleteSpecTemplate**](LibraryAPI.md#DeleteSpecTemplate) | **Delete** /api/library/spec-templates/{id} | Delete a Spec Template
[**DownloadVirtualImage**](LibraryAPI.md#DownloadVirtualImage) | **Get** /api/virtual-images/{virtualImageId}/download | Download Virtual Image
[**GetFileTemplate**](LibraryAPI.md#GetFileTemplate) | **Get** /api/library/container-templates/{id} | Get a Specific File Template
[**GetInput**](LibraryAPI.md#GetInput) | **Get** /api/library/option-types/{id} | Get A Specific Input
[**GetInstanceType**](LibraryAPI.md#GetInstanceType) | **Get** /api/library/instance-types/{instanceTypeId} | Get a Specific Instance Type
[**GetLayout**](LibraryAPI.md#GetLayout) | **Get** /api/library/layouts/{id} | Get a Specific Layout
[**GetNodeType**](LibraryAPI.md#GetNodeType) | **Get** /api/library/container-types/{id} | Get a Specific Node Type
[**GetOptionForm**](LibraryAPI.md#GetOptionForm) | **Get** /api/library/option-type-forms/{id} | Get a Specific Option Form
[**GetOptionList**](LibraryAPI.md#GetOptionList) | **Get** /api/library/option-type-lists/{id} | Get a Specific Option List
[**GetOptionListItems**](LibraryAPI.md#GetOptionListItems) | **Get** /api/library/option-type-lists/{id}/items | List Items for a Specific Option List
[**GetOsType**](LibraryAPI.md#GetOsType) | **Get** /api/library/operating-systems/os-types/{id} | Get an OsType
[**GetOsTypeImage**](LibraryAPI.md#GetOsTypeImage) | **Get** /api/library/operating-systems/os-types/images/{id} | Get an OsTypeImage
[**GetScript**](LibraryAPI.md#GetScript) | **Get** /api/library/container-scripts/{id} | Get a Specific Script
[**GetSecurityPackageType**](LibraryAPI.md#GetSecurityPackageType) | **Get** /api/security-package-types/{id} | Retrieves a Specific Security Package Type
[**GetSpecTemplate**](LibraryAPI.md#GetSpecTemplate) | **Get** /api/library/spec-templates/{id} | Get a Specific Spec Template
[**GetVirtualImage**](LibraryAPI.md#GetVirtualImage) | **Get** /api/virtual-images/{virtualImageId} | Get a Specific Virtual Image
[**ListFileTemplates**](LibraryAPI.md#ListFileTemplates) | **Get** /api/library/container-templates | Get All File Templates
[**ListInputs**](LibraryAPI.md#ListInputs) | **Get** /api/library/option-types | Get All Inputs
[**ListInstanceTypes**](LibraryAPI.md#ListInstanceTypes) | **Get** /api/library/instance-types | Get All Instance Types
[**ListLayouts**](LibraryAPI.md#ListLayouts) | **Get** /api/library/layouts | Get All Layouts
[**ListLayoutsForInstanceType**](LibraryAPI.md#ListLayoutsForInstanceType) | **Get** /api/library/instance-types/{instanceTypeId}/layouts | Get All Layouts For an Instance Type
[**ListNodeTypes**](LibraryAPI.md#ListNodeTypes) | **Get** /api/library/container-types | Get All Node Types
[**ListOptionForms**](LibraryAPI.md#ListOptionForms) | **Get** /api/library/option-type-forms | Get All Option Forms
[**ListOptionLists**](LibraryAPI.md#ListOptionLists) | **Get** /api/library/option-type-lists | Get All Option Lists
[**ListOsTypes**](LibraryAPI.md#ListOsTypes) | **Get** /api/library/operating-systems/os-types | Retrieves all OsTypes
[**ListScripts**](LibraryAPI.md#ListScripts) | **Get** /api/library/container-scripts | Get All Scripts
[**ListSecurityPackageTypes**](LibraryAPI.md#ListSecurityPackageTypes) | **Get** /api/security-package-types | Retrieves all Security Package Types
[**ListSpecTemplates**](LibraryAPI.md#ListSpecTemplates) | **Get** /api/library/spec-templates | Get All Spec Templates
[**ListVirtualImageLocations**](LibraryAPI.md#ListVirtualImageLocations) | **Get** /api/virtual-images/{virtualImageId}/locations | Get a List of Virtual Image Locations
[**ListVirtualImages**](LibraryAPI.md#ListVirtualImages) | **Get** /api/virtual-images | Get List of Virtual Images
[**RemoveSecurityScans**](LibraryAPI.md#RemoveSecurityScans) | **Delete** /api/security-scans/{id} | Deletes a Security Scan
[**RemoveVirtualImage**](LibraryAPI.md#RemoveVirtualImage) | **Delete** /api/virtual-images/{virtualImageId} | Delete a Virtual Image
[**RemoveVirtualImageFile**](LibraryAPI.md#RemoveVirtualImageFile) | **Delete** /api/virtual-images/{virtualImageId}/files | Remove Virtual Image File
[**RemoveVirtualImageLocation**](LibraryAPI.md#RemoveVirtualImageLocation) | **Delete** /api/virtual-images/{virtualImageId}/locations/{id} | Delete a Virtual Image Location
[**SetInstanceTypeFeatured**](LibraryAPI.md#SetInstanceTypeFeatured) | **Put** /api/library/instance-types/{instanceTypeId}/toggle-featured | Toggle Featured For Instance Type
[**UpdateFileTemplate**](LibraryAPI.md#UpdateFileTemplate) | **Put** /api/library/container-templates/{id} | Update a File Template
[**UpdateInstanceType**](LibraryAPI.md#UpdateInstanceType) | **Put** /api/library/instance-types/{instanceTypeId} | Update an Instance Type
[**UpdateInstanceTypeLogo**](LibraryAPI.md#UpdateInstanceTypeLogo) | **Post** /api/library/instance-types/{instanceTypeId}/update-logo | Update Logo For Instance Type
[**UpdateLayout**](LibraryAPI.md#UpdateLayout) | **Put** /api/library/layouts/{id} | Update a Layout
[**UpdateLayoutPermissions**](LibraryAPI.md#UpdateLayoutPermissions) | **Post** /api/library/layouts/{id}/permissions | Update Layout Permissions
[**UpdateNodeType**](LibraryAPI.md#UpdateNodeType) | **Put** /api/library/container-types/{id} | Update a Node Type
[**UpdateOptionForm**](LibraryAPI.md#UpdateOptionForm) | **Put** /api/library/option-type-forms/{id} | Update an Option Form
[**UpdateOptionList**](LibraryAPI.md#UpdateOptionList) | **Put** /api/library/option-type-lists/{id} | Update an Option List
[**UpdateOptionType**](LibraryAPI.md#UpdateOptionType) | **Put** /api/library/option-types/{id} | Update an Input
[**UpdateOsType**](LibraryAPI.md#UpdateOsType) | **Put** /api/library/operating-systems/os-types/{id} | Update an OsType
[**UpdateScript**](LibraryAPI.md#UpdateScript) | **Put** /api/library/container-scripts/{id} | Update a Script
[**UpdateSpecTemplate**](LibraryAPI.md#UpdateSpecTemplate) | **Put** /api/library/spec-templates/{id} | Update a Spec Template
[**UpdateVirtualImage**](LibraryAPI.md#UpdateVirtualImage) | **Put** /api/virtual-images/{virtualImageId} | Update a Virtual Image



## AddFileTemplate

> AddFileTemplate200Response AddFileTemplate(ctx).AddFileTemplateRequest(addFileTemplateRequest).Execute()

Create a File Template



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
	addFileTemplateRequest := *openapiclient.NewAddFileTemplateRequest() // AddFileTemplateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddFileTemplate(context.Background()).AddFileTemplateRequest(addFileTemplateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddFileTemplate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddFileTemplate`: AddFileTemplate200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddFileTemplate`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddFileTemplateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addFileTemplateRequest** | [**AddFileTemplateRequest**](AddFileTemplateRequest.md) |  | 

### Return type

[**AddFileTemplate200Response**](AddFileTemplate200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddInstanceType

> AddInstanceType200Response AddInstanceType(ctx).AddInstanceTypeRequest(addInstanceTypeRequest).Execute()

Create an Instance Type



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
	addInstanceTypeRequest := *openapiclient.NewAddInstanceTypeRequest() // AddInstanceTypeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddInstanceType(context.Background()).AddInstanceTypeRequest(addInstanceTypeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddInstanceType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddInstanceType`: AddInstanceType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddInstanceType`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddInstanceTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addInstanceTypeRequest** | [**AddInstanceTypeRequest**](AddInstanceTypeRequest.md) |  | 

### Return type

[**AddInstanceType200Response**](AddInstanceType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddLayout

> AddLayout200Response AddLayout(ctx, instanceTypeId).AddLayoutRequest(addLayoutRequest).Execute()

Create a Layout



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
	instanceTypeId := int64(2) // int64 | The ID of the instance type
	addLayoutRequest := *openapiclient.NewAddLayoutRequest() // AddLayoutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddLayout(context.Background(), instanceTypeId).AddLayoutRequest(addLayoutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddLayout``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddLayout`: AddLayout200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddLayout`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceTypeId** | **int64** | The ID of the instance type | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddLayoutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addLayoutRequest** | [**AddLayoutRequest**](AddLayoutRequest.md) |  | 

### Return type

[**AddLayout200Response**](AddLayout200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddNodeType

> AddNodeType200Response AddNodeType(ctx).AddNodeTypeRequest(addNodeTypeRequest).Execute()

Create a Node Type



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
	addNodeTypeRequest := *openapiclient.NewAddNodeTypeRequest() // AddNodeTypeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddNodeType(context.Background()).AddNodeTypeRequest(addNodeTypeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddNodeType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddNodeType`: AddNodeType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddNodeType`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddNodeTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addNodeTypeRequest** | [**AddNodeTypeRequest**](AddNodeTypeRequest.md) |  | 

### Return type

[**AddNodeType200Response**](AddNodeType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddOptionForm

> AddOptionForm200Response AddOptionForm(ctx).AddOptionFormRequest(addOptionFormRequest).Execute()

Create an Option Form



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
	addOptionFormRequest := *openapiclient.NewAddOptionFormRequest() // AddOptionFormRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddOptionForm(context.Background()).AddOptionFormRequest(addOptionFormRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddOptionForm``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddOptionForm`: AddOptionForm200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddOptionForm`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddOptionFormRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addOptionFormRequest** | [**AddOptionFormRequest**](AddOptionFormRequest.md) |  | 

### Return type

[**AddOptionForm200Response**](AddOptionForm200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddOptionList

> AddOptionList200Response AddOptionList(ctx).AddOptionListRequest(addOptionListRequest).Execute()

Create an Option List



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
	addOptionListRequest := *openapiclient.NewAddOptionListRequest() // AddOptionListRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddOptionList(context.Background()).AddOptionListRequest(addOptionListRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddOptionList``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddOptionList`: AddOptionList200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddOptionList`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddOptionListRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addOptionListRequest** | [**AddOptionListRequest**](AddOptionListRequest.md) |  | 

### Return type

[**AddOptionList200Response**](AddOptionList200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddOptionType

> AddOptionType200Response AddOptionType(ctx).AddOptionTypeRequest(addOptionTypeRequest).Execute()

Create an Input



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
	addOptionTypeRequest := *openapiclient.NewAddOptionTypeRequest() // AddOptionTypeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddOptionType(context.Background()).AddOptionTypeRequest(addOptionTypeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddOptionType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddOptionType`: AddOptionType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddOptionType`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddOptionTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addOptionTypeRequest** | [**AddOptionTypeRequest**](AddOptionTypeRequest.md) |  | 

### Return type

[**AddOptionType200Response**](AddOptionType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddOsTypeImage

> AddOsTypeImage(ctx).AddOsTypeImageRequest(addOsTypeImageRequest).Execute()

Create an OsTypeImage



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
	addOsTypeImageRequest := *openapiclient.NewAddOsTypeImageRequest() // AddOsTypeImageRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.LibraryAPI.AddOsTypeImage(context.Background()).AddOsTypeImageRequest(addOsTypeImageRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddOsTypeImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddOsTypeImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addOsTypeImageRequest** | [**AddOsTypeImageRequest**](AddOsTypeImageRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddOsTypes

> AddOsTypes200Response AddOsTypes(ctx).AddOsTypesRequest(addOsTypesRequest).Execute()

Create an OsType



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
	addOsTypesRequest := *openapiclient.NewAddOsTypesRequest() // AddOsTypesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddOsTypes(context.Background()).AddOsTypesRequest(addOsTypesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddOsTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddOsTypes`: AddOsTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddOsTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddOsTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addOsTypesRequest** | [**AddOsTypesRequest**](AddOsTypesRequest.md) |  | 

### Return type

[**AddOsTypes200Response**](AddOsTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddScript

> AddScript200Response AddScript(ctx).AddScriptRequest(addScriptRequest).Execute()

Create a Script



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
	addScriptRequest := *openapiclient.NewAddScriptRequest() // AddScriptRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddScript(context.Background()).AddScriptRequest(addScriptRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddScript`: AddScript200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddScript`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addScriptRequest** | [**AddScriptRequest**](AddScriptRequest.md) |  | 

### Return type

[**AddScript200Response**](AddScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddSpecTemplate

> AddSpecTemplate200Response AddSpecTemplate(ctx).AddSpecTemplateRequest(addSpecTemplateRequest).Execute()

Create a Spec Template



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
	addSpecTemplateRequest := *openapiclient.NewAddSpecTemplateRequest() // AddSpecTemplateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddSpecTemplate(context.Background()).AddSpecTemplateRequest(addSpecTemplateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddSpecTemplate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddSpecTemplate`: AddSpecTemplate200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddSpecTemplate`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddSpecTemplateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addSpecTemplateRequest** | [**AddSpecTemplateRequest**](AddSpecTemplateRequest.md) |  | 

### Return type

[**AddSpecTemplate200Response**](AddSpecTemplate200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddVirtualImage

> AddVirtualImage200Response AddVirtualImage(ctx).AddVirtualImageRequest(addVirtualImageRequest).Execute()

Create a Virtual Image



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
	addVirtualImageRequest := *openapiclient.NewAddVirtualImageRequest(*openapiclient.NewAddVirtualImageRequestVirtualImage()) // AddVirtualImageRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddVirtualImage(context.Background()).AddVirtualImageRequest(addVirtualImageRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddVirtualImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddVirtualImage`: AddVirtualImage200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddVirtualImage`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddVirtualImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addVirtualImageRequest** | [**AddVirtualImageRequest**](AddVirtualImageRequest.md) |  | 

### Return type

[**AddVirtualImage200Response**](AddVirtualImage200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddVirtualImageFile

> RemoveExecuteSchedules200Response AddVirtualImageFile(ctx, virtualImageId).Filename(filename).Url(url).Body(body).Execute()

Upload Virtual Image File



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
	virtualImageId := int64(4) // int64 | Virtual Image ID
	filename := "testimage.ovf" // string | The name of the file (optional)
	url := "https://example.com/testimage.ovf" // string | Download the file from a remote url. This can be used instead of uploading a local file. (optional)
	body := os.NewFile(1234, "some_file") // *os.File |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.AddVirtualImageFile(context.Background(), virtualImageId).Filename(filename).Url(url).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.AddVirtualImageFile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddVirtualImageFile`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.AddVirtualImageFile`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**virtualImageId** | **int64** | Virtual Image ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddVirtualImageFileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **filename** | **string** | The name of the file | 
 **url** | **string** | Download the file from a remote url. This can be used instead of uploading a local file. | 
 **body** | ***os.File** |  | 

### Return type

[**RemoveExecuteSchedules200Response**](RemoveExecuteSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/octet-stream
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ConvertImage

> RemoveExecuteSchedules200Response ConvertImage(ctx, virtualImageId).ConvertImageRequest(convertImageRequest).Execute()

Convert Virtual Image File



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
	virtualImageId := int64(4) // int64 | Virtual Image ID
	convertImageRequest := *openapiclient.NewConvertImageRequest() // ConvertImageRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ConvertImage(context.Background(), virtualImageId).ConvertImageRequest(convertImageRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ConvertImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ConvertImage`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ConvertImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**virtualImageId** | **int64** | Virtual Image ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiConvertImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **convertImageRequest** | [**ConvertImageRequest**](ConvertImageRequest.md) |  | 

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


## DeleteFileTemplate

> DeleteFileTemplate200Response DeleteFileTemplate(ctx, id).Execute()

Delete a File Template



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
	resp, r, err := apiClient.LibraryAPI.DeleteFileTemplate(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DeleteFileTemplate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteFileTemplate`: DeleteFileTemplate200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.DeleteFileTemplate`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteFileTemplateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteFileTemplate200Response**](DeleteFileTemplate200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteInstanceType

> DeleteInstanceType200Response DeleteInstanceType(ctx, instanceTypeId).Execute()

Delete an Instance Type



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
	instanceTypeId := int64(2) // int64 | The ID of the instance type

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.DeleteInstanceType(context.Background(), instanceTypeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DeleteInstanceType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteInstanceType`: DeleteInstanceType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.DeleteInstanceType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceTypeId** | **int64** | The ID of the instance type | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteInstanceTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteInstanceType200Response**](DeleteInstanceType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteLayout

> DeleteLayout200Response DeleteLayout(ctx, id).Execute()

Delete a Layout



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
	resp, r, err := apiClient.LibraryAPI.DeleteLayout(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DeleteLayout``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteLayout`: DeleteLayout200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.DeleteLayout`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteLayoutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteLayout200Response**](DeleteLayout200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteNodeType

> DeleteNodeType200Response DeleteNodeType(ctx, id).Execute()

Delete a Node Type



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
	resp, r, err := apiClient.LibraryAPI.DeleteNodeType(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DeleteNodeType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteNodeType`: DeleteNodeType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.DeleteNodeType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteNodeTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteNodeType200Response**](DeleteNodeType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOptionForm

> DeleteOptionForm200Response DeleteOptionForm(ctx, id).Execute()

Delete an Option Form



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
	resp, r, err := apiClient.LibraryAPI.DeleteOptionForm(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DeleteOptionForm``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteOptionForm`: DeleteOptionForm200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.DeleteOptionForm`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOptionFormRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteOptionForm200Response**](DeleteOptionForm200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOptionList

> DeleteOptionList200Response DeleteOptionList(ctx, id).Execute()

Delete an Option List



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
	resp, r, err := apiClient.LibraryAPI.DeleteOptionList(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DeleteOptionList``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteOptionList`: DeleteOptionList200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.DeleteOptionList`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOptionListRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteOptionList200Response**](DeleteOptionList200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOptionType

> DeleteOptionType200Response DeleteOptionType(ctx, id).Execute()

Delete an Input



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
	resp, r, err := apiClient.LibraryAPI.DeleteOptionType(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DeleteOptionType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteOptionType`: DeleteOptionType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.DeleteOptionType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOptionTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteOptionType200Response**](DeleteOptionType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOsType

> DeleteOsType200Response DeleteOsType(ctx, id).Execute()

Delete an OsType



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
	resp, r, err := apiClient.LibraryAPI.DeleteOsType(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DeleteOsType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteOsType`: DeleteOsType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.DeleteOsType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOsTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteOsType200Response**](DeleteOsType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOsTypeImage

> DeleteOsTypeImage200Response DeleteOsTypeImage(ctx, id).Execute()

Delete an OsTypeImage



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
	resp, r, err := apiClient.LibraryAPI.DeleteOsTypeImage(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DeleteOsTypeImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteOsTypeImage`: DeleteOsTypeImage200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.DeleteOsTypeImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOsTypeImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteOsTypeImage200Response**](DeleteOsTypeImage200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteScript

> DeleteScript200Response DeleteScript(ctx, id).Execute()

Delete a Script



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
	resp, r, err := apiClient.LibraryAPI.DeleteScript(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DeleteScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteScript`: DeleteScript200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.DeleteScript`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteScript200Response**](DeleteScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteSpecTemplate

> DeleteSpecTemplate200Response DeleteSpecTemplate(ctx, id).Execute()

Delete a Spec Template



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
	resp, r, err := apiClient.LibraryAPI.DeleteSpecTemplate(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DeleteSpecTemplate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteSpecTemplate`: DeleteSpecTemplate200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.DeleteSpecTemplate`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteSpecTemplateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteSpecTemplate200Response**](DeleteSpecTemplate200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DownloadVirtualImage

> DownloadVirtualImage(ctx, virtualImageId).Execute()

Download Virtual Image



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
	virtualImageId := int64(4) // int64 | Virtual Image ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.LibraryAPI.DownloadVirtualImage(context.Background(), virtualImageId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.DownloadVirtualImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**virtualImageId** | **int64** | Virtual Image ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDownloadVirtualImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetFileTemplate

> GetFileTemplate200Response GetFileTemplate(ctx, id).Execute()

Get a Specific File Template



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
	resp, r, err := apiClient.LibraryAPI.GetFileTemplate(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetFileTemplate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetFileTemplate`: GetFileTemplate200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetFileTemplate`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetFileTemplateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetFileTemplate200Response**](GetFileTemplate200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInput

> GetInput200Response GetInput(ctx, id).Execute()

Get A Specific Input



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
	resp, r, err := apiClient.LibraryAPI.GetInput(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetInput``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInput`: GetInput200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetInput`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInputRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetInput200Response**](GetInput200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceType

> GetInstanceType200Response GetInstanceType(ctx, instanceTypeId).Execute()

Get a Specific Instance Type



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
	instanceTypeId := int64(2) // int64 | The ID of the instance type

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.GetInstanceType(context.Background(), instanceTypeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetInstanceType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceType`: GetInstanceType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetInstanceType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceTypeId** | **int64** | The ID of the instance type | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetInstanceType200Response**](GetInstanceType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLayout

> GetLayout200Response GetLayout(ctx, id).Execute()

Get a Specific Layout



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
	resp, r, err := apiClient.LibraryAPI.GetLayout(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetLayout``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLayout`: GetLayout200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetLayout`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetLayoutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetLayout200Response**](GetLayout200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNodeType

> GetNodeType200Response GetNodeType(ctx, id).Execute()

Get a Specific Node Type



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
	resp, r, err := apiClient.LibraryAPI.GetNodeType(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetNodeType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodeType`: GetNodeType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetNodeType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodeTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetNodeType200Response**](GetNodeType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOptionForm

> GetOptionForm200Response GetOptionForm(ctx, id).Execute()

Get a Specific Option Form



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
	resp, r, err := apiClient.LibraryAPI.GetOptionForm(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetOptionForm``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOptionForm`: GetOptionForm200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetOptionForm`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOptionFormRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetOptionForm200Response**](GetOptionForm200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOptionList

> GetOptionList200Response GetOptionList(ctx, id).Execute()

Get a Specific Option List



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
	resp, r, err := apiClient.LibraryAPI.GetOptionList(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetOptionList``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOptionList`: GetOptionList200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetOptionList`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOptionListRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetOptionList200Response**](GetOptionList200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOptionListItems

> GetOptionListItems200Response GetOptionListItems(ctx, id).Execute()

List Items for a Specific Option List



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
	resp, r, err := apiClient.LibraryAPI.GetOptionListItems(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetOptionListItems``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOptionListItems`: GetOptionListItems200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetOptionListItems`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOptionListItemsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetOptionListItems200Response**](GetOptionListItems200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOsType

> GetOsType200Response GetOsType(ctx, id).Execute()

Get an OsType



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
	resp, r, err := apiClient.LibraryAPI.GetOsType(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetOsType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOsType`: GetOsType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetOsType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOsTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetOsType200Response**](GetOsType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOsTypeImage

> GetOsTypeImage200Response GetOsTypeImage(ctx, id).Execute()

Get an OsTypeImage



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
	resp, r, err := apiClient.LibraryAPI.GetOsTypeImage(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetOsTypeImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOsTypeImage`: GetOsTypeImage200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetOsTypeImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOsTypeImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetOsTypeImage200Response**](GetOsTypeImage200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetScript

> GetScript200Response GetScript(ctx, id).Execute()

Get a Specific Script



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
	resp, r, err := apiClient.LibraryAPI.GetScript(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetScript`: GetScript200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetScript`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetScript200Response**](GetScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSecurityPackageType

> GetSecurityPackageType200Response GetSecurityPackageType(ctx, id).Execute()

Retrieves a Specific Security Package Type



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
	resp, r, err := apiClient.LibraryAPI.GetSecurityPackageType(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetSecurityPackageType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSecurityPackageType`: GetSecurityPackageType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetSecurityPackageType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSecurityPackageTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetSecurityPackageType200Response**](GetSecurityPackageType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSpecTemplate

> GetSpecTemplate200Response GetSpecTemplate(ctx, id).Execute()

Get a Specific Spec Template



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
	resp, r, err := apiClient.LibraryAPI.GetSpecTemplate(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetSpecTemplate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSpecTemplate`: GetSpecTemplate200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetSpecTemplate`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSpecTemplateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetSpecTemplate200Response**](GetSpecTemplate200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVirtualImage

> GetVirtualImage200Response GetVirtualImage(ctx, virtualImageId).Execute()

Get a Specific Virtual Image



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
	virtualImageId := int64(4) // int64 | Virtual Image ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.GetVirtualImage(context.Background(), virtualImageId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.GetVirtualImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVirtualImage`: GetVirtualImage200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.GetVirtualImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**virtualImageId** | **int64** | Virtual Image ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVirtualImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetVirtualImage200Response**](GetVirtualImage200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListFileTemplates

> ListFileTemplates200Response ListFileTemplates(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Labels(labels).AllLabels(allLabels).FileName(fileName).Execute()

Get All File Templates



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
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)
	fileName := "testtext.txt" // string | Filename filter, restricts query to only load file template matching fileName specified (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListFileTemplates(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Labels(labels).AllLabels(allLabels).FileName(fileName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListFileTemplates``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListFileTemplates`: ListFileTemplates200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListFileTemplates`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListFileTemplatesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 
 **fileName** | **string** | Filename filter, restricts query to only load file template matching fileName specified | 

### Return type

[**ListFileTemplates200Response**](ListFileTemplates200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListInputs

> ListInputs200Response ListInputs(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).Labels(labels).AllLabels(allLabels).FieldName(fieldName).FieldContext(fieldContext).FieldLabel(fieldLabel).Execute()

Get All Inputs



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
	code := "azr" // string | If specified will return an exact match on code (optional)
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)
	fieldName := "cloudName" // string | Field Name filter, restricts query to only load type matching fieldName specified (optional)
	fieldContext := "config.customOptions" // string | Field Context filter, restricts query to only load type matching fieldContext specified (optional)
	fieldLabel := "DB Version" // string | Field Label filter, restricts query to only load type matching fieldLabel specified (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListInputs(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).Labels(labels).AllLabels(allLabels).FieldName(fieldName).FieldContext(fieldContext).FieldLabel(fieldLabel).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListInputs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListInputs`: ListInputs200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListInputs`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListInputsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 
 **fieldName** | **string** | Field Name filter, restricts query to only load type matching fieldName specified | 
 **fieldContext** | **string** | Field Context filter, restricts query to only load type matching fieldContext specified | 
 **fieldLabel** | **string** | Field Label filter, restricts query to only load type matching fieldLabel specified | 

### Return type

[**ListInputs200Response**](ListInputs200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListInstanceTypes

> ListInstanceTypes200Response ListInstanceTypes(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).Featured(featured).Labels(labels).AllLabels(allLabels).Details(details).Execute()

Get All Instance Types



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
	code := "azr" // string | If specified will return an exact match on code (optional)
	featured := false // bool | Filter by featured (optional)
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)
	details := true // bool | Load full details including optionTypes, environmentVariables, etc. These properties are excluded by default. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListInstanceTypes(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).Featured(featured).Labels(labels).AllLabels(allLabels).Details(details).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListInstanceTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListInstanceTypes`: ListInstanceTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListInstanceTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListInstanceTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **featured** | **bool** | Filter by featured | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 
 **details** | **bool** | Load full details including optionTypes, environmentVariables, etc. These properties are excluded by default. | 

### Return type

[**ListInstanceTypes200Response**](ListInstanceTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListLayouts

> ListLayouts200Response ListLayouts(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).ProvisionType(provisionType).Labels(labels).AllLabels(allLabels).Execute()

Get All Layouts



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
	code := "azr" // string | If specified will return an exact match on code (optional)
	provisionType := "provisionType_example" // string | Filter by `Provision Type` code. Refer to `Provision Types` API for up to date listings.  (optional)
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListLayouts(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).ProvisionType(provisionType).Labels(labels).AllLabels(allLabels).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListLayouts``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLayouts`: ListLayouts200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListLayouts`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListLayoutsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **provisionType** | **string** | Filter by &#x60;Provision Type&#x60; code. Refer to &#x60;Provision Types&#x60; API for up to date listings.  | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 

### Return type

[**ListLayouts200Response**](ListLayouts200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListLayoutsForInstanceType

> ListLayoutsForInstanceType200Response ListLayoutsForInstanceType(ctx, instanceTypeId).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).ProvisionType(provisionType).Labels(labels).AllLabels(allLabels).Execute()

Get All Layouts For an Instance Type



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
	instanceTypeId := int64(2) // int64 | The ID of the instance type
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	code := "azr" // string | If specified will return an exact match on code (optional)
	provisionType := "provisionType_example" // string | Filter by `Provision Type` code. Refer to `Provision Types` API for up to date listings.  (optional)
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListLayoutsForInstanceType(context.Background(), instanceTypeId).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).ProvisionType(provisionType).Labels(labels).AllLabels(allLabels).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListLayoutsForInstanceType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLayoutsForInstanceType`: ListLayoutsForInstanceType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListLayoutsForInstanceType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceTypeId** | **int64** | The ID of the instance type | 

### Other Parameters

Other parameters are passed through a pointer to a apiListLayoutsForInstanceTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **provisionType** | **string** | Filter by &#x60;Provision Type&#x60; code. Refer to &#x60;Provision Types&#x60; API for up to date listings.  | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 

### Return type

[**ListLayoutsForInstanceType200Response**](ListLayoutsForInstanceType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodeTypes

> ListNodeTypes200Response ListNodeTypes(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).ProvisionType(provisionType).Labels(labels).AllLabels(allLabels).Execute()

Get All Node Types



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
	code := "azr" // string | If specified will return an exact match on code (optional)
	provisionType := "provisionType_example" // string | Filter by `Provision Type` code. Refer to `Provision Types` API for up to date listings.  (optional)
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListNodeTypes(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).ProvisionType(provisionType).Labels(labels).AllLabels(allLabels).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListNodeTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodeTypes`: ListNodeTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListNodeTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListNodeTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **provisionType** | **string** | Filter by &#x60;Provision Type&#x60; code. Refer to &#x60;Provision Types&#x60; API for up to date listings.  | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 

### Return type

[**ListNodeTypes200Response**](ListNodeTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOptionForms

> ListOptionForms200Response ListOptionForms(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).Labels(labels).AllLabels(allLabels).Execute()

Get All Option Forms



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
	code := "azr" // string | If specified will return an exact match on code (optional)
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListOptionForms(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).Labels(labels).AllLabels(allLabels).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListOptionForms``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOptionForms`: ListOptionForms200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListOptionForms`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListOptionFormsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 

### Return type

[**ListOptionForms200Response**](ListOptionForms200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOptionLists

> ListOptionLists200Response ListOptionLists(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Labels(labels).AllLabels(allLabels).Execute()

Get All Option Lists



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
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListOptionLists(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Labels(labels).AllLabels(allLabels).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListOptionLists``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOptionLists`: ListOptionLists200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListOptionLists`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListOptionListsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 

### Return type

[**ListOptionLists200Response**](ListOptionLists200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOsTypes

> ListOsTypes200Response ListOsTypes(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Retrieves all OsTypes



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
	resp, r, err := apiClient.LibraryAPI.ListOsTypes(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListOsTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOsTypes`: ListOsTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListOsTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListOsTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListOsTypes200Response**](ListOsTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListScripts

> ListScripts200Response ListScripts(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Labels(labels).AllLabels(allLabels).ScriptType(scriptType).ScriptPhase(scriptPhase).Execute()

Get All Scripts



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
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)
	scriptType := "scriptType_example" // string | Script type code filter, restricts query to only load scripts of specified type (optional)
	scriptPhase := "scriptPhase_example" // string | Script phase filter, restricts query to only load scripts of specified phase (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListScripts(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Labels(labels).AllLabels(allLabels).ScriptType(scriptType).ScriptPhase(scriptPhase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListScripts``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListScripts`: ListScripts200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListScripts`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListScriptsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 
 **scriptType** | **string** | Script type code filter, restricts query to only load scripts of specified type | 
 **scriptPhase** | **string** | Script phase filter, restricts query to only load scripts of specified phase | 

### Return type

[**ListScripts200Response**](ListScripts200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSecurityPackageTypes

> ListSecurityPackageTypes200Response ListSecurityPackageTypes(ctx).Execute()

Retrieves all Security Package Types



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
	resp, r, err := apiClient.LibraryAPI.ListSecurityPackageTypes(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListSecurityPackageTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListSecurityPackageTypes`: ListSecurityPackageTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListSecurityPackageTypes`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListSecurityPackageTypesRequest struct via the builder pattern


### Return type

[**ListSecurityPackageTypes200Response**](ListSecurityPackageTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSpecTemplates

> ListSpecTemplates200Response ListSpecTemplates(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Labels(labels).AllLabels(allLabels).Execute()

Get All Spec Templates



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
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListSpecTemplates(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Labels(labels).AllLabels(allLabels).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListSpecTemplates``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListSpecTemplates`: ListSpecTemplates200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListSpecTemplates`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListSpecTemplatesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 

### Return type

[**ListSpecTemplates200Response**](ListSpecTemplates200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVirtualImageLocations

> ListVirtualImageLocations200Response ListVirtualImageLocations(ctx, virtualImageId).Execute()

Get a List of Virtual Image Locations



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
	virtualImageId := int64(4) // int64 | Virtual Image ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListVirtualImageLocations(context.Background(), virtualImageId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListVirtualImageLocations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListVirtualImageLocations`: ListVirtualImageLocations200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListVirtualImageLocations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**virtualImageId** | **int64** | Virtual Image ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListVirtualImageLocationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListVirtualImageLocations200Response**](ListVirtualImageLocations200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVirtualImages

> ListVirtualImages200Response ListVirtualImages(ctx).Max(max).Offset(offset).Name(name).Phrase(phrase).LastUpdated(lastUpdated).FilterType(filterType).ImageType(imageType).TagsName(tagsName).Labels(labels).AllLabels(allLabels).Execute()

Get List of Virtual Images



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
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	lastUpdated := time.Now() // time.Time | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) (optional)
	filterType := "System" // string | Filter by type, \"User\", \"System\", \"Synced\", or \"All\" (optional) (default to "User")
	imageType := "vmware" // string | Filter by image type code, \"vmware\", \"ami\", etc (optional)
	tagsName := "value" // string | Filter by tags (metadata). This allows filtering by a tag name and value(s)  (optional)
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.ListVirtualImages(context.Background()).Max(max).Offset(offset).Name(name).Phrase(phrase).LastUpdated(lastUpdated).FilterType(filterType).ImageType(imageType).TagsName(tagsName).Labels(labels).AllLabels(allLabels).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.ListVirtualImages``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListVirtualImages`: ListVirtualImages200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.ListVirtualImages`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListVirtualImagesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **lastUpdated** | **time.Time** | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) | 
 **filterType** | **string** | Filter by type, \&quot;User\&quot;, \&quot;System\&quot;, \&quot;Synced\&quot;, or \&quot;All\&quot; | [default to &quot;User&quot;]
 **imageType** | **string** | Filter by image type code, \&quot;vmware\&quot;, \&quot;ami\&quot;, etc | 
 **tagsName** | **string** | Filter by tags (metadata). This allows filtering by a tag name and value(s)  | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 

### Return type

[**ListVirtualImages200Response**](ListVirtualImages200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveSecurityScans

> RemoveExecuteSchedules200Response RemoveSecurityScans(ctx, id).Execute()

Deletes a Security Scan



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
	resp, r, err := apiClient.LibraryAPI.RemoveSecurityScans(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.RemoveSecurityScans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveSecurityScans`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.RemoveSecurityScans`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveSecurityScansRequest struct via the builder pattern


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


## RemoveVirtualImage

> RemoveExecuteSchedules200Response RemoveVirtualImage(ctx, virtualImageId).RemoveFromCloud(removeFromCloud).Execute()

Delete a Virtual Image



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
	virtualImageId := int64(4) // int64 | Virtual Image ID
	removeFromCloud := true // bool | Include removeFromCloud=true to delete the image location files from all clouds. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.RemoveVirtualImage(context.Background(), virtualImageId).RemoveFromCloud(removeFromCloud).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.RemoveVirtualImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveVirtualImage`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.RemoveVirtualImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**virtualImageId** | **int64** | Virtual Image ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveVirtualImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **removeFromCloud** | **bool** | Include removeFromCloud&#x3D;true to delete the image location files from all clouds. | 

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


## RemoveVirtualImageFile

> RemoveExecuteSchedules200Response RemoveVirtualImageFile(ctx, virtualImageId).Filename(filename).Execute()

Remove Virtual Image File



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
	virtualImageId := int64(4) // int64 | Virtual Image ID
	filename := "testimage.ovf" // string | The name of the file (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.RemoveVirtualImageFile(context.Background(), virtualImageId).Filename(filename).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.RemoveVirtualImageFile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveVirtualImageFile`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.RemoveVirtualImageFile`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**virtualImageId** | **int64** | Virtual Image ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveVirtualImageFileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **filename** | **string** | The name of the file | 

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


## RemoveVirtualImageLocation

> RemoveExecuteSchedules200Response RemoveVirtualImageLocation(ctx, virtualImageId, id).RemoveFromCloud(removeFromCloud).Execute()

Delete a Virtual Image Location



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
	virtualImageId := int64(4) // int64 | Virtual Image ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	removeFromCloud := true // bool | Include removeFromCloud=true to delete the image location files in the cloud as well. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.RemoveVirtualImageLocation(context.Background(), virtualImageId, id).RemoveFromCloud(removeFromCloud).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.RemoveVirtualImageLocation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveVirtualImageLocation`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.RemoveVirtualImageLocation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**virtualImageId** | **int64** | Virtual Image ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveVirtualImageLocationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **removeFromCloud** | **bool** | Include removeFromCloud&#x3D;true to delete the image location files in the cloud as well. | 

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


## SetInstanceTypeFeatured

> SetInstanceTypeFeatured200Response SetInstanceTypeFeatured(ctx, instanceTypeId).Execute()

Toggle Featured For Instance Type



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
	instanceTypeId := int64(2) // int64 | The ID of the instance type

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.SetInstanceTypeFeatured(context.Background(), instanceTypeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.SetInstanceTypeFeatured``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetInstanceTypeFeatured`: SetInstanceTypeFeatured200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.SetInstanceTypeFeatured`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceTypeId** | **int64** | The ID of the instance type | 

### Other Parameters

Other parameters are passed through a pointer to a apiSetInstanceTypeFeaturedRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**SetInstanceTypeFeatured200Response**](SetInstanceTypeFeatured200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateFileTemplate

> UpdateFileTemplate200Response UpdateFileTemplate(ctx, id).UpdateFileTemplateRequest(updateFileTemplateRequest).Execute()

Update a File Template



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
	updateFileTemplateRequest := *openapiclient.NewUpdateFileTemplateRequest() // UpdateFileTemplateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateFileTemplate(context.Background(), id).UpdateFileTemplateRequest(updateFileTemplateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateFileTemplate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateFileTemplate`: UpdateFileTemplate200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateFileTemplate`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateFileTemplateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateFileTemplateRequest** | [**UpdateFileTemplateRequest**](UpdateFileTemplateRequest.md) |  | 

### Return type

[**UpdateFileTemplate200Response**](UpdateFileTemplate200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateInstanceType

> UpdateInstanceType200Response UpdateInstanceType(ctx, instanceTypeId).UpdateInstanceTypeRequest(updateInstanceTypeRequest).Execute()

Update an Instance Type



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
	instanceTypeId := int64(2) // int64 | The ID of the instance type
	updateInstanceTypeRequest := *openapiclient.NewUpdateInstanceTypeRequest() // UpdateInstanceTypeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateInstanceType(context.Background(), instanceTypeId).UpdateInstanceTypeRequest(updateInstanceTypeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateInstanceType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateInstanceType`: UpdateInstanceType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateInstanceType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceTypeId** | **int64** | The ID of the instance type | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateInstanceTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateInstanceTypeRequest** | [**UpdateInstanceTypeRequest**](UpdateInstanceTypeRequest.md) |  | 

### Return type

[**UpdateInstanceType200Response**](UpdateInstanceType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateInstanceTypeLogo

> UpdateInstanceTypeLogo200Response UpdateInstanceTypeLogo(ctx, instanceTypeId).Logo(logo).DarkLogo(darkLogo).Execute()

Update Logo For Instance Type



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
	instanceTypeId := int64(2) // int64 | The ID of the instance type
	logo := os.NewFile(1234, "some_file") // *os.File | Logo File png,jpg,svg (optional)
	darkLogo := os.NewFile(1234, "some_file") // *os.File | Dark Logo File png,jpg,svg (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateInstanceTypeLogo(context.Background(), instanceTypeId).Logo(logo).DarkLogo(darkLogo).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateInstanceTypeLogo``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateInstanceTypeLogo`: UpdateInstanceTypeLogo200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateInstanceTypeLogo`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**instanceTypeId** | **int64** | The ID of the instance type | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateInstanceTypeLogoRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **logo** | ***os.File** | Logo File png,jpg,svg | 
 **darkLogo** | ***os.File** | Dark Logo File png,jpg,svg | 

### Return type

[**UpdateInstanceTypeLogo200Response**](UpdateInstanceTypeLogo200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateLayout

> UpdateLayout200Response UpdateLayout(ctx, id).UpdateLayoutRequest(updateLayoutRequest).Execute()

Update a Layout



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
	updateLayoutRequest := *openapiclient.NewUpdateLayoutRequest() // UpdateLayoutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateLayout(context.Background(), id).UpdateLayoutRequest(updateLayoutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateLayout``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateLayout`: UpdateLayout200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateLayout`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateLayoutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateLayoutRequest** | [**UpdateLayoutRequest**](UpdateLayoutRequest.md) |  | 

### Return type

[**UpdateLayout200Response**](UpdateLayout200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateLayoutPermissions

> UpdateLayoutPermissions200Response UpdateLayoutPermissions(ctx, id).UpdateLayoutPermissionsRequest(updateLayoutPermissionsRequest).Execute()

Update Layout Permissions



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
	updateLayoutPermissionsRequest := *openapiclient.NewUpdateLayoutPermissionsRequest() // UpdateLayoutPermissionsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateLayoutPermissions(context.Background(), id).UpdateLayoutPermissionsRequest(updateLayoutPermissionsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateLayoutPermissions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateLayoutPermissions`: UpdateLayoutPermissions200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateLayoutPermissions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateLayoutPermissionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateLayoutPermissionsRequest** | [**UpdateLayoutPermissionsRequest**](UpdateLayoutPermissionsRequest.md) |  | 

### Return type

[**UpdateLayoutPermissions200Response**](UpdateLayoutPermissions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateNodeType

> UpdateNodeType200Response UpdateNodeType(ctx, id).UpdateNodeTypeRequest(updateNodeTypeRequest).Execute()

Update a Node Type



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
	updateNodeTypeRequest := *openapiclient.NewUpdateNodeTypeRequest() // UpdateNodeTypeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateNodeType(context.Background(), id).UpdateNodeTypeRequest(updateNodeTypeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateNodeType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateNodeType`: UpdateNodeType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateNodeType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateNodeTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateNodeTypeRequest** | [**UpdateNodeTypeRequest**](UpdateNodeTypeRequest.md) |  | 

### Return type

[**UpdateNodeType200Response**](UpdateNodeType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateOptionForm

> UpdateOptionForm200Response UpdateOptionForm(ctx, id).UpdateOptionFormRequest(updateOptionFormRequest).Execute()

Update an Option Form



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
	updateOptionFormRequest := *openapiclient.NewUpdateOptionFormRequest() // UpdateOptionFormRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateOptionForm(context.Background(), id).UpdateOptionFormRequest(updateOptionFormRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateOptionForm``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateOptionForm`: UpdateOptionForm200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateOptionForm`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateOptionFormRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateOptionFormRequest** | [**UpdateOptionFormRequest**](UpdateOptionFormRequest.md) |  | 

### Return type

[**UpdateOptionForm200Response**](UpdateOptionForm200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateOptionList

> UpdateOptionList200Response UpdateOptionList(ctx, id).UpdateOptionListRequest(updateOptionListRequest).Execute()

Update an Option List



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
	updateOptionListRequest := *openapiclient.NewUpdateOptionListRequest() // UpdateOptionListRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateOptionList(context.Background(), id).UpdateOptionListRequest(updateOptionListRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateOptionList``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateOptionList`: UpdateOptionList200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateOptionList`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateOptionListRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateOptionListRequest** | [**UpdateOptionListRequest**](UpdateOptionListRequest.md) |  | 

### Return type

[**UpdateOptionList200Response**](UpdateOptionList200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateOptionType

> UpdateOptionType200Response UpdateOptionType(ctx, id).UpdateOptionTypeRequest(updateOptionTypeRequest).Execute()

Update an Input



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
	updateOptionTypeRequest := *openapiclient.NewUpdateOptionTypeRequest() // UpdateOptionTypeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateOptionType(context.Background(), id).UpdateOptionTypeRequest(updateOptionTypeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateOptionType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateOptionType`: UpdateOptionType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateOptionType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateOptionTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateOptionTypeRequest** | [**UpdateOptionTypeRequest**](UpdateOptionTypeRequest.md) |  | 

### Return type

[**UpdateOptionType200Response**](UpdateOptionType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateOsType

> UpdateOsType200Response UpdateOsType(ctx, id).UpdateOsTypeRequest(updateOsTypeRequest).Execute()

Update an OsType



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
	updateOsTypeRequest := *openapiclient.NewUpdateOsTypeRequest() // UpdateOsTypeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateOsType(context.Background(), id).UpdateOsTypeRequest(updateOsTypeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateOsType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateOsType`: UpdateOsType200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateOsType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateOsTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateOsTypeRequest** | [**UpdateOsTypeRequest**](UpdateOsTypeRequest.md) |  | 

### Return type

[**UpdateOsType200Response**](UpdateOsType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateScript

> UpdateScript200Response UpdateScript(ctx, id).UpdateScriptRequest(updateScriptRequest).Execute()

Update a Script



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
	updateScriptRequest := *openapiclient.NewUpdateScriptRequest() // UpdateScriptRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateScript(context.Background(), id).UpdateScriptRequest(updateScriptRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateScript`: UpdateScript200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateScript`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateScriptRequest** | [**UpdateScriptRequest**](UpdateScriptRequest.md) |  | 

### Return type

[**UpdateScript200Response**](UpdateScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateSpecTemplate

> UpdateSpecTemplate200Response UpdateSpecTemplate(ctx, id).UpdateSpecTemplateRequest(updateSpecTemplateRequest).Execute()

Update a Spec Template



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
	updateSpecTemplateRequest := *openapiclient.NewUpdateSpecTemplateRequest() // UpdateSpecTemplateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateSpecTemplate(context.Background(), id).UpdateSpecTemplateRequest(updateSpecTemplateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateSpecTemplate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateSpecTemplate`: UpdateSpecTemplate200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateSpecTemplate`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateSpecTemplateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateSpecTemplateRequest** | [**UpdateSpecTemplateRequest**](UpdateSpecTemplateRequest.md) |  | 

### Return type

[**UpdateSpecTemplate200Response**](UpdateSpecTemplate200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateVirtualImage

> UpdateVirtualImage200Response UpdateVirtualImage(ctx, virtualImageId).UpdateVirtualImageRequest(updateVirtualImageRequest).Execute()

Update a Virtual Image



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
	virtualImageId := int64(4) // int64 | Virtual Image ID
	updateVirtualImageRequest := *openapiclient.NewUpdateVirtualImageRequest(*openapiclient.NewUpdateVirtualImageRequestVirtualImage()) // UpdateVirtualImageRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LibraryAPI.UpdateVirtualImage(context.Background(), virtualImageId).UpdateVirtualImageRequest(updateVirtualImageRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LibraryAPI.UpdateVirtualImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateVirtualImage`: UpdateVirtualImage200Response
	fmt.Fprintf(os.Stdout, "Response from `LibraryAPI.UpdateVirtualImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**virtualImageId** | **int64** | Virtual Image ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateVirtualImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateVirtualImageRequest** | [**UpdateVirtualImageRequest**](UpdateVirtualImageRequest.md) |  | 

### Return type

[**UpdateVirtualImage200Response**](UpdateVirtualImage200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

