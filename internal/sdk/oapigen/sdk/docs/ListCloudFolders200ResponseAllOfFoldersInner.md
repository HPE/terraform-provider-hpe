# ListCloudFolders200ResponseAllOfFoldersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Zone** | Pointer to [**NullableListCloudFolders200ResponseAllOfFoldersInnerZone**](ListCloudFolders200ResponseAllOfFoldersInnerZone.md) |  | [optional] 
**Parent** | Pointer to [**ListCloudFolders200ResponseAllOfFoldersInnerParent**](ListCloudFolders200ResponseAllOfFoldersInnerParent.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**ReadOnly** | Pointer to **bool** |  | [optional] 
**DefaultFolder** | Pointer to **bool** |  | [optional] 
**DefaultStore** | Pointer to **bool** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Tenants** | Pointer to [**[]ListCloudFolders200ResponseAllOfFoldersInnerTenantsInner**](ListCloudFolders200ResponseAllOfFoldersInnerTenantsInner.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**ListCloudFolders200ResponseAllOfFoldersInnerResourcePermissions**](ListCloudFolders200ResponseAllOfFoldersInnerResourcePermissions.md) |  | [optional] 
**Depth** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListCloudFolders200ResponseAllOfFoldersInner{
    // Set fields directly
}
```

### Zone (Nullable)

Use the Nullable wrapper methods:
- `obj.Zone.IsSet()` — check if set
- `obj.Zone.Get()` — get the inner value (returns pointer)
- `obj.Zone.Set(&val)` — set the value
- `obj.Zone.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


