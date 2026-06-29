# ReportTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**Visible** | Pointer to **bool** |  | [optional] 
**MasterOnly** | Pointer to **bool** |  | [optional] 
**OwnerOnly** | Pointer to **bool** |  | [optional] 
**SupportsAllZoneTypes** | Pointer to **bool** |  | [optional] 
**IsPlugin** | Pointer to **NullableBool** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**OptionTypes** | Pointer to [**[]ReportTypesInnerOptionTypesInner**](ReportTypesInnerOptionTypesInner.md) |  | [optional] 
**SupportedZoneTypes** | Pointer to [**[]ReportTypesInnerSupportedZoneTypesInner**](ReportTypesInnerSupportedZoneTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ReportTypesInner{
    // Set fields directly
}
```

### IsPlugin (Nullable)

Use the Nullable wrapper methods:
- `obj.IsPlugin.IsSet()` — check if set
- `obj.IsPlugin.Get()` — get the inner value (returns pointer)
- `obj.IsPlugin.Set(&val)` — set the value
- `obj.IsPlugin.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


