# FastUDPLoadBalancerProfileConfig5

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProfileType** | Pointer to **string** | Profile category. Automatically set from serviceType (application-profile). | [optional] 
**FastUdpIdleTimeout** | Pointer to **NullableInt64** | Idle timeout in seconds before an inactive flow is closed. Defaults to 300. | [optional] [default to 300]
**HaFlowMirroring** | Pointer to **bool** | Whether all flows to the active member are mirrored to the standby member. | [optional] 
**Tags** | Pointer to [**[]LoadBalancerProfileTag42**](LoadBalancerProfileTag42.md) | NSX-T tags applied to the profile. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &FastUDPLoadBalancerProfileConfig5{
    // Set fields directly
}
```

### FastUdpIdleTimeout (Nullable)

Use the Nullable wrapper methods:
- `obj.FastUdpIdleTimeout.IsSet()` — check if set
- `obj.FastUdpIdleTimeout.Get()` — get the inner value (returns pointer)
- `obj.FastUdpIdleTimeout.Set(&val)` — set the value
- `obj.FastUdpIdleTimeout.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


