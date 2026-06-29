# BackupCreateResponseBackupSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RetentionCount** | Pointer to **NullableInt32** | Default retention count | [optional] 
**DefaultBackupSchedule** | Pointer to **NullableInt64** | Default backup schedule ID (execute schedule type) | [optional] 
**DefaultSyntheticFullBackupsEnabled** | Pointer to **NullableBool** | Whether synthetic full backups are enabled by default | [optional] 
**DefaultSyntheticFullBackupSchedule** | Pointer to **NullableInt64** | Default synthetic full backup schedule ID | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BackupCreateResponseBackupSettings{
    // Set fields directly
}
```

### RetentionCount (Nullable)

Use the Nullable wrapper methods:
- `obj.RetentionCount.IsSet()` — check if set
- `obj.RetentionCount.Get()` — get the inner value (returns pointer)
- `obj.RetentionCount.Set(&val)` — set the value
- `obj.RetentionCount.Unset()` — clear the value
### DefaultBackupSchedule (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultBackupSchedule.IsSet()` — check if set
- `obj.DefaultBackupSchedule.Get()` — get the inner value (returns pointer)
- `obj.DefaultBackupSchedule.Set(&val)` — set the value
- `obj.DefaultBackupSchedule.Unset()` — clear the value
### DefaultSyntheticFullBackupsEnabled (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultSyntheticFullBackupsEnabled.IsSet()` — check if set
- `obj.DefaultSyntheticFullBackupsEnabled.Get()` — get the inner value (returns pointer)
- `obj.DefaultSyntheticFullBackupsEnabled.Set(&val)` — set the value
- `obj.DefaultSyntheticFullBackupsEnabled.Unset()` — clear the value
### DefaultSyntheticFullBackupSchedule (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultSyntheticFullBackupSchedule.IsSet()` — check if set
- `obj.DefaultSyntheticFullBackupSchedule.Get()` — get the inner value (returns pointer)
- `obj.DefaultSyntheticFullBackupSchedule.Set(&val)` — set the value
- `obj.DefaultSyntheticFullBackupSchedule.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


