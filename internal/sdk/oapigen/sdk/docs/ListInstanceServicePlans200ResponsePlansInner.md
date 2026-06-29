# ListInstanceServicePlans200ResponsePlansInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int32** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Value** | Pointer to **int32** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**MaxMemory** | Pointer to **int32** |  | [optional] 
**MaxCpu** | Pointer to **int32** |  | [optional] 
**MaxCores** | Pointer to **int32** |  | [optional] 
**CustomCpu** | Pointer to **bool** |  | [optional] 
**CustomMaxMemory** | Pointer to **bool** |  | [optional] 
**CustomMaxStorage** | Pointer to **bool** |  | [optional] 
**CustomMaxDataStorage** | Pointer to **bool** |  | [optional] 
**CustomCoresPerSocket** | Pointer to **bool** |  | [optional] 
**CoresPerSocket** | Pointer to **int32** |  | [optional] 
**StorageTypes** | Pointer to [**[]ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner**](ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner.md) |  | [optional] 
**RootStorageTypes** | Pointer to [**[]ListInstanceServicePlans200ResponsePlansInnerRootStorageTypesInner**](ListInstanceServicePlans200ResponsePlansInnerRootStorageTypesInner.md) |  | [optional] 
**AddVolumes** | Pointer to **bool** |  | [optional] 
**CustomizeVolume** | Pointer to **bool** |  | [optional] 
**RootDiskCustomizable** | Pointer to **bool** |  | [optional] 
**NoDisks** | Pointer to **bool** |  | [optional] 
**HasDatastore** | Pointer to **bool** |  | [optional] 
**MinDisk** | Pointer to **int32** |  | [optional] 
**MaxDisk** | Pointer to **NullableString** |  | [optional] 
**LvmSupported** | Pointer to **bool** |  | [optional] 
**Datastores** | Pointer to [**ListInstanceServicePlans200ResponsePlansInnerDatastores**](ListInstanceServicePlans200ResponsePlansInnerDatastores.md) |  | [optional] 
**SupportsAutoDatastore** | Pointer to **bool** |  | [optional] 
**AutoOptions** | Pointer to [**[]ListInstanceServicePlans200ResponsePlansInnerAutoOptionsInner**](ListInstanceServicePlans200ResponsePlansInnerAutoOptionsInner.md) |  | [optional] 
**CpuOptions** | Pointer to **[]map[string]interface{}** |  | [optional] 
**CoreOptions** | Pointer to **[]map[string]interface{}** |  | [optional] 
**MemoryOptions** | Pointer to **[]map[string]interface{}** |  | [optional] 
**RootCustomSizeOptions** | Pointer to **map[string]interface{}** |  | [optional] 
**CustomSizeOptions** | Pointer to **map[string]interface{}** |  | [optional] 
**CustomCores** | Pointer to **bool** |  | [optional] 
**MaxDisks** | Pointer to **NullableString** |  | [optional] 
**MemorySizeType** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListInstanceServicePlans200ResponsePlansInner{
    // Set fields directly
}
```

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
### CoreOptions (Nullable)

Use the Nullable wrapper methods:
- `obj.CoreOptions.IsSet()` — check if set
- `obj.CoreOptions.Get()` — get the inner value (returns pointer)
- `obj.CoreOptions.Set(&val)` — set the value
- `obj.CoreOptions.Unset()` — clear the value
### MemoryOptions (Nullable)

Use the Nullable wrapper methods:
- `obj.MemoryOptions.IsSet()` — check if set
- `obj.MemoryOptions.Get()` — get the inner value (returns pointer)
- `obj.MemoryOptions.Set(&val)` — set the value
- `obj.MemoryOptions.Unset()` — clear the value
### RootCustomSizeOptions (Nullable)

Use the Nullable wrapper methods:
- `obj.RootCustomSizeOptions.IsSet()` — check if set
- `obj.RootCustomSizeOptions.Get()` — get the inner value (returns pointer)
- `obj.RootCustomSizeOptions.Set(&val)` — set the value
- `obj.RootCustomSizeOptions.Unset()` — clear the value
### CustomSizeOptions (Nullable)

Use the Nullable wrapper methods:
- `obj.CustomSizeOptions.IsSet()` — check if set
- `obj.CustomSizeOptions.Get()` — get the inner value (returns pointer)
- `obj.CustomSizeOptions.Set(&val)` — set the value
- `obj.CustomSizeOptions.Unset()` — clear the value
### MaxDisks (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxDisks.IsSet()` — check if set
- `obj.MaxDisks.Get()` — get the inner value (returns pointer)
- `obj.MaxDisks.Set(&val)` — set the value
- `obj.MaxDisks.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


