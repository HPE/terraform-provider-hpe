# AddPowerSchedulesRequestSchedule

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A name for the power schedule | 
**Description** | Pointer to **string** | A description for the power schedule | [optional] 
**ScheduleType** | Pointer to **string** | Type of schedule &#x60;power&#x60; on or &#x60;power off&#x60; | [optional] 
**ScheduleTimezone** | Pointer to **string** | Time Zone eg. America/New_York, Europe/Amsterdam, etc. | [optional] [default to "UTC"]
**Enabled** | Pointer to **bool** | Is the power schedule enabled | [optional] [default to true]
**MondayOnTime** | Pointer to **string** | Monday Start time of the day in 24-hour format | [optional] [default to "00:00"]
**MondayOffTime** | Pointer to **string** | Monday Off time of the day in 24-hour format | [optional] [default to "24:00"]
**TuesdayOnTime** | Pointer to **string** | Tuesday Start time of the day in 24-hour format | [optional] [default to "00:00"]
**TuesdayOffTime** | Pointer to **string** | Tuesday Off time of the day in 24-hour format | [optional] [default to "24:00"]
**WednesdayOnTime** | Pointer to **string** | Wednesday Start time of the day in 24-hour format | [optional] [default to "00:00"]
**WednesdayOffTime** | Pointer to **string** | Wednesday Off time of the day in 24-hour format | [optional] [default to "24:00"]
**ThursdayOnTime** | Pointer to **string** | Thursday Start time of the day in 24-hour format | [optional] [default to "00:00"]
**ThursdayOffTime** | Pointer to **string** | Thursday Off time of the day in 24-hour format | [optional] [default to "24:00"]
**FridayOnTime** | Pointer to **string** | Friday Start time of the day in 24-hour format | [optional] [default to "00:00"]
**FridayOffTime** | Pointer to **string** | Friday Off time of the day in 24-hour format | [optional] [default to "24:00"]
**SaturdayOnTime** | Pointer to **string** | Saturday Start time of the day in 24-hour format | [optional] [default to "00:00"]
**SaturdayOffTime** | Pointer to **string** | Saturday Off time of the day in 24-hour format | [optional] [default to "24:00"]
**SundayOnTime** | Pointer to **string** | Sunday Start time of the day in 24-hour format | [optional] [default to "00:00"]
**SundayOffTime** | Pointer to **string** | Sunday Off time of the day in 24-hour format | [optional] [default to "24:00"]

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddPowerSchedulesRequestSchedule{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


