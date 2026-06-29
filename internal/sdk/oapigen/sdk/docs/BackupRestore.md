# BackupRestore

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Backup Result ID | [optional] 
**BackupResultId** | Pointer to **int64** |  | [optional] 
**BackupId** | Pointer to **int64** |  | [optional] 
**Backup** | Pointer to [**GetBackupRestores200ResponseRestoreBackup**](GetBackupRestores200ResponseRestoreBackup.md) |  | [optional] 
**ContainerId** | Pointer to **NullableInt64** |  | [optional] 
**Container** | Pointer to [**GetBackupRestores200ResponseRestoreContainer**](GetBackupRestores200ResponseRestoreContainer.md) |  | [optional] 
**Instance** | Pointer to [**GetBackupRestores200ResponseRestoreInstance**](GetBackupRestores200ResponseRestoreInstance.md) |  | [optional] 
**RestoreToNew** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **NullableString** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**StartDate** | Pointer to **NullableTime** |  | [optional] 
**EndDate** | Pointer to **NullableTime** |  | [optional] 
**DurationMillis** | Pointer to **NullableInt64** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalStatusRef** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** | Date Created | [optional] 
**LastUpdated** | Pointer to **time.Time** | Last Updated | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BackupRestore{
    // Set fields directly
}
```

### ContainerId (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerId.IsSet()` — check if set
- `obj.ContainerId.Get()` — get the inner value (returns pointer)
- `obj.ContainerId.Set(&val)` — set the value
- `obj.ContainerId.Unset()` — clear the value
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
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### ExternalStatusRef (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalStatusRef.IsSet()` — check if set
- `obj.ExternalStatusRef.Get()` — get the inner value (returns pointer)
- `obj.ExternalStatusRef.Set(&val)` — set the value
- `obj.ExternalStatusRef.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


