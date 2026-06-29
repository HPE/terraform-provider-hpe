# \AutomationAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddExecuteSchedules**](AutomationAPI.md#AddExecuteSchedules) | **Post** /api/execute-schedules | Creates a Execute Schedule
[**AddPowerScheduleInstances**](AutomationAPI.md#AddPowerScheduleInstances) | **Put** /api/power-schedules/{id}/add-instances | Add Instances to a Power Schedule
[**AddPowerScheduleServers**](AutomationAPI.md#AddPowerScheduleServers) | **Put** /api/power-schedules/{id}/add-servers | Add Servers to a Power Schedule
[**AddPowerSchedules**](AutomationAPI.md#AddPowerSchedules) | **Post** /api/power-schedules | Creates a Power Schedule
[**AddScaleThresholds**](AutomationAPI.md#AddScaleThresholds) | **Post** /api/scale-thresholds | Creates a Scale Threshold
[**AddTasks**](AutomationAPI.md#AddTasks) | **Post** /api/tasks | Creates a Task
[**AddWorkflows**](AutomationAPI.md#AddWorkflows) | **Post** /api/task-sets | Creates a Workflow
[**DeletePowerScheduleInstances**](AutomationAPI.md#DeletePowerScheduleInstances) | **Put** /api/power-schedules/{id}/remove-instances | Remove Instances from a Power Schedule
[**DeletePowerScheduleServers**](AutomationAPI.md#DeletePowerScheduleServers) | **Put** /api/power-schedules/{id}/remove-servers | Remove Servers from a Power Schedule
[**ExecuteExecutionRequest**](AutomationAPI.md#ExecuteExecutionRequest) | **Post** /api/execution-request/execute | Executes an Execution Request
[**ExecuteTasks**](AutomationAPI.md#ExecuteTasks) | **Post** /api/tasks/{id}/execute | Executes a Task
[**ExecuteWorkflows**](AutomationAPI.md#ExecuteWorkflows) | **Post** /api/task-sets/{id}/execute | Executes a Workflow
[**GetExecuteSchedules**](AutomationAPI.md#GetExecuteSchedules) | **Get** /api/execute-schedules/{id} | Retrieves a Specific Execute Schedule
[**GetExecutionRequest**](AutomationAPI.md#GetExecutionRequest) | **Get** /api/execution-request/{uniqueId} | Retrieves a Specific Execution Request
[**GetPowerSchedules**](AutomationAPI.md#GetPowerSchedules) | **Get** /api/power-schedules/{id} | Retrieves a Specific Power Schedule
[**GetScaleThresholds**](AutomationAPI.md#GetScaleThresholds) | **Get** /api/scale-thresholds/{id} | Retrieves a Specific Scale Threshold
[**GetTaskTypes**](AutomationAPI.md#GetTaskTypes) | **Get** /api/task-types/{id} | Retrieves a Specific Task Type
[**GetTasks**](AutomationAPI.md#GetTasks) | **Get** /api/tasks/{id} | Retrieves a Specific Task
[**GetWorkflows**](AutomationAPI.md#GetWorkflows) | **Get** /api/task-sets/{id} | Retrieves a Specific Workflow
[**ListExecuteSchedules**](AutomationAPI.md#ListExecuteSchedules) | **Get** /api/execute-schedules | Retrieves all Execute Schedules
[**ListPowerSchedules**](AutomationAPI.md#ListPowerSchedules) | **Get** /api/power-schedules | Retrieves all Power Schedules
[**ListScaleThresholds**](AutomationAPI.md#ListScaleThresholds) | **Get** /api/scale-thresholds | Retrieves all Scale Thresholds
[**ListTaskTypes**](AutomationAPI.md#ListTaskTypes) | **Get** /api/task-types | Retrieves all Task Types
[**ListTasks**](AutomationAPI.md#ListTasks) | **Get** /api/tasks | Retrieves all Tasks
[**ListWorkflows**](AutomationAPI.md#ListWorkflows) | **Get** /api/task-sets | Retrieves all Workflows
[**RemoveExecuteSchedules**](AutomationAPI.md#RemoveExecuteSchedules) | **Delete** /api/execute-schedules/{id} | Deletes a Execute Schedule
[**RemovePowerSchedules**](AutomationAPI.md#RemovePowerSchedules) | **Delete** /api/power-schedules/{id} | Deletes a Power Schedule
[**RemoveScaleThresholds**](AutomationAPI.md#RemoveScaleThresholds) | **Delete** /api/scale-thresholds/{id} | Deletes a Scale Threshold
[**RemoveTasks**](AutomationAPI.md#RemoveTasks) | **Delete** /api/tasks/{id} | Deletes a Task
[**RemoveWorkflows**](AutomationAPI.md#RemoveWorkflows) | **Delete** /api/task-sets/{id} | Deletes a Workflow
[**UpdateExecuteSchedules**](AutomationAPI.md#UpdateExecuteSchedules) | **Put** /api/execute-schedules/{id} | Updates a Execute Schedule
[**UpdatePowerSchedules**](AutomationAPI.md#UpdatePowerSchedules) | **Put** /api/power-schedules/{id} | Updates a Power Schedule
[**UpdateScaleThresholds**](AutomationAPI.md#UpdateScaleThresholds) | **Put** /api/scale-thresholds/{id} | Updates a Scale Threshold
[**UpdateTasks**](AutomationAPI.md#UpdateTasks) | **Put** /api/tasks/{id} | Updates a Task
[**UpdateWorkflows**](AutomationAPI.md#UpdateWorkflows) | **Put** /api/task-sets/{id} | Updates a Workflow



## AddExecuteSchedules

> AddExecuteSchedules200Response AddExecuteSchedules(ctx).AddExecuteSchedulesRequest(addExecuteSchedulesRequest).Execute()

Creates a Execute Schedule



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
	addExecuteSchedulesRequest := *openapiclient.NewAddExecuteSchedulesRequest(*openapiclient.NewAddExecuteSchedulesRequestSchedule("Sample Execution")) // AddExecuteSchedulesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.AddExecuteSchedules(context.Background()).AddExecuteSchedulesRequest(addExecuteSchedulesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.AddExecuteSchedules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddExecuteSchedules`: AddExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.AddExecuteSchedules`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddExecuteSchedulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addExecuteSchedulesRequest** | [**AddExecuteSchedulesRequest**](AddExecuteSchedulesRequest.md) |  | 

### Return type

[**AddExecuteSchedules200Response**](AddExecuteSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddPowerScheduleInstances

> AddPowerScheduleInstances200Response AddPowerScheduleInstances(ctx, id).AddPowerScheduleInstancesRequest(addPowerScheduleInstancesRequest).Execute()

Add Instances to a Power Schedule



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
	addPowerScheduleInstancesRequest := *openapiclient.NewAddPowerScheduleInstancesRequest([]int64{int64(123)}) // AddPowerScheduleInstancesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.AddPowerScheduleInstances(context.Background(), id).AddPowerScheduleInstancesRequest(addPowerScheduleInstancesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.AddPowerScheduleInstances``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddPowerScheduleInstances`: AddPowerScheduleInstances200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.AddPowerScheduleInstances`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddPowerScheduleInstancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addPowerScheduleInstancesRequest** | [**AddPowerScheduleInstancesRequest**](AddPowerScheduleInstancesRequest.md) |  | 

### Return type

[**AddPowerScheduleInstances200Response**](AddPowerScheduleInstances200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddPowerScheduleServers

> AddPowerScheduleServers200Response AddPowerScheduleServers(ctx, id).AddPowerScheduleServersRequest(addPowerScheduleServersRequest).Execute()

Add Servers to a Power Schedule



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
	addPowerScheduleServersRequest := *openapiclient.NewAddPowerScheduleServersRequest([]int64{int64(123)}) // AddPowerScheduleServersRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.AddPowerScheduleServers(context.Background(), id).AddPowerScheduleServersRequest(addPowerScheduleServersRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.AddPowerScheduleServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddPowerScheduleServers`: AddPowerScheduleServers200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.AddPowerScheduleServers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddPowerScheduleServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addPowerScheduleServersRequest** | [**AddPowerScheduleServersRequest**](AddPowerScheduleServersRequest.md) |  | 

### Return type

[**AddPowerScheduleServers200Response**](AddPowerScheduleServers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddPowerSchedules

> AddPowerSchedules200Response AddPowerSchedules(ctx).AddPowerSchedulesRequest(addPowerSchedulesRequest).Execute()

Creates a Power Schedule



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
	addPowerSchedulesRequest := *openapiclient.NewAddPowerSchedulesRequest(*openapiclient.NewAddPowerSchedulesRequestSchedule("Sample Threshold")) // AddPowerSchedulesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.AddPowerSchedules(context.Background()).AddPowerSchedulesRequest(addPowerSchedulesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.AddPowerSchedules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddPowerSchedules`: AddPowerSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.AddPowerSchedules`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddPowerSchedulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addPowerSchedulesRequest** | [**AddPowerSchedulesRequest**](AddPowerSchedulesRequest.md) |  | 

### Return type

[**AddPowerSchedules200Response**](AddPowerSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddScaleThresholds

> AddScaleThresholds200Response AddScaleThresholds(ctx).AddScaleThresholdsRequest(addScaleThresholdsRequest).Execute()

Creates a Scale Threshold



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
	addScaleThresholdsRequest := *openapiclient.NewAddScaleThresholdsRequest(*openapiclient.NewAddScaleThresholdsRequestScaleThreshold("Sample Threshold")) // AddScaleThresholdsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.AddScaleThresholds(context.Background()).AddScaleThresholdsRequest(addScaleThresholdsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.AddScaleThresholds``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddScaleThresholds`: AddScaleThresholds200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.AddScaleThresholds`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddScaleThresholdsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addScaleThresholdsRequest** | [**AddScaleThresholdsRequest**](AddScaleThresholdsRequest.md) |  | 

### Return type

[**AddScaleThresholds200Response**](AddScaleThresholds200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddTasks

> AddTasks200Response AddTasks(ctx).AddTasksRequest(addTasksRequest).Execute()

Creates a Task



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
	addTasksRequest := *openapiclient.NewAddTasksRequest(*openapiclient.NewAddTasksRequestTask("Sample Task", *openapiclient.NewAddTasksRequestTaskTaskType("Code_example"), "ExecuteTarget_example")) // AddTasksRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.AddTasks(context.Background()).AddTasksRequest(addTasksRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.AddTasks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddTasks`: AddTasks200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.AddTasks`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddTasksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addTasksRequest** | [**AddTasksRequest**](AddTasksRequest.md) |  | 

### Return type

[**AddTasks200Response**](AddTasks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddWorkflows

> AddWorkflows200Response AddWorkflows(ctx).AddWorkflowsRequest(addWorkflowsRequest).Execute()

Creates a Workflow



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
	addWorkflowsRequest := *openapiclient.NewAddWorkflowsRequest(*openapiclient.NewAddWorkflowsRequestTaskSet("Sample Workflow")) // AddWorkflowsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.AddWorkflows(context.Background()).AddWorkflowsRequest(addWorkflowsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.AddWorkflows``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddWorkflows`: AddWorkflows200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.AddWorkflows`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddWorkflowsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addWorkflowsRequest** | [**AddWorkflowsRequest**](AddWorkflowsRequest.md) |  | 

### Return type

[**AddWorkflows200Response**](AddWorkflows200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeletePowerScheduleInstances

> AddPowerScheduleServers200Response DeletePowerScheduleInstances(ctx, id).DeletePowerScheduleInstancesRequest(deletePowerScheduleInstancesRequest).Execute()

Remove Instances from a Power Schedule



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
	deletePowerScheduleInstancesRequest := *openapiclient.NewDeletePowerScheduleInstancesRequest([]int64{int64(123)}) // DeletePowerScheduleInstancesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.DeletePowerScheduleInstances(context.Background(), id).DeletePowerScheduleInstancesRequest(deletePowerScheduleInstancesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.DeletePowerScheduleInstances``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeletePowerScheduleInstances`: AddPowerScheduleServers200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.DeletePowerScheduleInstances`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeletePowerScheduleInstancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **deletePowerScheduleInstancesRequest** | [**DeletePowerScheduleInstancesRequest**](DeletePowerScheduleInstancesRequest.md) |  | 

### Return type

[**AddPowerScheduleServers200Response**](AddPowerScheduleServers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeletePowerScheduleServers

> AddPowerScheduleServers200Response DeletePowerScheduleServers(ctx, id).DeletePowerScheduleServersRequest(deletePowerScheduleServersRequest).Execute()

Remove Servers from a Power Schedule



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
	deletePowerScheduleServersRequest := *openapiclient.NewDeletePowerScheduleServersRequest([]int64{int64(123)}) // DeletePowerScheduleServersRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.DeletePowerScheduleServers(context.Background(), id).DeletePowerScheduleServersRequest(deletePowerScheduleServersRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.DeletePowerScheduleServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeletePowerScheduleServers`: AddPowerScheduleServers200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.DeletePowerScheduleServers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeletePowerScheduleServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **deletePowerScheduleServersRequest** | [**DeletePowerScheduleServersRequest**](DeletePowerScheduleServersRequest.md) |  | 

### Return type

[**AddPowerScheduleServers200Response**](AddPowerScheduleServers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExecuteExecutionRequest

> ExecuteExecutionRequest200Response ExecuteExecutionRequest(ctx).InstanceId(instanceId).ContainerId(containerId).ServerId(serverId).ExecuteExecutionRequestRequest(executeExecutionRequestRequest).Execute()

Executes an Execution Request



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
	instanceId := int64(94) // int64 | The Instance ID for Filtering (optional)
	containerId := int64(135) // int64 | The Container ID for Filtering (optional)
	serverId := int64(97) // int64 | The Server ID for Filtering (optional)
	executeExecutionRequestRequest := *openapiclient.NewExecuteExecutionRequestRequest("uname -a") // ExecuteExecutionRequestRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.ExecuteExecutionRequest(context.Background()).InstanceId(instanceId).ContainerId(containerId).ServerId(serverId).ExecuteExecutionRequestRequest(executeExecutionRequestRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.ExecuteExecutionRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecuteExecutionRequest`: ExecuteExecutionRequest200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.ExecuteExecutionRequest`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiExecuteExecutionRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **instanceId** | **int64** | The Instance ID for Filtering | 
 **containerId** | **int64** | The Container ID for Filtering | 
 **serverId** | **int64** | The Server ID for Filtering | 
 **executeExecutionRequestRequest** | [**ExecuteExecutionRequestRequest**](ExecuteExecutionRequestRequest.md) |  | 

### Return type

[**ExecuteExecutionRequest200Response**](ExecuteExecutionRequest200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExecuteTasks

> ExecuteTasks200Response ExecuteTasks(ctx, id).ExecuteTasksRequest(executeTasksRequest).Execute()

Executes a Task



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
	executeTasksRequest := *openapiclient.NewExecuteTasksRequest(*openapiclient.NewExecuteTasksRequestJob()) // ExecuteTasksRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.ExecuteTasks(context.Background(), id).ExecuteTasksRequest(executeTasksRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.ExecuteTasks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecuteTasks`: ExecuteTasks200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.ExecuteTasks`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiExecuteTasksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **executeTasksRequest** | [**ExecuteTasksRequest**](ExecuteTasksRequest.md) |  | 

### Return type

[**ExecuteTasks200Response**](ExecuteTasks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExecuteWorkflows

> ExecuteWorkflows200Response ExecuteWorkflows(ctx, id).ExecuteWorkflowsRequest(executeWorkflowsRequest).Execute()

Executes a Workflow



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
	executeWorkflowsRequest := *openapiclient.NewExecuteWorkflowsRequest(*openapiclient.NewExecuteWorkflowsRequestJob()) // ExecuteWorkflowsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.ExecuteWorkflows(context.Background(), id).ExecuteWorkflowsRequest(executeWorkflowsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.ExecuteWorkflows``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecuteWorkflows`: ExecuteWorkflows200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.ExecuteWorkflows`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiExecuteWorkflowsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **executeWorkflowsRequest** | [**ExecuteWorkflowsRequest**](ExecuteWorkflowsRequest.md) |  | 

### Return type

[**ExecuteWorkflows200Response**](ExecuteWorkflows200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetExecuteSchedules

> GetExecuteSchedules200Response GetExecuteSchedules(ctx, id).Execute()

Retrieves a Specific Execute Schedule



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
	resp, r, err := apiClient.AutomationAPI.GetExecuteSchedules(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.GetExecuteSchedules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetExecuteSchedules`: GetExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.GetExecuteSchedules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetExecuteSchedulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetExecuteSchedules200Response**](GetExecuteSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetExecutionRequest

> GetExecutionRequest200Response GetExecutionRequest(ctx, uniqueId).Execute()

Retrieves a Specific Execution Request



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
	uniqueId := "c56f75d0-a59a-4566-92e3-4dc716c5d076" // string | The Unique ID of the execution request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.GetExecutionRequest(context.Background(), uniqueId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.GetExecutionRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetExecutionRequest`: GetExecutionRequest200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.GetExecutionRequest`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uniqueId** | **string** | The Unique ID of the execution request | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetExecutionRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetExecutionRequest200Response**](GetExecutionRequest200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPowerSchedules

> GetPowerSchedules200Response GetPowerSchedules(ctx, id).Execute()

Retrieves a Specific Power Schedule



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
	resp, r, err := apiClient.AutomationAPI.GetPowerSchedules(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.GetPowerSchedules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPowerSchedules`: GetPowerSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.GetPowerSchedules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPowerSchedulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetPowerSchedules200Response**](GetPowerSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetScaleThresholds

> GetScaleThresholds200Response GetScaleThresholds(ctx, id).Execute()

Retrieves a Specific Scale Threshold



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
	resp, r, err := apiClient.AutomationAPI.GetScaleThresholds(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.GetScaleThresholds``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetScaleThresholds`: GetScaleThresholds200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.GetScaleThresholds`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetScaleThresholdsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetScaleThresholds200Response**](GetScaleThresholds200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTaskTypes

> GetTaskTypes200Response GetTaskTypes(ctx, id).Execute()

Retrieves a Specific Task Type



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
	resp, r, err := apiClient.AutomationAPI.GetTaskTypes(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.GetTaskTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTaskTypes`: GetTaskTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.GetTaskTypes`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTaskTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetTaskTypes200Response**](GetTaskTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTasks

> GetTasks200Response GetTasks(ctx, id).Execute()

Retrieves a Specific Task



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
	resp, r, err := apiClient.AutomationAPI.GetTasks(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.GetTasks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTasks`: GetTasks200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.GetTasks`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTasksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetTasks200Response**](GetTasks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetWorkflows

> GetWorkflows200Response GetWorkflows(ctx, id).Execute()

Retrieves a Specific Workflow



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
	resp, r, err := apiClient.AutomationAPI.GetWorkflows(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.GetWorkflows``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetWorkflows`: GetWorkflows200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.GetWorkflows`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetWorkflowsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetWorkflows200Response**](GetWorkflows200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListExecuteSchedules

> ListExecuteSchedules200Response ListExecuteSchedules(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Retrieves all Execute Schedules



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
	resp, r, err := apiClient.AutomationAPI.ListExecuteSchedules(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.ListExecuteSchedules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListExecuteSchedules`: ListExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.ListExecuteSchedules`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListExecuteSchedulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListExecuteSchedules200Response**](ListExecuteSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPowerSchedules

> ListPowerSchedules200Response ListPowerSchedules(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Retrieves all Power Schedules



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
	resp, r, err := apiClient.AutomationAPI.ListPowerSchedules(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.ListPowerSchedules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPowerSchedules`: ListPowerSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.ListPowerSchedules`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListPowerSchedulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListPowerSchedules200Response**](ListPowerSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListScaleThresholds

> ListScaleThresholds200Response ListScaleThresholds(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Retrieves all Scale Thresholds



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
	resp, r, err := apiClient.AutomationAPI.ListScaleThresholds(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.ListScaleThresholds``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListScaleThresholds`: ListScaleThresholds200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.ListScaleThresholds`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListScaleThresholdsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListScaleThresholds200Response**](ListScaleThresholds200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListTaskTypes

> ListTaskTypes200Response ListTaskTypes(ctx).Name(name).Code(code).Execute()

Retrieves all Task Types



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
	code := "azr" // string | If specified will return an exact match on code (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.ListTaskTypes(context.Background()).Name(name).Code(code).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.ListTaskTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListTaskTypes`: ListTaskTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.ListTaskTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListTaskTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 

### Return type

[**ListTaskTypes200Response**](ListTaskTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListTasks

> ListTasks200Response ListTasks(ctx).Type_(type_).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Labels(labels).AllLabels(allLabels).Execute()

Retrieves all Tasks



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
	type_ := "type__example" // string | If specified will return all tasks by `task type` code. Refer to `Task Types` API for up to date listings.  (optional)
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
	resp, r, err := apiClient.AutomationAPI.ListTasks(context.Background()).Type_(type_).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Labels(labels).AllLabels(allLabels).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.ListTasks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListTasks`: ListTasks200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.ListTasks`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListTasksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **type_** | **string** | If specified will return all tasks by &#x60;task type&#x60; code. Refer to &#x60;Task Types&#x60; API for up to date listings.  | 
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 

### Return type

[**ListTasks200Response**](ListTasks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListWorkflows

> ListWorkflows200Response ListWorkflows(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Labels(labels).AllLabels(allLabels).Type_(type_).Execute()

Retrieves all Workflows



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
	type_ := "provision" // string | Filter by Workflow Type (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.ListWorkflows(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Labels(labels).AllLabels(allLabels).Type_(type_).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.ListWorkflows``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListWorkflows`: ListWorkflows200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.ListWorkflows`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListWorkflowsRequest struct via the builder pattern


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
 **type_** | **string** | Filter by Workflow Type | 

### Return type

[**ListWorkflows200Response**](ListWorkflows200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveExecuteSchedules

> RemoveExecuteSchedules200Response RemoveExecuteSchedules(ctx, id).Execute()

Deletes a Execute Schedule



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
	resp, r, err := apiClient.AutomationAPI.RemoveExecuteSchedules(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.RemoveExecuteSchedules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveExecuteSchedules`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.RemoveExecuteSchedules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveExecuteSchedulesRequest struct via the builder pattern


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


## RemovePowerSchedules

> RemoveExecuteSchedules200Response RemovePowerSchedules(ctx, id).Execute()

Deletes a Power Schedule



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
	resp, r, err := apiClient.AutomationAPI.RemovePowerSchedules(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.RemovePowerSchedules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemovePowerSchedules`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.RemovePowerSchedules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemovePowerSchedulesRequest struct via the builder pattern


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


## RemoveScaleThresholds

> RemoveExecuteSchedules200Response RemoveScaleThresholds(ctx, id).Execute()

Deletes a Scale Threshold



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
	resp, r, err := apiClient.AutomationAPI.RemoveScaleThresholds(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.RemoveScaleThresholds``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveScaleThresholds`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.RemoveScaleThresholds`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveScaleThresholdsRequest struct via the builder pattern


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


## RemoveTasks

> RemoveExecuteSchedules200Response RemoveTasks(ctx, id).Execute()

Deletes a Task



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
	resp, r, err := apiClient.AutomationAPI.RemoveTasks(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.RemoveTasks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveTasks`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.RemoveTasks`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveTasksRequest struct via the builder pattern


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


## RemoveWorkflows

> RemoveExecuteSchedules200Response RemoveWorkflows(ctx, id).Execute()

Deletes a Workflow



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
	resp, r, err := apiClient.AutomationAPI.RemoveWorkflows(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.RemoveWorkflows``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveWorkflows`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.RemoveWorkflows`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveWorkflowsRequest struct via the builder pattern


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


## UpdateExecuteSchedules

> UpdateExecuteSchedules200Response UpdateExecuteSchedules(ctx, id).UpdateExecuteSchedulesRequest(updateExecuteSchedulesRequest).Execute()

Updates a Execute Schedule



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
	updateExecuteSchedulesRequest := *openapiclient.NewUpdateExecuteSchedulesRequest(*openapiclient.NewUpdateExecuteSchedulesRequestSchedule()) // UpdateExecuteSchedulesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.UpdateExecuteSchedules(context.Background(), id).UpdateExecuteSchedulesRequest(updateExecuteSchedulesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.UpdateExecuteSchedules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateExecuteSchedules`: UpdateExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.UpdateExecuteSchedules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateExecuteSchedulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateExecuteSchedulesRequest** | [**UpdateExecuteSchedulesRequest**](UpdateExecuteSchedulesRequest.md) |  | 

### Return type

[**UpdateExecuteSchedules200Response**](UpdateExecuteSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdatePowerSchedules

> UpdatePowerSchedules200Response UpdatePowerSchedules(ctx, id).UpdatePowerSchedulesRequest(updatePowerSchedulesRequest).Execute()

Updates a Power Schedule



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
	updatePowerSchedulesRequest := *openapiclient.NewUpdatePowerSchedulesRequest(*openapiclient.NewUpdatePowerSchedulesRequestSchedule()) // UpdatePowerSchedulesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.UpdatePowerSchedules(context.Background(), id).UpdatePowerSchedulesRequest(updatePowerSchedulesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.UpdatePowerSchedules``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdatePowerSchedules`: UpdatePowerSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.UpdatePowerSchedules`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdatePowerSchedulesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updatePowerSchedulesRequest** | [**UpdatePowerSchedulesRequest**](UpdatePowerSchedulesRequest.md) |  | 

### Return type

[**UpdatePowerSchedules200Response**](UpdatePowerSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateScaleThresholds

> UpdateScaleThresholds200Response UpdateScaleThresholds(ctx, id).UpdateScaleThresholdsRequest(updateScaleThresholdsRequest).Execute()

Updates a Scale Threshold



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
	updateScaleThresholdsRequest := *openapiclient.NewUpdateScaleThresholdsRequest(*openapiclient.NewUpdateScaleThresholdsRequestScaleThreshold()) // UpdateScaleThresholdsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.UpdateScaleThresholds(context.Background(), id).UpdateScaleThresholdsRequest(updateScaleThresholdsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.UpdateScaleThresholds``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateScaleThresholds`: UpdateScaleThresholds200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.UpdateScaleThresholds`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateScaleThresholdsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateScaleThresholdsRequest** | [**UpdateScaleThresholdsRequest**](UpdateScaleThresholdsRequest.md) |  | 

### Return type

[**UpdateScaleThresholds200Response**](UpdateScaleThresholds200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateTasks

> UpdateTasks200Response UpdateTasks(ctx, id).UpdateTasksRequest(updateTasksRequest).Execute()

Updates a Task



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
	updateTasksRequest := *openapiclient.NewUpdateTasksRequest(*openapiclient.NewUpdateTasksRequestTask()) // UpdateTasksRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.UpdateTasks(context.Background(), id).UpdateTasksRequest(updateTasksRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.UpdateTasks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateTasks`: UpdateTasks200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.UpdateTasks`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateTasksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateTasksRequest** | [**UpdateTasksRequest**](UpdateTasksRequest.md) |  | 

### Return type

[**UpdateTasks200Response**](UpdateTasks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateWorkflows

> UpdateWorkflows200Response UpdateWorkflows(ctx, id).UpdateWorkflowsRequest(updateWorkflowsRequest).Execute()

Updates a Workflow



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
	updateWorkflowsRequest := *openapiclient.NewUpdateWorkflowsRequest(*openapiclient.NewUpdateWorkflowsRequestTaskSet()) // UpdateWorkflowsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AutomationAPI.UpdateWorkflows(context.Background(), id).UpdateWorkflowsRequest(updateWorkflowsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AutomationAPI.UpdateWorkflows``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateWorkflows`: UpdateWorkflows200Response
	fmt.Fprintf(os.Stdout, "Response from `AutomationAPI.UpdateWorkflows`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateWorkflowsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateWorkflowsRequest** | [**UpdateWorkflowsRequest**](UpdateWorkflowsRequest.md) |  | 

### Return type

[**UpdateWorkflows200Response**](UpdateWorkflows200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

