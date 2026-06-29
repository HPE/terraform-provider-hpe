# BillingServerUsagesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cost** | Pointer to **float32** |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**CreatedByUser** | Pointer to **string** |  | [optional] 
**CreatedByUserId** | Pointer to **int64** |  | [optional] 
**SiteId** | Pointer to **NullableString** |  | [optional] 
**SiteName** | Pointer to **NullableString** |  | [optional] 
**SiteUUID** | Pointer to **NullableString** |  | [optional] 
**SiteCode** | Pointer to **NullableString** |  | [optional] 
**Currency** | Pointer to **string** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Tags** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ApplicablePrices** | Pointer to [**[]BillingServerUsagesInnerApplicablePricesInner**](BillingServerUsagesInnerApplicablePricesInner.md) |  | [optional] 
**ServicePlanId** | Pointer to **int64** |  | [optional] 
**ServicePlanName** | Pointer to **string** |  | [optional] 
**ResourcePoolId** | Pointer to **int64** |  | [optional] 
**ResourcePoolName** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BillingServerUsagesInner{
    // Set fields directly
}
```

### SiteId (Nullable)

Use the Nullable wrapper methods:
- `obj.SiteId.IsSet()` — check if set
- `obj.SiteId.Get()` — get the inner value (returns pointer)
- `obj.SiteId.Set(&val)` — set the value
- `obj.SiteId.Unset()` — clear the value
### SiteName (Nullable)

Use the Nullable wrapper methods:
- `obj.SiteName.IsSet()` — check if set
- `obj.SiteName.Get()` — get the inner value (returns pointer)
- `obj.SiteName.Set(&val)` — set the value
- `obj.SiteName.Unset()` — clear the value
### SiteUUID (Nullable)

Use the Nullable wrapper methods:
- `obj.SiteUUID.IsSet()` — check if set
- `obj.SiteUUID.Get()` — get the inner value (returns pointer)
- `obj.SiteUUID.Set(&val)` — set the value
- `obj.SiteUUID.Unset()` — clear the value
### SiteCode (Nullable)

Use the Nullable wrapper methods:
- `obj.SiteCode.IsSet()` — check if set
- `obj.SiteCode.Get()` — get the inner value (returns pointer)
- `obj.SiteCode.Set(&val)` — set the value
- `obj.SiteCode.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


