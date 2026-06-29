# GetLoadBalancerMonitor200ResponseLoadBalancerMonitorConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Monitor** | Pointer to [**GetLoadBalancerMonitor200ResponseLoadBalancerMonitorConfigMonitor**](GetLoadBalancerMonitor200ResponseLoadBalancerMonitorConfigMonitor.md) |  | [optional] 
**MonitorConfig** | Pointer to **NullableString** | Free-form monitor configuration content. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetLoadBalancerMonitor200ResponseLoadBalancerMonitorConfig{
    // Set fields directly
}
```

### MonitorConfig (Nullable)

Use the Nullable wrapper methods:
- `obj.MonitorConfig.IsSet()` — check if set
- `obj.MonitorConfig.Get()` — get the inner value (returns pointer)
- `obj.MonitorConfig.Set(&val)` — set the value
- `obj.MonitorConfig.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


