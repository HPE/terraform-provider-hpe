# ListInvoices200ResponseAllOfInvoicesInnerLineItemsInner

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
**ItemId** | Pointer to **string** |  | [optional] 
**ItemType** | Pointer to **NullableString** |  | [optional] 
**ItemName** | Pointer to **string** |  | [optional] 
**ItemDescription** | Pointer to **string** |  | [optional] 
**ProductId** | Pointer to **NullableString** |  | [optional] 
**ProductCode** | Pointer to **string** |  | [optional] 
**ProductName** | Pointer to **string** |  | [optional] 
**ItemSeller** | Pointer to **NullableString** |  | [optional] 
**ItemAction** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**RateId** | Pointer to **string** |  | [optional] 
**RateClass** | Pointer to **NullableString** |  | [optional] 
**RateUnit** | Pointer to **string** |  | [optional] 
**RateTerm** | Pointer to **NullableString** |  | [optional] 
**UsageType** | Pointer to **string** |  | [optional] 
**UsageCategory** | Pointer to **string** |  | [optional] 
**UsageService** | Pointer to **string** |  | [optional] 
**ItemUsage** | Pointer to **float32** |  | [optional] 
**ItemRate** | Pointer to **float32** |  | [optional] 
**ItemCost** | Pointer to **float32** |  | [optional] 
**ItemPrice** | Pointer to **float32** |  | [optional] 
**ItemTax** | Pointer to **int64** |  | [optional] 
**ItemTerm** | Pointer to **NullableString** |  | [optional] 
**TaxType** | Pointer to **NullableString** |  | [optional] 
**RegionCode** | Pointer to **string** |  | [optional] 
**Currency** | Pointer to **string** |  | [optional] 
**ConversionRate** | Pointer to **int64** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListInvoices200ResponseAllOfInvoicesInnerLineItemsInner{
    // Set fields directly
}
```

### ItemType (Nullable)

Use the Nullable wrapper methods:
- `obj.ItemType.IsSet()` — check if set
- `obj.ItemType.Get()` — get the inner value (returns pointer)
- `obj.ItemType.Set(&val)` — set the value
- `obj.ItemType.Unset()` — clear the value
### ProductId (Nullable)

Use the Nullable wrapper methods:
- `obj.ProductId.IsSet()` — check if set
- `obj.ProductId.Get()` — get the inner value (returns pointer)
- `obj.ProductId.Set(&val)` — set the value
- `obj.ProductId.Unset()` — clear the value
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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


