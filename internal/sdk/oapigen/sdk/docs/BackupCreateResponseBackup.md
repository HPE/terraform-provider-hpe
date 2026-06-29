# BackupCreateResponseBackup

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
**TargetPort** | Pointer to **NullableInt32** | Default target port, computed for database backup types | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BackupCreateResponseBackup{
    // Set fields directly
}
```

### TargetPort (Nullable)

Use the Nullable wrapper methods:
- `obj.TargetPort.IsSet()` — check if set
- `obj.TargetPort.Get()` — get the inner value (returns pointer)
- `obj.TargetPort.Set(&val)` — set the value
- `obj.TargetPort.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


