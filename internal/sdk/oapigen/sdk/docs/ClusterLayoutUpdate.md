# ClusterLayoutUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Cluster layout name | [optional] 
**Description** | Pointer to **NullableString** | Cluster layout description | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**ComputeVersion** | Pointer to **string** | Version of the cluster layout | [optional] 
**Creatable** | Pointer to **bool** | Can be used to enable / disable the creatability of the cluster layout. | [optional] [default to true]
**HasAutoScale** | Pointer to **bool** | Can be used to enable / disable the horizontal scaling. | [optional] [default to false]
**InstallContainerRuntime** | Pointer to **bool** | Install Docker (container runtime). | [optional] [default to false]
**MemoryRequirement** | Pointer to **int64** | Memory requirement in bytes | [optional] 
**GroupType** | Pointer to [**ClusterLayoutUpdateGroupType**](ClusterLayoutUpdateGroupType.md) |  | [optional] 
**ProvisionType** | Pointer to [**ClusterLayoutUpdateProvisionType**](ClusterLayoutUpdateProvisionType.md) |  | [optional] 
**OptionTypes** | Pointer to [**[]ClusterLayoutUpdateOptionTypesInner**](ClusterLayoutUpdateOptionTypesInner.md) | Array of cluster layout option types | [optional] 
**TaskSets** | Pointer to [**[]ClusterLayoutUpdateTaskSetsInner**](ClusterLayoutUpdateTaskSetsInner.md) | Array of cluster layout task sets | [optional] 
**EnvironmentVariables** | Pointer to [**[]ClusterLayoutUpdateEnvironmentVariablesInner**](ClusterLayoutUpdateEnvironmentVariablesInner.md) | Array of cluster layout env variables | [optional] 
**Masters** | Pointer to [**[]ClusterLayoutUpdateMastersInner**](ClusterLayoutUpdateMastersInner.md) | Array of cluster layout master nodes | [optional] 
**Workers** | Pointer to [**[]ClusterLayoutUpdateWorkersInner**](ClusterLayoutUpdateWorkersInner.md) | Array of cluster layout worker nodes | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterLayoutUpdate{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


