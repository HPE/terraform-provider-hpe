# UpdateWorkflows200ResponseAllOfTaskSetTaskSetTasksInnerTask

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**TaskType** | Pointer to [**UpdateWorkflows200ResponseAllOfTaskSetTaskSetTasksInnerTaskTaskType**](UpdateWorkflows200ResponseAllOfTaskSetTaskSetTasksInnerTaskTaskType.md) |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**TaskOptions** | Pointer to [**UpdateWorkflows200ResponseAllOfTaskSetTaskSetTasksInnerTaskTaskOptions**](UpdateWorkflows200ResponseAllOfTaskSetTaskSetTasksInnerTaskTaskOptions.md) |  | [optional] 
**File** | Pointer to [**UpdateWorkflows200ResponseAllOfTaskSetTaskSetTasksInnerTaskFile**](UpdateWorkflows200ResponseAllOfTaskSetTaskSetTasksInnerTaskFile.md) |  | [optional] 
**ResultType** | Pointer to **NullableString** |  | [optional] 
**ExecuteTarget** | Pointer to **string** |  | [optional] 
**Retryable** | Pointer to **bool** |  | [optional] 
**RetryCount** | Pointer to **int64** |  | [optional] 
**RetryDelaySeconds** | Pointer to **int64** |  | [optional] 
**AllowCustomConfig** | Pointer to **bool** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateWorkflows200ResponseAllOfTaskSetTaskSetTasksInnerTask{
    // Set fields directly
}
```

### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### ResultType (Nullable)

Use the Nullable wrapper methods:
- `obj.ResultType.IsSet()` — check if set
- `obj.ResultType.Get()` — get the inner value (returns pointer)
- `obj.ResultType.Set(&val)` — set the value
- `obj.ResultType.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


