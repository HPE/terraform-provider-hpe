# BackupResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Backup Result ID | [optional] 
**Backup** | Pointer to [**GetBackupResults200ResponseResultBackup**](GetBackupResults200ResponseResultBackup.md) |  | [optional] 
**BackupSetId** | Pointer to **NullableString** |  | [optional] 
**InstanceId** | Pointer to **NullableInt64** |  | [optional] 
**ContainerId** | Pointer to **NullableInt64** |  | [optional] 
**ServerId** | Pointer to **NullableInt64** |  | [optional] 
**Status** | Pointer to **NullableString** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**StartDate** | Pointer to **NullableTime** |  | [optional] 
**EndDate** | Pointer to **NullableTime** |  | [optional] 
**DurationMillis** | Pointer to **NullableInt64** |  | [optional] 
**SizeInBytes** | Pointer to **NullableInt64** |  | [optional] 
**SizeInMb** | Pointer to **NullableInt64** |  | [optional] 
**VolumePath** | Pointer to **NullableString** |  | [optional] 
**ResultArchive** | Pointer to **NullableString** |  | [optional] 
**ResultPath** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**SnapshotId** | Pointer to **NullableString** |  | [optional] 
**SnapshotExternalId** | Pointer to **NullableString** |  | [optional] 
**CreatedBy** | Pointer to [**GetBackupResults200ResponseResultCreatedBy**](GetBackupResults200ResponseResultCreatedBy.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** | Date Created | [optional] 
**LastUpdated** | Pointer to **time.Time** | Last Updated | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BackupResult{
    // Set fields directly
}
```

### BackupSetId (Nullable)

Use the Nullable wrapper methods:
- `obj.BackupSetId.IsSet()` — check if set
- `obj.BackupSetId.Get()` — get the inner value (returns pointer)
- `obj.BackupSetId.Set(&val)` — set the value
- `obj.BackupSetId.Unset()` — clear the value
### InstanceId (Nullable)

Use the Nullable wrapper methods:
- `obj.InstanceId.IsSet()` — check if set
- `obj.InstanceId.Get()` — get the inner value (returns pointer)
- `obj.InstanceId.Set(&val)` — set the value
- `obj.InstanceId.Unset()` — clear the value
### ContainerId (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerId.IsSet()` — check if set
- `obj.ContainerId.Get()` — get the inner value (returns pointer)
- `obj.ContainerId.Set(&val)` — set the value
- `obj.ContainerId.Unset()` — clear the value
### ServerId (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerId.IsSet()` — check if set
- `obj.ServerId.Get()` — get the inner value (returns pointer)
- `obj.ServerId.Set(&val)` — set the value
- `obj.ServerId.Unset()` — clear the value
### Status (Nullable)

Use the Nullable wrapper methods:
- `obj.Status.IsSet()` — check if set
- `obj.Status.Get()` — get the inner value (returns pointer)
- `obj.Status.Set(&val)` — set the value
- `obj.Status.Unset()` — clear the value
### ErrorMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.ErrorMessage.IsSet()` — check if set
- `obj.ErrorMessage.Get()` — get the inner value (returns pointer)
- `obj.ErrorMessage.Set(&val)` — set the value
- `obj.ErrorMessage.Unset()` — clear the value
### StartDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StartDate.IsSet()` — check if set
- `obj.StartDate.Get()` — get the inner value (returns pointer)
- `obj.StartDate.Set(&val)` — set the value
- `obj.StartDate.Unset()` — clear the value
### EndDate (Nullable)

Use the Nullable wrapper methods:
- `obj.EndDate.IsSet()` — check if set
- `obj.EndDate.Get()` — get the inner value (returns pointer)
- `obj.EndDate.Set(&val)` — set the value
- `obj.EndDate.Unset()` — clear the value
### DurationMillis (Nullable)

Use the Nullable wrapper methods:
- `obj.DurationMillis.IsSet()` — check if set
- `obj.DurationMillis.Get()` — get the inner value (returns pointer)
- `obj.DurationMillis.Set(&val)` — set the value
- `obj.DurationMillis.Unset()` — clear the value
### SizeInBytes (Nullable)

Use the Nullable wrapper methods:
- `obj.SizeInBytes.IsSet()` — check if set
- `obj.SizeInBytes.Get()` — get the inner value (returns pointer)
- `obj.SizeInBytes.Set(&val)` — set the value
- `obj.SizeInBytes.Unset()` — clear the value
### SizeInMb (Nullable)

Use the Nullable wrapper methods:
- `obj.SizeInMb.IsSet()` — check if set
- `obj.SizeInMb.Get()` — get the inner value (returns pointer)
- `obj.SizeInMb.Set(&val)` — set the value
- `obj.SizeInMb.Unset()` — clear the value
### VolumePath (Nullable)

Use the Nullable wrapper methods:
- `obj.VolumePath.IsSet()` — check if set
- `obj.VolumePath.Get()` — get the inner value (returns pointer)
- `obj.VolumePath.Set(&val)` — set the value
- `obj.VolumePath.Unset()` — clear the value
### ResultArchive (Nullable)

Use the Nullable wrapper methods:
- `obj.ResultArchive.IsSet()` — check if set
- `obj.ResultArchive.Get()` — get the inner value (returns pointer)
- `obj.ResultArchive.Set(&val)` — set the value
- `obj.ResultArchive.Unset()` — clear the value
### ResultPath (Nullable)

Use the Nullable wrapper methods:
- `obj.ResultPath.IsSet()` — check if set
- `obj.ResultPath.Get()` — get the inner value (returns pointer)
- `obj.ResultPath.Set(&val)` — set the value
- `obj.ResultPath.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### SnapshotId (Nullable)

Use the Nullable wrapper methods:
- `obj.SnapshotId.IsSet()` — check if set
- `obj.SnapshotId.Get()` — get the inner value (returns pointer)
- `obj.SnapshotId.Set(&val)` — set the value
- `obj.SnapshotId.Unset()` — clear the value
### SnapshotExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.SnapshotExternalId.IsSet()` — check if set
- `obj.SnapshotExternalId.Get()` — get the inner value (returns pointer)
- `obj.SnapshotExternalId.Set(&val)` — set the value
- `obj.SnapshotExternalId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


