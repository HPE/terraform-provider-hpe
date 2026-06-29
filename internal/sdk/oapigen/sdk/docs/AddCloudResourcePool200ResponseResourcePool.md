# AddCloudResourcePool200ResponseResourcePool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Zone** | Pointer to [**NullableAddCloudResourcePool200ResponseResourcePoolAllOfZone**](AddCloudResourcePool200ResponseResourcePoolAllOfZone.md) |  | [optional] 
**Parent** | Pointer to [**AddCloudResourcePool200ResponseResourcePoolAllOfParent**](AddCloudResourcePool200ResponseResourcePoolAllOfParent.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**RegionCode** | Pointer to **NullableString** |  | [optional] 
**IacId** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**ReadOnly** | Pointer to **bool** |  | [optional] 
**DefaultPool** | Pointer to **bool** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Inventory** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to [**AddCloudResourcePool200ResponseResourcePoolAllOfConfig**](AddCloudResourcePool200ResponseResourcePoolAllOfConfig.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**DisplayName** | Pointer to **NullableString** |  | [optional] 
**Tenants** | Pointer to [**[]AddCloudResourcePool200ResponseResourcePoolAllOfTenantsInner**](AddCloudResourcePool200ResponseResourcePoolAllOfTenantsInner.md) |  | [optional] 
**ResourcePermission** | Pointer to [**AddCloudResourcePool200ResponseResourcePoolAllOfResourcePermission**](AddCloudResourcePool200ResponseResourcePoolAllOfResourcePermission.md) |  | [optional] 
**Depth** | Pointer to **int64** |  | [optional] 
**Success** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCloudResourcePool200ResponseResourcePool{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Zone (Nullable)

Use the Nullable wrapper methods:
- `obj.Zone.IsSet()` — check if set
- `obj.Zone.Get()` — get the inner value (returns pointer)
- `obj.Zone.Set(&val)` — set the value
- `obj.Zone.Unset()` — clear the value
### RegionCode (Nullable)

Use the Nullable wrapper methods:
- `obj.RegionCode.IsSet()` — check if set
- `obj.RegionCode.Get()` — get the inner value (returns pointer)
- `obj.RegionCode.Set(&val)` — set the value
- `obj.RegionCode.Unset()` — clear the value
### IacId (Nullable)

Use the Nullable wrapper methods:
- `obj.IacId.IsSet()` — check if set
- `obj.IacId.Get()` — get the inner value (returns pointer)
- `obj.IacId.Set(&val)` — set the value
- `obj.IacId.Unset()` — clear the value
### DisplayName (Nullable)

Use the Nullable wrapper methods:
- `obj.DisplayName.IsSet()` — check if set
- `obj.DisplayName.Get()` — get the inner value (returns pointer)
- `obj.DisplayName.Set(&val)` — set the value
- `obj.DisplayName.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


