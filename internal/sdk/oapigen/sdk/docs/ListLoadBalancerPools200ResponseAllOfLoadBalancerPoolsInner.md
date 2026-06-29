# ListLoadBalancerPools200ResponseAllOfLoadBalancerPoolsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**LoadBalancer** | Pointer to [**ListLoadBalancerPools200ResponseAllOfLoadBalancerPoolsInnerLoadBalancer**](ListLoadBalancerPools200ResponseAllOfLoadBalancerPoolsInnerLoadBalancer.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**VipSticky** | Pointer to **NullableString** |  | [optional] 
**VipBalance** | Pointer to **string** |  | [optional] 
**AllowNat** | Pointer to **NullableString** |  | [optional] 
**AllowSnat** | Pointer to **NullableString** |  | [optional] 
**VipClientIpMode** | Pointer to **NullableString** |  | [optional] 
**VipServerIpMode** | Pointer to **NullableString** |  | [optional] 
**MinActive** | Pointer to **int64** |  | [optional] 
**MinInService** | Pointer to **NullableString** |  | [optional] 
**MinUpMonitor** | Pointer to **NullableString** |  | [optional] 
**MinUpAction** | Pointer to **NullableString** |  | [optional] 
**MaxQueueDepth** | Pointer to **NullableString** |  | [optional] 
**MaxQueueTime** | Pointer to **NullableString** |  | [optional] 
**NumberActive** | Pointer to **int64** |  | [optional] 
**NumberInService** | Pointer to **int64** |  | [optional] 
**HealthScore** | Pointer to **int64** |  | [optional] 
**PerformanceScore** | Pointer to **int64** |  | [optional] 
**HealthPenalty** | Pointer to **int64** |  | [optional] 
**SecurityPenalty** | Pointer to **int64** |  | [optional] 
**ErrorPenalty** | Pointer to **int64** |  | [optional] 
**DownAction** | Pointer to **NullableString** |  | [optional] 
**RampTime** | Pointer to **NullableString** |  | [optional] 
**Port** | Pointer to **NullableString** |  | [optional] 
**PortType** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Nodes** | Pointer to [**[]ListLoadBalancerPools200ResponseAllOfLoadBalancerPoolsInnerNodesInner**](ListLoadBalancerPools200ResponseAllOfLoadBalancerPoolsInnerNodesInner.md) |  | [optional] 
**Monitors** | Pointer to [**[]ListLoadBalancerPools200ResponseAllOfLoadBalancerPoolsInnerMonitorsInner**](ListLoadBalancerPools200ResponseAllOfLoadBalancerPoolsInnerMonitorsInner.md) |  | [optional] 
**Members** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**CreatedBy** | Pointer to [**ListLoadBalancerPools200ResponseAllOfLoadBalancerPoolsInnerCreatedBy**](ListLoadBalancerPools200ResponseAllOfLoadBalancerPoolsInnerCreatedBy.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListLoadBalancerPools200ResponseAllOfLoadBalancerPoolsInner{
    // Set fields directly
}
```

### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### VipSticky (Nullable)

Use the Nullable wrapper methods:
- `obj.VipSticky.IsSet()` — check if set
- `obj.VipSticky.Get()` — get the inner value (returns pointer)
- `obj.VipSticky.Set(&val)` — set the value
- `obj.VipSticky.Unset()` — clear the value
### AllowNat (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowNat.IsSet()` — check if set
- `obj.AllowNat.Get()` — get the inner value (returns pointer)
- `obj.AllowNat.Set(&val)` — set the value
- `obj.AllowNat.Unset()` — clear the value
### AllowSnat (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowSnat.IsSet()` — check if set
- `obj.AllowSnat.Get()` — get the inner value (returns pointer)
- `obj.AllowSnat.Set(&val)` — set the value
- `obj.AllowSnat.Unset()` — clear the value
### VipClientIpMode (Nullable)

Use the Nullable wrapper methods:
- `obj.VipClientIpMode.IsSet()` — check if set
- `obj.VipClientIpMode.Get()` — get the inner value (returns pointer)
- `obj.VipClientIpMode.Set(&val)` — set the value
- `obj.VipClientIpMode.Unset()` — clear the value
### VipServerIpMode (Nullable)

Use the Nullable wrapper methods:
- `obj.VipServerIpMode.IsSet()` — check if set
- `obj.VipServerIpMode.Get()` — get the inner value (returns pointer)
- `obj.VipServerIpMode.Set(&val)` — set the value
- `obj.VipServerIpMode.Unset()` — clear the value
### MinInService (Nullable)

Use the Nullable wrapper methods:
- `obj.MinInService.IsSet()` — check if set
- `obj.MinInService.Get()` — get the inner value (returns pointer)
- `obj.MinInService.Set(&val)` — set the value
- `obj.MinInService.Unset()` — clear the value
### MinUpMonitor (Nullable)

Use the Nullable wrapper methods:
- `obj.MinUpMonitor.IsSet()` — check if set
- `obj.MinUpMonitor.Get()` — get the inner value (returns pointer)
- `obj.MinUpMonitor.Set(&val)` — set the value
- `obj.MinUpMonitor.Unset()` — clear the value
### MinUpAction (Nullable)

Use the Nullable wrapper methods:
- `obj.MinUpAction.IsSet()` — check if set
- `obj.MinUpAction.Get()` — get the inner value (returns pointer)
- `obj.MinUpAction.Set(&val)` — set the value
- `obj.MinUpAction.Unset()` — clear the value
### MaxQueueDepth (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxQueueDepth.IsSet()` — check if set
- `obj.MaxQueueDepth.Get()` — get the inner value (returns pointer)
- `obj.MaxQueueDepth.Set(&val)` — set the value
- `obj.MaxQueueDepth.Unset()` — clear the value
### MaxQueueTime (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxQueueTime.IsSet()` — check if set
- `obj.MaxQueueTime.Get()` — get the inner value (returns pointer)
- `obj.MaxQueueTime.Set(&val)` — set the value
- `obj.MaxQueueTime.Unset()` — clear the value
### DownAction (Nullable)

Use the Nullable wrapper methods:
- `obj.DownAction.IsSet()` — check if set
- `obj.DownAction.Get()` — get the inner value (returns pointer)
- `obj.DownAction.Set(&val)` — set the value
- `obj.DownAction.Unset()` — clear the value
### RampTime (Nullable)

Use the Nullable wrapper methods:
- `obj.RampTime.IsSet()` — check if set
- `obj.RampTime.Get()` — get the inner value (returns pointer)
- `obj.RampTime.Set(&val)` — set the value
- `obj.RampTime.Unset()` — clear the value
### Port (Nullable)

Use the Nullable wrapper methods:
- `obj.Port.IsSet()` — check if set
- `obj.Port.Get()` — get the inner value (returns pointer)
- `obj.Port.Set(&val)` — set the value
- `obj.Port.Unset()` — clear the value
### PortType (Nullable)

Use the Nullable wrapper methods:
- `obj.PortType.IsSet()` — check if set
- `obj.PortType.Get()` — get the inner value (returns pointer)
- `obj.PortType.Set(&val)` — set the value
- `obj.PortType.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


