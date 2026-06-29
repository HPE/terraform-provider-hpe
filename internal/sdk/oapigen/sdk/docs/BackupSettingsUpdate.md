# BackupSettingsUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BackupsEnabled** | Pointer to **bool** | Use this to enable / disable scheduled backups | [optional] 
**RetentionCount** | Pointer to **int64** | Maximum number of successful backups to retain | [optional] 
**CreateBackups** | Pointer to **bool** | Use this to enable / disable create backups | [optional] 
**BackupAppliance** | Pointer to **bool** | When enabled, a Backup will be created to backup the Morpheus appliance database | [optional] 
**UpdateExisting** | Pointer to **bool** | Use this to update existing backups with new settings | [optional] 
**DefaultSchedule** | Pointer to [**BackupSettingsUpdateDefaultSchedule**](BackupSettingsUpdateDefaultSchedule.md) |  | [optional] 
**ClearDefaultSchedule** | Pointer to **bool** | Use this to clear existing default backup schedule | [optional] 
**DefaultStorageBucket** | Pointer to [**BackupSettingsUpdateDefaultStorageBucket**](BackupSettingsUpdateDefaultStorageBucket.md) |  | [optional] 
**ClearDefaultStorageBucket** | Pointer to **bool** | Use this to clear default store bucket | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BackupSettingsUpdate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


