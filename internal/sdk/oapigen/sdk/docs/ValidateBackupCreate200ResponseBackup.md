# ValidateBackupCreate200ResponseBackup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LocationType** | Pointer to **string** | The resolved location type | [optional] 
**Name** | Pointer to **string** | The backup name | [optional] 
**InstanceId** | Pointer to **int64** |  | [optional] 
**ContainerId** | Pointer to **int64** |  | [optional] 
**ServerId** | Pointer to **int64** |  | [optional] 
**BackupType** | Pointer to **string** | The computed or submitted backup type code | [optional] 
**JobAction** | Pointer to **string** |  | [optional] 
**JobId** | Pointer to **int64** |  | [optional] 
**JobName** | Pointer to **string** |  | [optional] 
**TargetPort** | Pointer to **int32** | Default target port, computed for database backup types | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ValidateBackupCreate200ResponseBackup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


