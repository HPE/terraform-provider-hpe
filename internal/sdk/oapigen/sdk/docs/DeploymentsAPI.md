# \DeploymentsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddDeploymentFile**](DeploymentsAPI.md#AddDeploymentFile) | **Post** /api/deployments/{deploymentId}/versions/{id}/files{filepath} | Upload a Deployment File
[**AddDeploymentVersion**](DeploymentsAPI.md#AddDeploymentVersion) | **Post** /api/deployments/{deploymentId}/versions | Create a new Deployment Version
[**AddDeployments**](DeploymentsAPI.md#AddDeployments) | **Post** /api/deployments | Create a new Deployment
[**DeleteDeployment**](DeploymentsAPI.md#DeleteDeployment) | **Delete** /api/deployments/{deploymentId} | Delete a Deployment
[**DeleteDeploymentFile**](DeploymentsAPI.md#DeleteDeploymentFile) | **Delete** /api/deployments/{deploymentId}/versions/{id}/files{filepath} | Delete a Deployment File
[**DeleteDeploymentVersion**](DeploymentsAPI.md#DeleteDeploymentVersion) | **Delete** /api/deployments/{deploymentId}/versions/{id} | Delete a Deployment Version
[**GetDeployment**](DeploymentsAPI.md#GetDeployment) | **Get** /api/deployments/{deploymentId} | Get a Specific Deployment
[**GetDeploymentVersion**](DeploymentsAPI.md#GetDeploymentVersion) | **Get** /api/deployments/{deploymentId}/versions/{id} | Get a Specific Deployment Version
[**ListDeploymentFiles**](DeploymentsAPI.md#ListDeploymentFiles) | **Get** /api/deployments/{deploymentId}/versions/{id}/files{filepath} | List Deployment Files
[**ListDeploymentVersions**](DeploymentsAPI.md#ListDeploymentVersions) | **Get** /api/deployments/{deploymentId}/versions | Get All Versions For a Deployment
[**ListDeployments**](DeploymentsAPI.md#ListDeployments) | **Get** /api/deployments | Get All Deployments
[**UpdateDeployment**](DeploymentsAPI.md#UpdateDeployment) | **Put** /api/deployments/{deploymentId} | Updating a Deployment
[**UpdateDeploymentVersion**](DeploymentsAPI.md#UpdateDeploymentVersion) | **Put** /api/deployments/{deploymentId}/versions/{id} | Updating a Deployment Version



## AddDeploymentFile

> AddDeploymentFile200Response AddDeploymentFile(ctx, deploymentId, id, filepath).File(file).Execute()

Upload a Deployment File



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
	deploymentId := int64(4) // int64 | Deployment ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	filepath := "/config/environments/" // string | The path to to search for files under. Default is the root directory /. (default to "/")
	file := os.NewFile(1234, "some_file") // *os.File |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.AddDeploymentFile(context.Background(), deploymentId, id, filepath).File(file).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.AddDeploymentFile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddDeploymentFile`: AddDeploymentFile200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.AddDeploymentFile`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **int64** | Deployment ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 
**filepath** | **string** | The path to to search for files under. Default is the root directory /. | [default to &quot;/&quot;]

### Other Parameters

Other parameters are passed through a pointer to a apiAddDeploymentFileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **file** | ***os.File** |  | 

### Return type

[**AddDeploymentFile200Response**](AddDeploymentFile200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddDeploymentVersion

> AddDeploymentVersion200Response AddDeploymentVersion(ctx, deploymentId).AddDeploymentVersionRequest(addDeploymentVersionRequest).Execute()

Create a new Deployment Version



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
	deploymentId := int64(4) // int64 | Deployment ID
	addDeploymentVersionRequest := *openapiclient.NewAddDeploymentVersionRequest() // AddDeploymentVersionRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.AddDeploymentVersion(context.Background(), deploymentId).AddDeploymentVersionRequest(addDeploymentVersionRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.AddDeploymentVersion``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddDeploymentVersion`: AddDeploymentVersion200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.AddDeploymentVersion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **int64** | Deployment ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddDeploymentVersionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addDeploymentVersionRequest** | [**AddDeploymentVersionRequest**](AddDeploymentVersionRequest.md) |  | 

### Return type

[**AddDeploymentVersion200Response**](AddDeploymentVersion200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddDeployments

> AddDeployments200Response AddDeployments(ctx).AddDeploymentsRequest(addDeploymentsRequest).Execute()

Create a new Deployment



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
	addDeploymentsRequest := *openapiclient.NewAddDeploymentsRequest() // AddDeploymentsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.AddDeployments(context.Background()).AddDeploymentsRequest(addDeploymentsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.AddDeployments``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddDeployments`: AddDeployments200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.AddDeployments`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddDeploymentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addDeploymentsRequest** | [**AddDeploymentsRequest**](AddDeploymentsRequest.md) |  | 

### Return type

[**AddDeployments200Response**](AddDeployments200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteDeployment

> DeleteDeployment200Response DeleteDeployment(ctx, deploymentId).Execute()

Delete a Deployment



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
	deploymentId := int64(4) // int64 | Deployment ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.DeleteDeployment(context.Background(), deploymentId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.DeleteDeployment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteDeployment`: DeleteDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.DeleteDeployment`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **int64** | Deployment ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDeploymentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteDeployment200Response**](DeleteDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteDeploymentFile

> DeleteDeploymentFile200Response DeleteDeploymentFile(ctx, deploymentId, id, filepath).Force(force).Execute()

Delete a Deployment File



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
	deploymentId := int64(4) // int64 | Deployment ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	filepath := "/config/environments/" // string | The path to to search for files under. Default is the root directory /. (default to "/")
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.DeleteDeploymentFile(context.Background(), deploymentId, id, filepath).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.DeleteDeploymentFile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteDeploymentFile`: DeleteDeploymentFile200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.DeleteDeploymentFile`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **int64** | Deployment ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 
**filepath** | **string** | The path to to search for files under. Default is the root directory /. | [default to &quot;/&quot;]

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDeploymentFileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**DeleteDeploymentFile200Response**](DeleteDeploymentFile200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteDeploymentVersion

> DeleteDeploymentVersion200Response DeleteDeploymentVersion(ctx, deploymentId, id).Execute()

Delete a Deployment Version



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
	deploymentId := int64(4) // int64 | Deployment ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.DeleteDeploymentVersion(context.Background(), deploymentId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.DeleteDeploymentVersion``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteDeploymentVersion`: DeleteDeploymentVersion200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.DeleteDeploymentVersion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **int64** | Deployment ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDeploymentVersionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DeleteDeploymentVersion200Response**](DeleteDeploymentVersion200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDeployment

> GetDeployment200Response GetDeployment(ctx, deploymentId).MaxVersions(maxVersions).Execute()

Get a Specific Deployment



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
	deploymentId := int64(4) // int64 | Deployment ID
	maxVersions := int64(6) // int64 | Max number of recent versions to return. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.GetDeployment(context.Background(), deploymentId).MaxVersions(maxVersions).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.GetDeployment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDeployment`: GetDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.GetDeployment`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **int64** | Deployment ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDeploymentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **maxVersions** | **int64** | Max number of recent versions to return. | 

### Return type

[**GetDeployment200Response**](GetDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDeploymentVersion

> GetDeploymentVersion200Response GetDeploymentVersion(ctx, deploymentId, id).Execute()

Get a Specific Deployment Version



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
	deploymentId := int64(4) // int64 | Deployment ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.GetDeploymentVersion(context.Background(), deploymentId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.GetDeploymentVersion``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDeploymentVersion`: GetDeploymentVersion200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.GetDeploymentVersion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **int64** | Deployment ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDeploymentVersionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetDeploymentVersion200Response**](GetDeploymentVersion200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListDeploymentFiles

> ListDeploymentFiles200Response ListDeploymentFiles(ctx, deploymentId, id, filepath).Max(max).Offset(offset).Phrase(phrase).Version(version).Execute()

List Deployment Files



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
	deploymentId := int64(4) // int64 | Deployment ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	filepath := "/config/environments/" // string | The path to to search for files under. Default is the root directory /. (default to "/")
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	version := int64(5) // int64 | Filter by version number (userVersion) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.ListDeploymentFiles(context.Background(), deploymentId, id, filepath).Max(max).Offset(offset).Phrase(phrase).Version(version).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.ListDeploymentFiles``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListDeploymentFiles`: ListDeploymentFiles200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.ListDeploymentFiles`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **int64** | Deployment ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 
**filepath** | **string** | The path to to search for files under. Default is the root directory /. | [default to &quot;/&quot;]

### Other Parameters

Other parameters are passed through a pointer to a apiListDeploymentFilesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **version** | **int64** | Filter by version number (userVersion) | 

### Return type

[**ListDeploymentFiles200Response**](ListDeploymentFiles200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListDeploymentVersions

> ListDeploymentVersions200Response ListDeploymentVersions(ctx, deploymentId).Max(max).Offset(offset).Phrase(phrase).Version(version).Type_(type_).DateCreated(dateCreated).LastUpdated(lastUpdated).Execute()

Get All Versions For a Deployment



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
	deploymentId := int64(4) // int64 | Deployment ID
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	version := int64(5) // int64 | Filter by version number (userVersion) (optional)
	type_ := "file" // string | Filter by type (deployType), file, git, fetch (optional)
	dateCreated := "2019-01-01" // string | Filter by dateCreated, the created timestamp is more recent or equal to the date specified (optional)
	lastUpdated := time.Now() // time.Time | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.ListDeploymentVersions(context.Background(), deploymentId).Max(max).Offset(offset).Phrase(phrase).Version(version).Type_(type_).DateCreated(dateCreated).LastUpdated(lastUpdated).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.ListDeploymentVersions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListDeploymentVersions`: ListDeploymentVersions200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.ListDeploymentVersions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **int64** | Deployment ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListDeploymentVersionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **version** | **int64** | Filter by version number (userVersion) | 
 **type_** | **string** | Filter by type (deployType), file, git, fetch | 
 **dateCreated** | **string** | Filter by dateCreated, the created timestamp is more recent or equal to the date specified | 
 **lastUpdated** | **time.Time** | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) | 

### Return type

[**ListDeploymentVersions200Response**](ListDeploymentVersions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListDeployments

> ListDeployments200Response ListDeployments(ctx).Max(max).Offset(offset).Phrase(phrase).Name(name).Description(description).DateCreated(dateCreated).LastUpdated(lastUpdated).Execute()

Get All Deployments



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
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	description := "The desription of my cool object" // string | Filter by description, wildcard may be specified as %. eg. `example-%` (optional)
	dateCreated := "2019-01-01" // string | Filter by dateCreated, the created timestamp is more recent or equal to the date specified (optional)
	lastUpdated := time.Now() // time.Time | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.ListDeployments(context.Background()).Max(max).Offset(offset).Phrase(phrase).Name(name).Description(description).DateCreated(dateCreated).LastUpdated(lastUpdated).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.ListDeployments``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListDeployments`: ListDeployments200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.ListDeployments`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListDeploymentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **description** | **string** | Filter by description, wildcard may be specified as %. eg. &#x60;example-%&#x60; | 
 **dateCreated** | **string** | Filter by dateCreated, the created timestamp is more recent or equal to the date specified | 
 **lastUpdated** | **time.Time** | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) | 

### Return type

[**ListDeployments200Response**](ListDeployments200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateDeployment

> UpdateDeployment200Response UpdateDeployment(ctx, deploymentId).UpdateDeploymentRequest(updateDeploymentRequest).Execute()

Updating a Deployment



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
	deploymentId := int64(4) // int64 | Deployment ID
	updateDeploymentRequest := *openapiclient.NewUpdateDeploymentRequest() // UpdateDeploymentRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.UpdateDeployment(context.Background(), deploymentId).UpdateDeploymentRequest(updateDeploymentRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.UpdateDeployment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateDeployment`: UpdateDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.UpdateDeployment`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **int64** | Deployment ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateDeploymentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateDeploymentRequest** | [**UpdateDeploymentRequest**](UpdateDeploymentRequest.md) |  | 

### Return type

[**UpdateDeployment200Response**](UpdateDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateDeploymentVersion

> UpdateDeploymentVersion200Response UpdateDeploymentVersion(ctx, deploymentId, id).UpdateDeploymentVersionRequest(updateDeploymentVersionRequest).Execute()

Updating a Deployment Version



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
	deploymentId := int64(4) // int64 | Deployment ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateDeploymentVersionRequest := *openapiclient.NewUpdateDeploymentVersionRequest() // UpdateDeploymentVersionRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploymentsAPI.UpdateDeploymentVersion(context.Background(), deploymentId, id).UpdateDeploymentVersionRequest(updateDeploymentVersionRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploymentsAPI.UpdateDeploymentVersion``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateDeploymentVersion`: UpdateDeploymentVersion200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploymentsAPI.UpdateDeploymentVersion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **int64** | Deployment ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateDeploymentVersionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateDeploymentVersionRequest** | [**UpdateDeploymentVersionRequest**](UpdateDeploymentVersionRequest.md) |  | 

### Return type

[**UpdateDeploymentVersion200Response**](UpdateDeploymentVersion200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

