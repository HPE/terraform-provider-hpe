# BillingServerUsagesInnerApplicablePricesInnerPricesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | Pointer to **string** |  | [optional] 
**PricePerUnit** | Pointer to **float32** |  | [optional] 
**CostPerUnit** | Pointer to **float32** |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**Quantity** | Pointer to **int64** |  | [optional] 
**DatastoreId** | Pointer to **NullableString** |  | [optional] 
**VolumeType** | Pointer to **string** |  | [optional] 
**Datastore** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BillingServerUsagesInnerApplicablePricesInnerPricesInner{
    // Set fields directly
}
```

### DatastoreId (Nullable)

Use the Nullable wrapper methods:
- `obj.DatastoreId.IsSet()` — check if set
- `obj.DatastoreId.Get()` — get the inner value (returns pointer)
- `obj.DatastoreId.Set(&val)` — set the value
- `obj.DatastoreId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


