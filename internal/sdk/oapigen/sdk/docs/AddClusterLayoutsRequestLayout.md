# AddClusterLayoutsRequestLayout

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Cluster layout name | 
**Description** | Pointer to **NullableString** | Cluster layout description | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**ComputeVersion** | **string** | Version of the cluster layout | 
**Creatable** | Pointer to **bool** | Can be used to enable / disable the creatability of the cluster layout. | [optional] [default to true]
**HasAutoScale** | Pointer to **bool** | Can be used to enable / disable the horizontal scaling. | [optional] [default to false]
**InstallContainerRuntime** | Pointer to **bool** | Install Docker (container runtime). | [optional] [default to false]
**MemoryRequirement** | Pointer to **int64** | Memory requirement in bytes | [optional] 
**GroupType** | [**AddClusterLayoutsRequestLayoutGroupType**](AddClusterLayoutsRequestLayoutGroupType.md) |  | 
**ProvisionType** | [**AddClusterLayoutsRequestLayoutProvisionType**](AddClusterLayoutsRequestLayoutProvisionType.md) |  | 
**OptionTypes** | Pointer to [**[]AddClusterLayoutsRequestLayoutOptionTypesInner**](AddClusterLayoutsRequestLayoutOptionTypesInner.md) | Array of cluster layout option types | [optional] 
**TaskSets** | Pointer to [**[]AddClusterLayoutsRequestLayoutTaskSetsInner**](AddClusterLayoutsRequestLayoutTaskSetsInner.md) | Array of cluster layout task sets | [optional] 
**EnvironmentVariables** | Pointer to [**[]AddClusterLayoutsRequestLayoutEnvironmentVariablesInner**](AddClusterLayoutsRequestLayoutEnvironmentVariablesInner.md) | Array of cluster layout env variables | [optional] 
**Masters** | Pointer to [**[]AddClusterLayoutsRequestLayoutMastersInner**](AddClusterLayoutsRequestLayoutMastersInner.md) | Array of cluster layout master nodes | [optional] 
**Workers** | Pointer to [**[]AddClusterLayoutsRequestLayoutWorkersInner**](AddClusterLayoutsRequestLayoutWorkersInner.md) | Array of cluster layout worker nodes | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddClusterLayoutsRequestLayout{
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


