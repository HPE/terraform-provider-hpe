# ClusterNamespacePermissionsResourcePermissions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllGroups** | Pointer to **bool** |  | [optional] 
**DefaultStore** | Pointer to **bool** |  | [optional] 
**AllPlans** | Pointer to **bool** |  | [optional] 
**DefaultTarget** | Pointer to **bool** |  | [optional] 
**MorpheusResourceType** | Pointer to **string** |  | [optional] 
**MorpheusResourceId** | Pointer to **int64** |  | [optional] 
**CanManage** | Pointer to **bool** |  | [optional] 
**All** | Pointer to **bool** |  | [optional] 
**Account** | Pointer to [**ClusterNamespacePermissionsResourcePermissionsAccount**](ClusterNamespacePermissionsResourcePermissionsAccount.md) |  | [optional] 
**Sites** | Pointer to [**[]ClusterNamespacePermissionsResourcePermissionsSitesInner**](ClusterNamespacePermissionsResourcePermissionsSitesInner.md) |  | [optional] 
**Plans** | Pointer to [**[]ClusterNamespacePermissionsResourcePermissionsPlansInner**](ClusterNamespacePermissionsResourcePermissionsPlansInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterNamespacePermissionsResourcePermissions{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


