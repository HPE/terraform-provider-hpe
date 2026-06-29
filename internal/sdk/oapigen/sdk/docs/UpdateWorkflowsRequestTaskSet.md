# UpdateWorkflowsRequestTaskSet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A unique name for the workflow | [optional] 
**Description** | Pointer to **string** | A description of the workflow | [optional] 
**Labels** | Pointer to **[]string** | An array of Category labels for filtering | [optional] 
**Type** | Pointer to **string** | Workflow type | [optional] [default to "provision"]
**OptionTypes** | Pointer to **[]int64** | List of option type IDs for use with operational workflow configuration. | [optional] 
**Tasks** | Pointer to [**UpdateWorkflowsRequestTaskSetTasks**](UpdateWorkflowsRequestTaskSetTasks.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateWorkflowsRequestTaskSet{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


