# UpdateBlueprintPermissionsRequestResourcePermission

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**All** | Pointer to **bool** | Set to true to grant access to all groups | [optional] 
**Sites** | Pointer to [**[]UpdateBlueprintPermissionsRequestResourcePermissionSitesInner**](UpdateBlueprintPermissionsRequestResourcePermissionSitesInner.md) | Array of objects identifying groups with access | [optional] 
**OwnerId** | Pointer to **int64** | User ID, can be used to change blueprint owner. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateBlueprintPermissionsRequestResourcePermission{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


