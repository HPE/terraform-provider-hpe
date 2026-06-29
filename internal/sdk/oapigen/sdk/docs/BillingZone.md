# BillingZone

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
**ComputeServers** | Pointer to [**GetBillingZoneIdentifier200ResponseAllOfBillingInfoComputeServers**](GetBillingZoneIdentifier200ResponseAllOfBillingInfoComputeServers.md) |  | [optional] 
**Instances** | Pointer to [**GetBillingZoneIdentifier200ResponseAllOfBillingInfoInstances**](GetBillingZoneIdentifier200ResponseAllOfBillingInfoInstances.md) |  | [optional] 
**DiscoveredServers** | Pointer to [**GetBillingZoneIdentifier200ResponseAllOfBillingInfoDiscoveredServers**](GetBillingZoneIdentifier200ResponseAllOfBillingInfoDiscoveredServers.md) |  | [optional] 
**LoadBalancers** | Pointer to [**GetBillingZoneIdentifier200ResponseAllOfBillingInfoLoadBalancers**](GetBillingZoneIdentifier200ResponseAllOfBillingInfoLoadBalancers.md) |  | [optional] 
**VirtualImages** | Pointer to [**GetBillingZoneIdentifier200ResponseAllOfBillingInfoVirtualImages**](GetBillingZoneIdentifier200ResponseAllOfBillingInfoVirtualImages.md) |  | [optional] 
**Snapshots** | Pointer to [**GetBillingZoneIdentifier200ResponseAllOfBillingInfoSnapshots**](GetBillingZoneIdentifier200ResponseAllOfBillingInfoSnapshots.md) |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BillingZone{
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


