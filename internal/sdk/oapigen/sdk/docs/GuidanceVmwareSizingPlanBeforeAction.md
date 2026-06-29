# GuidanceVmwareSizingPlanBeforeAction

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**SortOrder** | Pointer to **int64** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**MaxCpu** | Pointer to **NullableString** |  | [optional] 
**MaxCores** | Pointer to **int64** |  | [optional] 
**MaxDisks** | Pointer to **NullableString** |  | [optional] 
**CoresPerSocket** | Pointer to **int64** |  | [optional] 
**CustomCpu** | Pointer to **bool** |  | [optional] 
**CustomCores** | Pointer to **bool** |  | [optional] 
**CustomMaxStorage** | Pointer to **bool** |  | [optional] 
**CustomMaxDataStorage** | Pointer to **bool** |  | [optional] 
**CustomMaxMemory** | Pointer to **bool** |  | [optional] 
**AddVolumes** | Pointer to **bool** |  | [optional] 
**MemoryOptionSource** | Pointer to **NullableString** |  | [optional] 
**CpuOptionSource** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**RegionCode** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Editable** | Pointer to **bool** |  | [optional] 
**ProvisionType** | Pointer to [**GuidanceVmwareSizingPlanBeforeActionProvisionType**](GuidanceVmwareSizingPlanBeforeActionProvisionType.md) |  | [optional] 
**Tenants** | Pointer to **string** |  | [optional] 
**PriceSets** | Pointer to [**[]GuidanceVmwareSizingPlanBeforeActionPriceSetsInner**](GuidanceVmwareSizingPlanBeforeActionPriceSetsInner.md) |  | [optional] 
**Config** | Pointer to [**GuidanceVmwareSizingPlanBeforeActionConfig**](GuidanceVmwareSizingPlanBeforeActionConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GuidanceVmwareSizingPlanBeforeAction{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### MaxCpu (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxCpu.IsSet()` — check if set
- `obj.MaxCpu.Get()` — get the inner value (returns pointer)
- `obj.MaxCpu.Set(&val)` — set the value
- `obj.MaxCpu.Unset()` — clear the value
### MaxDisks (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxDisks.IsSet()` — check if set
- `obj.MaxDisks.Get()` — get the inner value (returns pointer)
- `obj.MaxDisks.Set(&val)` — set the value
- `obj.MaxDisks.Unset()` — clear the value
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
### RegionCode (Nullable)

Use the Nullable wrapper methods:
- `obj.RegionCode.IsSet()` — check if set
- `obj.RegionCode.Get()` — get the inner value (returns pointer)
- `obj.RegionCode.Set(&val)` — set the value
- `obj.RegionCode.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


