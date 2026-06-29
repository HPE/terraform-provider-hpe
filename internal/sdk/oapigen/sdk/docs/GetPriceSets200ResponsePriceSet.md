# GetPriceSets200ResponsePriceSet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**PriceUnit** | Pointer to **string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**RegionCode** | Pointer to **string** |  | [optional] 
**SystemCreated** | Pointer to **bool** |  | [optional] 
**Zone** | Pointer to **NullableString** |  | [optional] 
**ZonePool** | Pointer to **NullableString** |  | [optional] 
**Account** | Pointer to **NullableString** |  | [optional] 
**Prices** | Pointer to [**[]AddPriceSets200ResponseAllOfBudgetPricesInner**](AddPriceSets200ResponseAllOfBudgetPricesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetPriceSets200ResponsePriceSet{
    // Set fields directly
}
```

### Zone (Nullable)

Use the Nullable wrapper methods:
- `obj.Zone.IsSet()` — check if set
- `obj.Zone.Get()` — get the inner value (returns pointer)
- `obj.Zone.Set(&val)` — set the value
- `obj.Zone.Unset()` — clear the value
### ZonePool (Nullable)

Use the Nullable wrapper methods:
- `obj.ZonePool.IsSet()` — check if set
- `obj.ZonePool.Get()` — get the inner value (returns pointer)
- `obj.ZonePool.Set(&val)` — set the value
- `obj.ZonePool.Unset()` — clear the value
### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


