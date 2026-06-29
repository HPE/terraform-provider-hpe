# InstanceContainerServerVolume

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**ControllerId** | Pointer to **NullableInt64** |  | [optional] 
**ControllerMountPoint** | Pointer to **NullableString** |  | [optional] 
**Resizeable** | Pointer to **bool** |  | [optional] 
**RootVolume** | Pointer to **bool** |  | [optional] 
**DeviceName** | Pointer to **string** |  | [optional] 
**DeviceDisplayName** | Pointer to **string** |  | [optional] 
**DiskMode** | Pointer to **NullableString** |  | [optional] 
**DiskType** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**StorageProfile** | Pointer to **string** | Storage Profile Code for the volume storage profile assignment. eg. &#x60;\&quot;kvm-cache-none\&quot;&#x60; or &#x60;\&quot;kvm-cache-directsync\&quot;&#x60;. Use &#x60;/api/provision-types?code&#x3D;kvm&#x60; to see the available &#x60;storageProfiles&#x60; for HVM and KVM. | [optional] 
**StorageServer** | Pointer to [**InstanceContainerServerVolumeStorageServer**](InstanceContainerServerVolumeStorageServer.md) |  | [optional] 
**ZoneId** | Pointer to **NullableInt64** |  | [optional] 
**Zone** | Pointer to [**InstanceContainerServerVolumeZone**](InstanceContainerServerVolumeZone.md) |  | [optional] 
**Datastore** | Pointer to [**InstanceContainerServerVolumeDatastore**](InstanceContainerServerVolumeDatastore.md) |  | [optional] 
**UnitNumber** | Pointer to **NullableString** |  | [optional] 
**TypeId** | Pointer to **int64** |  | [optional] 
**ConfigurableIOPS** | Pointer to **bool** |  | [optional] 
**DatastoreId** | Pointer to **NullableInt64** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**MaxIOPS** | Pointer to **NullableString** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**UniqueId** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceContainerServerVolume{
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
### DiskMode (Nullable)

Use the Nullable wrapper methods:
- `obj.DiskMode.IsSet()` — check if set
- `obj.DiskMode.Get()` — get the inner value (returns pointer)
- `obj.DiskMode.Set(&val)` — set the value
- `obj.DiskMode.Unset()` — clear the value
### DiskType (Nullable)

Use the Nullable wrapper methods:
- `obj.DiskType.IsSet()` — check if set
- `obj.DiskType.Get()` — get the inner value (returns pointer)
- `obj.DiskType.Set(&val)` — set the value
- `obj.DiskType.Unset()` — clear the value
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### ZoneId (Nullable)

Use the Nullable wrapper methods:
- `obj.ZoneId.IsSet()` — check if set
- `obj.ZoneId.Get()` — get the inner value (returns pointer)
- `obj.ZoneId.Set(&val)` — set the value
- `obj.ZoneId.Unset()` — clear the value
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
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


