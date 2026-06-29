# RunReports200ResponseAllOfReportResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Type** | Pointer to [**RunReports200ResponseAllOfReportResultType**](RunReports200ResponseAllOfReportResultType.md) |  | [optional] 
**ReportTitle** | Pointer to **NullableString** |  | [optional] 
**FilterTitle** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**StartDate** | Pointer to **NullableTime** |  | [optional] 
**EndDate** | Pointer to **NullableTime** |  | [optional] 
**Config** | Pointer to [**RunReports200ResponseAllOfReportResultConfig**](RunReports200ResponseAllOfReportResultConfig.md) |  | [optional] 
**CreatedBy** | Pointer to [**RunReports200ResponseAllOfReportResultCreatedBy**](RunReports200ResponseAllOfReportResultCreatedBy.md) |  | [optional] 
**Rows** | Pointer to [**[]RunReports200ResponseAllOfReportResultRowsInner**](RunReports200ResponseAllOfReportResultRowsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &RunReports200ResponseAllOfReportResult{
    // Set fields directly
}
```

### ReportTitle (Nullable)

Use the Nullable wrapper methods:
- `obj.ReportTitle.IsSet()` — check if set
- `obj.ReportTitle.Get()` — get the inner value (returns pointer)
- `obj.ReportTitle.Set(&val)` — set the value
- `obj.ReportTitle.Unset()` — clear the value
### FilterTitle (Nullable)

Use the Nullable wrapper methods:
- `obj.FilterTitle.IsSet()` — check if set
- `obj.FilterTitle.Get()` — get the inner value (returns pointer)
- `obj.FilterTitle.Set(&val)` — set the value
- `obj.FilterTitle.Unset()` — clear the value
### StartDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StartDate.IsSet()` — check if set
- `obj.StartDate.Get()` — get the inner value (returns pointer)
- `obj.StartDate.Set(&val)` — set the value
- `obj.StartDate.Unset()` — clear the value
### EndDate (Nullable)

Use the Nullable wrapper methods:
- `obj.EndDate.IsSet()` — check if set
- `obj.EndDate.Get()` — get the inner value (returns pointer)
- `obj.EndDate.Set(&val)` — set the value
- `obj.EndDate.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


