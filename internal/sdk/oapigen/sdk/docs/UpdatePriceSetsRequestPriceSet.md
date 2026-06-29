# UpdatePriceSetsRequestPriceSet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Price set name | [optional] 
**Code** | Pointer to **string** | Price set code. Must be unique. | [optional] 
**RegionCode** | Pointer to **string** | Price set region code | [optional] 
**Zone** | Pointer to [**UpdatePriceSetsRequestPriceSetZone**](UpdatePriceSetsRequestPriceSetZone.md) |  | [optional] 
**ZonePool** | Pointer to [**UpdatePriceSetsRequestPriceSetZonePool**](UpdatePriceSetsRequestPriceSetZonePool.md) |  | [optional] 
**PriceUnit** | Pointer to **string** | Price Unit | [optional] 
**Type** | Pointer to **string** | Price set type | [optional] 
**Prices** | Pointer to **[]int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdatePriceSetsRequestPriceSet{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


