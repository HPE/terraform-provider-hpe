# ListBillingAccount200ResponseAllOfBillingInfoZonesInner

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
**ComputeServers** | Pointer to [**ListBillingAccount200ResponseAllOfBillingInfoZonesInnerComputeServers**](ListBillingAccount200ResponseAllOfBillingInfoZonesInnerComputeServers.md) |  | [optional] 
**Instances** | Pointer to [**ListBillingAccount200ResponseAllOfBillingInfoZonesInnerInstances**](ListBillingAccount200ResponseAllOfBillingInfoZonesInnerInstances.md) |  | [optional] 
**DiscoveredServers** | Pointer to [**ListBillingAccount200ResponseAllOfBillingInfoZonesInnerDiscoveredServers**](ListBillingAccount200ResponseAllOfBillingInfoZonesInnerDiscoveredServers.md) |  | [optional] 
**LoadBalancers** | Pointer to [**ListBillingAccount200ResponseAllOfBillingInfoZonesInnerLoadBalancers**](ListBillingAccount200ResponseAllOfBillingInfoZonesInnerLoadBalancers.md) |  | [optional] 
**VirtualImages** | Pointer to [**ListBillingAccount200ResponseAllOfBillingInfoZonesInnerVirtualImages**](ListBillingAccount200ResponseAllOfBillingInfoZonesInnerVirtualImages.md) |  | [optional] 
**Snapshots** | Pointer to [**ListBillingAccount200ResponseAllOfBillingInfoZonesInnerSnapshots**](ListBillingAccount200ResponseAllOfBillingInfoZonesInnerSnapshots.md) |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListBillingAccount200ResponseAllOfBillingInfoZonesInner{
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


