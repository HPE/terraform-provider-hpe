# UpdateServicePlansRequestServicePlanConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StorageSizeType** | Pointer to **string** | Specifies range min / max storage multiplier | [optional] [default to "gb"]
**MemorySizeType** | Pointer to **string** | Specifies range min / max memory multiplier | [optional] [default to "mb"]
**Ranges** | Pointer to [**UpdateServicePlansRequestServicePlanConfigRanges**](UpdateServicePlansRequestServicePlanConfigRanges.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateServicePlansRequestServicePlanConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


