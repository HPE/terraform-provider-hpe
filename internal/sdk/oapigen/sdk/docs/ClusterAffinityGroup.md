# ClusterAffinityGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** | Name | [optional] 
**AffinityType** | Pointer to **string** | Affinity Type | [optional] 
**Source** | Pointer to **string** | Source | [optional] 
**RefType** | Pointer to **string** | Reference Type for the Affinity Group. Can be ComputeZone or ComputeServerGroup for cloud or cluster respectively | [optional] 
**RefId** | Pointer to **int64** | Reference ID for the Affinity Group. The ID of the Cloud or Clusterfor cloud or cluster respectively | [optional] 
**Active** | Pointer to **bool** | Active | [optional] 
**Pool** | Pointer to [**SaveClusterAffinityGroup200ResponseAllOfAffinityGroupPool**](SaveClusterAffinityGroup200ResponseAllOfAffinityGroupPool.md) |  | [optional] 
**Servers** | Pointer to [**[]SaveClusterAffinityGroup200ResponseAllOfAffinityGroupServersInner**](SaveClusterAffinityGroup200ResponseAllOfAffinityGroupServersInner.md) | List of Servers in the Affinity Group | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Tenants** | Pointer to [**[]SaveClusterAffinityGroup200ResponseAllOfAffinityGroupTenantsInner**](SaveClusterAffinityGroup200ResponseAllOfAffinityGroupTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 
**ResourcePermissions** | Pointer to [**SaveClusterAffinityGroup200ResponseAllOfAffinityGroupResourcePermissions**](SaveClusterAffinityGroup200ResponseAllOfAffinityGroupResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterAffinityGroup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


