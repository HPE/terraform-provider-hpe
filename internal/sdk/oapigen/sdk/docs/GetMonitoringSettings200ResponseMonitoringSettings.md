# GetMonitoringSettings200ResponseMonitoringSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AutoManageChecks** | Pointer to **bool** | Auto Create Checks | [optional] 
**AvailabilityTimeFrame** | Pointer to **NullableInt32** | Availability Time Frame. The number of days availability should be calculated for. Changes will not take effect until your checks have passed their check interval. | [optional] 
**AvailabilityPrecision** | Pointer to **NullableInt32** | Availability Precision. The number of decimal places availability should be displayed in. Can be anywhere between 0 and 5. | [optional] 
**DefaultCheckInterval** | Pointer to **NullableInt32** | Default Check Interval. The number of minutes to use as the default interval to use when creating new checks. | [optional] 
**ServiceNow** | Pointer to [**GetMonitoringSettings200ResponseMonitoringSettingsServiceNow**](GetMonitoringSettings200ResponseMonitoringSettingsServiceNow.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetMonitoringSettings200ResponseMonitoringSettings{
    // Set fields directly
}
```

### AvailabilityTimeFrame (Nullable)

Use the Nullable wrapper methods:
- `obj.AvailabilityTimeFrame.IsSet()` — check if set
- `obj.AvailabilityTimeFrame.Get()` — get the inner value (returns pointer)
- `obj.AvailabilityTimeFrame.Set(&val)` — set the value
- `obj.AvailabilityTimeFrame.Unset()` — clear the value
### AvailabilityPrecision (Nullable)

Use the Nullable wrapper methods:
- `obj.AvailabilityPrecision.IsSet()` — check if set
- `obj.AvailabilityPrecision.Get()` — get the inner value (returns pointer)
- `obj.AvailabilityPrecision.Set(&val)` — set the value
- `obj.AvailabilityPrecision.Unset()` — clear the value
### DefaultCheckInterval (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultCheckInterval.IsSet()` — check if set
- `obj.DefaultCheckInterval.Get()` — get the inner value (returns pointer)
- `obj.DefaultCheckInterval.Set(&val)` — set the value
- `obj.DefaultCheckInterval.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


