# InstanceTypeInstanceTypeLayoutsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**InstanceType** | Pointer to [**InstanceTypeInstanceTypeLayoutsInnerInstanceType**](InstanceTypeInstanceTypeLayoutsInnerInstanceType.md) |  | [optional] 
**Account** | Pointer to [**InstanceTypeInstanceTypeLayoutsInnerAccount**](InstanceTypeInstanceTypeLayoutsInnerAccount.md) |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**InstanceVersion** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Creatable** | Pointer to **bool** |  | [optional] 
**MemoryRequirement** | Pointer to **NullableInt64** |  | [optional] 
**SortOrder** | Pointer to **int64** |  | [optional] 
**SupportsConvertToManaged** | Pointer to **NullableBool** |  | [optional] 
**ProvisionType** | Pointer to [**InstanceTypeInstanceTypeLayoutsInnerProvisionType**](InstanceTypeInstanceTypeLayoutsInnerProvisionType.md) |  | [optional] 
**TaskSets** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ContainerTypes** | Pointer to [**[]InstanceTypeInstanceTypeLayoutsInnerContainerTypesInner**](InstanceTypeInstanceTypeLayoutsInnerContainerTypesInner.md) |  | [optional] 
**Mounts** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Ports** | Pointer to **[]map[string]interface{}** |  | [optional] 
**OptionTypes** | Pointer to **[]map[string]interface{}** |  | [optional] 
**EnvironmentVariables** | Pointer to **[]map[string]interface{}** |  | [optional] 
**PriceSets** | Pointer to **[]map[string]interface{}** |  | [optional] 
**SpecTemplates** | Pointer to **[]map[string]interface{}** |  | [optional] 
**TfvarSecret** | Pointer to **NullableString** |  | [optional] 
**Permissions** | Pointer to [**InstanceTypeInstanceTypeLayoutsInnerPermissions**](InstanceTypeInstanceTypeLayoutsInnerPermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceTypeInstanceTypeLayoutsInner{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### MemoryRequirement (Nullable)

Use the Nullable wrapper methods:
- `obj.MemoryRequirement.IsSet()` — check if set
- `obj.MemoryRequirement.Get()` — get the inner value (returns pointer)
- `obj.MemoryRequirement.Set(&val)` — set the value
- `obj.MemoryRequirement.Unset()` — clear the value
### SupportsConvertToManaged (Nullable)

Use the Nullable wrapper methods:
- `obj.SupportsConvertToManaged.IsSet()` — check if set
- `obj.SupportsConvertToManaged.Get()` — get the inner value (returns pointer)
- `obj.SupportsConvertToManaged.Set(&val)` — set the value
- `obj.SupportsConvertToManaged.Unset()` — clear the value
### TaskSets (Nullable)

Use the Nullable wrapper methods:
- `obj.TaskSets.IsSet()` — check if set
- `obj.TaskSets.Get()` — get the inner value (returns pointer)
- `obj.TaskSets.Set(&val)` — set the value
- `obj.TaskSets.Unset()` — clear the value
### Mounts (Nullable)

Use the Nullable wrapper methods:
- `obj.Mounts.IsSet()` — check if set
- `obj.Mounts.Get()` — get the inner value (returns pointer)
- `obj.Mounts.Set(&val)` — set the value
- `obj.Mounts.Unset()` — clear the value
### Ports (Nullable)

Use the Nullable wrapper methods:
- `obj.Ports.IsSet()` — check if set
- `obj.Ports.Get()` — get the inner value (returns pointer)
- `obj.Ports.Set(&val)` — set the value
- `obj.Ports.Unset()` — clear the value
### OptionTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.OptionTypes.IsSet()` — check if set
- `obj.OptionTypes.Get()` — get the inner value (returns pointer)
- `obj.OptionTypes.Set(&val)` — set the value
- `obj.OptionTypes.Unset()` — clear the value
### EnvironmentVariables (Nullable)

Use the Nullable wrapper methods:
- `obj.EnvironmentVariables.IsSet()` — check if set
- `obj.EnvironmentVariables.Get()` — get the inner value (returns pointer)
- `obj.EnvironmentVariables.Set(&val)` — set the value
- `obj.EnvironmentVariables.Unset()` — clear the value
### PriceSets (Nullable)

Use the Nullable wrapper methods:
- `obj.PriceSets.IsSet()` — check if set
- `obj.PriceSets.Get()` — get the inner value (returns pointer)
- `obj.PriceSets.Set(&val)` — set the value
- `obj.PriceSets.Unset()` — clear the value
### SpecTemplates (Nullable)

Use the Nullable wrapper methods:
- `obj.SpecTemplates.IsSet()` — check if set
- `obj.SpecTemplates.Get()` — get the inner value (returns pointer)
- `obj.SpecTemplates.Set(&val)` — set the value
- `obj.SpecTemplates.Unset()` — clear the value
### TfvarSecret (Nullable)

Use the Nullable wrapper methods:
- `obj.TfvarSecret.IsSet()` — check if set
- `obj.TfvarSecret.Get()` — get the inner value (returns pointer)
- `obj.TfvarSecret.Set(&val)` — set the value
- `obj.TfvarSecret.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


