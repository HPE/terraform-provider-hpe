# ListBillingZone200ResponseAllOfBillingInfoZonesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ZoneName** | Pointer to **string** |  | [optional] 
**ZoneId** | Pointer to **int64** |  | [optional] 
**ZoneUUID** | Pointer to **string** |  | [optional] 
**ZoneCode** | Pointer to **NullableString** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**PriceUnit** | Pointer to **string** |  | [optional] 
**ComputeServers** | Pointer to [**ListBillingZone200ResponseAllOfBillingInfoZonesInnerComputeServers**](ListBillingZone200ResponseAllOfBillingInfoZonesInnerComputeServers.md) |  | [optional] 
**Instances** | Pointer to [**ListBillingZone200ResponseAllOfBillingInfoZonesInnerInstances**](ListBillingZone200ResponseAllOfBillingInfoZonesInnerInstances.md) |  | [optional] 
**DiscoveredServers** | Pointer to [**ListBillingZone200ResponseAllOfBillingInfoZonesInnerDiscoveredServers**](ListBillingZone200ResponseAllOfBillingInfoZonesInnerDiscoveredServers.md) |  | [optional] 
**LoadBalancers** | Pointer to [**ListBillingZone200ResponseAllOfBillingInfoZonesInnerLoadBalancers**](ListBillingZone200ResponseAllOfBillingInfoZonesInnerLoadBalancers.md) |  | [optional] 
**VirtualImages** | Pointer to [**ListBillingZone200ResponseAllOfBillingInfoZonesInnerVirtualImages**](ListBillingZone200ResponseAllOfBillingInfoZonesInnerVirtualImages.md) |  | [optional] 
**Snapshots** | Pointer to [**ListBillingZone200ResponseAllOfBillingInfoZonesInnerSnapshots**](ListBillingZone200ResponseAllOfBillingInfoZonesInnerSnapshots.md) |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListBillingZone200ResponseAllOfBillingInfoZonesInner{
    // Set fields directly
}
```

### ZoneCode (Nullable)

Use the Nullable wrapper methods:
- `obj.ZoneCode.IsSet()` — check if set
- `obj.ZoneCode.Get()` — get the inner value (returns pointer)
- `obj.ZoneCode.Set(&val)` — set the value
- `obj.ZoneCode.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


