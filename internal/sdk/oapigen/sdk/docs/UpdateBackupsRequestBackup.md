# UpdateBackupsRequestBackup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the backup | [optional] 
**Enabled** | Pointer to **bool** | Can be used to enable or disable the backup | [optional] 
**StorageProviderId** | Pointer to **NullableInt64** | The ID of the storage bucket. Set to null to clear. | [optional] 
**BackupJobId** | Pointer to **int64** | The Backup Job ID to assign the backup to. This determines when the backup is run. | [optional] 
**CronExpression** | Pointer to **string** | Cron expression for custom schedule (overrides job schedule). | [optional] 
**TargetPath** | Pointer to **string** | The file or directory path on the target host. Applies to server/host backups (fileBackup, directoryBackup). | [optional] 
**TargetHost** | Pointer to **string** | Target host address. Applies to server/host and database backups. | [optional] 
**TargetPort** | Pointer to **int32** | Target port. Applies to server/host and database backups. | [optional] 
**TargetUsername** | Pointer to **string** | Target username for connecting to the backup target. | [optional] 
**TargetPassword** | Pointer to **string** | Target password for connecting to the backup target. | [optional] 
**TargetAll** | Pointer to **bool** | Whether to backup all databases. Applies to database backup types. | [optional] 
**TargetName** | Pointer to **string** | Target database name. Applies when targetAll is false. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateBackupsRequestBackup{
    // Set fields directly
}
```

### StorageProviderId (Nullable)

Use the Nullable wrapper methods:
- `obj.StorageProviderId.IsSet()` — check if set
- `obj.StorageProviderId.Get()` — get the inner value (returns pointer)
- `obj.StorageProviderId.Set(&val)` — set the value
- `obj.StorageProviderId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


