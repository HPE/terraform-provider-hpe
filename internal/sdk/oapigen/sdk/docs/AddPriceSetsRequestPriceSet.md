# AddPriceSetsRequestPriceSet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Price set name | 
**Code** | **string** | Price set code. Must be unique. | 
**RegionCode** | Pointer to **string** | Price set region code | [optional] 
**Zone** | Pointer to [**AddPriceSetsRequestPriceSetZone**](AddPriceSetsRequestPriceSetZone.md) |  | [optional] 
**ZonePool** | Pointer to [**AddPriceSetsRequestPriceSetZonePool**](AddPriceSetsRequestPriceSetZonePool.md) |  | [optional] 
**PriceUnit** | **string** | Price Unit | 
**Type** | **string** | Price set type | 
**Prices** | Pointer to **[]int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddPriceSetsRequestPriceSet{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


