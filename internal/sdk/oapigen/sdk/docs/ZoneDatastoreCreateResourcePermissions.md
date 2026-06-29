# ZoneDatastoreCreateResourcePermissions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllGroups** | Pointer to **bool** |  | [optional] 
**DefaultStore** | Pointer to **bool** |  | [optional] 
**AllPlans** | Pointer to **bool** |  | [optional] 
**DefaultTarget** | Pointer to **bool** |  | [optional] 
**MorpheusResourceType** | Pointer to **string** |  | [optional] 
**MorpheusResourceId** | Pointer to **int64** |  | [optional] 
**CanManage** | Pointer to **bool** |  | [optional] 
**All** | Pointer to **bool** |  | [optional] 
**Account** | Pointer to [**ZoneDatastoreCreateResourcePermissionsAccount**](ZoneDatastoreCreateResourcePermissionsAccount.md) |  | [optional] 
**Sites** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Plans** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ZoneDatastoreCreateResourcePermissions{
    // Set fields directly
}
```

### Sites (Nullable)

Use the Nullable wrapper methods:
- `obj.Sites.IsSet()` — check if set
- `obj.Sites.Get()` — get the inner value (returns pointer)
- `obj.Sites.Set(&val)` — set the value
- `obj.Sites.Unset()` — clear the value
### Plans (Nullable)

Use the Nullable wrapper methods:
- `obj.Plans.IsSet()` — check if set
- `obj.Plans.Get()` — get the inner value (returns pointer)
- `obj.Plans.Set(&val)` — set the value
- `obj.Plans.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


