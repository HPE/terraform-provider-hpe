# GetVDIPools200ResponseVdiPoolConfigVolumesDisplayInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Storage** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Controller** | Pointer to **NullableString** |  | [optional] 
**Datastore** | Pointer to **string** |  | [optional] 
**DisplayOrder** | Pointer to **NullableString** |  | [optional] 
**Size** | Pointer to **int64** |  | [optional] 
**MountPoint** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetVDIPools200ResponseVdiPoolConfigVolumesDisplayInner{
    // Set fields directly
}
```

### Controller (Nullable)

Use the Nullable wrapper methods:
- `obj.Controller.IsSet()` — check if set
- `obj.Controller.Get()` — get the inner value (returns pointer)
- `obj.Controller.Set(&val)` — set the value
- `obj.Controller.Unset()` — clear the value
### DisplayOrder (Nullable)

Use the Nullable wrapper methods:
- `obj.DisplayOrder.IsSet()` — check if set
- `obj.DisplayOrder.Get()` — get the inner value (returns pointer)
- `obj.DisplayOrder.Set(&val)` — set the value
- `obj.DisplayOrder.Unset()` — clear the value
### MountPoint (Nullable)

Use the Nullable wrapper methods:
- `obj.MountPoint.IsSet()` — check if set
- `obj.MountPoint.Get()` — get the inner value (returns pointer)
- `obj.MountPoint.Set(&val)` — set the value
- `obj.MountPoint.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


