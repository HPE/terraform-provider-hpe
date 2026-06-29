# ClusterMastersVolumesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**ControllerId** | Pointer to **NullableString** |  | [optional] 
**ControllerMountPoint** | Pointer to **NullableString** |  | [optional] 
**Resizeable** | Pointer to **bool** |  | [optional] 
**PlanResizable** | Pointer to **bool** |  | [optional] 
**RootVolume** | Pointer to **bool** |  | [optional] 
**UnitNumber** | Pointer to **NullableString** |  | [optional] 
**TypeId** | Pointer to **int64** |  | [optional] 
**ConfigurableIOPS** | Pointer to **bool** |  | [optional] 
**DatastoreId** | Pointer to **NullableString** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**DisplayOrder** | Pointer to **int64** |  | [optional] 
**MaxIOPS** | Pointer to **NullableString** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterMastersVolumesInner{
    // Set fields directly
}
```

### ControllerId (Nullable)

Use the Nullable wrapper methods:
- `obj.ControllerId.IsSet()` — check if set
- `obj.ControllerId.Get()` — get the inner value (returns pointer)
- `obj.ControllerId.Set(&val)` — set the value
- `obj.ControllerId.Unset()` — clear the value
### ControllerMountPoint (Nullable)

Use the Nullable wrapper methods:
- `obj.ControllerMountPoint.IsSet()` — check if set
- `obj.ControllerMountPoint.Get()` — get the inner value (returns pointer)
- `obj.ControllerMountPoint.Set(&val)` — set the value
- `obj.ControllerMountPoint.Unset()` — clear the value
### UnitNumber (Nullable)

Use the Nullable wrapper methods:
- `obj.UnitNumber.IsSet()` — check if set
- `obj.UnitNumber.Get()` — get the inner value (returns pointer)
- `obj.UnitNumber.Set(&val)` — set the value
- `obj.UnitNumber.Unset()` — clear the value
### DatastoreId (Nullable)

Use the Nullable wrapper methods:
- `obj.DatastoreId.IsSet()` — check if set
- `obj.DatastoreId.Get()` — get the inner value (returns pointer)
- `obj.DatastoreId.Set(&val)` — set the value
- `obj.DatastoreId.Unset()` — clear the value
### MaxIOPS (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxIOPS.IsSet()` — check if set
- `obj.MaxIOPS.Get()` — get the inner value (returns pointer)
- `obj.MaxIOPS.Set(&val)` — set the value
- `obj.MaxIOPS.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


