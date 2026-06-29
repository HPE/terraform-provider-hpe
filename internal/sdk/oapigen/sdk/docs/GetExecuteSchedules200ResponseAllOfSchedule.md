# GetExecuteSchedules200ResponseAllOfSchedule

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**ScheduleType** | Pointer to **string** |  | [optional] 
**ScheduleTimezone** | Pointer to **NullableString** |  | [optional] 
**Cron** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetExecuteSchedules200ResponseAllOfSchedule{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Visibility (Nullable)

Use the Nullable wrapper methods:
- `obj.Visibility.IsSet()` — check if set
- `obj.Visibility.Get()` — get the inner value (returns pointer)
- `obj.Visibility.Set(&val)` — set the value
- `obj.Visibility.Unset()` — clear the value
### ScheduleTimezone (Nullable)

Use the Nullable wrapper methods:
- `obj.ScheduleTimezone.IsSet()` — check if set
- `obj.ScheduleTimezone.Get()` — get the inner value (returns pointer)
- `obj.ScheduleTimezone.Set(&val)` — set the value
- `obj.ScheduleTimezone.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


