# LoadBalancerInstanceUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**VipName** | Pointer to **string** | VIP Name | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**VipAddress** | Pointer to **string** | VIP Address | [optional] 
**VipPort** | Pointer to **int64** | VIP Port | [optional] 
**VipProtocol** | Pointer to **string** | VIP Protocol. For NSX-T load balancers, this is the virtual server type. Valid values are &#x60;http&#x60;, &#x60;tcp&#x60;, &#x60;udp&#x60;. | [optional] 
**VipHostname** | Pointer to **string** | VIP Hostname | [optional] 
**VipPool** | Pointer to **NullableInt64** | Network Pool ID for automatic VIP address allocation. When set, a VIP address will be leased from this pool and &#x60;vipAddress&#x60; does not need to be specified. | [optional] 
**Pool** | Pointer to **int64** | Backend server pool ID. Can also be specified in &#x60;config.pool&#x60; for NSX-T load balancers. | [optional] 
**SslCert** | Pointer to **int64** | SSL Client Certificate ID. Use &#x60;0&#x60; for none. | [optional] 
**SslServerCert** | Pointer to **int64** | SSL Server Certificate ID. Use &#x60;0&#x60; for none. | [optional] 
**Config** | Pointer to [**LoadBalancerInstanceUpdateConfig**](LoadBalancerInstanceUpdateConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &LoadBalancerInstanceUpdate{
    // Set fields directly
}
```

### VipPool (Nullable)

Use the Nullable wrapper methods:
- `obj.VipPool.IsSet()` — check if set
- `obj.VipPool.Get()` — get the inner value (returns pointer)
- `obj.VipPool.Set(&val)` — set the value
- `obj.VipPool.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


