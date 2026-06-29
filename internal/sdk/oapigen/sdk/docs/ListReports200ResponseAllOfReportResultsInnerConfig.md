# ListReports200ResponseAllOfReportResultsInnerConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ReportType** | Pointer to **NullableString** |  | [optional] 
**StartDate** | Pointer to **string** |  | [optional] 
**EndDate** | Pointer to **string** |  | [optional] 
**CloudId** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListReports200ResponseAllOfReportResultsInnerConfig{
    // Set fields directly
}
```

### ReportType (Nullable)

Use the Nullable wrapper methods:
- `obj.ReportType.IsSet()` — check if set
- `obj.ReportType.Get()` — get the inner value (returns pointer)
- `obj.ReportType.Set(&val)` — set the value
- `obj.ReportType.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


