# UsersAvailableRolesRolesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Authority** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**RoleType** | Pointer to **string** |  | [optional] 
**Owner** | Pointer to [**UsersAvailableRolesRolesInnerOwner**](UsersAvailableRolesRolesInnerOwner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UsersAvailableRolesRolesInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


