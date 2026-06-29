# RunWorkflowInstanceRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TaskSet** | Pointer to [**RunWorkflowInstanceRequestTaskSet**](RunWorkflowInstanceRequestTaskSet.md) |  | [optional] 
**TaskPhase** | Pointer to **string** | Task Phase to run for Provisioning workflows. The default is &#x60;provision&#x60;. | [optional] [default to "provision"]

## Usage

Instantiate with a Go composite literal:

```go
obj := &RunWorkflowInstanceRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


