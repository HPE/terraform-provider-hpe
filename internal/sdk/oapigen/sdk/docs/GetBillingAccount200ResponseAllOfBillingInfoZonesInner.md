# GetBillingAccount200ResponseAllOfBillingInfoZonesInner

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
**ComputeServers** | Pointer to [**GetBillingAccount200ResponseAllOfBillingInfoZonesInnerComputeServers**](GetBillingAccount200ResponseAllOfBillingInfoZonesInnerComputeServers.md) |  | [optional] 
**Instances** | Pointer to [**GetBillingAccount200ResponseAllOfBillingInfoZonesInnerInstances**](GetBillingAccount200ResponseAllOfBillingInfoZonesInnerInstances.md) |  | [optional] 
**DiscoveredServers** | Pointer to [**GetBillingAccount200ResponseAllOfBillingInfoZonesInnerDiscoveredServers**](GetBillingAccount200ResponseAllOfBillingInfoZonesInnerDiscoveredServers.md) |  | [optional] 
**LoadBalancers** | Pointer to [**GetBillingAccount200ResponseAllOfBillingInfoZonesInnerLoadBalancers**](GetBillingAccount200ResponseAllOfBillingInfoZonesInnerLoadBalancers.md) |  | [optional] 
**VirtualImages** | Pointer to [**GetBillingAccount200ResponseAllOfBillingInfoZonesInnerVirtualImages**](GetBillingAccount200ResponseAllOfBillingInfoZonesInnerVirtualImages.md) |  | [optional] 
**Snapshots** | Pointer to [**GetBillingAccount200ResponseAllOfBillingInfoZonesInnerSnapshots**](GetBillingAccount200ResponseAllOfBillingInfoZonesInnerSnapshots.md) |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetBillingAccount200ResponseAllOfBillingInfoZonesInner{
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


