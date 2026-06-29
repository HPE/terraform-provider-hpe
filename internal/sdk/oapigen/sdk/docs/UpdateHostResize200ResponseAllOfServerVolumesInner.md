# UpdateHostResize200ResponseAllOfServerVolumesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Controller** | Pointer to [**UpdateHostResize200ResponseAllOfServerVolumesInnerController**](UpdateHostResize200ResponseAllOfServerVolumesInnerController.md) |  | [optional] 
**ControllerId** | Pointer to **NullableInt64** |  | [optional] 
**ControllerMountPoint** | Pointer to **NullableString** |  | [optional] 
**Resizeable** | Pointer to **NullableBool** |  | [optional] 
**RootVolume** | Pointer to **bool** |  | [optional] 
**UnitNumber** | Pointer to **NullableString** |  | [optional] 
**DeviceName** | Pointer to **string** |  | [optional] 
**DeviceDisplayName** | Pointer to **string** |  | [optional] 
**VolumeName** | Pointer to **string** |  | [optional] 
**VolumePath** | Pointer to **string** |  | [optional] 
**VolumeType** | Pointer to **string** |  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**DiskMode** | Pointer to **string** |  | [optional] 
**DiskType** | Pointer to **string** |  | [optional] 
**Type** | Pointer to [**UpdateHostResize200ResponseAllOfServerVolumesInnerType**](UpdateHostResize200ResponseAllOfServerVolumesInnerType.md) |  | [optional] 
**TypeId** | Pointer to **int64** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**ConfigurableIOPS** | Pointer to **bool** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**UsedStorage** | Pointer to **int64** |  | [optional] 
**DisplayOrder** | Pointer to **int64** |  | [optional] 
**MaxIOPS** | Pointer to **NullableString** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**ReadOnly** | Pointer to **bool** |  | [optional] 
**Removable** | Pointer to **bool** |  | [optional] 
**PoolName** | Pointer to **string** |  | [optional] 
**Zone** | Pointer to [**UpdateHostResize200ResponseAllOfServerVolumesInnerZone**](UpdateHostResize200ResponseAllOfServerVolumesInnerZone.md) |  | [optional] 
**ZoneId** | Pointer to **int64** |  | [optional] 
**Datastore** | Pointer to [**UpdateHostResize200ResponseAllOfServerVolumesInnerDatastore**](UpdateHostResize200ResponseAllOfServerVolumesInnerDatastore.md) |  | [optional] 
**DatastoreId** | Pointer to **NullableInt64** |  | [optional] 
**DatastoreOption** | Pointer to **string** |  | [optional] 
**StorageGroup** | Pointer to **NullableString** |  | [optional] 
**Namespace** | Pointer to **NullableString** |  | [optional] 
**StorageServer** | Pointer to **map[string]interface{}** |  | [optional] 
**Source** | Pointer to **string** |  | [optional] 
**UniqueId** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**ProvisionType** | Pointer to **NullableString** |  | [optional] 
**CopyType** | Pointer to **NullableString** |  | [optional] 
**FiberWwn** | Pointer to **NullableString** |  | [optional] 
**FileName** | Pointer to **NullableString** |  | [optional] 
**ClaimName** | Pointer to **NullableString** |  | [optional] 
**SharePath** | Pointer to **NullableString** |  | [optional] 
**SourceId** | Pointer to **string** |  | [optional] 
**SourceImage** | Pointer to **string** |  | [optional] 
**ImageType** | Pointer to **string** |  | [optional] 
**Online** | Pointer to **bool** |  | [optional] 
**RawData** | Pointer to **string** |  | [optional] 
**CreateForMultiAttach** | Pointer to **bool** |  | [optional] 
**IsMultiAttach** | Pointer to **bool** |  | [optional] 
**StorageProfile** | Pointer to **NullableString** | Storage Profile Code for the volume storage profile assignment. eg. &#x60;\&quot;kvm-cache-none\&quot;&#x60; or &#x60;\&quot;kvm-cache-directsync\&quot;&#x60;. Use &#x60;/api/provision-types?code&#x3D;kvm&#x60; to see the available &#x60;storageProfiles&#x60; for HVM and KVM. | [optional] 
**Account** | Pointer to [**UpdateHostResize200ResponseAllOfServerVolumesInnerAccount**](UpdateHostResize200ResponseAllOfServerVolumesInnerAccount.md) |  | [optional] 
**Owner** | Pointer to [**UpdateHostResize200ResponseAllOfServerVolumesInnerOwner**](UpdateHostResize200ResponseAllOfServerVolumesInnerOwner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateHostResize200ResponseAllOfServerVolumesInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
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
### Resizeable (Nullable)

Use the Nullable wrapper methods:
- `obj.Resizeable.IsSet()` — check if set
- `obj.Resizeable.Get()` — get the inner value (returns pointer)
- `obj.Resizeable.Set(&val)` — set the value
- `obj.Resizeable.Unset()` — clear the value
### UnitNumber (Nullable)

Use the Nullable wrapper methods:
- `obj.UnitNumber.IsSet()` — check if set
- `obj.UnitNumber.Get()` — get the inner value (returns pointer)
- `obj.UnitNumber.Set(&val)` — set the value
- `obj.UnitNumber.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### MaxIOPS (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxIOPS.IsSet()` — check if set
- `obj.MaxIOPS.Get()` — get the inner value (returns pointer)
- `obj.MaxIOPS.Set(&val)` — set the value
- `obj.MaxIOPS.Unset()` — clear the value
### DatastoreId (Nullable)

Use the Nullable wrapper methods:
- `obj.DatastoreId.IsSet()` — check if set
- `obj.DatastoreId.Get()` — get the inner value (returns pointer)
- `obj.DatastoreId.Set(&val)` — set the value
- `obj.DatastoreId.Unset()` — clear the value
### StorageGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.StorageGroup.IsSet()` — check if set
- `obj.StorageGroup.Get()` — get the inner value (returns pointer)
- `obj.StorageGroup.Set(&val)` — set the value
- `obj.StorageGroup.Unset()` — clear the value
### Namespace (Nullable)

Use the Nullable wrapper methods:
- `obj.Namespace.IsSet()` — check if set
- `obj.Namespace.Get()` — get the inner value (returns pointer)
- `obj.Namespace.Set(&val)` — set the value
- `obj.Namespace.Unset()` — clear the value
### UniqueId (Nullable)

Use the Nullable wrapper methods:
- `obj.UniqueId.IsSet()` — check if set
- `obj.UniqueId.Get()` — get the inner value (returns pointer)
- `obj.UniqueId.Set(&val)` — set the value
- `obj.UniqueId.Unset()` — clear the value
### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### ProvisionType (Nullable)

Use the Nullable wrapper methods:
- `obj.ProvisionType.IsSet()` — check if set
- `obj.ProvisionType.Get()` — get the inner value (returns pointer)
- `obj.ProvisionType.Set(&val)` — set the value
- `obj.ProvisionType.Unset()` — clear the value
### CopyType (Nullable)

Use the Nullable wrapper methods:
- `obj.CopyType.IsSet()` — check if set
- `obj.CopyType.Get()` — get the inner value (returns pointer)
- `obj.CopyType.Set(&val)` — set the value
- `obj.CopyType.Unset()` — clear the value
### FiberWwn (Nullable)

Use the Nullable wrapper methods:
- `obj.FiberWwn.IsSet()` — check if set
- `obj.FiberWwn.Get()` — get the inner value (returns pointer)
- `obj.FiberWwn.Set(&val)` — set the value
- `obj.FiberWwn.Unset()` — clear the value
### FileName (Nullable)

Use the Nullable wrapper methods:
- `obj.FileName.IsSet()` — check if set
- `obj.FileName.Get()` — get the inner value (returns pointer)
- `obj.FileName.Set(&val)` — set the value
- `obj.FileName.Unset()` — clear the value
### ClaimName (Nullable)

Use the Nullable wrapper methods:
- `obj.ClaimName.IsSet()` — check if set
- `obj.ClaimName.Get()` — get the inner value (returns pointer)
- `obj.ClaimName.Set(&val)` — set the value
- `obj.ClaimName.Unset()` — clear the value
### SharePath (Nullable)

Use the Nullable wrapper methods:
- `obj.SharePath.IsSet()` — check if set
- `obj.SharePath.Get()` — get the inner value (returns pointer)
- `obj.SharePath.Set(&val)` — set the value
- `obj.SharePath.Unset()` — clear the value
### StorageProfile (Nullable)

Use the Nullable wrapper methods:
- `obj.StorageProfile.IsSet()` — check if set
- `obj.StorageProfile.Get()` — get the inner value (returns pointer)
- `obj.StorageProfile.Set(&val)` — set the value
- `obj.StorageProfile.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


