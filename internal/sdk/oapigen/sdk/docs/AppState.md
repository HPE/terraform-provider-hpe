# AppState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Workloads** | Pointer to [**[]AppStateWorkloadsInner**](AppStateWorkloadsInner.md) |  | [optional] 
**IacDrift** | Pointer to **bool** |  | [optional] 
**PlanResources** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Specs** | Pointer to [**[]AppStateSpecsInner**](AppStateSpecsInner.md) |  | [optional] 
**StateData** | Pointer to **string** |  | [optional] 
**PlanData** | Pointer to **string** |  | [optional] 
**Input** | Pointer to [**AppStateInput**](AppStateInput.md) |  | [optional] 
**Output** | Pointer to [**AppStateOutput**](AppStateOutput.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AppState{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


