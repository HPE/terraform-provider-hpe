# UpdateSecurityGroups200ResponseSecurityGroup

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
**Zone** | Pointer to [**UpdateSecurityGroups200ResponseSecurityGroupAllOfZone**](UpdateSecurityGroups200ResponseSecurityGroupAllOfZone.md) |  | [optional] 
**Locations** | Pointer to [**[]UpdateSecurityGroups200ResponseSecurityGroupAllOfLocationsInner**](UpdateSecurityGroups200ResponseSecurityGroupAllOfLocationsInner.md) |  | [optional] 
**Rules** | Pointer to [**[]UpdateSecurityGroups200ResponseSecurityGroupAllOfRulesInner**](UpdateSecurityGroups200ResponseSecurityGroupAllOfRulesInner.md) |  | [optional] 
**Tenants** | Pointer to [**[]UpdateSecurityGroups200ResponseSecurityGroupAllOfTenantsInner**](UpdateSecurityGroups200ResponseSecurityGroupAllOfTenantsInner.md) |  | [optional] 
**ResourcePermission** | Pointer to [**UpdateSecurityGroups200ResponseSecurityGroupAllOfResourcePermission**](UpdateSecurityGroups200ResponseSecurityGroupAllOfResourcePermission.md) |  | [optional] 
**Success** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateSecurityGroups200ResponseSecurityGroup{
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


