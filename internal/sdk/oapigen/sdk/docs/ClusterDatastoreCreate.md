# ClusterDatastoreCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**DatastoreType** | Pointer to [**ClusterDatastoreCreateDatastoreType**](ClusterDatastoreCreateDatastoreType.md) |  | [optional] 
**StorageServer** | Pointer to [**ClusterDatastoreCreateStorageServer**](ClusterDatastoreCreateStorageServer.md) |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**DefaultStore** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to [**ClusterDatastoreCreateConfig**](ClusterDatastoreCreateConfig.md) |  | [optional] 
**Tenants** | Pointer to [**[]ClusterDatastoreCreateTenantsInner**](ClusterDatastoreCreateTenantsInner.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**ClusterDatastoreCreateResourcePermissions**](ClusterDatastoreCreateResourcePermissions.md) |  | [optional] 
**Datastores** | Pointer to **[]map[string]interface{}** |  | [optional] 
**HeartbeatTarget** | Pointer to **bool** | Heartbeat Target | [optional] 
**SupportsVmSecureMetadata** | Pointer to **bool** | When &#x60;true&#x60;, designates this datastore to hold NVRAM and swtpm state for TPM/SecureBoot VMs in the cluster, enabling live migration and HA failover. Only one datastore per cluster scope may be designated at a time; setting this to &#x60;true&#x60; automatically clears the flag on any previously designated datastore. Only applicable to GFS2 and NFS datastore types.  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterDatastoreCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


