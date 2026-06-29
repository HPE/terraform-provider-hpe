# Budget

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Account** | Pointer to [**BudgetAccount**](BudgetAccount.md) |  | [optional] 
**ForecastType** | Pointer to [**BudgetForecastType**](BudgetForecastType.md) |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**RefScope** | Pointer to **string** |  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**RefName** | Pointer to **string** |  | [optional] 
**Period** | Pointer to **string** |  | [optional] 
**Year** | Pointer to **string** |  | [optional] 
**ResourceType** | Pointer to **string** |  | [optional] 
**Timezone** | Pointer to **string** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**Interval** | Pointer to **string** |  | [optional] 
**Costs** | Pointer to **[]int64** |  | [optional] 
**IsFiscal** | Pointer to **bool** |  | [optional] 
**AverageCost** | Pointer to **int64** |  | [optional] 
**TotalCost** | Pointer to **int64** |  | [optional] 
**Currency** | Pointer to **string** |  | [optional] 
**Rollover** | Pointer to **bool** |  | [optional] 
**WarningLimit** | Pointer to **NullableString** |  | [optional] 
**OverLimit** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**CreatedById** | Pointer to **int64** |  | [optional] 
**CreatedByName** | Pointer to **string** |  | [optional] 
**UpdatedById** | Pointer to **NullableString** |  | [optional] 
**UpdatedByName** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Stats** | Pointer to [**BudgetStats**](BudgetStats.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &Budget{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### WarningLimit (Nullable)

Use the Nullable wrapper methods:
- `obj.WarningLimit.IsSet()` — check if set
- `obj.WarningLimit.Get()` — get the inner value (returns pointer)
- `obj.WarningLimit.Set(&val)` — set the value
- `obj.WarningLimit.Unset()` — clear the value
### OverLimit (Nullable)

Use the Nullable wrapper methods:
- `obj.OverLimit.IsSet()` — check if set
- `obj.OverLimit.Get()` — get the inner value (returns pointer)
- `obj.OverLimit.Set(&val)` — set the value
- `obj.OverLimit.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### UpdatedById (Nullable)

Use the Nullable wrapper methods:
- `obj.UpdatedById.IsSet()` — check if set
- `obj.UpdatedById.Get()` — get the inner value (returns pointer)
- `obj.UpdatedById.Set(&val)` — set the value
- `obj.UpdatedById.Unset()` — clear the value
### UpdatedByName (Nullable)

Use the Nullable wrapper methods:
- `obj.UpdatedByName.IsSet()` — check if set
- `obj.UpdatedByName.Get()` — get the inner value (returns pointer)
- `obj.UpdatedByName.Set(&val)` — set the value
- `obj.UpdatedByName.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


