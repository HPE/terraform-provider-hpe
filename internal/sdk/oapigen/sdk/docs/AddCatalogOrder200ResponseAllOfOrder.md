# AddCatalogOrder200ResponseAllOfOrder

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **NullableString** |  | [optional] 
**Items** | Pointer to [**[]AddCatalogOrder200ResponseAllOfOrderItemsInner**](AddCatalogOrder200ResponseAllOfOrderItemsInner.md) |  | [optional] 
**Stats** | Pointer to [**AddCatalogOrder200ResponseAllOfOrderStats**](AddCatalogOrder200ResponseAllOfOrderStats.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCatalogOrder200ResponseAllOfOrder{
    // Set fields directly
}
```

### Name (Nullable)

Use the Nullable wrapper methods:
- `obj.Name.IsSet()` — check if set
- `obj.Name.Get()` — get the inner value (returns pointer)
- `obj.Name.Set(&val)` — set the value
- `obj.Name.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


