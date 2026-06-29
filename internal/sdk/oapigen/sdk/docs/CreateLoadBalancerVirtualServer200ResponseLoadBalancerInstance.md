# CreateLoadBalancerVirtualServer200ResponseLoadBalancerInstance

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**LoadBalancer** | Pointer to [**CreateLoadBalancerVirtualServer200ResponseLoadBalancerInstanceLoadBalancer**](CreateLoadBalancerVirtualServer200ResponseLoadBalancerInstanceLoadBalancer.md) |  | [optional] 
**Instance** | Pointer to **NullableString** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** | Configuration object with parameters that vary by type. | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Sticky** | Pointer to **bool** |  | [optional] 
**SslEnabled** | Pointer to **NullableString** |  | [optional] 
**ExternalAddress** | Pointer to **bool** |  | [optional] 
**BackendPort** | Pointer to **NullableString** |  | [optional] 
**VipType** | Pointer to **NullableString** |  | [optional] 
**VipAddress** | Pointer to **string** |  | [optional] 
**VipHostname** | Pointer to **NullableString** |  | [optional] 
**VipProtocol** | Pointer to **string** |  | [optional] 
**VipScheme** | Pointer to **NullableString** |  | [optional] 
**VipMode** | Pointer to **NullableString** |  | [optional] 
**VipName** | Pointer to **string** |  | [optional] 
**VipPort** | Pointer to **int64** |  | [optional] 
**VipPool** | Pointer to **map[string]interface{}** | Network Pool used for VIP address allocation, if set. | [optional] 
**VipSticky** | Pointer to **NullableString** |  | [optional] 
**VipBalance** | Pointer to **NullableString** |  | [optional] 
**ServicePort** | Pointer to **NullableString** |  | [optional] 
**SourceAddress** | Pointer to **NullableString** |  | [optional] 
**SslCert** | Pointer to [**CreateLoadBalancerVirtualServer200ResponseLoadBalancerInstanceSslCert**](CreateLoadBalancerVirtualServer200ResponseLoadBalancerInstanceSslCert.md) |  | [optional] 
**SslServerCert** | Pointer to [**CreateLoadBalancerVirtualServer200ResponseLoadBalancerInstanceSslServerCert**](CreateLoadBalancerVirtualServer200ResponseLoadBalancerInstanceSslServerCert.md) |  | [optional] 
**SslMode** | Pointer to **NullableString** |  | [optional] 
**SslRedirectMode** | Pointer to **NullableString** |  | [optional] 
**VipShared** | Pointer to **bool** |  | [optional] 
**VipDirectAddress** | Pointer to **NullableString** |  | [optional] 
**ServerName** | Pointer to **NullableString** |  | [optional] 
**PoolName** | Pointer to **NullableString** |  | [optional] 
**Removing** | Pointer to **bool** |  | [optional] 
**Pool** | Pointer to [**CreateLoadBalancerVirtualServer200ResponseLoadBalancerInstancePool**](CreateLoadBalancerVirtualServer200ResponseLoadBalancerInstancePool.md) |  | [optional] 
**VipSource** | Pointer to **string** |  | [optional] 
**ExtraConfig** | Pointer to **NullableString** |  | [optional] 
**ServiceAccess** | Pointer to **NullableString** |  | [optional] 
**NetworkId** | Pointer to **NullableString** |  | [optional] 
**SubnetId** | Pointer to **NullableString** |  | [optional] 
**ExternalPortId** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**VipStatus** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateLoadBalancerVirtualServer200ResponseLoadBalancerInstance{
    // Set fields directly
}
```

### Instance (Nullable)

Use the Nullable wrapper methods:
- `obj.Instance.IsSet()` — check if set
- `obj.Instance.Get()` — get the inner value (returns pointer)
- `obj.Instance.Set(&val)` — set the value
- `obj.Instance.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### SslEnabled (Nullable)

Use the Nullable wrapper methods:
- `obj.SslEnabled.IsSet()` — check if set
- `obj.SslEnabled.Get()` — get the inner value (returns pointer)
- `obj.SslEnabled.Set(&val)` — set the value
- `obj.SslEnabled.Unset()` — clear the value
### BackendPort (Nullable)

Use the Nullable wrapper methods:
- `obj.BackendPort.IsSet()` — check if set
- `obj.BackendPort.Get()` — get the inner value (returns pointer)
- `obj.BackendPort.Set(&val)` — set the value
- `obj.BackendPort.Unset()` — clear the value
### VipType (Nullable)

Use the Nullable wrapper methods:
- `obj.VipType.IsSet()` — check if set
- `obj.VipType.Get()` — get the inner value (returns pointer)
- `obj.VipType.Set(&val)` — set the value
- `obj.VipType.Unset()` — clear the value
### VipHostname (Nullable)

Use the Nullable wrapper methods:
- `obj.VipHostname.IsSet()` — check if set
- `obj.VipHostname.Get()` — get the inner value (returns pointer)
- `obj.VipHostname.Set(&val)` — set the value
- `obj.VipHostname.Unset()` — clear the value
### VipScheme (Nullable)

