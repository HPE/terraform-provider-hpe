# BackupSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BackupsEnabled** | Pointer to **bool** |  | [optional] 
**CreateBackups** | Pointer to **bool** |  | [optional] 
**BackupAppliance** | Pointer to **bool** |  | [optional] 
**DefaultStorageBucket** | Pointer to [**BackupSettingsDefaultStorageBucket**](BackupSettingsDefaultStorageBucket.md) |  | [optional] 
**DefaultSchedule** | Pointer to [**BackupSettingsDefaultSchedule**](BackupSettingsDefaultSchedule.md) |  | [optional] 
**RetentionCount** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BackupSettings{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


