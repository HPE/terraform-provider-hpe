# BillingServersServersInnerUsagesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**ZoneName** | Pointer to **string** |  | [optional] 
**AccountName** | Pointer to **string** |  | [optional] 
**Volumes** | Pointer to [**[]BillingServersServersInnerUsagesInnerVolumesInner**](BillingServersServersInnerUsagesInnerVolumesInner.md) |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**MaxCpu** | Pointer to **NullableString** |  | [optional] 
**MaxCores** | Pointer to **int64** |  | [optional] 
**ServerExternalId** | Pointer to **NullableString** |  | [optional] 
**ServerInternalId** | Pointer to **NullableString** |  | [optional] 
**PlanName** | Pointer to **string** |  | [optional] 
**HourlyPrice** | Pointer to **float32** |  | [optional] 
**HourlyCost** | Pointer to **float32** |  | [optional] 
**Currency** | Pointer to **string** |  | [optional] 
**PricesUsed** | Pointer to [**[]BillingServersServersInnerUsagesInnerPricesUsedInner**](BillingServersServersInnerUsagesInnerPricesUsedInner.md) |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**CreatedByUser** | Pointer to **string** |  | [optional] 
**CreatedByUserId** | Pointer to **int64** |  | [optional] 
**SiteId** | Pointer to **NullableString** |  | [optional] 
**SiteName** | Pointer to **NullableString** |  | [optional] 
**SiteUUID** | Pointer to **NullableString** |  | [optional] 
**SiteCode** | Pointer to **NullableString** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Tags** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ApplicablePrices** | Pointer to [**[]BillingServersServersInnerUsagesInnerApplicablePricesInner**](BillingServersServersInnerUsagesInnerApplicablePricesInner.md) |  | [optional] 
**ServicePlanId** | Pointer to **int64** |  | [optional] 
**ServicePlanName** | Pointer to **string** |  | [optional] 
**ResourcePoolId** | Pointer to **NullableString** |  | [optional] 
**ResourcePoolName** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BillingServersServersInnerUsagesInner{
    // Set fields directly
}
```

### MaxCpu (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxCpu.IsSet()` — check if set
- `obj.MaxCpu.Get()` — get the inner value (returns pointer)
- `obj.MaxCpu.Set(&val)` — set the value
- `obj.MaxCpu.Unset()` — clear the value
### ServerExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerExternalId.IsSet()` — check if set
- `obj.ServerExternalId.Get()` — get the inner value (returns pointer)
- `obj.ServerExternalId.Set(&val)` — set the value
- `obj.ServerExternalId.Unset()` — clear the value
### ServerInternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerInternalId.IsSet()` — check if set
- `obj.ServerInternalId.Get()` — get the inner value (returns pointer)
- `obj.ServerInternalId.Set(&val)` — set the value
- `obj.ServerInternalId.Unset()` — clear the value
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
### ResourcePoolId (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourcePoolId.IsSet()` — check if set
- `obj.ResourcePoolId.Get()` — get the inner value (returns pointer)
- `obj.ResourcePoolId.Set(&val)` — set the value
- `obj.ResourcePoolId.Unset()` — clear the value
### ResourcePoolName (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourcePoolName.IsSet()` — check if set
- `obj.ResourcePoolName.Get()` — get the inner value (returns pointer)
- `obj.ResourcePoolName.Set(&val)` — set the value
- `obj.ResourcePoolName.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


