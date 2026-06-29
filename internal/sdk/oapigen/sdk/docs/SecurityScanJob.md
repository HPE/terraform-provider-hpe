# SecurityScanJob

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A name for the Job | 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Enabled** | Pointer to **bool** | Use this to set enabled state | [optional] [default to true]
**SecurityPackage** | [**SecurityScanJobSecurityPackage**](SecurityScanJobSecurityPackage.md) |  | 
**ScanPath** | Pointer to **NullableString** | Scan Checklist | [optional] 
**SecurityProfile** | Pointer to **NullableString** | Security Profile | [optional] 
**TargetType** | **string** | Target type where job will execute | 
**Targets** | [**[]SecurityScanJobTargetsInner**](SecurityScanJobTargetsInner.md) |  | 
**ScheduleMode** | [**SecurityScanJobScheduleMode**](SecurityScanJobScheduleMode.md) |  | 
**CustomOptions** | Pointer to **map[string]interface{}** | Map of options to be used as values in the workflow tasks. These correspond to option types. | [optional] 
**CustomConfig** | Pointer to **string** | Job custom configuration (String in JSON format) | [optional] 
**DateTime** | Pointer to **time.Time** | Date and Time to execute the job. Use UTC time in the format 2020-02-15T05:00:00Z. Required when scheduleMode is &#39;dateTime&#39;. | [optional] 
**Run** | Pointer to **bool** | If true, executes job | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SecurityScanJob{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### ScanPath (Nullable)

Use the Nullable wrapper methods:
- `obj.ScanPath.IsSet()` — check if set
- `obj.ScanPath.Get()` — get the inner value (returns pointer)
- `obj.ScanPath.Set(&val)` — set the value
- `obj.ScanPath.Unset()` — clear the value
### SecurityProfile (Nullable)

Use the Nullable wrapper methods:
- `obj.SecurityProfile.IsSet()` — check if set
- `obj.SecurityProfile.Get()` — get the inner value (returns pointer)
- `obj.SecurityProfile.Set(&val)` — set the value
- `obj.SecurityProfile.Unset()` — clear the value
### Targets (Nullable)

Use the Nullable wrapper methods:
- `obj.Targets.IsSet()` — check if set
- `obj.Targets.Get()` — get the inner value (returns pointer)
- `obj.Targets.Set(&val)` — set the value
- `obj.Targets.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


