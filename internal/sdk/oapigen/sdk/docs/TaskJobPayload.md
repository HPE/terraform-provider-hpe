# TaskJobPayload

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A name for the Job | 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Enabled** | Pointer to **bool** | Use this to set enabled state | [optional] [default to true]
**Task** | [**TaskJobPayloadTask**](TaskJobPayloadTask.md) |  | 
**TargetType** | **string** | Target type where job will execute | 
**Targets** | Pointer to [**[]TaskJobPayloadTargetsInner**](TaskJobPayloadTargetsInner.md) |  | [optional] 
**ScheduleMode** | [**TaskJobPayloadScheduleMode**](TaskJobPayloadScheduleMode.md) |  | 
**CustomOptions** | Pointer to **map[string]interface{}** | Map of options to be used as values in the workflow tasks. These correspond to option types. | [optional] 
**CustomConfig** | Pointer to **string** | Job custom configuration (String in JSON format) | [optional] 
**DateTime** | Pointer to **time.Time** | Date and Time to execute the job. Use UTC time in the format 2020-02-15T05:00:00Z. Required when scheduleMode is &#39;dateTime&#39;. | [optional] 
**Run** | Pointer to **bool** | If true, executes job | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &TaskJobPayload{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### Targets (Nullable)

Use the Nullable wrapper methods:
- `obj.Targets.IsSet()` — check if set
- `obj.Targets.Get()` — get the inner value (returns pointer)
- `obj.Targets.Set(&val)` — set the value
- `obj.Targets.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


