# ExecuteTasksRequestJob

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the execution job. Can be used to find execution results with &#x60;/api/processes?name&#x3D;&#x60; | [optional] 
**TargetType** | Pointer to **string** | The target context for task execution. This is required for tasks with &#x60;executeTarget&#x60; set to &#x60;resource&#x60;. | [optional] 
**Instances** | Pointer to **[]int64** | Array of Instance IDs. Only applicable if &#x60;targetType&#x60; is instance. | [optional] 
**Servers** | Pointer to **[]int64** | Array of Server IDs. Only applicable if &#x60;targetType&#x60; is &#x60;server&#x60;. | [optional] 
**InstanceLabel** | Pointer to **string** | Instance Label. Only applicable if &#x60;targetType&#x60; is &#x60;instance-label&#x60;. | [optional] 
**ServerLabel** | Pointer to **string** | Server Label. Only applicable if &#x60;targetType&#x60; is &#x60;server-label&#x60;. | [optional] 
**CustomOptions** | Pointer to **map[string]interface{}** | Map of options to be used as values in the task. These correspond to option types. | [optional] 
**CustomConfig** | Pointer to **string** | String of custom configuration values as JSON. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ExecuteTasksRequestJob{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


