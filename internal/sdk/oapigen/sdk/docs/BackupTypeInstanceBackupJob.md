# BackupTypeInstanceBackupJob

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SyntheticFullEnabled** | Pointer to **bool** | Enable synthetic full backups on this backup job. Only applies to backup types whose job type supports synthetic full backups. | [optional] 
**SyntheticFullSchedule** | Pointer to **int64** | The ID of the execute schedule for the synthetic full backup. Required when &#x60;syntheticFullEnabled&#x60; is true. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BackupTypeInstanceBackupJob{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


