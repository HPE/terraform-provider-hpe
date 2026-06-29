# \InstancesAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddInstance**](InstancesAPI.md#AddInstance) | **Post** /api/instances | Create an Instance
[**BackupInstance**](InstancesAPI.md#BackupInstance) | **Put** /api/instances/{id}/backup | Backup an instance
[**BackupsInstance**](InstancesAPI.md#BackupsInstance) | **Get** /api/instances/{id}/backups | Get list of backups for an Instance
[**CancelExpirationInstance**](InstancesAPI.md#CancelExpirationInstance) | **Put** /api/instances/{id}/cancel-expiration | Cancel Expiration of an Instance
[**CancelRemovalInstance**](InstancesAPI.md#CancelRemovalInstance) | **Put** /api/instances/{id}/cancel-removal | Cancel Removal of an Instance
[**CancelShutdownInstance**](InstancesAPI.md#CancelShutdownInstance) | **Put** /api/instances/{id}/cancel-shutdown | Cancel Shutdown of an Instance
[**CloneImageInstance**](InstancesAPI.md#CloneImageInstance) | **Put** /api/instances/{id}/clone-image | Clone to Image
[**CloneInstance**](InstancesAPI.md#CloneInstance) | **Put** /api/instances/{id}/clone | Clone an Instance
[**CreateInstanceSchedule**](InstancesAPI.md#CreateInstanceSchedule) | **Post** /api/instances/{id}/schedules | Create a new Instance Schedule
[**DeleteAllSnapshotsInstance**](InstancesAPI.md#DeleteAllSnapshotsInstance) | **Delete** /api/instances/{id}/delete-all-snapshots | Delete All Snapshots of Instance
[**DeleteAllSnapshotsInstanceContainer**](InstancesAPI.md#DeleteAllSnapshotsInstanceContainer) | **Delete** /api/instances/{id}/delete-container-snapshots/{containerId} | Delete All Snapshots of Instance Container
[**DeleteInstance**](InstancesAPI.md#DeleteInstance) | **Delete** /api/instances/{id} | Delete an instance
[**DeleteInstanceSchedule**](InstancesAPI.md#DeleteInstanceSchedule) | **Delete** /api/instances/{id}/schedules/{scheduleId} | Deletes an Instance Schedule
[**DeleteSnapshotInstance**](InstancesAPI.md#DeleteSnapshotInstance) | **Delete** /api/snapshots/{id} | Delete Snapshot of Instance
[**EjectInstance**](InstancesAPI.md#EjectInstance) | **Put** /api/instances/{id}/eject | Eject an instance
[**ExecuteInstanceAction**](InstancesAPI.md#ExecuteInstanceAction) | **Put** /api/instances/{id}/action | Execute Instance Action
[**ExtendExpirationInstance**](InstancesAPI.md#ExtendExpirationInstance) | **Put** /api/instances/{id}/extend-expiration | Extend Expiration of an Instance
[**ExtendShutdownInstance**](InstancesAPI.md#ExtendShutdownInstance) | **Put** /api/instances/{id}/extend-shutdown | Extend Shutdown of an Instance
[**GetEnvVariables**](InstancesAPI.md#GetEnvVariables) | **Get** /api/instances/{id}/envs | Get Env Variables
[**GetInstance**](InstancesAPI.md#GetInstance) | **Get** /api/instances/{id} | Retrieves a Specific Instance
[**GetInstanceActions**](InstancesAPI.md#GetInstanceActions) | **Get** /api/instances/{id}/actions | List Instance Actions
[**GetInstanceContainers**](InstancesAPI.md#GetInstanceContainers) | **Get** /api/instances/{id}/containers | Get Container Details
[**GetInstanceHistory**](InstancesAPI.md#GetInstanceHistory) | **Get** /api/instances/{id}/history | Get Instance History
[**GetInstanceSchedule**](InstancesAPI.md#GetInstanceSchedule) | **Get** /api/instances/{id}/schedules/{scheduleId} | Get a Specific Instance Schedule
[**GetInstanceSchedules**](InstancesAPI.md#GetInstanceSchedules) | **Get** /api/instances/{id}/schedules | Get all Instance Schedules
[**GetInstanceStats**](InstancesAPI.md#GetInstanceStats) | **Get** /api/instances/stats | Get Statistics for all Instances
[**GetInstanceThreshold**](InstancesAPI.md#GetInstanceThreshold) | **Get** /api/instances/{id}/threshold | Get an Instance Scale Threshold
[**GetInstanceTypeProvisioning**](InstancesAPI.md#GetInstanceTypeProvisioning) | **Get** /api/instance-types/{id} | Get Specific Instance Type for Provisioning
[**GetPrepareApplyInstance**](InstancesAPI.md#GetPrepareApplyInstance) | **Get** /api/instances/{id}/prepare-apply | Prepare To Apply an Instance
[**GetSnapshotInstance**](InstancesAPI.md#GetSnapshotInstance) | **Get** /api/snapshots/{id} | Get a Specific Snapshot
[**GetStateInstance**](InstancesAPI.md#GetStateInstance) | **Get** /api/instances/{id}/state | Get State of an Instance
[**GetValidateApplyInstance**](InstancesAPI.md#GetValidateApplyInstance) | **Post** /api/instances/{id}/validate-apply | Validate Apply State for an Instance
[**ImportSnapshotInstance**](InstancesAPI.md#ImportSnapshotInstance) | **Put** /api/instances/{id}/import-snapshot | Import Snapshot of an Instance
[**LinkedCloneSnapshotInstance**](InstancesAPI.md#LinkedCloneSnapshotInstance) | **Put** /api/instances/{id}/linked-clone/{snapshotId} | Create Linked Clone of Instance Snapshot
[**ListInstanceServicePlans**](InstancesAPI.md#ListInstanceServicePlans) | **Get** /api/instances/service-plans | Get Available Service Plans for an Instance
[**ListInstanceTypesProvisioning**](InstancesAPI.md#ListInstanceTypesProvisioning) | **Get** /api/instance-types | Get All Instance Types for Provisioning
[**ListInstances**](InstancesAPI.md#ListInstances) | **Get** /api/instances | Get All Instances
[**ListSecurityGroupsInstance**](InstancesAPI.md#ListSecurityGroupsInstance) | **Get** /api/instances/{id}/security-groups | Get Security Groups for an Instance
[**LockInstance**](InstancesAPI.md#LockInstance) | **Put** /api/instances/{id}/lock | Lock an Instance
[**RefreshStateInstance**](InstancesAPI.md#RefreshStateInstance) | **Post** /api/instances/{id}/refresh | Refresh State of an Instance
[**RemoveInstancesFromControl**](InstancesAPI.md#RemoveInstancesFromControl) | **Delete** /api/instances/remove-from-control | Remove From Control
[**ResizeInstance**](InstancesAPI.md#ResizeInstance) | **Put** /api/instances/{id}/resize | Resize an Instance
[**RestartInstance**](InstancesAPI.md#RestartInstance) | **Put** /api/instances/{id}/restart | Restart an instance
[**RevertSnapshotInstance**](InstancesAPI.md#RevertSnapshotInstance) | **Put** /api/instances/{id}/revert-snapshot/{snapshotId} | Revert Instance to Snapshot
[**RunWorkflowInstance**](InstancesAPI.md#RunWorkflowInstance) | **Put** /api/instances/{id}/workflow | Run Workflow on an Instance
[**SetApplyInstance**](InstancesAPI.md#SetApplyInstance) | **Post** /api/instances/{id}/apply | Apply State of an Instance
[**SetInstanceSecurityGroups**](InstancesAPI.md#SetInstanceSecurityGroups) | **Post** /api/instances/{id}/security-groups | Set Security Groups for an Instance
[**SnapshotInstance**](InstancesAPI.md#SnapshotInstance) | **Put** /api/instances/{id}/snapshot | Snapshot an Instance
[**SnapshotsInstance**](InstancesAPI.md#SnapshotsInstance) | **Get** /api/instances/{id}/snapshots | Get list of snapshots for an Instance
[**StartInstance**](InstancesAPI.md#StartInstance) | **Put** /api/instances/{id}/start | Start an instance
[**StopInstance**](InstancesAPI.md#StopInstance) | **Put** /api/instances/{id}/stop | Stop an instance
[**SuspendInstance**](InstancesAPI.md#SuspendInstance) | **Put** /api/instances/{id}/suspend | Suspend an instance
[**UnlockInstance**](InstancesAPI.md#UnlockInstance) | **Put** /api/instances/{id}/unlock | Unlock an Instance
[**UpdateInstance**](InstancesAPI.md#UpdateInstance) | **Put** /api/instances/{id} | Updating an Instance
[**UpdateInstanceNetworkInterface**](InstancesAPI.md#UpdateInstanceNetworkInterface) | **Put** /api/instances/{id}/networkInterfaces/{networkInterfaceId} | Updating a label for an Instance&#39;s Network
[**UpdateInstanceSchedule**](InstancesAPI.md#UpdateInstanceSchedule) | **Put** /api/instances/{id}/schedules/{scheduleId} | Updating an Instance Schedule
[**UpdateInstanceThreshold**](InstancesAPI.md#UpdateInstanceThreshold) | **Put** /api/instances/{id}/threshold | Updates an Instance Scale Threshold



## AddInstance

> AddInstance200Response AddInstance(ctx).AddInstanceRequest(addInstanceRequest).Execute()

Create an Instance



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
	addInstanceRequest := *openapiclient.NewAddInstanceRequest(*openapiclient.NewAddInstanceRequestInstance("myinstance", *openapiclient.NewAddInstanceRequestInstanceSite(int64(2)), *openapiclient.NewAddInstanceRequestInstanceInstanceType("Ubuntu"), *openapiclient.NewAddInstanceRequestInstanceLayout(int64(105)), *openapiclient.NewAddInstanceRequestInstancePlan(int64(75))), *openapiclient.NewAddInstanceRequestConfig()) // AddInstanceRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.AddInstance(context.Background()).AddInstanceRequest(addInstanceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.AddInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddInstance`: AddInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.AddInstance`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addInstanceRequest** | [**AddInstanceRequest**](AddInstanceRequest.md) |  | 

### Return type

[**AddInstance200Response**](AddInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BackupInstance

> BackupInstance200Response BackupInstance(ctx, id).Execute()

Backup an instance



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
	resp, r, err := apiClient.InstancesAPI.BackupInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.BackupInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BackupInstance`: BackupInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.BackupInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiBackupInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**BackupInstance200Response**](BackupInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BackupsInstance

> BackupsInstance200Response BackupsInstance(ctx, id).Execute()

Get list of backups for an Instance



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
	resp, r, err := apiClient.InstancesAPI.BackupsInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.BackupsInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BackupsInstance`: BackupsInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.BackupsInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiBackupsInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**BackupsInstance200Response**](BackupsInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CancelExpirationInstance

> CancelExpirationInstance200Response CancelExpirationInstance(ctx, id).Execute()

Cancel Expiration of an Instance



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
	resp, r, err := apiClient.InstancesAPI.CancelExpirationInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.CancelExpirationInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CancelExpirationInstance`: CancelExpirationInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.CancelExpirationInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiCancelExpirationInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CancelExpirationInstance200Response**](CancelExpirationInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CancelRemovalInstance

> CancelRemovalInstance200Response CancelRemovalInstance(ctx, id).Execute()

Cancel Removal of an Instance



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
	resp, r, err := apiClient.InstancesAPI.CancelRemovalInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.CancelRemovalInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CancelRemovalInstance`: CancelRemovalInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.CancelRemovalInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiCancelRemovalInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CancelRemovalInstance200Response**](CancelRemovalInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CancelShutdownInstance

> CancelShutdownInstance200Response CancelShutdownInstance(ctx, id).Execute()

Cancel Shutdown of an Instance



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
	resp, r, err := apiClient.InstancesAPI.CancelShutdownInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.CancelShutdownInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CancelShutdownInstance`: CancelShutdownInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.CancelShutdownInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiCancelShutdownInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CancelShutdownInstance200Response**](CancelShutdownInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CloneImageInstance

> CloneImageInstance200Response CloneImageInstance(ctx, id).CloneImageInstanceRequest(cloneImageInstanceRequest).Execute()

Clone to Image



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
	cloneImageInstanceRequest := *openapiclient.NewCloneImageInstanceRequest() // CloneImageInstanceRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.CloneImageInstance(context.Background(), id).CloneImageInstanceRequest(cloneImageInstanceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.CloneImageInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CloneImageInstance`: CloneImageInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.CloneImageInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiCloneImageInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **cloneImageInstanceRequest** | [**CloneImageInstanceRequest**](CloneImageInstanceRequest.md) |  | 

### Return type

[**CloneImageInstance200Response**](CloneImageInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CloneInstance

> CloneInstance200Response CloneInstance(ctx, id).CloneInstanceRequest(cloneInstanceRequest).Execute()

Clone an Instance



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
	cloneInstanceRequest := *openapiclient.NewCloneInstanceRequest() // CloneInstanceRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.CloneInstance(context.Background(), id).CloneInstanceRequest(cloneInstanceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.CloneInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CloneInstance`: CloneInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.CloneInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiCloneInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **cloneInstanceRequest** | [**CloneInstanceRequest**](CloneInstanceRequest.md) |  | 

### Return type

[**CloneInstance200Response**](CloneInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateInstanceSchedule

> CreateInstanceSchedule200Response CreateInstanceSchedule(ctx, id).CreateInstanceScheduleRequest(createInstanceScheduleRequest).Execute()

Create a new Instance Schedule



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
	id := int64(42) // int64 | Instance ID
	createInstanceScheduleRequest := *openapiclient.NewCreateInstanceScheduleRequest() // CreateInstanceScheduleRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.CreateInstanceSchedule(context.Background(), id).CreateInstanceScheduleRequest(createInstanceScheduleRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.CreateInstanceSchedule``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateInstanceSchedule`: CreateInstanceSchedule200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.CreateInstanceSchedule`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Instance ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateInstanceScheduleRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createInstanceScheduleRequest** | [**CreateInstanceScheduleRequest**](CreateInstanceScheduleRequest.md) |  | 

### Return type

[**CreateInstanceSchedule200Response**](CreateInstanceSchedule200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteAllSnapshotsInstance

> DeleteAllSnapshotsInstance200Response DeleteAllSnapshotsInstance(ctx, id).Execute()

Delete All Snapshots of Instance



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
	resp, r, err := apiClient.InstancesAPI.DeleteAllSnapshotsInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.DeleteAllSnapshotsInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteAllSnapshotsInstance`: DeleteAllSnapshotsInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.DeleteAllSnapshotsInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteAllSnapshotsInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteAllSnapshotsInstance200Response**](DeleteAllSnapshotsInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteAllSnapshotsInstanceContainer

> DeleteAllSnapshotsInstanceContainer200Response DeleteAllSnapshotsInstanceContainer(ctx, id, containerId).Execute()

Delete All Snapshots of Instance Container



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
	containerId := float32(4) // float32 | Container ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.DeleteAllSnapshotsInstanceContainer(context.Background(), id, containerId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.DeleteAllSnapshotsInstanceContainer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteAllSnapshotsInstanceContainer`: DeleteAllSnapshotsInstanceContainer200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.DeleteAllSnapshotsInstanceContainer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 
**containerId** | **float32** | Container ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteAllSnapshotsInstanceContainerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DeleteAllSnapshotsInstanceContainer200Response**](DeleteAllSnapshotsInstanceContainer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteInstance

> DeleteInstance200Response DeleteInstance(ctx, id).PreserveVolumes(preserveVolumes).KeepBackups(keepBackups).ReleaseFloatingIps(releaseFloatingIps).ReleaseEIPs(releaseEIPs).Force(force).RemoveVolumes(removeVolumes).Execute()

Delete an instance



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
	preserveVolumes := "on" // string | Preserve Volumes (optional) (default to "off")
	keepBackups := "on" // string | Preserve copy of backups (optional) (default to "off")
	releaseFloatingIps := "off" // string | Release Floating IPs (optional) (default to "on")
	releaseEIPs := "off" // string | Alias for releaseFloatingIps (optional) (default to "on")
	force := "on" // string | Force Delete (optional) (default to "off")
	removeVolumes := "on" // string | Remove Volumes (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.DeleteInstance(context.Background(), id).PreserveVolumes(preserveVolumes).KeepBackups(keepBackups).ReleaseFloatingIps(releaseFloatingIps).ReleaseEIPs(releaseEIPs).Force(force).RemoveVolumes(removeVolumes).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.DeleteInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteInstance`: DeleteInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.DeleteInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **preserveVolumes** | **string** | Preserve Volumes | [default to &quot;off&quot;]
 **keepBackups** | **string** | Preserve copy of backups | [default to &quot;off&quot;]
 **releaseFloatingIps** | **string** | Release Floating IPs | [default to &quot;on&quot;]
 **releaseEIPs** | **string** | Alias for releaseFloatingIps | [default to &quot;on&quot;]
 **force** | **string** | Force Delete | [default to &quot;off&quot;]
 **removeVolumes** | **string** | Remove Volumes | [default to &quot;off&quot;]

### Return type

[**DeleteInstance200Response**](DeleteInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteInstanceSchedule

> RemoveExecuteSchedules200Response DeleteInstanceSchedule(ctx, id, scheduleId).Execute()

Deletes an Instance Schedule



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
	id := int64(42) // int64 | Instance ID
	scheduleId := int64(1) // int64 | Instance Schedule ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.DeleteInstanceSchedule(context.Background(), id, scheduleId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.DeleteInstanceSchedule``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteInstanceSchedule`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.DeleteInstanceSchedule`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Instance ID | 
**scheduleId** | **int64** | Instance Schedule ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteInstanceScheduleRequest struct via the builder pattern


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


## DeleteSnapshotInstance

> DeleteSnapshotInstance200Response DeleteSnapshotInstance(ctx, id).Execute()

Delete Snapshot of Instance



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
	resp, r, err := apiClient.InstancesAPI.DeleteSnapshotInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.DeleteSnapshotInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteSnapshotInstance`: DeleteSnapshotInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.DeleteSnapshotInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteSnapshotInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteSnapshotInstance200Response**](DeleteSnapshotInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EjectInstance

> EjectInstance200Response EjectInstance(ctx, id).Execute()

Eject an instance



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
	resp, r, err := apiClient.InstancesAPI.EjectInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.EjectInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `EjectInstance`: EjectInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.EjectInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiEjectInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**EjectInstance200Response**](EjectInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExecuteInstanceAction

> ExecuteInstanceAction200Response ExecuteInstanceAction(ctx, id).Code(code).ExecuteInstanceActionRequest(executeInstanceActionRequest).Execute()

Execute Instance Action



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
	code := "generic-add-node" // string | Action code to be executed on a specific instance.
	executeInstanceActionRequest := *openapiclient.NewExecuteInstanceActionRequest() // ExecuteInstanceActionRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.ExecuteInstanceAction(context.Background(), id).Code(code).ExecuteInstanceActionRequest(executeInstanceActionRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.ExecuteInstanceAction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecuteInstanceAction`: ExecuteInstanceAction200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.ExecuteInstanceAction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiExecuteInstanceActionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **code** | **string** | Action code to be executed on a specific instance. | 
 **executeInstanceActionRequest** | [**ExecuteInstanceActionRequest**](ExecuteInstanceActionRequest.md) |  | 

### Return type

[**ExecuteInstanceAction200Response**](ExecuteInstanceAction200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExtendExpirationInstance

> ExtendExpirationInstance200Response ExtendExpirationInstance(ctx, id).Execute()

Extend Expiration of an Instance



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
	resp, r, err := apiClient.InstancesAPI.ExtendExpirationInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.ExtendExpirationInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExtendExpirationInstance`: ExtendExpirationInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.ExtendExpirationInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiExtendExpirationInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ExtendExpirationInstance200Response**](ExtendExpirationInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExtendShutdownInstance

> ExtendShutdownInstance200Response ExtendShutdownInstance(ctx, id).Execute()

Extend Shutdown of an Instance



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
	resp, r, err := apiClient.InstancesAPI.ExtendShutdownInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.ExtendShutdownInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExtendShutdownInstance`: ExtendShutdownInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.ExtendShutdownInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiExtendShutdownInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ExtendShutdownInstance200Response**](ExtendShutdownInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetEnvVariables

> GetEnvVariables200Response GetEnvVariables(ctx, id).Execute()

Get Env Variables



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
	resp, r, err := apiClient.InstancesAPI.GetEnvVariables(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetEnvVariables``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetEnvVariables`: GetEnvVariables200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetEnvVariables`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEnvVariablesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetEnvVariables200Response**](GetEnvVariables200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstance

> GetInstance200Response GetInstance(ctx, id).Details(details).Execute()

Retrieves a Specific Instance



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
	details := true // bool | Include details=false to exclude extra details about the instance and its container(s) and server(s), ie. instancePrice, apps, powerSchedule, isScalable, instanceThreshold, containerDetails.server. (optional) (default to true)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.GetInstance(context.Background(), id).Details(details).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstance`: GetInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **details** | **bool** | Include details&#x3D;false to exclude extra details about the instance and its container(s) and server(s), ie. instancePrice, apps, powerSchedule, isScalable, instanceThreshold, containerDetails.server. | [default to true]

### Return type

[**GetInstance200Response**](GetInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceActions

> GetInstanceActions200Response GetInstanceActions(ctx, id).Execute()

List Instance Actions



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
	resp, r, err := apiClient.InstancesAPI.GetInstanceActions(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetInstanceActions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceActions`: GetInstanceActions200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetInstanceActions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceActionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetInstanceActions200Response**](GetInstanceActions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceContainers

> GetInstanceContainers200Response GetInstanceContainers(ctx, id).Execute()

Get Container Details



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
	resp, r, err := apiClient.InstancesAPI.GetInstanceContainers(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetInstanceContainers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceContainers`: GetInstanceContainers200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetInstanceContainers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceContainersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetInstanceContainers200Response**](GetInstanceContainers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceHistory

> GetInstanceHistory200Response GetInstanceHistory(ctx, id).Phrase(phrase).ContainerId(containerId).ServerId(serverId).ZoneId(zoneId).Execute()

Get Instance History



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
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	containerId := int64(135) // int64 | The Container ID for Filtering (optional)
	serverId := int64(97) // int64 | The Server ID for Filtering (optional)
	zoneId := int64(3) // int64 | The Cloud ID (Zone ID) for Filtering (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.GetInstanceHistory(context.Background(), id).Phrase(phrase).ContainerId(containerId).ServerId(serverId).ZoneId(zoneId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetInstanceHistory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceHistory`: GetInstanceHistory200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetInstanceHistory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceHistoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **containerId** | **int64** | The Container ID for Filtering | 
 **serverId** | **int64** | The Server ID for Filtering | 
 **zoneId** | **int64** | The Cloud ID (Zone ID) for Filtering | 

### Return type

[**GetInstanceHistory200Response**](GetInstanceHistory200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceSchedule

> GetInstanceSchedule200Response GetInstanceSchedule(ctx, id, scheduleId).Execute()

Get a Specific Instance Schedule



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
	id := int64(42) // int64 | Instance ID
	scheduleId := int64(1) // int64 | Instance Schedule ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.GetInstanceSchedule(context.Background(), id, scheduleId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetInstanceSchedule``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceSchedule`: GetInstanceSchedule200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetInstanceSchedule`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Instance ID | 
**scheduleId** | **int64** | Instance Schedule ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceScheduleRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetInstanceSchedule200Response**](GetInstanceSchedule200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceSchedules

> GetInstanceSchedules200Response GetInstanceSchedules(ctx, id).Execute()

Get all Instance Schedules



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
	id := int64(42) // int64 | Instance ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.GetInstanceSchedules(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetInstanceSchedules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceSchedules`: GetInstanceSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetInstanceSchedules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Instance ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceSchedulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetInstanceSchedules200Response**](GetInstanceSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceStats

> GetInstanceStats200Response GetInstanceStats(ctx).ZoneId(zoneId).Execute()

Get Statistics for all Instances



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
	zoneId := int64(3) // int64 | The ID of the cloud(s) for filtering (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.GetInstanceStats(context.Background()).ZoneId(zoneId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetInstanceStats``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceStats`: GetInstanceStats200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetInstanceStats`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceStatsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zoneId** | **int64** | The ID of the cloud(s) for filtering | 

### Return type

[**GetInstanceStats200Response**](GetInstanceStats200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceThreshold

> GetInstanceThreshold200Response GetInstanceThreshold(ctx, id).Execute()

Get an Instance Scale Threshold



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
	resp, r, err := apiClient.InstancesAPI.GetInstanceThreshold(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetInstanceThreshold``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceThreshold`: GetInstanceThreshold200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetInstanceThreshold`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceThresholdRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetInstanceThreshold200Response**](GetInstanceThreshold200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceTypeProvisioning

> GetInstanceTypeProvisioning200Response GetInstanceTypeProvisioning(ctx, id).Execute()

Get Specific Instance Type for Provisioning



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
	resp, r, err := apiClient.InstancesAPI.GetInstanceTypeProvisioning(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetInstanceTypeProvisioning``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceTypeProvisioning`: GetInstanceTypeProvisioning200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetInstanceTypeProvisioning`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceTypeProvisioningRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetInstanceTypeProvisioning200Response**](GetInstanceTypeProvisioning200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPrepareApplyInstance

> GetPrepareApplyInstance200Response GetPrepareApplyInstance(ctx, id).Execute()

Prepare To Apply an Instance



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
	resp, r, err := apiClient.InstancesAPI.GetPrepareApplyInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetPrepareApplyInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPrepareApplyInstance`: GetPrepareApplyInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetPrepareApplyInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPrepareApplyInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetPrepareApplyInstance200Response**](GetPrepareApplyInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSnapshotInstance

> GetSnapshotInstance200Response GetSnapshotInstance(ctx, id).Execute()

Get a Specific Snapshot



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
	resp, r, err := apiClient.InstancesAPI.GetSnapshotInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetSnapshotInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSnapshotInstance`: GetSnapshotInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetSnapshotInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSnapshotInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetSnapshotInstance200Response**](GetSnapshotInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetStateInstance

> GetStateInstance200Response GetStateInstance(ctx, id).Execute()

Get State of an Instance



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
	resp, r, err := apiClient.InstancesAPI.GetStateInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetStateInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetStateInstance`: GetStateInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetStateInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetStateInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetStateInstance200Response**](GetStateInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetValidateApplyInstance

> GetValidateApplyInstance200Response GetValidateApplyInstance(ctx, id).ValidateAppStateRequest(validateAppStateRequest).Execute()

Validate Apply State for an Instance



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
	validateAppStateRequest := *openapiclient.NewValidateAppStateRequest() // ValidateAppStateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.GetValidateApplyInstance(context.Background(), id).ValidateAppStateRequest(validateAppStateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.GetValidateApplyInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetValidateApplyInstance`: GetValidateApplyInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.GetValidateApplyInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetValidateApplyInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **validateAppStateRequest** | [**ValidateAppStateRequest**](ValidateAppStateRequest.md) |  | 

### Return type

[**GetValidateApplyInstance200Response**](GetValidateApplyInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ImportSnapshotInstance

> ImportSnapshotInstance200Response ImportSnapshotInstance(ctx, id).ImportSnapshotInstanceRequest(importSnapshotInstanceRequest).Execute()

Import Snapshot of an Instance



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
	importSnapshotInstanceRequest := *openapiclient.NewImportSnapshotInstanceRequest() // ImportSnapshotInstanceRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.ImportSnapshotInstance(context.Background(), id).ImportSnapshotInstanceRequest(importSnapshotInstanceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.ImportSnapshotInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ImportSnapshotInstance`: ImportSnapshotInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.ImportSnapshotInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiImportSnapshotInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **importSnapshotInstanceRequest** | [**ImportSnapshotInstanceRequest**](ImportSnapshotInstanceRequest.md) |  | 

### Return type

[**ImportSnapshotInstance200Response**](ImportSnapshotInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## LinkedCloneSnapshotInstance

> LinkedCloneSnapshotInstance200Response LinkedCloneSnapshotInstance(ctx, id, snapshotId).Execute()

Create Linked Clone of Instance Snapshot



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
	snapshotId := float32(4) // float32 | Snapshot ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.LinkedCloneSnapshotInstance(context.Background(), id, snapshotId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.LinkedCloneSnapshotInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `LinkedCloneSnapshotInstance`: LinkedCloneSnapshotInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.LinkedCloneSnapshotInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 
**snapshotId** | **float32** | Snapshot ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiLinkedCloneSnapshotInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**LinkedCloneSnapshotInstance200Response**](LinkedCloneSnapshotInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListInstanceServicePlans

> ListInstanceServicePlans200Response ListInstanceServicePlans(ctx).ZoneId(zoneId).LayoutId(layoutId).SiteId(siteId).Execute()

Get Available Service Plans for an Instance



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
	zoneId := int64(3) // int64 | The Cloud ID (Zone ID) for Filtering (optional)
	layoutId := int64(98) // int64 | The Layout ID for Filtering (optional)
	siteId := int64(7) // int64 | The Group ID (Site ID) for Filtering (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.ListInstanceServicePlans(context.Background()).ZoneId(zoneId).LayoutId(layoutId).SiteId(siteId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.ListInstanceServicePlans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListInstanceServicePlans`: ListInstanceServicePlans200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.ListInstanceServicePlans`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListInstanceServicePlansRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zoneId** | **int64** | The Cloud ID (Zone ID) for Filtering | 
 **layoutId** | **int64** | The Layout ID for Filtering | 
 **siteId** | **int64** | The Group ID (Site ID) for Filtering | 

### Return type

[**ListInstanceServicePlans200Response**](ListInstanceServicePlans200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListInstanceTypesProvisioning

> ListInstanceTypesProvisioning200Response ListInstanceTypesProvisioning(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).Featured(featured).Details(details).Execute()

Get All Instance Types for Provisioning



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
	details := true // bool | Load full details including optionTypes, environmentVariables, etc. These properties are excluded by default. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.ListInstanceTypesProvisioning(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).Featured(featured).Details(details).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.ListInstanceTypesProvisioning``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListInstanceTypesProvisioning`: ListInstanceTypesProvisioning200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.ListInstanceTypesProvisioning`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListInstanceTypesProvisioningRequest struct via the builder pattern


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
 **details** | **bool** | Load full details including optionTypes, environmentVariables, etc. These properties are excluded by default. | 

### Return type

[**ListInstanceTypesProvisioning200Response**](ListInstanceTypesProvisioning200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListInstances

> ListInstances200Response ListInstances(ctx).Max(max).Offset(offset).Name(name).Phrase(phrase).ZoneId(zoneId).SiteId(siteId).PlanId(planId).ServerId(serverId).ParentServerId(parentServerId).InstanceType(instanceType).LastUpdated(lastUpdated).CreatedBy(createdBy).AgentInstalled(agentInstalled).Status(status).Environment(environment).ShowDeleted(showDeleted).Deleted(deleted).ExpireDate(expireDate).ExpireDateMin(expireDateMin).ExpireDays(expireDays).ExpireDaysMin(expireDaysMin).ShutdownDate(shutdownDate).ShutdownDateMin(shutdownDateMin).ShutdownDays(shutdownDays).ShutdownDaysMin(shutdownDaysMin).Labels(labels).AllLabels(allLabels).TagsName(tagsName).Details(details).IncludeTenants(includeTenants).TenantId(tenantId).Execute()

Get All Instances



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
	zoneId := int64(3) // int64 | The Cloud ID (Zone ID) for Filtering (optional)
	siteId := int64(7) // int64 | The Group ID (Site ID) for Filtering (optional)
	planId := int64(1) // int64 | The Plan ID for Filtering (optional)
	serverId := int64(97) // int64 | The Server ID for Filtering (optional)
	parentServerId := int64(36) // int64 | The Parent Server (Hypervisor) ID for Filtering (optional)
	instanceType := "ubuntu" // string | The Instance Type Code for Filtering (optional)
	lastUpdated := time.Now() // time.Time | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) (optional)
	createdBy := int64(63) // int64 | The User ID for Filtering (optional)
	agentInstalled := true // bool | Filter instances by if agent is installed or not on the associated servers. (optional)
	status := "running" // string | The instance status for filtering. (optional)
	environment := "lab" // string | The environment for filtering. (optional)
	showDeleted := true // bool | If true, includes instances in pending removal status. (optional) (default to false)
	deleted := true // bool | If true, only deleted resources or instances in pending removal status are returned. (optional)
	expireDate := "2019-01-01" // string | Filter by expireDate less than or equal to specified date (optional)
	expireDateMin := "2019-01-01" // string | Filter expireDate greater than or equal to the specified date (optional)
	expireDays := "20" // string | Filter by expireDays less than or equal to the specified value (optional)
	expireDaysMin := "20" // string | Filter by expireDays greater than or equal to the specified value (optional)
	shutdownDate := "2019-01-01" // string | Filter by shutdownDate less than equal to the specified date (optional)
	shutdownDateMin := "2019-01-01" // string | Filter by shutdownDate greater than or equal to the specified date (optional)
	shutdownDays := "20" // string | Filter by shutdownDays less than or equal to the specified value (optional)
	shutdownDaysMin := "20" // string | Filter by shutdownDays greater than or equal to the specified value (optional)
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)
	tagsName := "value" // string | Filter by tags (metadata). This allows filtering by a tag name and value(s)  (optional)
	details := true // bool | Include details=true to return more details about the instance and its container(s) and server(s), ie. instancePrice, apps, powerSchedule, isScalable, instanceThreshold, containerDetails.server. (optional) (default to false)
	includeTenants := true // bool | Optional ability to include all subtenant resources when calling from a master tenant user (optional) (default to false)
	tenantId := float32(3) // float32 | Filter by Tenant ID. Only available to the master account. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.ListInstances(context.Background()).Max(max).Offset(offset).Name(name).Phrase(phrase).ZoneId(zoneId).SiteId(siteId).PlanId(planId).ServerId(serverId).ParentServerId(parentServerId).InstanceType(instanceType).LastUpdated(lastUpdated).CreatedBy(createdBy).AgentInstalled(agentInstalled).Status(status).Environment(environment).ShowDeleted(showDeleted).Deleted(deleted).ExpireDate(expireDate).ExpireDateMin(expireDateMin).ExpireDays(expireDays).ExpireDaysMin(expireDaysMin).ShutdownDate(shutdownDate).ShutdownDateMin(shutdownDateMin).ShutdownDays(shutdownDays).ShutdownDaysMin(shutdownDaysMin).Labels(labels).AllLabels(allLabels).TagsName(tagsName).Details(details).IncludeTenants(includeTenants).TenantId(tenantId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.ListInstances``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListInstances`: ListInstances200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.ListInstances`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListInstancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **zoneId** | **int64** | The Cloud ID (Zone ID) for Filtering | 
 **siteId** | **int64** | The Group ID (Site ID) for Filtering | 
 **planId** | **int64** | The Plan ID for Filtering | 
 **serverId** | **int64** | The Server ID for Filtering | 
 **parentServerId** | **int64** | The Parent Server (Hypervisor) ID for Filtering | 
 **instanceType** | **string** | The Instance Type Code for Filtering | 
 **lastUpdated** | **time.Time** | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) | 
 **createdBy** | **int64** | The User ID for Filtering | 
 **agentInstalled** | **bool** | Filter instances by if agent is installed or not on the associated servers. | 
 **status** | **string** | The instance status for filtering. | 
 **environment** | **string** | The environment for filtering. | 
 **showDeleted** | **bool** | If true, includes instances in pending removal status. | [default to false]
 **deleted** | **bool** | If true, only deleted resources or instances in pending removal status are returned. | 
 **expireDate** | **string** | Filter by expireDate less than or equal to specified date | 
 **expireDateMin** | **string** | Filter expireDate greater than or equal to the specified date | 
 **expireDays** | **string** | Filter by expireDays less than or equal to the specified value | 
 **expireDaysMin** | **string** | Filter by expireDays greater than or equal to the specified value | 
 **shutdownDate** | **string** | Filter by shutdownDate less than equal to the specified date | 
 **shutdownDateMin** | **string** | Filter by shutdownDate greater than or equal to the specified date | 
 **shutdownDays** | **string** | Filter by shutdownDays less than or equal to the specified value | 
 **shutdownDaysMin** | **string** | Filter by shutdownDays greater than or equal to the specified value | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 
 **tagsName** | **string** | Filter by tags (metadata). This allows filtering by a tag name and value(s)  | 
 **details** | **bool** | Include details&#x3D;true to return more details about the instance and its container(s) and server(s), ie. instancePrice, apps, powerSchedule, isScalable, instanceThreshold, containerDetails.server. | [default to false]
 **includeTenants** | **bool** | Optional ability to include all subtenant resources when calling from a master tenant user | [default to false]
 **tenantId** | **float32** | Filter by Tenant ID. Only available to the master account. | 

### Return type

[**ListInstances200Response**](ListInstances200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSecurityGroupsInstance

> ListSecurityGroupsInstance200Response ListSecurityGroupsInstance(ctx, id).Execute()

Get Security Groups for an Instance



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
	resp, r, err := apiClient.InstancesAPI.ListSecurityGroupsInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.ListSecurityGroupsInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListSecurityGroupsInstance`: ListSecurityGroupsInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.ListSecurityGroupsInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiListSecurityGroupsInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListSecurityGroupsInstance200Response**](ListSecurityGroupsInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## LockInstance

> LockInstance200Response LockInstance(ctx, id).Execute()

Lock an Instance



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
	resp, r, err := apiClient.InstancesAPI.LockInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.LockInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `LockInstance`: LockInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.LockInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiLockInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**LockInstance200Response**](LockInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RefreshStateInstance

> ContainersAttachFloatingIp200Response RefreshStateInstance(ctx, id).Execute()

Refresh State of an Instance



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
	resp, r, err := apiClient.InstancesAPI.RefreshStateInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.RefreshStateInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RefreshStateInstance`: ContainersAttachFloatingIp200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.RefreshStateInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRefreshStateInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ContainersAttachFloatingIp200Response**](ContainersAttachFloatingIp200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveInstancesFromControl

> RemoveInstancesFromControl200Response RemoveInstancesFromControl(ctx).RemoveInstancesFromControlRequest(removeInstancesFromControlRequest).Execute()

Remove From Control



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
	removeInstancesFromControlRequest := *openapiclient.NewRemoveInstancesFromControlRequest() // RemoveInstancesFromControlRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.RemoveInstancesFromControl(context.Background()).RemoveInstancesFromControlRequest(removeInstancesFromControlRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.RemoveInstancesFromControl``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveInstancesFromControl`: RemoveInstancesFromControl200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.RemoveInstancesFromControl`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRemoveInstancesFromControlRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **removeInstancesFromControlRequest** | [**RemoveInstancesFromControlRequest**](RemoveInstancesFromControlRequest.md) |  | 

### Return type

[**RemoveInstancesFromControl200Response**](RemoveInstancesFromControl200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResizeInstance

> ResizeInstance200Response ResizeInstance(ctx, id).ResizeInstanceRequest(resizeInstanceRequest).Execute()

Resize an Instance



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
	resizeInstanceRequest := *openapiclient.NewResizeInstanceRequest() // ResizeInstanceRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.ResizeInstance(context.Background(), id).ResizeInstanceRequest(resizeInstanceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.ResizeInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ResizeInstance`: ResizeInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.ResizeInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiResizeInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **resizeInstanceRequest** | [**ResizeInstanceRequest**](ResizeInstanceRequest.md) |  | 

### Return type

[**ResizeInstance200Response**](ResizeInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RestartInstance

> RestartInstance200Response RestartInstance(ctx, id).Server(server).Execute()

Restart an instance



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
	server := true // bool | Restart the underlying server(s) as well by passing `true`. By default only the service will be restarted. (optional) (default to true)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.RestartInstance(context.Background(), id).Server(server).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.RestartInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RestartInstance`: RestartInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.RestartInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRestartInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **server** | **bool** | Restart the underlying server(s) as well by passing &#x60;true&#x60;. By default only the service will be restarted. | [default to true]

### Return type

[**RestartInstance200Response**](RestartInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RevertSnapshotInstance

> RevertSnapshotInstance200Response RevertSnapshotInstance(ctx, id, snapshotId).Execute()

Revert Instance to Snapshot



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
	snapshotId := float32(4) // float32 | Snapshot ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.RevertSnapshotInstance(context.Background(), id, snapshotId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.RevertSnapshotInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RevertSnapshotInstance`: RevertSnapshotInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.RevertSnapshotInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 
**snapshotId** | **float32** | Snapshot ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiRevertSnapshotInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**RevertSnapshotInstance200Response**](RevertSnapshotInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RunWorkflowInstance

> RunWorkflowInstance200Response RunWorkflowInstance(ctx, id).WorkflowId(workflowId).WorkflowName(workflowName).RunWorkflowInstanceRequest(runWorkflowInstanceRequest).Execute()

Run Workflow on an Instance



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
	workflowId := int64(3) // int64 | ID of the workflow to execute (optional)
	workflowName := "myworkflow" // string | Name of the workflow to execute (optional)
	runWorkflowInstanceRequest := *openapiclient.NewRunWorkflowInstanceRequest() // RunWorkflowInstanceRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.RunWorkflowInstance(context.Background(), id).WorkflowId(workflowId).WorkflowName(workflowName).RunWorkflowInstanceRequest(runWorkflowInstanceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.RunWorkflowInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RunWorkflowInstance`: RunWorkflowInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.RunWorkflowInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRunWorkflowInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **workflowId** | **int64** | ID of the workflow to execute | 
 **workflowName** | **string** | Name of the workflow to execute | 
 **runWorkflowInstanceRequest** | [**RunWorkflowInstanceRequest**](RunWorkflowInstanceRequest.md) |  | 

### Return type

[**RunWorkflowInstance200Response**](RunWorkflowInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetApplyInstance

> SetApplyInstance200Response SetApplyInstance(ctx, id).ValidateAppStateRequest(validateAppStateRequest).Execute()

Apply State of an Instance



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
	validateAppStateRequest := *openapiclient.NewValidateAppStateRequest() // ValidateAppStateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.SetApplyInstance(context.Background(), id).ValidateAppStateRequest(validateAppStateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.SetApplyInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetApplyInstance`: SetApplyInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.SetApplyInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiSetApplyInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **validateAppStateRequest** | [**ValidateAppStateRequest**](ValidateAppStateRequest.md) |  | 

### Return type

[**SetApplyInstance200Response**](SetApplyInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetInstanceSecurityGroups

> SetInstanceSecurityGroups200Response SetInstanceSecurityGroups(ctx, id).SetInstanceSecurityGroupsRequest(setInstanceSecurityGroupsRequest).Execute()

Set Security Groups for an Instance



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
	setInstanceSecurityGroupsRequest := *openapiclient.NewSetInstanceSecurityGroupsRequest() // SetInstanceSecurityGroupsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.SetInstanceSecurityGroups(context.Background(), id).SetInstanceSecurityGroupsRequest(setInstanceSecurityGroupsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.SetInstanceSecurityGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetInstanceSecurityGroups`: SetInstanceSecurityGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.SetInstanceSecurityGroups`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiSetInstanceSecurityGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **setInstanceSecurityGroupsRequest** | [**SetInstanceSecurityGroupsRequest**](SetInstanceSecurityGroupsRequest.md) |  | 

### Return type

[**SetInstanceSecurityGroups200Response**](SetInstanceSecurityGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SnapshotInstance

> SnapshotInstance200Response SnapshotInstance(ctx, id).SnapshotInstanceRequest(snapshotInstanceRequest).Execute()

Snapshot an Instance



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
	snapshotInstanceRequest := *openapiclient.NewSnapshotInstanceRequest() // SnapshotInstanceRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.SnapshotInstance(context.Background(), id).SnapshotInstanceRequest(snapshotInstanceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.SnapshotInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SnapshotInstance`: SnapshotInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.SnapshotInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiSnapshotInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **snapshotInstanceRequest** | [**SnapshotInstanceRequest**](SnapshotInstanceRequest.md) |  | 

### Return type

[**SnapshotInstance200Response**](SnapshotInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SnapshotsInstance

> SnapshotsInstance200Response SnapshotsInstance(ctx, id).Execute()

Get list of snapshots for an Instance



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
	resp, r, err := apiClient.InstancesAPI.SnapshotsInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.SnapshotsInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SnapshotsInstance`: SnapshotsInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.SnapshotsInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiSnapshotsInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**SnapshotsInstance200Response**](SnapshotsInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StartInstance

> StartInstance200Response StartInstance(ctx, id).Server(server).Execute()

Start an instance



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
	server := true // bool | Start the underlying server(s) as well by passing `true`. By default only the service will be started. (optional) (default to true)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.StartInstance(context.Background(), id).Server(server).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.StartInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `StartInstance`: StartInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.StartInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiStartInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **server** | **bool** | Start the underlying server(s) as well by passing &#x60;true&#x60;. By default only the service will be started. | [default to true]

### Return type

[**StartInstance200Response**](StartInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StopInstance

> StopInstance200Response StopInstance(ctx, id).Server(server).MuteMonitoring(muteMonitoring).Execute()

Stop an instance



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
	server := true // bool | Stop the underlying server(s) as well by passing `true`. By default only the service will be stopped. (optional) (default to true)
	muteMonitoring := true // bool | Mute monitoring checks for the instance by passing `true`. (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.StopInstance(context.Background(), id).Server(server).MuteMonitoring(muteMonitoring).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.StopInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `StopInstance`: StopInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.StopInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiStopInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **server** | **bool** | Stop the underlying server(s) as well by passing &#x60;true&#x60;. By default only the service will be stopped. | [default to true]
 **muteMonitoring** | **bool** | Mute monitoring checks for the instance by passing &#x60;true&#x60;. | [default to false]

### Return type

[**StopInstance200Response**](StopInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SuspendInstance

> SuspendInstance200Response SuspendInstance(ctx, id).Execute()

Suspend an instance



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
	resp, r, err := apiClient.InstancesAPI.SuspendInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.SuspendInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SuspendInstance`: SuspendInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.SuspendInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiSuspendInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**SuspendInstance200Response**](SuspendInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UnlockInstance

> UnlockInstance200Response UnlockInstance(ctx, id).Execute()

Unlock an Instance



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
	resp, r, err := apiClient.InstancesAPI.UnlockInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.UnlockInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UnlockInstance`: UnlockInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.UnlockInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUnlockInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**UnlockInstance200Response**](UnlockInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateInstance

> UpdateInstance200Response UpdateInstance(ctx, id).UpdateInstanceRequest(updateInstanceRequest).Execute()

Updating an Instance



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
	updateInstanceRequest := *openapiclient.NewUpdateInstanceRequest() // UpdateInstanceRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.UpdateInstance(context.Background(), id).UpdateInstanceRequest(updateInstanceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.UpdateInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateInstance`: UpdateInstance200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.UpdateInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateInstanceRequest** | [**UpdateInstanceRequest**](UpdateInstanceRequest.md) |  | 

### Return type

[**UpdateInstance200Response**](UpdateInstance200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateInstanceNetworkInterface

> UpdateInstanceNetworkInterface200Response UpdateInstanceNetworkInterface(ctx, id, networkInterfaceId).UpdateInstanceNetworkInterfaceRequest(updateInstanceNetworkInterfaceRequest).Execute()

Updating a label for an Instance's Network



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
	networkInterfaceId := float32(7) // float32 | NetworkInterface (ComputeServerInterface) ID of the Object being created or referenced
	updateInstanceNetworkInterfaceRequest := *openapiclient.NewUpdateInstanceNetworkInterfaceRequest() // UpdateInstanceNetworkInterfaceRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.UpdateInstanceNetworkInterface(context.Background(), id, networkInterfaceId).UpdateInstanceNetworkInterfaceRequest(updateInstanceNetworkInterfaceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.UpdateInstanceNetworkInterface``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateInstanceNetworkInterface`: UpdateInstanceNetworkInterface200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.UpdateInstanceNetworkInterface`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 
**networkInterfaceId** | **float32** | NetworkInterface (ComputeServerInterface) ID of the Object being created or referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateInstanceNetworkInterfaceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateInstanceNetworkInterfaceRequest** | [**UpdateInstanceNetworkInterfaceRequest**](UpdateInstanceNetworkInterfaceRequest.md) |  | 

### Return type

[**UpdateInstanceNetworkInterface200Response**](UpdateInstanceNetworkInterface200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateInstanceSchedule

> UpdateInstanceSchedule200Response UpdateInstanceSchedule(ctx, id, scheduleId).UpdateInstanceScheduleRequest(updateInstanceScheduleRequest).Execute()

Updating an Instance Schedule



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
	id := int64(9) // int64 | Instance ID
	scheduleId := int64(1) // int64 | Instance Schedule ID
	updateInstanceScheduleRequest := *openapiclient.NewUpdateInstanceScheduleRequest(*openapiclient.NewUpdateInstanceScheduleRequestInstanceSchedule()) // UpdateInstanceScheduleRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.UpdateInstanceSchedule(context.Background(), id, scheduleId).UpdateInstanceScheduleRequest(updateInstanceScheduleRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.UpdateInstanceSchedule``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateInstanceSchedule`: UpdateInstanceSchedule200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.UpdateInstanceSchedule`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Instance ID | 
**scheduleId** | **int64** | Instance Schedule ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateInstanceScheduleRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateInstanceScheduleRequest** | [**UpdateInstanceScheduleRequest**](UpdateInstanceScheduleRequest.md) |  | 

### Return type

[**UpdateInstanceSchedule200Response**](UpdateInstanceSchedule200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateInstanceThreshold

> UpdateInstanceThreshold200Response UpdateInstanceThreshold(ctx, id).UpdateInstanceThresholdRequest(updateInstanceThresholdRequest).Execute()

Updates an Instance Scale Threshold



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
	updateInstanceThresholdRequest := *openapiclient.NewUpdateInstanceThresholdRequest(*openapiclient.NewUpdateInstanceThresholdRequestInstanceThreshold()) // UpdateInstanceThresholdRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstancesAPI.UpdateInstanceThreshold(context.Background(), id).UpdateInstanceThresholdRequest(updateInstanceThresholdRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstancesAPI.UpdateInstanceThreshold``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateInstanceThreshold`: UpdateInstanceThreshold200Response
	fmt.Fprintf(os.Stdout, "Response from `InstancesAPI.UpdateInstanceThreshold`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateInstanceThresholdRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateInstanceThresholdRequest** | [**UpdateInstanceThresholdRequest**](UpdateInstanceThresholdRequest.md) |  | 

### Return type

[**UpdateInstanceThreshold200Response**](UpdateInstanceThreshold200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

