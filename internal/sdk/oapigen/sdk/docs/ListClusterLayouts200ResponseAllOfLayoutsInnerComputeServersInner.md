# ListClusterLayouts200ResponseAllOfLayoutsInnerComputeServersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**PriorityOrder** | Pointer to **int64** |  | [optional] 
**NodeCount** | Pointer to **int64** |  | [optional] 
**NodeType** | Pointer to **string** |  | [optional] 
**MinNodeCount** | Pointer to **int64** |  | [optional] 
**MaxNodeCount** | Pointer to **NullableString** |  | [optional] 
**DynamicCount** | Pointer to **bool** |  | [optional] 
**InstallContainerRuntime** | Pointer to **bool** |  | [optional] 
**InstallStorageRuntime** | Pointer to **bool** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**Config** | Pointer to **NullableString** |  | [optional] 
**ContainerType** | Pointer to [**ListClusterLayouts200ResponseAllOfLayoutsInnerComputeServersInnerContainerType**](ListClusterLayouts200ResponseAllOfLayoutsInnerComputeServersInnerContainerType.md) |  | [optional] 
**ComputeServerType** | Pointer to [**ListClusterLayouts200ResponseAllOfLayoutsInnerComputeServersInnerComputeServerType**](ListClusterLayouts200ResponseAllOfLayoutsInnerComputeServersInnerComputeServerType.md) |  | [optional] 
**ProvisionService** | Pointer to **NullableString** |  | [optional] 
**PlanCategory** | Pointer to **NullableString** |  | [optional] 
**NamePrefix** | Pointer to **NullableString** |  | [optional] 
**NameSuffix** | Pointer to **NullableString** |  | [optional] 
**ForceNameIndex** | Pointer to **bool** |  | [optional] 
**LoadBalance** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListClusterLayouts200ResponseAllOfLayoutsInnerComputeServersInner{
    // Set fields directly
}
```

### MaxNodeCount (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxNodeCount.IsSet()` — check if set
- `obj.MaxNodeCount.Get()` — get the inner value (returns pointer)
- `obj.MaxNodeCount.Set(&val)` — set the value
- `obj.MaxNodeCount.Unset()` — clear the value
### Config (Nullable)

Use the Nullable wrapper methods:
- `obj.Config.IsSet()` — check if set
- `obj.Config.Get()` — get the inner value (returns pointer)
- `obj.Config.Set(&val)` — set the value
- `obj.Config.Unset()` — clear the value
### ProvisionService (Nullable)

Use the Nullable wrapper methods:
- `obj.ProvisionService.IsSet()` — check if set
- `obj.ProvisionService.Get()` — get the inner value (returns pointer)
- `obj.ProvisionService.Set(&val)` — set the value
- `obj.ProvisionService.Unset()` — clear the value
### PlanCategory (Nullable)

Use the Nullable wrapper methods:
- `obj.PlanCategory.IsSet()` — check if set
- `obj.PlanCategory.Get()` — get the inner value (returns pointer)
- `obj.PlanCategory.Set(&val)` — set the value
- `obj.PlanCategory.Unset()` — clear the value
### NamePrefix (Nullable)

Use the Nullable wrapper methods:
- `obj.NamePrefix.IsSet()` — check if set
- `obj.NamePrefix.Get()` — get the inner value (returns pointer)
- `obj.NamePrefix.Set(&val)` — set the value
- `obj.NamePrefix.Unset()` — clear the value
### NameSuffix (Nullable)

Use the Nullable wrapper methods:
- `obj.NameSuffix.IsSet()` — check if set
- `obj.NameSuffix.Get()` — get the inner value (returns pointer)
- `obj.NameSuffix.Set(&val)` — set the value
- `obj.NameSuffix.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


