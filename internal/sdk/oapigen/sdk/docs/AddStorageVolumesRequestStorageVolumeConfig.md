# AddStorageVolumesRequestStorageVolumeConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**HpeStorageDatastore** | **int64** | ID of the Alletra MP BMaaS data store (pool) in which to create the volume. | 
**HpeStorageVolumeShared** | Pointer to **string** | Whether the volume is shared (multi-attach). When &#39;off&#39; the volume may be exported to a single compute server; when &#39;on&#39; it may be exported to one or more instances.  | [optional] 
**HpeStorageComputeServer** | Pointer to **string** | Compute server to export a non-shared volume to, in the form \&quot;[id: &lt;computeServerId&gt;]\&quot;. Mutually exclusive with hpe_storage_instances.  | [optional] 
**HpeStorageInstances** | Pointer to **[]string** | Instances to export a shared volume to, each in the form \&quot;[id: &lt;instanceId&gt;]\&quot;. Mutually exclusive with hpe_storage_compute_server.  | [optional] 
**HpeStorageRemotecopytargetId** | Pointer to **string** | Remote copy (replication) target ID. Required for replicated LUN volume types (Peer Persistence: Active or Classic).  | [optional] 
**HpeStorageExistingVolumeSet** | Pointer to **string** | Whether to add the volume to an existing, exported volume set rather than creating a new one. | [optional] 
**HpeStorageVolumesetId** | Pointer to **string** | ID of an existing volume set to add the volume to (when using an existing volume set). | [optional] 
**HpeStorageVolumeSetName** | Pointer to **string** | Base name for a new volume set (a unique suffix is always appended). Used when not using an existing volume set.  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddStorageVolumesRequestStorageVolumeConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


