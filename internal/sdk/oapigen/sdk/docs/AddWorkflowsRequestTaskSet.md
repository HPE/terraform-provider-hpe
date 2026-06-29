# AddWorkflowsRequestTaskSet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A unique name for the workflow | 
**Description** | Pointer to **string** | A description of the workflow | [optional] 
**Labels** | Pointer to **[]string** | An array of Category labels for filtering | [optional] 
**Type** | Pointer to **string** | Workflow type | [optional] [default to "provision"]
**Visibility** | Pointer to **string** | private or public | [optional] [default to "private"]
**OptionTypes** | Pointer to **[]int64** | List of option type IDs for use with operational workflow configuration. | [optional] 
**Tasks** | Pointer to [**[]AddWorkflowsRequestTaskSetTasksInner**](AddWorkflowsRequestTaskSetTasksInner.md) | List of task objects in order | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddWorkflowsRequestTaskSet{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


