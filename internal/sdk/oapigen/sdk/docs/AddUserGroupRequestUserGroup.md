# AddUserGroupRequestUserGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**SudoUser** | Pointer to **bool** |  | [optional] 
**ServerGroup** | Pointer to **NullableString** |  | [optional] 
**Users** | Pointer to **[]int64** | A list of IDs of users that are in the user group | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddUserGroupRequestUserGroup{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ServerGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerGroup.IsSet()` — check if set
- `obj.ServerGroup.Get()` — get the inner value (returns pointer)
- `obj.ServerGroup.Set(&val)` — set the value
- `obj.ServerGroup.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


