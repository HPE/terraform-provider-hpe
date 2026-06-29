# ServerServicePlans

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Value** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**MaxCpu** | Pointer to **int64** |  | [optional] 
**MaxCores** | Pointer to **int64** |  | [optional] 
**MaxDataStorage** | Pointer to **int64** |  | [optional] 
**CustomCpu** | Pointer to **bool** |  | [optional] 
**CustomMaxMemory** | Pointer to **bool** |  | [optional] 
**CustomMaxStorage** | Pointer to **bool** |  | [optional] 
**CustomMaxDataStorage** | Pointer to **bool** |  | [optional] 
**CustomCoresPerSocket** | Pointer to **bool** |  | [optional] 
**CoresPerSocket** | Pointer to **NullableInt64** |  | [optional] 
**StorageTypes** | Pointer to **[]map[string]interface{}** |  | [optional] 
**RootStorageTypes** | Pointer to **[]map[string]interface{}** |  | [optional] 
**AddVolumes** | Pointer to **bool** |  | [optional] 
**CustomizeVolume** | Pointer to **bool** |  | [optional] 
**RootDiskCustomizable** | Pointer to **bool** |  | [optional] 
**HostDiskMode** | Pointer to **NullableString** |  | [optional] 
**HasDatastore** | Pointer to **NullableString** |  | [optional] 
**LvmSupported** | Pointer to **NullableString** |  | [optional] 
**MinDisk** | Pointer to **NullableString** |  | [optional] 
**MaxDisk** | Pointer to **NullableString** |  | [optional] 
**Datastores** | Pointer to [**ServerServicePlansDatastores**](ServerServicePlansDatastores.md) |  | [optional] 
**SupportsAutoDatastore** | Pointer to **bool** |  | [optional] 
**AutoOptions** | Pointer to [**[]ServerServicePlansAutoOptionsInner**](ServerServicePlansAutoOptionsInner.md) |  | [optional] 
**CpuOptions** | Pointer to **[]map[string]interface{}** |  | [optional] 
**MemoryOptions** | Pointer to **[]map[string]interface{}** |  | [optional] 
**RootCustomSizeOptions** | Pointer to **map[string]interface{}** |  | [optional] 
**CustomSizeOptions** | Pointer to **map[string]interface{}** |  | [optional] 
**CustomCores** | Pointer to **bool** |  | [optional] 
**MaxDisks** | Pointer to **NullableString** |  | [optional] 
**MemorySizeType** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ServerServicePlans{
    // Set fields directly
}
```

### CoresPerSocket (Nullable)

Use the Nullable wrapper methods:
- `obj.CoresPerSocket.IsSet()` — check if set
- `obj.CoresPerSocket.Get()` — get the inner value (returns pointer)
- `obj.CoresPerSocket.Set(&val)` — set the value
- `obj.CoresPerSocket.Unset()` — clear the value
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
### HostDiskMode (Nullable)

Use the Nullable wrapper methods:
- `obj.HostDiskMode.IsSet()` — check if set
- `obj.HostDiskMode.Get()` — get the inner value (returns pointer)
- `obj.HostDiskMode.Set(&val)` — set the value
- `obj.HostDiskMode.Unset()` — clear the value
### HasDatastore (Nullable)

Use the Nullable wrapper methods:
- `obj.HasDatastore.IsSet()` — check if set
- `obj.HasDatastore.Get()` — get the inner value (returns pointer)
- `obj.HasDatastore.Set(&val)` — set the value
- `obj.HasDatastore.Unset()` — clear the value
### LvmSupported (Nullable)

Use the Nullable wrapper methods:
- `obj.LvmSupported.IsSet()` — check if set
- `obj.LvmSupported.Get()` — get the inner value (returns pointer)
- `obj.LvmSupported.Set(&val)` — set the value
- `obj.LvmSupported.Unset()` — clear the value
### MinDisk (Nullable)

Use the Nullable wrapper methods:
- `obj.MinDisk.IsSet()` — check if set
- `obj.MinDisk.Get()` — get the inner value (returns pointer)
- `obj.MinDisk.Set(&val)` — set the value
- `obj.MinDisk.Unset()` — clear the value
### MaxDisk (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxDisk.IsSet()` — check if set
- `obj.MaxDisk.Get()` — get the inner value (returns pointer)
- `obj.MaxDisk.Set(&val)` — set the value
- `obj.MaxDisk.Unset()` — clear the value
### CpuOptions (Nullable)

Use the Nullable wrapper methods:
- `obj.CpuOptions.IsSet()` — check if set
- `obj.CpuOptions.Get()` — get the inner value (returns pointer)
- `obj.CpuOptions.Set(&val)` — set the value
- `obj.CpuOptions.Unset()` — clear the value
### MemoryOptions (Nullable)

Use the Nullable wrapper methods:
- `obj.MemoryOptions.IsSet()` — check if set
- `obj.MemoryOptions.Get()` — get the inner value (returns pointer)
- `obj.MemoryOptions.Set(&val)` — set the value
- `obj.MemoryOptions.Unset()` — clear the value
### MaxDisks (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxDisks.IsSet()` — check if set
- `obj.MaxDisks.Get()` — get the inner value (returns pointer)
- `obj.MaxDisks.Set(&val)` — set the value
- `obj.MaxDisks.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


