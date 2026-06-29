# AddBackupJobsRequestJob

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A name for the backup job | 
**Code** | Pointer to **string** | A code for the backup job | [optional] 
**ScheduleId** | Pointer to **NullableInt64** | Execute Schedule ID to use for the backup job | [optional] 
**RetentionCount** | Pointer to **int64** | Retention Count. By default the backup settings value will be used. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddBackupJobsRequestJob{
    // Set fields directly
}
```

### ScheduleId (Nullable)

Use the Nullable wrapper methods:
- `obj.ScheduleId.IsSet()` — check if set
- `obj.ScheduleId.Get()` — get the inner value (returns pointer)
- `obj.ScheduleId.Set(&val)` — set the value
- `obj.ScheduleId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


