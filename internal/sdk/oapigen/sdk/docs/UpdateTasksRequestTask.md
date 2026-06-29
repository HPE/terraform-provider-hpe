# UpdateTasksRequestTask

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllowCustomConfig** | Pointer to **bool** | When enabled, a text area is provided at Task execution time to allow the user to pass extra variables or specify extra configuration | [optional] 
**Name** | Pointer to **string** | A unique name for the task | [optional] 
**Code** | Pointer to **string** | A unique code for the task | [optional] 
**Visibility** | Pointer to **string** | Visibility | [optional] [default to "private"]
**TaskType** | Pointer to [**UpdateTasksRequestTaskTaskType**](UpdateTasksRequestTaskTaskType.md) |  | [optional] 
**Labels** | Pointer to **[]string** | An array of Category labels for filtering | [optional] 
**TaskOptions** | Pointer to [**UpdateTasksRequestTaskTaskOptions**](UpdateTasksRequestTaskTaskOptions.md) |  | [optional] 
**ResultType** | Pointer to **NullableString** |  | [optional] 
**ExecuteTarget** | Pointer to **string** | The execution target. eg. local,remote,resource. The default value varies by task type.  | [optional] 
**Retryable** | Pointer to **bool** | If the task should be retried or not. | [optional] [default to false]
**RetryCount** | Pointer to **int64** | The number of times to retry. | [optional] 
**RetryDelaySeconds** | Pointer to **int64** | The delay, between retries. | [optional] 
**File** | Pointer to [**UpdateTasksRequestTaskFile**](UpdateTasksRequestTaskFile.md) |  | [optional] 
**Credential** | Pointer to [**UpdateTasksRequestTaskCredential**](UpdateTasksRequestTaskCredential.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateTasksRequestTask{
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


