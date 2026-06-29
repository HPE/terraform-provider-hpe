# ExecuteBackups200ResponseAllOfBackup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Backup ID | [optional] 
**Name** | Pointer to **string** | Name | [optional] 
**LocationType** | Pointer to **string** | Source Type (instance, server, storage) | [optional] 
**Location** | Pointer to **NullableString** | Location label, typically the storage provider name | [optional] 
**Instance** | Pointer to [**ExecuteBackups200ResponseAllOfBackupInstance**](ExecuteBackups200ResponseAllOfBackupInstance.md) |  | [optional] 
**ContainerId** | Pointer to **NullableInt64** | Present when locationType is &#x60;instance&#x60; | [optional] 
**Server** | Pointer to [**ExecuteBackups200ResponseAllOfBackupServer**](ExecuteBackups200ResponseAllOfBackupServer.md) |  | [optional] 
**VolumePath** | Pointer to **NullableString** | Volume path. Present when backupType is &#x60;lvmSnapshot&#x60;. | [optional] 
**TargetPath** | Pointer to **NullableString** | The file or directory path on the target host. Present when backupType is &#x60;fileBackup&#x60;, &#x60;directoryBackup&#x60;, or &#x60;tarDirectoryBackup&#x60;. | [optional] 
**TargetHost** | Pointer to **NullableString** | Target host. Present when backupType is a database type (MySQL, SqlServer, Postgres, MongoDB) or lvmSnapshot. | [optional] 
**TargetPort** | Pointer to **NullableInt32** | Target port. Present when backupType is a database type (MySQL, SqlServer, Postgres, MongoDB). | [optional] 
**TargetAll** | Pointer to **NullableBool** | Whether to backup all databases. Present when backupType is a database type. | [optional] 
**TargetName** | Pointer to **NullableString** | Target database name. Present when backupType is a database type and &#x60;targetAll&#x60; is false. | [optional] 
**TargetUsername** | Pointer to **NullableString** | Target username. Present when backupType is a database type. | [optional] 
**TargetPassword** | Pointer to **NullableString** | Target password (masked). Present when backupType is a database type. | [optional] 
**TargetPasswordHash** | Pointer to **NullableString** | Target password hash. Present when backupType is a database type. | [optional] 
**Job** | Pointer to [**ExecuteBackups200ResponseAllOfBackupJob**](ExecuteBackups200ResponseAllOfBackupJob.md) |  | [optional] 
**Schedule** | Pointer to [**ExecuteBackups200ResponseAllOfBackupSchedule**](ExecuteBackups200ResponseAllOfBackupSchedule.md) |  | [optional] 
**RetentionCount** | Pointer to **NullableInt64** |  | [optional] 
**BackupType** | Pointer to [**ExecuteBackups200ResponseAllOfBackupBackupType**](ExecuteBackups200ResponseAllOfBackupBackupType.md) |  | [optional] 
**BackupProvider** | Pointer to [**ExecuteBackups200ResponseAllOfBackupBackupProvider**](ExecuteBackups200ResponseAllOfBackupBackupProvider.md) |  | [optional] 
**StorageProvider** | Pointer to [**ExecuteBackups200ResponseAllOfBackupStorageProvider**](ExecuteBackups200ResponseAllOfBackupStorageProvider.md) |  | [optional] 
**BackupRepository** | Pointer to [**ExecuteBackups200ResponseAllOfBackupBackupRepository**](ExecuteBackups200ResponseAllOfBackupBackupRepository.md) |  | [optional] 
**CronExpression** | Pointer to **NullableString** | Cron Expression | [optional] 
**NextFire** | Pointer to **NullableTime** | Next Fire | [optional] 
**LastStatus** | Pointer to **NullableString** | Last Status | [optional] 
**LastResult** | Pointer to [**ExecuteBackups200ResponseAllOfBackupLastResult**](ExecuteBackups200ResponseAllOfBackupLastResult.md) |  | [optional] 
**Stats** | Pointer to [**ExecuteBackups200ResponseAllOfBackupStats**](ExecuteBackups200ResponseAllOfBackupStats.md) |  | [optional] 
**Enabled** | Pointer to **bool** | Enabled | [optional] 
**DateCreated** | Pointer to **time.Time** | Date Created | [optional] 
**LastUpdated** | Pointer to **time.Time** | Last Updated | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ExecuteBackups200ResponseAllOfBackup{
    // Set fields directly
}
```

### Location (Nullable)

Use the Nullable wrapper methods:
- `obj.Location.IsSet()` — check if set
- `obj.Location.Get()` — get the inner value (returns pointer)
- `obj.Location.Set(&val)` — set the value
- `obj.Location.Unset()` — clear the value
### ContainerId (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerId.IsSet()` — check if set
- `obj.ContainerId.Get()` — get the inner value (returns pointer)
- `obj.ContainerId.Set(&val)` — set the value
- `obj.ContainerId.Unset()` — clear the value
### VolumePath (Nullable)

