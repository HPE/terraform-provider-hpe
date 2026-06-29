# CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitorConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Monitor** | Pointer to [**CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitorConfigMonitor**](CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitorConfigMonitor.md) |  | [optional] 
**MonitorConfig** | Pointer to **NullableString** | Free-form monitor configuration content. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitorConfig{
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


