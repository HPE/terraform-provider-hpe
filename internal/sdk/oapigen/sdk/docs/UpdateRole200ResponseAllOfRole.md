# UpdateRole200ResponseAllOfRole

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** | a unique name of the role | [optional] 
**Authority** | Pointer to **string** | Alias for name | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**LandingUrl** | Pointer to **NullableString** | An optional override for the default landing page after login for a user. | [optional] 
**Scope** | Pointer to **string** |  | [optional] 
**RoleType** | Pointer to **string** |  | [optional] 
**Multitenant** | Pointer to **bool** |  | [optional] 
**MultitenantLocked** | Pointer to **bool** |  | [optional] 
**ParentRoleId** | Pointer to **NullableString** |  | [optional] 
**Diverged** | Pointer to **bool** |  | [optional] 
**OwnerId** | Pointer to **int64** |  | [optional] 
**Owner** | Pointer to [**UpdateRole200ResponseAllOfRoleOwner**](UpdateRole200ResponseAllOfRoleOwner.md) |  | [optional] 
**DefaultPersona** | Pointer to [**UpdateRole200ResponseAllOfRoleDefaultPersona**](UpdateRole200ResponseAllOfRoleDefaultPersona.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateRole200ResponseAllOfRole{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### LandingUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.LandingUrl.IsSet()` — check if set
- `obj.LandingUrl.Get()` — get the inner value (returns pointer)
- `obj.LandingUrl.Set(&val)` — set the value
- `obj.LandingUrl.Unset()` — clear the value
### ParentRoleId (Nullable)

Use the Nullable wrapper methods:
- `obj.ParentRoleId.IsSet()` — check if set
- `obj.ParentRoleId.Get()` — get the inner value (returns pointer)
- `obj.ParentRoleId.Set(&val)` — set the value
- `obj.ParentRoleId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


