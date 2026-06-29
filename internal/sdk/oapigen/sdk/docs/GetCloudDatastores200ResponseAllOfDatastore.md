# GetCloudDatastores200ResponseAllOfDatastore

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Zone** | Pointer to [**GetCloudDatastores200ResponseAllOfDatastoreZone**](GetCloudDatastores200ResponseAllOfDatastoreZone.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**FreeSpace** | Pointer to **int64** |  | [optional] 
**Online** | Pointer to **bool** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Tenants** | Pointer to [**[]GetCloudDatastores200ResponseAllOfDatastoreTenantsInner**](GetCloudDatastores200ResponseAllOfDatastoreTenantsInner.md) |  | [optional] 
**ResourcePermission** | Pointer to [**GetCloudDatastores200ResponseAllOfDatastoreResourcePermission**](GetCloudDatastores200ResponseAllOfDatastoreResourcePermission.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetCloudDatastores200ResponseAllOfDatastore{
    // Set fields directly
}
```

### Tenants (Nullable)

Use the Nullable wrapper methods:
- `obj.Tenants.IsSet()` — check if set
- `obj.Tenants.Get()` — get the inner value (returns pointer)
- `obj.Tenants.Set(&val)` — set the value
- `obj.Tenants.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


