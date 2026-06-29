# AddCloudResourcePoolRequestResourcePoolResourcePermissions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**All** | Pointer to **bool** | Pass &#x60;true&#x60; to allow access all groups | [optional] [default to true]
**Sites** | Pointer to [**[]UpdateCloudFoldersRequestFolderResourcePermissionsSitesInner**](UpdateCloudFoldersRequestFolderResourcePermissionsSitesInner.md) | Array of groups that are allowed access | [optional] 
**AllPlans** | Pointer to **bool** | Pass true to allow access all plans | [optional] [default to true]
**Plans** | Pointer to [**[]UpdateCloudFoldersRequestFolderResourcePermissionsPlansInner**](UpdateCloudFoldersRequestFolderResourcePermissionsPlansInner.md) | Array of plans that are allowed access | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCloudResourcePoolRequestResourcePoolResourcePermissions{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


