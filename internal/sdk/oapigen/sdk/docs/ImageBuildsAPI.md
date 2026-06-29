# \ImageBuildsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddBootScript**](ImageBuildsAPI.md#AddBootScript) | **Post** /api/boot-scripts | Create a Boot Script
[**AddImageBuild**](ImageBuildsAPI.md#AddImageBuild) | **Post** /api/image-builds | Create an Image Build
[**AddPreseedScript**](ImageBuildsAPI.md#AddPreseedScript) | **Post** /api/preseed-scripts | Create a Preseed Script
[**DeleteBootScript**](ImageBuildsAPI.md#DeleteBootScript) | **Delete** /api/boot-scripts/{id} | Delete a Boot Script
[**DeleteImageBuild**](ImageBuildsAPI.md#DeleteImageBuild) | **Delete** /api/image-builds/{id} | Delete an Image Build
[**DeletePreseedScript**](ImageBuildsAPI.md#DeletePreseedScript) | **Delete** /api/preseed-scripts/{id} | Delete a Preseed Script
[**ExecuteImageBuild**](ImageBuildsAPI.md#ExecuteImageBuild) | **Post** /api/image-builds/{id}/run | Run an Image Build
[**GetBootScript**](ImageBuildsAPI.md#GetBootScript) | **Get** /api/boot-scripts/{id} | Get a Specific Boot Script
[**GetImageBuild**](ImageBuildsAPI.md#GetImageBuild) | **Get** /api/image-builds/{id} | Get a Specific Image Build
[**GetImageBuildExecutions**](ImageBuildsAPI.md#GetImageBuildExecutions) | **Get** /api/image-builds/{id}/list-executions | List Image Build Executions
[**GetPreseedScript**](ImageBuildsAPI.md#GetPreseedScript) | **Get** /api/preseed-scripts/{id} | Get a Specific Preseed Script
[**ListBootScripts**](ImageBuildsAPI.md#ListBootScripts) | **Get** /api/boot-scripts | Boot Scripts
[**ListImageBuilds**](ImageBuildsAPI.md#ListImageBuilds) | **Get** /api/image-builds | Get All Image Builds
[**ListPreseedScripts**](ImageBuildsAPI.md#ListPreseedScripts) | **Get** /api/preseed-scripts | Preseed Scripts
[**UpdateBootScript**](ImageBuildsAPI.md#UpdateBootScript) | **Put** /api/boot-scripts/{id} | Update a Boot Script
[**UpdateImageBuild**](ImageBuildsAPI.md#UpdateImageBuild) | **Put** /api/image-builds/{id} | Update an Image Build
[**UpdatePreseedScript**](ImageBuildsAPI.md#UpdatePreseedScript) | **Put** /api/preseed-scripts/{id} | Update a Preseed Script



## AddBootScript

> AddBootScript200Response AddBootScript(ctx).AddBootScriptRequest(addBootScriptRequest).Execute()

Create a Boot Script



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
	addBootScriptRequest := *openapiclient.NewAddBootScriptRequest() // AddBootScriptRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ImageBuildsAPI.AddBootScript(context.Background()).AddBootScriptRequest(addBootScriptRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.AddBootScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddBootScript`: AddBootScript200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.AddBootScript`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddBootScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addBootScriptRequest** | [**AddBootScriptRequest**](AddBootScriptRequest.md) |  | 

### Return type

[**AddBootScript200Response**](AddBootScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddImageBuild

> AddImageBuild200Response AddImageBuild(ctx).AddImageBuildRequest(addImageBuildRequest).Execute()

Create an Image Build



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
	addImageBuildRequest := *openapiclient.NewAddImageBuildRequest() // AddImageBuildRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ImageBuildsAPI.AddImageBuild(context.Background()).AddImageBuildRequest(addImageBuildRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.AddImageBuild``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddImageBuild`: AddImageBuild200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.AddImageBuild`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddImageBuildRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addImageBuildRequest** | [**AddImageBuildRequest**](AddImageBuildRequest.md) |  | 

### Return type

[**AddImageBuild200Response**](AddImageBuild200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddPreseedScript

> AddPreseedScript200Response AddPreseedScript(ctx).AddPreseedScriptRequest(addPreseedScriptRequest).Execute()

Create a Preseed Script



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
	addPreseedScriptRequest := *openapiclient.NewAddPreseedScriptRequest() // AddPreseedScriptRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ImageBuildsAPI.AddPreseedScript(context.Background()).AddPreseedScriptRequest(addPreseedScriptRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.AddPreseedScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddPreseedScript`: AddPreseedScript200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.AddPreseedScript`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddPreseedScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addPreseedScriptRequest** | [**AddPreseedScriptRequest**](AddPreseedScriptRequest.md) |  | 

### Return type

[**AddPreseedScript200Response**](AddPreseedScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteBootScript

> DeleteBootScript200Response DeleteBootScript(ctx, id).Execute()

Delete a Boot Script



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
	resp, r, err := apiClient.ImageBuildsAPI.DeleteBootScript(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.DeleteBootScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteBootScript`: DeleteBootScript200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.DeleteBootScript`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteBootScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteBootScript200Response**](DeleteBootScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteImageBuild

> DeleteImageBuild200Response DeleteImageBuild(ctx, id).Execute()

Delete an Image Build



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
	resp, r, err := apiClient.ImageBuildsAPI.DeleteImageBuild(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.DeleteImageBuild``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteImageBuild`: DeleteImageBuild200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.DeleteImageBuild`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteImageBuildRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteImageBuild200Response**](DeleteImageBuild200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeletePreseedScript

> DeletePreseedScript200Response DeletePreseedScript(ctx, id).Execute()

Delete a Preseed Script



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
	resp, r, err := apiClient.ImageBuildsAPI.DeletePreseedScript(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.DeletePreseedScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeletePreseedScript`: DeletePreseedScript200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.DeletePreseedScript`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeletePreseedScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeletePreseedScript200Response**](DeletePreseedScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExecuteImageBuild

> ExecuteImageBuild200Response ExecuteImageBuild(ctx, id).Execute()

Run an Image Build



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
	resp, r, err := apiClient.ImageBuildsAPI.ExecuteImageBuild(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.ExecuteImageBuild``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecuteImageBuild`: ExecuteImageBuild200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.ExecuteImageBuild`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiExecuteImageBuildRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ExecuteImageBuild200Response**](ExecuteImageBuild200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetBootScript

> GetBootScript200Response GetBootScript(ctx, id).Execute()

Get a Specific Boot Script



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
	resp, r, err := apiClient.ImageBuildsAPI.GetBootScript(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.GetBootScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetBootScript`: GetBootScript200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.GetBootScript`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetBootScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetBootScript200Response**](GetBootScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetImageBuild

> GetImageBuild200Response GetImageBuild(ctx, id).Execute()

Get a Specific Image Build



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
	resp, r, err := apiClient.ImageBuildsAPI.GetImageBuild(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.GetImageBuild``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetImageBuild`: GetImageBuild200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.GetImageBuild`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetImageBuildRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetImageBuild200Response**](GetImageBuild200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetImageBuildExecutions

> GetImageBuildExecutions200Response GetImageBuildExecutions(ctx, id).BuildNumber(buildNumber).Status(status).Execute()

List Image Build Executions



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
	buildNumber := int64(4) // int64 | If specified will return an exact match on buildNumber (optional)
	status := "running" // string | Filter by status (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ImageBuildsAPI.GetImageBuildExecutions(context.Background(), id).BuildNumber(buildNumber).Status(status).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.GetImageBuildExecutions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetImageBuildExecutions`: GetImageBuildExecutions200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.GetImageBuildExecutions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetImageBuildExecutionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **buildNumber** | **int64** | If specified will return an exact match on buildNumber | 
 **status** | **string** | Filter by status | 

### Return type

[**GetImageBuildExecutions200Response**](GetImageBuildExecutions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPreseedScript

> GetPreseedScript200Response GetPreseedScript(ctx, id).Execute()

Get a Specific Preseed Script



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
	resp, r, err := apiClient.ImageBuildsAPI.GetPreseedScript(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.GetPreseedScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPreseedScript`: GetPreseedScript200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.GetPreseedScript`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPreseedScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetPreseedScript200Response**](GetPreseedScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListBootScripts

> ListBootScripts200Response ListBootScripts(ctx).Name(name).Phrase(phrase).Execute()

Boot Scripts



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
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ImageBuildsAPI.ListBootScripts(context.Background()).Name(name).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.ListBootScripts``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListBootScripts`: ListBootScripts200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.ListBootScripts`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListBootScriptsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListBootScripts200Response**](ListBootScripts200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListImageBuilds

> ListImageBuilds200Response ListImageBuilds(ctx).Name(name).Phrase(phrase).Execute()

Get All Image Builds



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
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ImageBuildsAPI.ListImageBuilds(context.Background()).Name(name).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.ListImageBuilds``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListImageBuilds`: ListImageBuilds200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.ListImageBuilds`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListImageBuildsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListImageBuilds200Response**](ListImageBuilds200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPreseedScripts

> ListPreseedScripts200Response ListPreseedScripts(ctx).Name(name).Phrase(phrase).Execute()

Preseed Scripts



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
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ImageBuildsAPI.ListPreseedScripts(context.Background()).Name(name).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.ListPreseedScripts``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPreseedScripts`: ListPreseedScripts200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.ListPreseedScripts`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListPreseedScriptsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListPreseedScripts200Response**](ListPreseedScripts200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateBootScript

> UpdateBootScript200Response UpdateBootScript(ctx, id).UpdateBootScriptRequest(updateBootScriptRequest).Execute()

Update a Boot Script



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
	updateBootScriptRequest := *openapiclient.NewUpdateBootScriptRequest() // UpdateBootScriptRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ImageBuildsAPI.UpdateBootScript(context.Background(), id).UpdateBootScriptRequest(updateBootScriptRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.UpdateBootScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateBootScript`: UpdateBootScript200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.UpdateBootScript`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateBootScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateBootScriptRequest** | [**UpdateBootScriptRequest**](UpdateBootScriptRequest.md) |  | 

### Return type

[**UpdateBootScript200Response**](UpdateBootScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateImageBuild

> UpdateImageBuild200Response UpdateImageBuild(ctx, id).UpdateImageBuildRequest(updateImageBuildRequest).Execute()

Update an Image Build



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
	updateImageBuildRequest := *openapiclient.NewUpdateImageBuildRequest() // UpdateImageBuildRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ImageBuildsAPI.UpdateImageBuild(context.Background(), id).UpdateImageBuildRequest(updateImageBuildRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.UpdateImageBuild``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateImageBuild`: UpdateImageBuild200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.UpdateImageBuild`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateImageBuildRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateImageBuildRequest** | [**UpdateImageBuildRequest**](UpdateImageBuildRequest.md) |  | 

### Return type

[**UpdateImageBuild200Response**](UpdateImageBuild200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdatePreseedScript

> UpdatePreseedScript200Response UpdatePreseedScript(ctx, id).UpdatePreseedScriptRequest(updatePreseedScriptRequest).Execute()

Update a Preseed Script



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
	updatePreseedScriptRequest := *openapiclient.NewUpdatePreseedScriptRequest() // UpdatePreseedScriptRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ImageBuildsAPI.UpdatePreseedScript(context.Background(), id).UpdatePreseedScriptRequest(updatePreseedScriptRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageBuildsAPI.UpdatePreseedScript``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdatePreseedScript`: UpdatePreseedScript200Response
	fmt.Fprintf(os.Stdout, "Response from `ImageBuildsAPI.UpdatePreseedScript`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdatePreseedScriptRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updatePreseedScriptRequest** | [**UpdatePreseedScriptRequest**](UpdatePreseedScriptRequest.md) |  | 

### Return type

[**UpdatePreseedScript200Response**](UpdatePreseedScript200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

