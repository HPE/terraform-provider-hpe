# ListServicePlans200ResponseAllOfServicePlansInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**SortOrder** | Pointer to **int64** |  | [optional] 
**RegionCode** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**MaxCpu** | Pointer to **NullableInt64** |  | [optional] 
**MaxGpus** | Pointer to **NullableInt64** |  | [optional] 
**MaxCores** | Pointer to **NullableInt64** |  | [optional] 
**MaxDisks** | Pointer to **NullableInt64** |  | [optional] 
**CoresPerSocket** | Pointer to **NullableInt64** |  | [optional] 
**CustomCpu** | Pointer to **bool** |  | [optional] 
**CustomCores** | Pointer to **bool** |  | [optional] 
**CustomMaxStorage** | Pointer to **NullableBool** |  | [optional] 
**CustomMaxDataStorage** | Pointer to **NullableBool** |  | [optional] 
**CustomMaxMemory** | Pointer to **NullableBool** |  | [optional] 
**AddVolumes** | Pointer to **NullableBool** |  | [optional] 
**MemoryOptionSource** | Pointer to **NullableString** |  | [optional] 
**CpuOptionSource** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Editable** | Pointer to **bool** |  | [optional] 
**ProvisionType** | Pointer to [**ListServicePlans200ResponseAllOfServicePlansInnerProvisionType**](ListServicePlans200ResponseAllOfServicePlansInnerProvisionType.md) |  | [optional] 
**Tenants** | Pointer to **string** |  | [optional] 
**PriceSets** | Pointer to [**[]ListServicePlans200ResponseAllOfServicePlansInnerPriceSetsInner**](ListServicePlans200ResponseAllOfServicePlansInnerPriceSetsInner.md) |  | [optional] 
**Config** | Pointer to [**ListServicePlans200ResponseAllOfServicePlansInnerConfig**](ListServicePlans200ResponseAllOfServicePlansInnerConfig.md) |  | [optional] 
**Zones** | Pointer to [**[]ListServicePlans200ResponseAllOfServicePlansInnerZonesInner**](ListServicePlans200ResponseAllOfServicePlansInnerZonesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListServicePlans200ResponseAllOfServicePlansInner{
    // Set fields directly
}
```

### RegionCode (Nullable)

Use the Nullable wrapper methods:
- `obj.RegionCode.IsSet()` — check if set
- `obj.RegionCode.Get()` — get the inner value (returns pointer)
- `obj.RegionCode.Set(&val)` — set the value
- `obj.RegionCode.Unset()` — clear the value
### MaxCpu (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxCpu.IsSet()` — check if set
- `obj.MaxCpu.Get()` — get the inner value (returns pointer)
- `obj.MaxCpu.Set(&val)` — set the value
- `obj.MaxCpu.Unset()` — clear the value
### MaxGpus (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxGpus.IsSet()` — check if set
- `obj.MaxGpus.Get()` — get the inner value (returns pointer)
- `obj.MaxGpus.Set(&val)` — set the value
- `obj.MaxGpus.Unset()` — clear the value
### MaxCores (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxCores.IsSet()` — check if set
- `obj.MaxCores.Get()` — get the inner value (returns pointer)
- `obj.MaxCores.Set(&val)` — set the value
- `obj.MaxCores.Unset()` — clear the value
### MaxDisks (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxDisks.IsSet()` — check if set
- `obj.MaxDisks.Get()` — get the inner value (returns pointer)
- `obj.MaxDisks.Set(&val)` — set the value
- `obj.MaxDisks.Unset()` — clear the value
### CoresPerSocket (Nullable)

Use the Nullable wrapper methods:
- `obj.CoresPerSocket.IsSet()` — check if set
- `obj.CoresPerSocket.Get()` — get the inner value (returns pointer)
- `obj.CoresPerSocket.Set(&val)` — set the value
- `obj.CoresPerSocket.Unset()` — clear the value
### CustomMaxStorage (Nullable)

Use the Nullable wrapper methods:
- `obj.CustomMaxStorage.IsSet()` — check if set
- `obj.CustomMaxStorage.Get()` — get the inner value (returns pointer)
- `obj.CustomMaxStorage.Set(&val)` — set the value
- `obj.CustomMaxStorage.Unset()` — clear the value
### CustomMaxDataStorage (Nullable)

Use the Nullable wrapper methods:
- `obj.CustomMaxDataStorage.IsSet()` — check if set
- `obj.CustomMaxDataStorage.Get()` — get the inner value (returns pointer)
- `obj.CustomMaxDataStorage.Set(&val)` — set the value
- `obj.CustomMaxDataStorage.Unset()` — clear the value
### CustomMaxMemory (Nullable)

Use the Nullable wrapper methods:
- `obj.CustomMaxMemory.IsSet()` — check if set
- `obj.CustomMaxMemory.Get()` — get the inner value (returns pointer)
- `obj.CustomMaxMemory.Set(&val)` — set the value
- `obj.CustomMaxMemory.Unset()` — clear the value
### AddVolumes (Nullable)

Use the Nullable wrapper methods:
- `obj.AddVolumes.IsSet()` — check if set
- `obj.AddVolumes.Get()` — get the inner value (returns pointer)
- `obj.AddVolumes.Set(&val)` — set the value
- `obj.AddVolumes.Unset()` — clear the value
### MemoryOptionSource (Nullable)

Use the Nullable wrapper methods:
- `obj.MemoryOptionSource.IsSet()` — check if set
- `obj.MemoryOptionSource.Get()` — get the inner value (returns pointer)
- `obj.MemoryOptionSource.Set(&val)` — set the value
- `obj.MemoryOptionSource.Unset()` — clear the value
### CpuOptionSource (Nullable)

Use the Nullable wrapper methods:
- `obj.CpuOptionSource.IsSet()` — check if set
- `obj.CpuOptionSource.Get()` — get the inner value (returns pointer)
- `obj.CpuOptionSource.Set(&val)` — set the value
- `obj.CpuOptionSource.Unset()` — clear the value
### PriceSets (Nullable)

Use the Nullable wrapper methods:
- `obj.PriceSets.IsSet()` — check if set
- `obj.PriceSets.Get()` — get the inner value (returns pointer)
- `obj.PriceSets.Set(&val)` — set the value
- `obj.PriceSets.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


