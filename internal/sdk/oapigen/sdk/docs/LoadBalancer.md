# LoadBalancer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Cloud** | Pointer to [**CreateLoadBalancer200ResponseLoadBalancerCloud**](CreateLoadBalancer200ResponseLoadBalancerCloud.md) |  | [optional] 
**Type** | Pointer to [**CreateLoadBalancer200ResponseLoadBalancerType**](CreateLoadBalancer200ResponseLoadBalancerType.md) |  | [optional] 
**Owner** | Pointer to [**CreateLoadBalancer200ResponseLoadBalancerOwner**](CreateLoadBalancer200ResponseLoadBalancerOwner.md) |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Host** | Pointer to **string** |  | [optional] 
**Port** | Pointer to **int64** |  | [optional] 
**Username** | Pointer to **NullableString** |  | [optional] 
**Password** | Pointer to **NullableString** |  | [optional] 
**PasswordHash** | Pointer to **NullableString** |  | [optional] 
**Ip** | Pointer to **string** |  | [optional] 
**InternalIp** | Pointer to **NullableString** |  | [optional] 
**ExternalIp** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**ApiPort** | Pointer to **NullableString** |  | [optional] 
**AdminPort** | Pointer to **NullableString** |  | [optional] 
**SslEnabled** | Pointer to **NullableBool** |  | [optional] 
**SslCert** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**AllowVipEntry** | Pointer to **bool** |  | [optional] 
**VipPools** | Pointer to **[]map[string]interface{}** |  | [optional] 
**VirtualServiceName** | Pointer to **NullableString** |  | [optional] 
**PoolName** | Pointer to **NullableString** |  | [optional] 
**ServerName** | Pointer to **NullableString** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Credential** | Pointer to [**CreateLoadBalancer200ResponseLoadBalancerCredential**](CreateLoadBalancer200ResponseLoadBalancerCredential.md) |  | [optional] 
**Tenants** | Pointer to [**[]CreateLoadBalancer200ResponseLoadBalancerTenantsInner**](CreateLoadBalancer200ResponseLoadBalancerTenantsInner.md) |  | [optional] 
**ResourcePermission** | Pointer to [**CreateLoadBalancer200ResponseLoadBalancerResourcePermission**](CreateLoadBalancer200ResponseLoadBalancerResourcePermission.md) |  | [optional] 
**InstancePrice** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &LoadBalancer{
    // Set fields directly
}
```

### Username (Nullable)

Use the Nullable wrapper methods:
- `obj.Username.IsSet()` — check if set
- `obj.Username.Get()` — get the inner value (returns pointer)
- `obj.Username.Set(&val)` — set the value
- `obj.Username.Unset()` — clear the value
### Password (Nullable)

Use the Nullable wrapper methods:
- `obj.Password.IsSet()` — check if set
- `obj.Password.Get()` — get the inner value (returns pointer)
- `obj.Password.Set(&val)` — set the value
- `obj.Password.Unset()` — clear the value
### PasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.PasswordHash.IsSet()` — check if set
- `obj.PasswordHash.Get()` — get the inner value (returns pointer)
- `obj.PasswordHash.Set(&val)` — set the value
- `obj.PasswordHash.Unset()` — clear the value
### InternalIp (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalIp.IsSet()` — check if set
- `obj.InternalIp.Get()` — get the inner value (returns pointer)
- `obj.InternalIp.Set(&val)` — set the value
- `obj.InternalIp.Unset()` — clear the value
### ExternalIp (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalIp.IsSet()` — check if set
- `obj.ExternalIp.Get()` — get the inner value (returns pointer)
- `obj.ExternalIp.Set(&val)` — set the value
- `obj.ExternalIp.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### ApiPort (Nullable)

Use the Nullable wrapper methods:
- `obj.ApiPort.IsSet()` — check if set
- `obj.ApiPort.Get()` — get the inner value (returns pointer)
- `obj.ApiPort.Set(&val)` — set the value
- `obj.ApiPort.Unset()` — clear the value
### AdminPort (Nullable)

Use the Nullable wrapper methods:
- `obj.AdminPort.IsSet()` — check if set
- `obj.AdminPort.Get()` — get the inner value (returns pointer)
- `obj.AdminPort.Set(&val)` — set the value
- `obj.AdminPort.Unset()` — clear the value
### SslEnabled (Nullable)

Use the Nullable wrapper methods:
- `obj.SslEnabled.IsSet()` — check if set
- `obj.SslEnabled.Get()` — get the inner value (returns pointer)
- `obj.SslEnabled.Set(&val)` — set the value
- `obj.SslEnabled.Unset()` — clear the value
### SslCert (Nullable)

Use the Nullable wrapper methods:
- `obj.SslCert.IsSet()` — check if set
- `obj.SslCert.Get()` — get the inner value (returns pointer)
- `obj.SslCert.Set(&val)` — set the value
- `obj.SslCert.Unset()` — clear the value
### VirtualServiceName (Nullable)

Use the Nullable wrapper methods:
- `obj.VirtualServiceName.IsSet()` — check if set
- `obj.VirtualServiceName.Get()` — get the inner value (returns pointer)
- `obj.VirtualServiceName.Set(&val)` — set the value
- `obj.VirtualServiceName.Unset()` — clear the value
### PoolName (Nullable)

Use the Nullable wrapper methods:
- `obj.PoolName.IsSet()` — check if set
- `obj.PoolName.Get()` — get the inner value (returns pointer)
- `obj.PoolName.Set(&val)` — set the value
- `obj.PoolName.Unset()` — clear the value
### ServerName (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerName.IsSet()` — check if set
- `obj.ServerName.Get()` — get the inner value (returns pointer)
- `obj.ServerName.Set(&val)` — set the value
- `obj.ServerName.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


