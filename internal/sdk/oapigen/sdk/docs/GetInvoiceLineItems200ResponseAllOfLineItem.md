# GetInvoiceLineItems200ResponseAllOfLineItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**InvoiceId** | Pointer to **int64** |  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**RefName** | Pointer to **string** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**ItemId** | Pointer to **NullableString** |  | [optional] 
**ItemType** | Pointer to **NullableString** |  | [optional] 
**ItemName** | Pointer to **NullableString** |  | [optional] 
**ItemDescription** | Pointer to **NullableString** |  | [optional] 
**ProductId** | Pointer to **NullableString** |  | [optional] 
**ProductCode** | Pointer to **NullableString** |  | [optional] 
**ProductName** | Pointer to **NullableString** |  | [optional] 
**ItemSeller** | Pointer to **NullableString** |  | [optional] 
**ItemAction** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**RateId** | Pointer to **NullableString** |  | [optional] 
**RateClass** | Pointer to **NullableString** |  | [optional] 
**RateUnit** | Pointer to **string** |  | [optional] 
**RateTerm** | Pointer to **NullableString** |  | [optional] 
**UsageType** | Pointer to **string** |  | [optional] 
**UsageCategory** | Pointer to **string** |  | [optional] 
**UsageService** | Pointer to **NullableString** |  | [optional] 
**ItemUsage** | Pointer to **int64** |  | [optional] 
**ItemRate** | Pointer to **float32** |  | [optional] 
**ItemCost** | Pointer to **float32** |  | [optional] 
**ItemPriceRate** | Pointer to **float32** |  | [optional] 
**ItemPrice** | Pointer to **float32** |  | [optional] 
**ItemTax** | Pointer to **int64** |  | [optional] 
**ItemTerm** | Pointer to **NullableString** |  | [optional] 
**TaxType** | Pointer to **NullableString** |  | [optional] 
**RegionCode** | Pointer to **NullableString** |  | [optional] 
**Currency** | Pointer to **string** |  | [optional] 
**ConversionRate** | Pointer to **int64** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetInvoiceLineItems200ResponseAllOfLineItem{
    // Set fields directly
}
```

### ItemId (Nullable)

Use the Nullable wrapper methods:
- `obj.ItemId.IsSet()` — check if set
- `obj.ItemId.Get()` — get the inner value (returns pointer)
- `obj.ItemId.Set(&val)` — set the value
- `obj.ItemId.Unset()` — clear the value
### ItemType (Nullable)

Use the Nullable wrapper methods:
- `obj.ItemType.IsSet()` — check if set
- `obj.ItemType.Get()` — get the inner value (returns pointer)
- `obj.ItemType.Set(&val)` — set the value
- `obj.ItemType.Unset()` — clear the value
### ItemName (Nullable)

Use the Nullable wrapper methods:
- `obj.ItemName.IsSet()` — check if set
- `obj.ItemName.Get()` — get the inner value (returns pointer)
- `obj.ItemName.Set(&val)` — set the value
- `obj.ItemName.Unset()` — clear the value
### ItemDescription (Nullable)

Use the Nullable wrapper methods:
- `obj.ItemDescription.IsSet()` — check if set
- `obj.ItemDescription.Get()` — get the inner value (returns pointer)
- `obj.ItemDescription.Set(&val)` — set the value
- `obj.ItemDescription.Unset()` — clear the value
### ProductId (Nullable)

Use the Nullable wrapper methods:
- `obj.ProductId.IsSet()` — check if set
- `obj.ProductId.Get()` — get the inner value (returns pointer)
- `obj.ProductId.Set(&val)` — set the value
- `obj.ProductId.Unset()` — clear the value
### ProductCode (Nullable)

Use the Nullable wrapper methods:
- `obj.ProductCode.IsSet()` — check if set
- `obj.ProductCode.Get()` — get the inner value (returns pointer)
- `obj.ProductCode.Set(&val)` — set the value
- `obj.ProductCode.Unset()` — clear the value
### ProductName (Nullable)

Use the Nullable wrapper methods:
- `obj.ProductName.IsSet()` — check if set
- `obj.ProductName.Get()` — get the inner value (returns pointer)
- `obj.ProductName.Set(&val)` — set the value
- `obj.ProductName.Unset()` — clear the value
### ItemSeller (Nullable)

Use the Nullable wrapper methods:
- `obj.ItemSeller.IsSet()` — check if set
- `obj.ItemSeller.Get()` — get the inner value (returns pointer)
- `obj.ItemSeller.Set(&val)` — set the value
- `obj.ItemSeller.Unset()` — clear the value
### ItemAction (Nullable)

Use the Nullable wrapper methods:
- `obj.ItemAction.IsSet()` — check if set
- `obj.ItemAction.Get()` — get the inner value (returns pointer)
- `obj.ItemAction.Set(&val)` — set the value
- `obj.ItemAction.Unset()` — clear the value
### RateId (Nullable)

Use the Nullable wrapper methods:
- `obj.RateId.IsSet()` — check if set
- `obj.RateId.Get()` — get the inner value (returns pointer)
- `obj.RateId.Set(&val)` — set the value
- `obj.RateId.Unset()` — clear the value
### RateClass (Nullable)

Use the Nullable wrapper methods:
- `obj.RateClass.IsSet()` — check if set
- `obj.RateClass.Get()` — get the inner value (returns pointer)
- `obj.RateClass.Set(&val)` — set the value
- `obj.RateClass.Unset()` — clear the value
### RateTerm (Nullable)

Use the Nullable wrapper methods:
- `obj.RateTerm.IsSet()` — check if set
- `obj.RateTerm.Get()` — get the inner value (returns pointer)
- `obj.RateTerm.Set(&val)` — set the value
- `obj.RateTerm.Unset()` — clear the value
### UsageService (Nullable)

Use the Nullable wrapper methods:
- `obj.UsageService.IsSet()` — check if set
- `obj.UsageService.Get()` — get the inner value (returns pointer)
- `obj.UsageService.Set(&val)` — set the value
- `obj.UsageService.Unset()` — clear the value
### ItemTerm (Nullable)

Use the Nullable wrapper methods:
- `obj.ItemTerm.IsSet()` — check if set
- `obj.ItemTerm.Get()` — get the inner value (returns pointer)
- `obj.ItemTerm.Set(&val)` — set the value
- `obj.ItemTerm.Unset()` — clear the value
### TaxType (Nullable)

Use the Nullable wrapper methods:
- `obj.TaxType.IsSet()` — check if set
- `obj.TaxType.Get()` — get the inner value (returns pointer)
- `obj.TaxType.Set(&val)` — set the value
- `obj.TaxType.Unset()` — clear the value
### RegionCode (Nullable)

Use the Nullable wrapper methods:
- `obj.RegionCode.IsSet()` — check if set
- `obj.RegionCode.Get()` — get the inner value (returns pointer)
- `obj.RegionCode.Set(&val)` — set the value
- `obj.RegionCode.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


