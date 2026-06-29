# GetAppState200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Workloads** | Pointer to [**[]GetAppState200ResponseAllOfWorkloadsInner**](GetAppState200ResponseAllOfWorkloadsInner.md) |  | [optional] 
**IacDrift** | Pointer to **bool** |  | [optional] 
**PlanResources** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Specs** | Pointer to [**[]GetAppState200ResponseAllOfSpecsInner**](GetAppState200ResponseAllOfSpecsInner.md) |  | [optional] 
**StateData** | Pointer to **string** |  | [optional] 
**PlanData** | Pointer to **string** |  | [optional] 
**Input** | Pointer to [**GetAppState200ResponseAllOfInput**](GetAppState200ResponseAllOfInput.md) |  | [optional] 
**Output** | Pointer to [**GetAppState200ResponseAllOfOutput**](GetAppState200ResponseAllOfOutput.md) |  | [optional] 
**Success** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetAppState200Response{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