Use the Nullable wrapper methods:
- `obj.VipScheme.IsSet()` — check if set
- `obj.VipScheme.Get()` — get the inner value (returns pointer)
- `obj.VipScheme.Set(&val)` — set the value
- `obj.VipScheme.Unset()` — clear the value
### VipMode (Nullable)

Use the Nullable wrapper methods:
- `obj.VipMode.IsSet()` — check if set
- `obj.VipMode.Get()` — get the inner value (returns pointer)
- `obj.VipMode.Set(&val)` — set the value
- `obj.VipMode.Unset()` — clear the value
### VipPool (Nullable)

Use the Nullable wrapper methods:
- `obj.VipPool.IsSet()` — check if set
- `obj.VipPool.Get()` — get the inner value (returns pointer)
- `obj.VipPool.Set(&val)` — set the value
- `obj.VipPool.Unset()` — clear the value
### VipSticky (Nullable)

Use the Nullable wrapper methods:
- `obj.VipSticky.IsSet()` — check if set
- `obj.VipSticky.Get()` — get the inner value (returns pointer)
- `obj.VipSticky.Set(&val)` — set the value
- `obj.VipSticky.Unset()` — clear the value
### VipBalance (Nullable)

Use the Nullable wrapper methods:
- `obj.VipBalance.IsSet()` — check if set
- `obj.VipBalance.Get()` — get the inner value (returns pointer)
- `obj.VipBalance.Set(&val)` — set the value
- `obj.VipBalance.Unset()` — clear the value
### ServicePort (Nullable)

Use the Nullable wrapper methods:
- `obj.ServicePort.IsSet()` — check if set
- `obj.ServicePort.Get()` — get the inner value (returns pointer)
- `obj.ServicePort.Set(&val)` — set the value
- `obj.ServicePort.Unset()` — clear the value
### SourceAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.SourceAddress.IsSet()` — check if set
- `obj.SourceAddress.Get()` — get the inner value (returns pointer)
- `obj.SourceAddress.Set(&val)` — set the value
- `obj.SourceAddress.Unset()` — clear the value
### SslMode (Nullable)

Use the Nullable wrapper methods:
- `obj.SslMode.IsSet()` — check if set
- `obj.SslMode.Get()` — get the inner value (returns pointer)
- `obj.SslMode.Set(&val)` — set the value
- `obj.SslMode.Unset()` — clear the value
### SslRedirectMode (Nullable)

Use the Nullable wrapper methods:
- `obj.SslRedirectMode.IsSet()` — check if set
- `obj.SslRedirectMode.Get()` — get the inner value (returns pointer)
- `obj.SslRedirectMode.Set(&val)` — set the value
- `obj.SslRedirectMode.Unset()` — clear the value
### VipDirectAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.VipDirectAddress.IsSet()` — check if set
- `obj.VipDirectAddress.Get()` — get the inner value (returns pointer)
- `obj.VipDirectAddress.Set(&val)` — set the value
- `obj.VipDirectAddress.Unset()` — clear the value
### ServerName (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerName.IsSet()` — check if set
- `obj.ServerName.Get()` — get the inner value (returns pointer)
- `obj.ServerName.Set(&val)` — set the value
- `obj.ServerName.Unset()` — clear the value
### PoolName (Nullable)

Use the Nullable wrapper methods:
- `obj.PoolName.IsSet()` — check if set
- `obj.PoolName.Get()` — get the inner value (returns pointer)
- `obj.PoolName.Set(&val)` — set the value
- `obj.PoolName.Unset()` — clear the value
### ExtraConfig (Nullable)

Use the Nullable wrapper methods:
- `obj.ExtraConfig.IsSet()` — check if set
- `obj.ExtraConfig.Get()` — get the inner value (returns pointer)
- `obj.ExtraConfig.Set(&val)` — set the value
- `obj.ExtraConfig.Unset()` — clear the value
### ServiceAccess (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceAccess.IsSet()` — check if set
- `obj.ServiceAccess.Get()` — get the inner value (returns pointer)
- `obj.ServiceAccess.Set(&val)` — set the value
- `obj.ServiceAccess.Unset()` — clear the value
### NetworkId (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkId.IsSet()` — check if set
- `obj.NetworkId.Get()` — get the inner value (returns pointer)
- `obj.NetworkId.Set(&val)` — set the value
- `obj.NetworkId.Unset()` — clear the value
### SubnetId (Nullable)

Use the Nullable wrapper methods:
- `obj.SubnetId.IsSet()` — check if set
- `obj.SubnetId.Get()` — get the inner value (returns pointer)
- `obj.SubnetId.Set(&val)` — set the value
- `obj.SubnetId.Unset()` — clear the value
### ExternalPortId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalPortId.IsSet()` — check if set
- `obj.ExternalPortId.Get()` — get the inner value (returns pointer)
- `obj.ExternalPortId.Set(&val)` — set the value
- `obj.ExternalPortId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


