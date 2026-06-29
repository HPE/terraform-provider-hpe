# ListBackupSettings200ResponseBackupSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BackupsEnabled** | Pointer to **bool** |  | [optional] 
**CreateBackups** | Pointer to **bool** |  | [optional] 
**BackupAppliance** | Pointer to **bool** |  | [optional] 
**DefaultStorageBucket** | Pointer to [**NullableListBackupSettings200ResponseBackupSettingsDefaultStorageBucket**](ListBackupSettings200ResponseBackupSettingsDefaultStorageBucket.md) |  | [optional] 
**DefaultSchedule** | Pointer to [**ListBackupSettings200ResponseBackupSettingsDefaultSchedule**](ListBackupSettings200ResponseBackupSettingsDefaultSchedule.md) |  | [optional] 
**RetentionCount** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListBackupSettings200ResponseBackupSettings{
    // Set fields directly
}
```

### DefaultStorageBucket (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultStorageBucket.IsSet()` — check if set
- `obj.DefaultStorageBucket.Get()` — get the inner value (returns pointer)
- `obj.DefaultStorageBucket.Set(&val)` — set the value
- `obj.DefaultStorageBucket.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


