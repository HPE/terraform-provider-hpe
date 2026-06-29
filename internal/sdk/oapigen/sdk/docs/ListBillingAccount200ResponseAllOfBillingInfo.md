# ListBillingAccount200ResponseAllOfBillingInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccountId** | Pointer to **int64** |  | [optional] 
**AccountUUID** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**PriceUnit** | Pointer to **string** |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 
**Zones** | Pointer to [**[]ListBillingAccount200ResponseAllOfBillingInfoZonesInner**](ListBillingAccount200ResponseAllOfBillingInfoZonesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListBillingAccount200ResponseAllOfBillingInfo{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


