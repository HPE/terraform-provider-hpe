# ListTenantsAvailableRoles200ResponseRolesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Authority** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**RoleType** | Pointer to **NullableString** |  | [optional] 
**Owner** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListTenantsAvailableRoles200ResponseRolesInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### RoleType (Nullable)

Use the Nullable wrapper methods:
- `obj.RoleType.IsSet()` — check if set
- `obj.RoleType.Get()` — get the inner value (returns pointer)
- `obj.RoleType.Set(&val)` — set the value
- `obj.RoleType.Unset()` — clear the value
### Owner (Nullable)

Use the Nullable wrapper methods:
- `obj.Owner.IsSet()` — check if set
- `obj.Owner.Get()` — get the inner value (returns pointer)
- `obj.Owner.Set(&val)` — set the value
- `obj.Owner.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


