# SaveDatastoreRequestDatastore

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the datastore to be created. | 
**DatastoreType** | **string** | The code of the datastoreType | 
**Config** | [**SaveDatastoreRequestDatastoreConfig**](SaveDatastoreRequestDatastoreConfig.md) |  | 
**RefType** | **string** | Type of the resource this datastore is associated with, can be &#39;ComputeZone&#39; (&#39;Cloud&#39;) or &#39;ComputeServerGroup&#39; (&#39;Cluster&#39;) | 
**RefId** | **int64** | The ID of the resource this datastore is associated with, e.g. ComputeZone, ComputeServerGroup | 
**StorageServer** | Pointer to [**SaveDatastoreRequestDatastoreStorageServer**](SaveDatastoreRequestDatastoreStorageServer.md) |  | [optional] 
**Visibility** | Pointer to **string** | Visibility level of the datastore, can be &#39;private&#39; or &#39;public&#39;. If not specified, defaults to &#39;private&#39;. | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**DefaultStore** | Pointer to **bool** |  | [optional] 
**TenantPermissions** | Pointer to [**SaveDatastoreRequestDatastoreTenantPermissions**](SaveDatastoreRequestDatastoreTenantPermissions.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**SaveDatastoreRequestDatastoreResourcePermissions**](SaveDatastoreRequestDatastoreResourcePermissions.md) |  | [optional] 
**Datastores** | Pointer to **[]map[string]interface{}** | List of datastores associated with this datastore, for use with vSphere clouds. | [optional] 
**HeartbeatTarget** | Pointer to **bool** | Heartbeat Target | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SaveDatastoreRequestDatastore{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


