# ListClusterDatastores200ResponseAllOfDatastoresInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**StorageSize** | Pointer to **NullableInt64** |  | [optional] 
**FreeSpace** | Pointer to **NullableInt64** |  | [optional] 
**DrsEnabled** | Pointer to **bool** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**AllowWrite** | Pointer to **bool** |  | [optional] 
**DefaultStore** | Pointer to **bool** |  | [optional] 
**Online** | Pointer to **bool** |  | [optional] 
**AllowRead** | Pointer to **bool** |  | [optional] 
**AllowProvision** | Pointer to **bool** |  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **string** |  | [optional] 
**Zone** | Pointer to [**ListClusterDatastores200ResponseAllOfDatastoresInnerZone**](ListClusterDatastores200ResponseAllOfDatastoresInnerZone.md) |  | [optional] 
**ZonePool** | Pointer to [**ListClusterDatastores200ResponseAllOfDatastoresInnerZonePool**](ListClusterDatastores200ResponseAllOfDatastoresInnerZonePool.md) |  | [optional] 
**Owner** | Pointer to [**ListClusterDatastores200ResponseAllOfDatastoresInnerOwner**](ListClusterDatastores200ResponseAllOfDatastoresInnerOwner.md) |  | [optional] 
**Tenants** | Pointer to [**[]ListClusterDatastores200ResponseAllOfDatastoresInnerTenantsInner**](ListClusterDatastores200ResponseAllOfDatastoresInnerTenantsInner.md) |  | [optional] 
**Datastores** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Locations** | Pointer to [**[]ListClusterDatastores200ResponseAllOfDatastoresInnerLocationsInner**](ListClusterDatastores200ResponseAllOfDatastoresInnerLocationsInner.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**ListClusterDatastores200ResponseAllOfDatastoresInnerResourcePermissions**](ListClusterDatastores200ResponseAllOfDatastoresInnerResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListClusterDatastores200ResponseAllOfDatastoresInner{
    // Set fields directly
}
```

### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
### StorageSize (Nullable)

Use the Nullable wrapper methods:
- `obj.StorageSize.IsSet()` — check if set
- `obj.StorageSize.Get()` — get the inner value (returns pointer)
- `obj.StorageSize.Set(&val)` — set the value
- `obj.StorageSize.Unset()` — clear the value
### FreeSpace (Nullable)

Use the Nullable wrapper methods:
- `obj.FreeSpace.IsSet()` — check if set
- `obj.FreeSpace.Get()` — get the inner value (returns pointer)
- `obj.FreeSpace.Set(&val)` — set the value
- `obj.FreeSpace.Unset()` — clear the value
### Datastores (Nullable)

Use the Nullable wrapper methods:
- `obj.Datastores.IsSet()` — check if set
- `obj.Datastores.Get()` — get the inner value (returns pointer)
- `obj.Datastores.Set(&val)` — set the value
- `obj.Datastores.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


