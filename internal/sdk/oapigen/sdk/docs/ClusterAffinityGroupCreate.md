# ClusterAffinityGroupCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** | Name | [optional] 
**AffinityType** | Pointer to **string** | Affinity Type | [optional] 
**Active** | Pointer to **bool** | Active | [optional] 
**Pool** | Pointer to [**ClusterAffinityGroupCreatePool**](ClusterAffinityGroupCreatePool.md) |  | [optional] 
**Servers** | Pointer to **[]int32** | List of Server IDs to include in the Affinity Group | [optional] 
**Visibility** | Pointer to **string** | Visibility - Set to public to allow all tenants | [optional] [default to "private"]
**Tenants** | Pointer to [**[]ClusterAffinityGroupCreateTenantsInner**](ClusterAffinityGroupCreateTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 
**ResourcePermissions** | Pointer to [**ClusterAffinityGroupCreateResourcePermissions**](ClusterAffinityGroupCreateResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterAffinityGroupCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


