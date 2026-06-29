# ConditionalWorkflowTaskConfig4

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConditionalScript** | Pointer to **string** | Allows the user to set JavaScript logic. If it resolves to true, HPE Morpheus Enterprise will run the Operational Workflow set as the “IF OPERATIONAL WORKFLOW” and if it resolves to false, HPE Morpheus Enterprise will run the “ELSE OPERATIONAL WORKFLOW” | [optional] 
**IfOperationalWorkflowId** | Pointer to **int64** | If Operational Workflow ID | [optional] 
**IfOperationalWorkflowName** | Pointer to **string** | If Operational Workflow Name | [optional] 
**ElseOperationalWorkflowId** | Pointer to **int64** | else Operational Workflow ID | [optional] 
**ElseOperationalWorkflowName** | Pointer to **string** | Else Operational Workflow Name | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ConditionalWorkflowTaskConfig4{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


