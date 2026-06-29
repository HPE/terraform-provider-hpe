# JobExecution

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Process** | Pointer to **NullableString** |  | [optional] 
**Job** | Pointer to [**GetJobExecutions200ResponseAllOfJobExecutionJob**](GetJobExecutions200ResponseAllOfJobExecutionJob.md) |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**Duration** | Pointer to **int64** |  | [optional] 
**ResultData** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**CreatedBy** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &JobExecution{
    // Set fields directly
}
```

### Process (Nullable)

Use the Nullable wrapper methods:
- `obj.Process.IsSet()` — check if set
- `obj.Process.Get()` — get the inner value (returns pointer)
- `obj.Process.Set(&val)` — set the value
- `obj.Process.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ResultData (Nullable)

Use the Nullable wrapper methods:
- `obj.ResultData.IsSet()` — check if set
- `obj.ResultData.Get()` — get the inner value (returns pointer)
- `obj.ResultData.Set(&val)` — set the value
- `obj.ResultData.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### CreatedBy (Nullable)

Use the Nullable wrapper methods:
- `obj.CreatedBy.IsSet()` — check if set
- `obj.CreatedBy.Get()` — get the inner value (returns pointer)
- `obj.CreatedBy.Set(&val)` — set the value
- `obj.CreatedBy.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


