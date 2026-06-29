# BillingZonesInner

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
**ComputeServers** | Pointer to [**BillingZonesInnerComputeServers**](BillingZonesInnerComputeServers.md) |  | [optional] 
**Instances** | Pointer to [**BillingZonesInnerInstances**](BillingZonesInnerInstances.md) |  | [optional] 
**DiscoveredServers** | Pointer to [**BillingZonesInnerDiscoveredServers**](BillingZonesInnerDiscoveredServers.md) |  | [optional] 
**LoadBalancers** | Pointer to [**BillingZonesInnerLoadBalancers**](BillingZonesInnerLoadBalancers.md) |  | [optional] 
**VirtualImages** | Pointer to [**BillingZonesInnerVirtualImages**](BillingZonesInnerVirtualImages.md) |  | [optional] 
**Snapshots** | Pointer to [**BillingZonesInnerSnapshots**](BillingZonesInnerSnapshots.md) |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BillingZonesInner{
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


