# GetClusterDatastore200ResponseDatastoreResourcePermissions

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
**Account** | Pointer to [**GetClusterDatastore200ResponseDatastoreResourcePermissionsAccount**](GetClusterDatastore200ResponseDatastoreResourcePermissionsAccount.md) |  | [optional] 
**Sites** | Pointer to [**[]GetClusterDatastore200ResponseDatastoreResourcePermissionsSitesInner**](GetClusterDatastore200ResponseDatastoreResourcePermissionsSitesInner.md) |  | [optional] 
**Plans** | Pointer to [**[]GetClusterDatastore200ResponseDatastoreResourcePermissionsPlansInner**](GetClusterDatastore200ResponseDatastoreResourcePermissionsPlansInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetClusterDatastore200ResponseDatastoreResourcePermissions{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


