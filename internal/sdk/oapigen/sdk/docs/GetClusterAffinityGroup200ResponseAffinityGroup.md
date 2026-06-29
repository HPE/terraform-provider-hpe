# GetClusterAffinityGroup200ResponseAffinityGroup

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
**Pool** | Pointer to [**GetClusterAffinityGroup200ResponseAffinityGroupPool**](GetClusterAffinityGroup200ResponseAffinityGroupPool.md) |  | [optional] 
**Servers** | Pointer to [**[]GetClusterAffinityGroup200ResponseAffinityGroupServersInner**](GetClusterAffinityGroup200ResponseAffinityGroupServersInner.md) | List of Servers in the Affinity Group | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Tenants** | Pointer to [**[]GetClusterAffinityGroup200ResponseAffinityGroupTenantsInner**](GetClusterAffinityGroup200ResponseAffinityGroupTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 
**ResourcePermissions** | Pointer to [**GetClusterAffinityGroup200ResponseAffinityGroupResourcePermissions**](GetClusterAffinityGroup200ResponseAffinityGroupResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetClusterAffinityGroup200ResponseAffinityGroup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


