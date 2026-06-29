# ListInstances200ResponseAllOfInstancesInnerVolumesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ControllerId** | Pointer to **NullableInt64** |  | [optional] 
**DatastoreId** | Pointer to **NullableString** |  | [optional] 
**DisplayOrder** | Pointer to **int64** |  | [optional] 
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**MaxIOPS** | Pointer to **NullableString** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**ShortName** | Pointer to **string** |  | [optional] 
**Resizeable** | Pointer to **bool** |  | [optional] 
**PlanResizable** | Pointer to **bool** |  | [optional] 
**RootVolume** | Pointer to **bool** |  | [optional] 
**Size** | Pointer to **int64** |  | [optional] 
**StorageType** | Pointer to **int64** |  | [optional] 
**UnitNumber** | Pointer to **NullableString** |  | [optional] 
**ControllerMountPoint** | Pointer to **NullableString** |  | [optional] 
**CreateForMultiAttach** | Pointer to **bool** |  | [optional] 
**StorageProfile** | Pointer to **NullableString** | Storage Profile Code for the volume storage profile assignment. eg. &#x60;\&quot;kvm-cache-none\&quot;&#x60; or &#x60;\&quot;kvm-cache-directsync\&quot;&#x60;. Use &#x60;/api/provision-types?code&#x3D;kvm&#x60; to see the available &#x60;storageProfiles&#x60; for HVM and KVM. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListInstances200ResponseAllOfInstancesInnerVolumesInner{
    // Set fields directly
}
```

### ControllerId (Nullable)

Use the Nullable wrapper methods:
- `obj.ControllerId.IsSet()` — check if set
- `obj.ControllerId.Get()` — get the inner value (returns pointer)
- `obj.ControllerId.Set(&val)` — set the value
- `obj.ControllerId.Unset()` — clear the value
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
### UnitNumber (Nullable)

Use the Nullable wrapper methods:
- `obj.UnitNumber.IsSet()` — check if set
- `obj.UnitNumber.Get()` — get the inner value (returns pointer)
- `obj.UnitNumber.Set(&val)` — set the value
- `obj.UnitNumber.Unset()` — clear the value
### ControllerMountPoint (Nullable)

Use the Nullable wrapper methods:
- `obj.ControllerMountPoint.IsSet()` — check if set
- `obj.ControllerMountPoint.Get()` — get the inner value (returns pointer)
- `obj.ControllerMountPoint.Set(&val)` — set the value
- `obj.ControllerMountPoint.Unset()` — clear the value
### StorageProfile (Nullable)

Use the Nullable wrapper methods:
- `obj.StorageProfile.IsSet()` — check if set
- `obj.StorageProfile.Get()` — get the inner value (returns pointer)
- `obj.StorageProfile.Set(&val)` — set the value
- `obj.StorageProfile.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


