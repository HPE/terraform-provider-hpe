# AddTasksRequestTask

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllowCustomConfig** | Pointer to **bool** | When enabled, a text area is provided at Task execution time to allow the user to pass extra variables or specify extra configuration | [optional] 
**Name** | **string** | A unique name for the task | 
**Code** | Pointer to **string** | A unique code for the task | [optional] 
**Visibility** | Pointer to **string** | private or public | [optional] [default to "private"]
**TaskType** | [**AddTasksRequestTaskTaskType**](AddTasksRequestTaskTaskType.md) |  | 
**Labels** | Pointer to **[]string** | An array of Category labels for filtering | [optional] 
**TaskOptions** | Pointer to [**AddTasksRequestTaskTaskOptions**](AddTasksRequestTaskTaskOptions.md) |  | [optional] 
**ResultType** | Pointer to **NullableString** |  | [optional] 
**ExecuteTarget** | **string** | The execution target. eg. local,remote,resource. The default value varies by task type.  | 
**Retryable** | Pointer to **bool** | If the task should be retried or not. | [optional] [default to false]
**RetryCount** | Pointer to **int64** | The number of times to retry. | [optional] 
**RetryDelaySeconds** | Pointer to **int64** | The delay, between retries. | [optional] 
**File** | Pointer to [**AddTasksRequestTaskFile**](AddTasksRequestTaskFile.md) |  | [optional] 
**Credential** | Pointer to [**AddTasksRequestTaskCredential**](AddTasksRequestTaskCredential.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddTasksRequestTask{
    // Set fields directly
}
```

### ResultType (Nullable)

Use the Nullable wrapper methods:
- `obj.ResultType.IsSet()` — check if set
- `obj.ResultType.Get()` — get the inner value (returns pointer)
- `obj.ResultType.Set(&val)` — set the value
- `obj.ResultType.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


