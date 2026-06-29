# ListPowerSchedules200ResponseAllOfSchedulesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**ScheduleType** | Pointer to **string** |  | [optional] 
**ScheduleTimezone** | Pointer to **string** |  | [optional] 
**MondayOn** | Pointer to **int64** |  | [optional] 
**MondayOnTime** | Pointer to **string** |  | [optional] 
**MondayOff** | Pointer to **int64** |  | [optional] 
**MondayOffTime** | Pointer to **string** |  | [optional] 
**TuesdayOn** | Pointer to **int64** |  | [optional] 
**TuesdayOnTime** | Pointer to **string** |  | [optional] 
**TuesdayOff** | Pointer to **int64** |  | [optional] 
**TuesdayOffTime** | Pointer to **string** |  | [optional] 
**WednesdayOn** | Pointer to **int64** |  | [optional] 
**WednesdayOnTime** | Pointer to **string** |  | [optional] 
**WednesdayOff** | Pointer to **int64** |  | [optional] 
**WednesdayOffTime** | Pointer to **string** |  | [optional] 
**ThursdayOn** | Pointer to **int64** |  | [optional] 
**ThursdayOnTime** | Pointer to **string** |  | [optional] 
**ThursdayOff** | Pointer to **int64** |  | [optional] 
**ThursdayOffTime** | Pointer to **string** |  | [optional] 
**FridayOn** | Pointer to **int64** |  | [optional] 
**FridayOnTime** | Pointer to **string** |  | [optional] 
**FridayOff** | Pointer to **int64** |  | [optional] 
**FridayOffTime** | Pointer to **string** |  | [optional] 
**SaturdayOn** | Pointer to **int64** |  | [optional] 
**SaturdayOnTime** | Pointer to **string** |  | [optional] 
**SaturdayOff** | Pointer to **int64** |  | [optional] 
**SaturdayOffTime** | Pointer to **string** |  | [optional] 
**SundayOn** | Pointer to **int64** |  | [optional] 
**SundayOnTime** | Pointer to **string** |  | [optional] 
**SundayOff** | Pointer to **int64** |  | [optional] 
**SundayOffTime** | Pointer to **string** |  | [optional] 
**TotalMonthlyHoursSaved** | Pointer to **float32** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListPowerSchedules200ResponseAllOfSchedulesInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


