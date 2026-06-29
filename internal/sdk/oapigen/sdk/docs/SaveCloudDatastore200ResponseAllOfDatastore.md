# SaveCloudDatastore200ResponseAllOfDatastore

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**DatastoreType** | Pointer to [**SaveCloudDatastore200ResponseAllOfDatastoreDatastoreType**](SaveCloudDatastore200ResponseAllOfDatastoreDatastoreType.md) |  | [optional] 
**StorageServer** | Pointer to [**SaveCloudDatastore200ResponseAllOfDatastoreStorageServer**](SaveCloudDatastore200ResponseAllOfDatastoreStorageServer.md) |  | [optional] 
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
**HeartbeatTarget** | Pointer to **bool** | Heartbeat Target | [optional] 
**SupportsVmSecureMetadata** | Pointer to **bool** | When &#x60;true&#x60;, this datastore is designated to hold NVRAM and swtpm state for TPM/SecureBoot VMs in the cluster, enabling live migration and HA failover. Only one datastore per cluster scope may have this set to &#x60;true&#x60;. Only applicable to GFS2 and NFS datastore types.  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**Zone** | Pointer to [**SaveCloudDatastore200ResponseAllOfDatastoreZone**](SaveCloudDatastore200ResponseAllOfDatastoreZone.md) |  | [optional] 
**ZonePool** | Pointer to [**SaveCloudDatastore200ResponseAllOfDatastoreZonePool**](SaveCloudDatastore200ResponseAllOfDatastoreZonePool.md) |  | [optional] 
**Owner** | Pointer to [**SaveCloudDatastore200ResponseAllOfDatastoreOwner**](SaveCloudDatastore200ResponseAllOfDatastoreOwner.md) |  | [optional] 
**Tenants** | Pointer to [**[]SaveCloudDatastore200ResponseAllOfDatastoreTenantsInner**](SaveCloudDatastore200ResponseAllOfDatastoreTenantsInner.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**SaveCloudDatastore200ResponseAllOfDatastoreResourcePermissions**](SaveCloudDatastore200ResponseAllOfDatastoreResourcePermissions.md) |  | [optional] 
**Datastores** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SaveCloudDatastore200ResponseAllOfDatastore{
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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


