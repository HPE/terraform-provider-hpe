# UpdateClusterDatastoreRequestDatastorePermissionsResourcePermissions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**All** | Pointer to **bool** | Pass true to allow access to all groups | [optional] 
**Sites** | Pointer to [**[]UpdateClusterDatastoreRequestDatastorePermissionsResourcePermissionsSitesInner**](UpdateClusterDatastoreRequestDatastorePermissionsResourcePermissionsSitesInner.md) | Array of groups that are allowed access | [optional] 
**AllPlans** | Pointer to **bool** | Pass true to allow access to all plans | [optional] 
**Plans** | Pointer to [**[]UpdateClusterDatastoreRequestDatastorePermissionsResourcePermissionsPlansInner**](UpdateClusterDatastoreRequestDatastorePermissionsResourcePermissionsPlansInner.md) | Array of plans that are allowed access | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateClusterDatastoreRequestDatastorePermissionsResourcePermissions{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


