# LoadBalancerProfileConfigFastTcp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProfileType** | Pointer to **string** | Profile category. Automatically set from serviceType (application-profile). | [optional] 
**FastTcpIdleTimeout** | Pointer to **NullableInt64** | Idle timeout in seconds before an inactive connection is closed. Defaults to 1800. | [optional] [default to 1800]
**HaFlowMirroring** | Pointer to **bool** | Whether all flows to the active member are mirrored to the standby member. | [optional] 
**ConnectionCloseTimeout** | Pointer to **NullableInt64** | Timeout in seconds before a closed connection is removed. Defaults to 8. Range 1-60. | [optional] [default to 8]
**Tags** | Pointer to [**[]LoadBalancerProfileTag49**](LoadBalancerProfileTag49.md) | NSX-T tags applied to the profile. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &LoadBalancerProfileConfigFastTcp{
    // Set fields directly
}
```

### FastTcpIdleTimeout (Nullable)

Use the Nullable wrapper methods:
- `obj.FastTcpIdleTimeout.IsSet()` — check if set
- `obj.FastTcpIdleTimeout.Get()` — get the inner value (returns pointer)
- `obj.FastTcpIdleTimeout.Set(&val)` — set the value
- `obj.FastTcpIdleTimeout.Unset()` — clear the value
### ConnectionCloseTimeout (Nullable)

Use the Nullable wrapper methods:
- `obj.ConnectionCloseTimeout.IsSet()` — check if set
- `obj.ConnectionCloseTimeout.Get()` — get the inner value (returns pointer)
- `obj.ConnectionCloseTimeout.Set(&val)` — set the value
- `obj.ConnectionCloseTimeout.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


