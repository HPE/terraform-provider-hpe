# ListSecurityGroups200ResponseAllOfSecurityGroupsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**GroupSource** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **NullableString** |  | [optional] 
**SyncSource** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Zone** | Pointer to [**ListSecurityGroups200ResponseAllOfSecurityGroupsInnerZone**](ListSecurityGroups200ResponseAllOfSecurityGroupsInnerZone.md) |  | [optional] 
**Locations** | Pointer to [**[]ListSecurityGroups200ResponseAllOfSecurityGroupsInnerLocationsInner**](ListSecurityGroups200ResponseAllOfSecurityGroupsInnerLocationsInner.md) |  | [optional] 
**Rules** | Pointer to [**[]ListSecurityGroups200ResponseAllOfSecurityGroupsInnerRulesInner**](ListSecurityGroups200ResponseAllOfSecurityGroupsInnerRulesInner.md) |  | [optional] 
**Tenants** | Pointer to [**[]ListSecurityGroups200ResponseAllOfSecurityGroupsInnerTenantsInner**](ListSecurityGroups200ResponseAllOfSecurityGroupsInnerTenantsInner.md) |  | [optional] 
**ResourcePermission** | Pointer to [**ListSecurityGroups200ResponseAllOfSecurityGroupsInnerResourcePermission**](ListSecurityGroups200ResponseAllOfSecurityGroupsInnerResourcePermission.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListSecurityGroups200ResponseAllOfSecurityGroupsInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### GroupSource (Nullable)

Use the Nullable wrapper methods:
- `obj.GroupSource.IsSet()` — check if set
- `obj.GroupSource.Get()` — get the inner value (returns pointer)
- `obj.GroupSource.Set(&val)` — set the value
- `obj.GroupSource.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### Enabled (Nullable)

Use the Nullable wrapper methods:
- `obj.Enabled.IsSet()` — check if set
- `obj.Enabled.Get()` — get the inner value (returns pointer)
- `obj.Enabled.Set(&val)` — set the value
- `obj.Enabled.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


