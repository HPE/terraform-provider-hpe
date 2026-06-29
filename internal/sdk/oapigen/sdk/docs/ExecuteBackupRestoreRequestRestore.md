# ExecuteBackupRestoreRequestRestore

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BackupResultId** | **int64** | Id of backup result | 
**RestoreInstance** | **string** | Type of restore | 
**InstanceId** | Pointer to **int64** | Id of instance | [optional] 
**SiteId** | Pointer to **int64** | Id of site for restore to new | [optional] 
**Config** | Pointer to **map[string]interface{}** | Additional config | [optional] 
**InstanceConfig** | Pointer to **map[string]interface{}** | Instance config for restore to new. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ExecuteBackupRestoreRequestRestore{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


