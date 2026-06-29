# UpdateCloudFoldersRequestFolder

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DefaultFolder** | Pointer to **bool** |  | [optional] [default to false]
**DefaultImage** | Pointer to **bool** |  | [optional] [default to false]
**Active** | Pointer to **bool** | Activate &#x60;true&#x60; or disable &#x60;false&#x60; the folder | [optional] 
**Visibility** | Pointer to **string** | Setting &#x60;private&#x60; or &#x60;public&#x60; | [optional] [default to "private"]
**TenantPermissions** | Pointer to [**UpdateCloudFoldersRequestFolderTenantPermissions**](UpdateCloudFoldersRequestFolderTenantPermissions.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**UpdateCloudFoldersRequestFolderResourcePermissions**](UpdateCloudFoldersRequestFolderResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateCloudFoldersRequestFolder{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


