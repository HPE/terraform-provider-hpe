# AddBackupsRequestBackup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LocationType** | **string** |  | 
**Name** | **string** | A name for the backup | 
**InstanceId** | **int64** | The ID of the instance to backup | 
**ContainerId** | **int64** | The ID of the container (workload) to backup. Can be obtained from the instance&#39;s containers list via &#x60;GET /api/instances/{instanceId}&#x60;. | 
**BackupType** | **string** | The backup type code, options vary by the type of cloud and source | 
**JobAction** | **string** | Create a new backup job, clone an existing job or add the new backup to an existing job | 
**JobId** | Pointer to **int64** | The ID of the job to clone or add to. Only applies to jobAction &#x60;clone&#x60; and &#x60;addTo&#x60;. | [optional] 
**JobName** | Pointer to **string** | Name for new job. Only applies to jobAction &#x60;new&#x60; and &#x60;clone&#x60;. | [optional] 
**JobSchedule** | Pointer to **int64** | The ID of the execute schedule for new job. See Execute Schedules. Only applies to jobAction &#x60;new&#x60; and &#x60;clone&#x60;. | [optional] 
**RetentionCount** | Pointer to **int64** | Retention Count for new job. By default the backup settings value will be used. Only applies to jobAction &#x60;new&#x60; and &#x60;clone&#x60;. | [optional] 
**Target** | Pointer to **int64** | The ID of the storage bucket to store the backup in. | [optional] 
**BackupRepository** | Pointer to **int64** | The ID of the backup repository (for provider-specific repositories such as Veeam or Commvault). | [optional] 
**ProviderBackupType** | Pointer to **int64** | The ID of the backup type when using an external backup provider (e.g. Veeam, Commvault). Overrides &#x60;backupType&#x60;. | [optional] 
**BackupJob** | Pointer to [**BackupInstanceBackupJob**](BackupInstanceBackupJob.md) |  | [optional] 
**ServerId** | **int64** | The ID of the server to backup | 
**TargetPath** | Pointer to **string** | The file or directory path on the target host to back up. | [optional] 
**SshUsername** | Pointer to **string** | SSH username for connecting to the server. | [optional] 
**SshPassword** | Pointer to **string** | SSH password for connecting to the server. | [optional] 
**SourceProviderId** | Pointer to **int64** | Source Object Store. The ID of the storage provider (bucket) to be backed up. | [optional] 
**StorageProviderId** | Pointer to **int64** | Target Object Store. The ID of the storage provider (bucket) the backup will be copied to. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddBackupsRequestBackup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


