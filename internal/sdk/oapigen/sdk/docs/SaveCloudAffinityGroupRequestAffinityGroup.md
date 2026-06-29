# SaveCloudAffinityGroupRequestAffinityGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**AffinityType** | Pointer to **string** | Affinity Type | [optional] 
**Active** | Pointer to **bool** | Active | [optional] 
**Pool** | Pointer to [**SaveCloudAffinityGroupRequestAffinityGroupPool**](SaveCloudAffinityGroupRequestAffinityGroupPool.md) |  | [optional] 
**Servers** | Pointer to **[]int32** | List of Server IDs to include in the Affinity Group | [optional] 
**Visibility** | Pointer to **string** | Visibility - Set to public to allow all tenants | [optional] [default to "private"]
**Tenants** | Pointer to [**[]SaveCloudAffinityGroupRequestAffinityGroupTenantsInner**](SaveCloudAffinityGroupRequestAffinityGroupTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 
**ResourcePermissions** | Pointer to [**SaveCloudAffinityGroupRequestAffinityGroupResourcePermissions**](SaveCloudAffinityGroupRequestAffinityGroupResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SaveCloudAffinityGroupRequestAffinityGroup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


