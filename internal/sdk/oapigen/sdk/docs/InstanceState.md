# InstanceState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Workloads** | Pointer to **[]map[string]interface{}** |  | [optional] 
**IacDrift** | Pointer to **bool** |  | [optional] 
**PlanResources** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Specs** | Pointer to **[]map[string]interface{}** |  | [optional] 
**StateData** | Pointer to **string** |  | [optional] 
**PlanData** | Pointer to **string** |  | [optional] 
**Input** | Pointer to [**InstanceStateInput**](InstanceStateInput.md) |  | [optional] 
**Output** | Pointer to [**InstanceStateOutput**](InstanceStateOutput.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceState{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


