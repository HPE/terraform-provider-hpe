# ListSecurityScans200ResponseAllOfSecurityScansInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**SecurityPackage** | Pointer to [**ListSecurityScans200ResponseAllOfSecurityScansInnerSecurityPackage**](ListSecurityScans200ResponseAllOfSecurityScansInnerSecurityPackage.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**ScanDate** | Pointer to **time.Time** |  | [optional] 
**ScanDuration** | Pointer to **int64** |  | [optional] 
**TestCount** | Pointer to **int64** |  | [optional] 
**RunCount** | Pointer to **int64** |  | [optional] 
**PassCount** | Pointer to **int64** |  | [optional] 
**FailCount** | Pointer to **int64** |  | [optional] 
**OtherCount** | Pointer to **int64** |  | [optional] 
**ScanScore** | Pointer to **float32** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**CreatedBy** | Pointer to **string** |  | [optional] 
**UpdatedBy** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Results** | Pointer to **map[string]interface{}** | Results Summary (only returned when using query parameter results&#x3D;true) | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListSecurityScans200ResponseAllOfSecurityScansInner{
    // Set fields directly
}
```

### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