Use the Nullable wrapper methods:
- `obj.VolumePath.IsSet()` — check if set
- `obj.VolumePath.Get()` — get the inner value (returns pointer)
- `obj.VolumePath.Set(&val)` — set the value
- `obj.VolumePath.Unset()` — clear the value
### TargetPath (Nullable)

Use the Nullable wrapper methods:
- `obj.TargetPath.IsSet()` — check if set
- `obj.TargetPath.Get()` — get the inner value (returns pointer)
- `obj.TargetPath.Set(&val)` — set the value
- `obj.TargetPath.Unset()` — clear the value
### TargetHost (Nullable)

Use the Nullable wrapper methods:
- `obj.TargetHost.IsSet()` — check if set
- `obj.TargetHost.Get()` — get the inner value (returns pointer)
- `obj.TargetHost.Set(&val)` — set the value
- `obj.TargetHost.Unset()` — clear the value
### TargetPort (Nullable)

Use the Nullable wrapper methods:
- `obj.TargetPort.IsSet()` — check if set
- `obj.TargetPort.Get()` — get the inner value (returns pointer)
- `obj.TargetPort.Set(&val)` — set the value
- `obj.TargetPort.Unset()` — clear the value
### TargetAll (Nullable)

Use the Nullable wrapper methods:
- `obj.TargetAll.IsSet()` — check if set
- `obj.TargetAll.Get()` — get the inner value (returns pointer)
- `obj.TargetAll.Set(&val)` — set the value
- `obj.TargetAll.Unset()` — clear the value
### TargetName (Nullable)

Use the Nullable wrapper methods:
- `obj.TargetName.IsSet()` — check if set
- `obj.TargetName.Get()` — get the inner value (returns pointer)
- `obj.TargetName.Set(&val)` — set the value
- `obj.TargetName.Unset()` — clear the value
### TargetUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.TargetUsername.IsSet()` — check if set
- `obj.TargetUsername.Get()` — get the inner value (returns pointer)
- `obj.TargetUsername.Set(&val)` — set the value
- `obj.TargetUsername.Unset()` — clear the value
### TargetPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.TargetPassword.IsSet()` — check if set
- `obj.TargetPassword.Get()` — get the inner value (returns pointer)
- `obj.TargetPassword.Set(&val)` — set the value
- `obj.TargetPassword.Unset()` — clear the value
### TargetPasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.TargetPasswordHash.IsSet()` — check if set
- `obj.TargetPasswordHash.Get()` — get the inner value (returns pointer)
- `obj.TargetPasswordHash.Set(&val)` — set the value
- `obj.TargetPasswordHash.Unset()` — clear the value
### RetentionCount (Nullable)

Use the Nullable wrapper methods:
- `obj.RetentionCount.IsSet()` — check if set
- `obj.RetentionCount.Get()` — get the inner value (returns pointer)
- `obj.RetentionCount.Set(&val)` — set the value
- `obj.RetentionCount.Unset()` — clear the value
### CronExpression (Nullable)

Use the Nullable wrapper methods:
- `obj.CronExpression.IsSet()` — check if set
- `obj.CronExpression.Get()` — get the inner value (returns pointer)
- `obj.CronExpression.Set(&val)` — set the value
- `obj.CronExpression.Unset()` — clear the value
### NextFire (Nullable)

Use the Nullable wrapper methods:
- `obj.NextFire.IsSet()` — check if set
- `obj.NextFire.Get()` — get the inner value (returns pointer)
- `obj.NextFire.Set(&val)` — set the value
- `obj.NextFire.Unset()` — clear the value
### LastStatus (Nullable)

Use the Nullable wrapper methods:
- `obj.LastStatus.IsSet()` — check if set
- `obj.LastStatus.Get()` — get the inner value (returns pointer)
- `obj.LastStatus.Set(&val)` — set the value
- `obj.LastStatus.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


