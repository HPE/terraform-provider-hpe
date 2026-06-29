# AddExecuteSchedulesRequestSchedule

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A name for the execute schedule | 
**Description** | Pointer to **string** | A description for the execute schedule | [optional] 
**ScheduleType** | Pointer to **string** | Type of schedule | [optional] 
**ScheduleTimezone** | Pointer to **string** | Time Zone eg. America/New_York, Europe/Amsterdam, etc. | [optional] [default to "UTC"]
**Cron** | Pointer to **string** | Cron Expression. The default is daily at midnight | [optional] [default to "0 0 * * *"]
**Enabled** | Pointer to **bool** | Is enabled | [optional] [default to true]

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddExecuteSchedulesRequestSchedule{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


