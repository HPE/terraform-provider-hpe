# \DeploysAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddInstanceDeploy**](DeploysAPI.md#AddInstanceDeploy) | **Post** /api/instances/{id}/deploys | Deploy to an Instance
[**Deletedeploy**](DeploysAPI.md#Deletedeploy) | **Delete** /api/deploys/{id} | Delete a Deploy
[**GetInstanceDeploys**](DeploysAPI.md#GetInstanceDeploys) | **Get** /api/instances/{id}/deploys | Get all Deploys for an Instance
[**ListDeploys**](DeploysAPI.md#ListDeploys) | **Get** /api/deploys | Get all Deploys
[**RunDeploy**](DeploysAPI.md#RunDeploy) | **Post** /api/deploys/{id}/deploy | Run a Deploy
[**UpdateDeploy**](DeploysAPI.md#UpdateDeploy) | **Put** /api/deploys/{id} | Update a Deploy



## AddInstanceDeploy

> AddInstanceDeploy200Response AddInstanceDeploy(ctx, id).AddInstanceDeployRequest(addInstanceDeployRequest).Execute()

Deploy to an Instance



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
	addInstanceDeployRequest := *openapiclient.NewAddInstanceDeployRequest() // AddInstanceDeployRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploysAPI.AddInstanceDeploy(context.Background(), id).AddInstanceDeployRequest(addInstanceDeployRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploysAPI.AddInstanceDeploy``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddInstanceDeploy`: AddInstanceDeploy200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploysAPI.AddInstanceDeploy`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddInstanceDeployRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addInstanceDeployRequest** | [**AddInstanceDeployRequest**](AddInstanceDeployRequest.md) |  | 

### Return type

[**AddInstanceDeploy200Response**](AddInstanceDeploy200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Deletedeploy

> Deletedeploy200Response Deletedeploy(ctx, id).Execute()

Delete a Deploy



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
	resp, r, err := apiClient.DeploysAPI.Deletedeploy(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploysAPI.Deletedeploy``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Deletedeploy`: Deletedeploy200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploysAPI.Deletedeploy`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeletedeployRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Deletedeploy200Response**](Deletedeploy200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceDeploys

> GetInstanceDeploys200Response GetInstanceDeploys(ctx, id).Max(max).Offset(offset).Phrase(phrase).Name(name).DeploymentId(deploymentId).InstanceName(instanceName).InstanceId(instanceId).Version(version).VersionId(versionId).CreatedById(createdById).DeployType(deployType).DateCreated(dateCreated).LastUpdated(lastUpdated).DeployDate(deployDate).Status(status).Execute()

Get all Deploys for an Instance



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
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	deploymentId := int64(5) // int64 | Filter by deployment id (optional)
	instanceName := "instance1" // string | Filter by instance name (optional)
	instanceId := int64(94) // int64 | The Instance ID for Filtering (optional)
	version := int64(5) // int64 | Filter by version number (userVersion) (optional)
	versionId := int64(5) // int64 | Filter by deployment version id (optional)
	createdById := int64(63) // int64 | Filter by owner (user) id (optional)
	deployType := "file" // string | Filter by type (deployType), file, git, fetch (optional)
	dateCreated := "2019-01-01" // string | Filter by dateCreated, the created timestamp is more recent or equal to the date specified (optional)
	lastUpdated := time.Now() // time.Time | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) (optional)
	deployDate := "2020-01-01" // string | Filter by deployDate, deployment completion timestamp is more recent or equal to the date specified (optional)
	status := "deploying" // string | Filter by status (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploysAPI.GetInstanceDeploys(context.Background(), id).Max(max).Offset(offset).Phrase(phrase).Name(name).DeploymentId(deploymentId).InstanceName(instanceName).InstanceId(instanceId).Version(version).VersionId(versionId).CreatedById(createdById).DeployType(deployType).DateCreated(dateCreated).LastUpdated(lastUpdated).DeployDate(deployDate).Status(status).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploysAPI.GetInstanceDeploys``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceDeploys`: GetInstanceDeploys200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploysAPI.GetInstanceDeploys`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceDeploysRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **deploymentId** | **int64** | Filter by deployment id | 
 **instanceName** | **string** | Filter by instance name | 
 **instanceId** | **int64** | The Instance ID for Filtering | 
 **version** | **int64** | Filter by version number (userVersion) | 
 **versionId** | **int64** | Filter by deployment version id | 
 **createdById** | **int64** | Filter by owner (user) id | 
 **deployType** | **string** | Filter by type (deployType), file, git, fetch | 
 **dateCreated** | **string** | Filter by dateCreated, the created timestamp is more recent or equal to the date specified | 
 **lastUpdated** | **time.Time** | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) | 
 **deployDate** | **string** | Filter by deployDate, deployment completion timestamp is more recent or equal to the date specified | 
 **status** | **string** | Filter by status | 

### Return type

[**GetInstanceDeploys200Response**](GetInstanceDeploys200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListDeploys

> ListDeploys200Response ListDeploys(ctx).Max(max).Offset(offset).Phrase(phrase).Name(name).DeploymentId(deploymentId).InstanceName(instanceName).InstanceId(instanceId).Version(version).VersionId(versionId).CreatedById(createdById).DeployType(deployType).DateCreated(dateCreated).LastUpdated(lastUpdated).DeployDate(deployDate).Status(status).Execute()

Get all Deploys



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
	deploymentId := int64(5) // int64 | Filter by deployment id (optional)
	instanceName := "instance1" // string | Filter by instance name (optional)
	instanceId := int64(94) // int64 | The Instance ID for Filtering (optional)
	version := int64(5) // int64 | Filter by version number (userVersion) (optional)
	versionId := int64(5) // int64 | Filter by deployment version id (optional)
	createdById := int64(63) // int64 | Filter by owner (user) id (optional)
	deployType := "file" // string | Filter by type (deployType), file, git, fetch (optional)
	dateCreated := "2019-01-01" // string | Filter by dateCreated, the created timestamp is more recent or equal to the date specified (optional)
	lastUpdated := time.Now() // time.Time | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) (optional)
	deployDate := "2020-01-01" // string | Filter by deployDate, deployment completion timestamp is more recent or equal to the date specified (optional)
	status := "deploying" // string | Filter by status (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploysAPI.ListDeploys(context.Background()).Max(max).Offset(offset).Phrase(phrase).Name(name).DeploymentId(deploymentId).InstanceName(instanceName).InstanceId(instanceId).Version(version).VersionId(versionId).CreatedById(createdById).DeployType(deployType).DateCreated(dateCreated).LastUpdated(lastUpdated).DeployDate(deployDate).Status(status).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploysAPI.ListDeploys``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListDeploys`: ListDeploys200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploysAPI.ListDeploys`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListDeploysRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **deploymentId** | **int64** | Filter by deployment id | 
 **instanceName** | **string** | Filter by instance name | 
 **instanceId** | **int64** | The Instance ID for Filtering | 
 **version** | **int64** | Filter by version number (userVersion) | 
 **versionId** | **int64** | Filter by deployment version id | 
 **createdById** | **int64** | Filter by owner (user) id | 
 **deployType** | **string** | Filter by type (deployType), file, git, fetch | 
 **dateCreated** | **string** | Filter by dateCreated, the created timestamp is more recent or equal to the date specified | 
 **lastUpdated** | **time.Time** | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) | 
 **deployDate** | **string** | Filter by deployDate, deployment completion timestamp is more recent or equal to the date specified | 
 **status** | **string** | Filter by status | 

### Return type

[**ListDeploys200Response**](ListDeploys200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RunDeploy

> RunDeploy200Response RunDeploy(ctx, id).RunDeployRequest(runDeployRequest).Execute()

Run a Deploy



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
	runDeployRequest := *openapiclient.NewRunDeployRequest() // RunDeployRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploysAPI.RunDeploy(context.Background(), id).RunDeployRequest(runDeployRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploysAPI.RunDeploy``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RunDeploy`: RunDeploy200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploysAPI.RunDeploy`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRunDeployRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **runDeployRequest** | [**RunDeployRequest**](RunDeployRequest.md) |  | 

### Return type

[**RunDeploy200Response**](RunDeploy200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateDeploy

> UpdateDeploy200Response UpdateDeploy(ctx, id).UpdateDeployRequest(updateDeployRequest).Execute()

Update a Deploy



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
	updateDeployRequest := *openapiclient.NewUpdateDeployRequest() // UpdateDeployRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DeploysAPI.UpdateDeploy(context.Background(), id).UpdateDeployRequest(updateDeployRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeploysAPI.UpdateDeploy``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateDeploy`: UpdateDeploy200Response
	fmt.Fprintf(os.Stdout, "Response from `DeploysAPI.UpdateDeploy`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateDeployRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateDeployRequest** | [**UpdateDeployRequest**](UpdateDeployRequest.md) |  | 

### Return type

[**UpdateDeploy200Response**](UpdateDeploy200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

