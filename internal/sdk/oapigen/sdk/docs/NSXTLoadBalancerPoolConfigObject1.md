# NSXTLoadBalancerPoolConfigObject1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ActiveMonitorPaths** | Pointer to **NullableInt64** | The ID of the active health monitor (NetworkLoadBalancerMonitor). The Options API &#x60;/api/options/nsxt/nsxtLBPoolActiveMonitor?loadBalancerId&#x3D;{id}&#x60; can be used to see which options are available. | [optional] 
**PassiveMonitorPath** | Pointer to **NullableInt64** | The ID of the passive health monitor (NetworkLoadBalancerMonitor). The Options API &#x60;/api/options/nsxt/nsxtLBPoolPassiveMonitor?loadBalancerId&#x3D;{id}&#x60; can be used to see which options are available. | [optional] 
**SnatTranslationType** | Pointer to **string** | SNAT translation type. Determines how source NAT is applied to pool traffic. | [optional] 
**SnatIpAddresses** | Pointer to **[]string** | List of SNAT IP addresses. Required when snatTranslationType is LBSnatIpPool. | [optional] 
**TcpMultiplexing** | Pointer to **bool** | Whether TCP multiplexing is enabled for the pool. | [optional] 
**TcpMultiplexingNumber** | Pointer to **NullableInt64** | Maximum number of TCP multiplexing connections. Defaults to 6. | [optional] 
**MemberGroup** | Pointer to [**NSXTLoadBalancerPoolConfigObject1MemberGroup**](NSXTLoadBalancerPoolConfigObject1MemberGroup.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NSXTLoadBalancerPoolConfigObject1{
    // Set fields directly
}
```

### ActiveMonitorPaths (Nullable)

Use the Nullable wrapper methods:
- `obj.ActiveMonitorPaths.IsSet()` — check if set
- `obj.ActiveMonitorPaths.Get()` — get the inner value (returns pointer)
- `obj.ActiveMonitorPaths.Set(&val)` — set the value
- `obj.ActiveMonitorPaths.Unset()` — clear the value
### PassiveMonitorPath (Nullable)

Use the Nullable wrapper methods:
- `obj.PassiveMonitorPath.IsSet()` — check if set
- `obj.PassiveMonitorPath.Get()` — get the inner value (returns pointer)
- `obj.PassiveMonitorPath.Set(&val)` — set the value
- `obj.PassiveMonitorPath.Unset()` — clear the value
### TcpMultiplexingNumber (Nullable)

Use the Nullable wrapper methods:
- `obj.TcpMultiplexingNumber.IsSet()` — check if set
- `obj.TcpMultiplexingNumber.Get()` — get the inner value (returns pointer)
- `obj.TcpMultiplexingNumber.Set(&val)` — set the value
- `obj.TcpMultiplexingNumber.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


