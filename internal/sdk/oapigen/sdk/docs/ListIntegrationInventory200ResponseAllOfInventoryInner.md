# ListIntegrationInventory200ResponseAllOfInventoryInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**Owner** | Pointer to [**NullableListIntegrationInventory200ResponseAllOfInventoryInnerOwner**](ListIntegrationInventory200ResponseAllOfInventoryInnerOwner.md) |  | [optional] 
**Tenants** | Pointer to [**[]ListIntegrationInventory200ResponseAllOfInventoryInnerTenantsInner**](ListIntegrationInventory200ResponseAllOfInventoryInnerTenantsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListIntegrationInventory200ResponseAllOfInventoryInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### Owner (Nullable)

Use the Nullable wrapper methods:
- `obj.Owner.IsSet()` — check if set
- `obj.Owner.Get()` — get the inner value (returns pointer)
- `obj.Owner.Set(&val)` — set the value
- `obj.Owner.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


