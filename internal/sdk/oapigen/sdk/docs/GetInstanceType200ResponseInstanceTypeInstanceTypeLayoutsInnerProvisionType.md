# GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**AclEnabled** | Pointer to **bool** |  | [optional] 
**MultiTenant** | Pointer to **bool** |  | [optional] 
**Managed** | Pointer to **bool** |  | [optional] 
**HostNetwork** | Pointer to **bool** |  | [optional] 
**CustomSupported** | Pointer to **bool** |  | [optional] 
**MapPorts** | Pointer to **bool** |  | [optional] 
**ExportServer** | Pointer to **bool** |  | [optional] 
**ViewSet** | Pointer to **NullableString** |  | [optional] 
**ServerType** | Pointer to **string** |  | [optional] 
**HostType** | Pointer to **string** |  | [optional] 
**AddVolumes** | Pointer to **bool** |  | [optional] 
**HasVolumes** | Pointer to **bool** |  | [optional] 
**HasDatastore** | Pointer to **bool** |  | [optional] 
**HasNetworks** | Pointer to **bool** |  | [optional] 
**MaxNetworks** | Pointer to **int64** |  | [optional] 
**CustomizeVolume** | Pointer to **bool** |  | [optional] 
**RootDiskCustomizable** | Pointer to **bool** |  | [optional] 
**RootDiskSizeKnown** | Pointer to **bool** |  | [optional] 
**RootDiskResizable** | Pointer to **bool** |  | [optional] 
**LvmSupported** | Pointer to **bool** |  | [optional] 
**HostDiskMode** | Pointer to **string** |  | [optional] 
**MinDisk** | Pointer to **int64** |  | [optional] 
**MaxDisk** | Pointer to **NullableString** |  | [optional] 
**ResizeCopiesVolumes** | Pointer to **bool** |  | [optional] 
**SupportsAutoDatastore** | Pointer to **bool** |  | [optional] 
**HasZonePools** | Pointer to **bool** |  | [optional] 
**HasSecurityGroups** | Pointer to **bool** |  | [optional] 
**HasParameters** | Pointer to **bool** |  | [optional] 
**CanEnforceTags** | Pointer to **NullableBool** |  | [optional] 
**DisableRootDatastore** | Pointer to **bool** |  | [optional] 
**HasSnapshots** | Pointer to **bool** |  | [optional] 
**HasMemorySnapshots** | Pointer to **bool** |  | [optional] 
**HasSpecTemplates** | Pointer to **bool** |  | [optional] 
**HasPreview** | Pointer to **bool** |  | [optional] 
**ZonePoolRequired** | Pointer to **bool** |  | [optional] 
**PlanRequiresPool** | Pointer to **bool** |  | [optional] 
**HasFolders** | Pointer to **NullableBool** |  | [optional] 
**OptionTypes** | Pointer to [**[]GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeOptionTypesInner**](GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeOptionTypesInner.md) |  | [optional] 
**CustomOptionTypes** | Pointer to [**[]GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeCustomOptionTypesInner**](GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeCustomOptionTypesInner.md) |  | [optional] 
**NetworkTypes** | Pointer to [**[]GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeNetworkTypesInner**](GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeNetworkTypesInner.md) |  | [optional] 
**StorageTypes** | Pointer to [**[]GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeStorageTypesInner**](GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeStorageTypesInner.md) |  | [optional] 
**RootStorageTypes** | Pointer to [**[]GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeRootStorageTypesInner**](GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeRootStorageTypesInner.md) |  | [optional] 
**ControllerTypes** | Pointer to [**[]GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeControllerTypesInner**](GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeControllerTypesInner.md) |  | [optional] 
**StorageProfiles** | Pointer to [**[]GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeStorageProfilesInner**](GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionTypeStorageProfilesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInnerProvisionType{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ViewSet (Nullable)

Use the Nullable wrapper methods:
- `obj.ViewSet.IsSet()` — check if set
- `obj.ViewSet.Get()` — get the inner value (returns pointer)
- `obj.ViewSet.Set(&val)` — set the value
- `obj.ViewSet.Unset()` — clear the value
### MaxDisk (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxDisk.IsSet()` — check if set
- `obj.MaxDisk.Get()` — get the inner value (returns pointer)
- `obj.MaxDisk.Set(&val)` — set the value
- `obj.MaxDisk.Unset()` — clear the value
### CanEnforceTags (Nullable)

Use the Nullable wrapper methods:
- `obj.CanEnforceTags.IsSet()` — check if set
- `obj.CanEnforceTags.Get()` — get the inner value (returns pointer)
- `obj.CanEnforceTags.Set(&val)` — set the value
- `obj.CanEnforceTags.Unset()` — clear the value
### HasFolders (Nullable)

Use the Nullable wrapper methods:
- `obj.HasFolders.IsSet()` — check if set
- `obj.HasFolders.Get()` — get the inner value (returns pointer)
- `obj.HasFolders.Set(&val)` — set the value
- `obj.HasFolders.Unset()` — clear the value
### OptionTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.OptionTypes.IsSet()` — check if set
- `obj.OptionTypes.Get()` — get the inner value (returns pointer)
- `obj.OptionTypes.Set(&val)` — set the value
- `obj.OptionTypes.Unset()` — clear the value
### CustomOptionTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.CustomOptionTypes.IsSet()` — check if set
- `obj.CustomOptionTypes.Get()` — get the inner value (returns pointer)
- `obj.CustomOptionTypes.Set(&val)` — set the value
- `obj.CustomOptionTypes.Unset()` — clear the value
### NetworkTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkTypes.IsSet()` — check if set
- `obj.NetworkTypes.Get()` — get the inner value (returns pointer)
- `obj.NetworkTypes.Set(&val)` — set the value
- `obj.NetworkTypes.Unset()` — clear the value
### StorageTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.StorageTypes.IsSet()` — check if set
- `obj.StorageTypes.Get()` — get the inner value (returns pointer)
- `obj.StorageTypes.Set(&val)` — set the value
- `obj.StorageTypes.Unset()` — clear the value
### RootStorageTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.RootStorageTypes.IsSet()` — check if set
- `obj.RootStorageTypes.Get()` — get the inner value (returns pointer)
- `obj.RootStorageTypes.Set(&val)` — set the value
- `obj.RootStorageTypes.Unset()` — clear the value
### ControllerTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.ControllerTypes.IsSet()` — check if set
- `obj.ControllerTypes.Get()` — get the inner value (returns pointer)
- `obj.ControllerTypes.Set(&val)` — set the value
- `obj.ControllerTypes.Unset()` — clear the value
### StorageProfiles (Nullable)

Use the Nullable wrapper methods:
- `obj.StorageProfiles.IsSet()` — check if set
- `obj.StorageProfiles.Get()` — get the inner value (returns pointer)
- `obj.StorageProfiles.Set(&val)` — set the value
- `obj.StorageProfiles.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


