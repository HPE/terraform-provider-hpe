# NetworkCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name | 
**DisplayName** | Pointer to **string** | Display Name | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Description** | Pointer to **NullableString** | Description | [optional] 
**Site** | [**NetworkCreateSite**](NetworkCreateSite.md) |  | 
**Zone** | [**NetworkCreateZone**](NetworkCreateZone.md) |  | 
**Type** | Pointer to [**NetworkCreateType**](NetworkCreateType.md) |  | [optional] 
**Ipv4Enabled** | Pointer to **bool** |  | [optional] 
**Ipv6Enabled** | Pointer to **bool** |  | [optional] 
**Cidr** | Pointer to **string** | CIDR Network | [optional] 
**Gateway** | Pointer to **string** | Network Gateway | [optional] 
**DnsPrimary** | Pointer to **string** | Primary DNS Server | [optional] 
**DnsSecondary** | Pointer to **string** | Secondary DNS Server | [optional] 
**GatewayIPv6** | Pointer to **NullableString** | IPv6 Network Gateway | [optional] 
**NetmaskIPv6** | Pointer to **NullableString** |  | [optional] 
**DnsPrimaryIPv6** | Pointer to **NullableString** | Primary IPv6 DNS Server | [optional] 
**DnsSecondaryIPv6** | Pointer to **NullableString** | Secondary IPv6 DNS Server | [optional] 
**CidrIPv6** | Pointer to **string** | IPv6 Network CIDR | [optional] 
**VlanId** | Pointer to **int64** |  | [optional] 
**SwitchId** | Pointer to **string** | Network switch identifier | [optional] 
**Pool** | Pointer to **NullableInt64** | Network Pool ID | [optional] 
**PoolIPv6** | Pointer to **NullableInt64** | IPv6 Network Pool ID | [optional] 
**ZonePool** | Pointer to [**NetworkCreateZonePool**](NetworkCreateZonePool.md) |  | [optional] 
**AllowStaticOverride** | Pointer to **bool** | Allow IP Override | [optional] 
**AssignPublicIp** | Pointer to **bool** | Assign Public IP | [optional] 
**Active** | Pointer to **bool** | Activate (true) or disable (false) the network | [optional] 
**DhcpServer** | Pointer to **bool** | DHCP Server enabled network | [optional] 
**DhcpServerIPv6** | Pointer to **bool** | IPv6 DHCP Server enabled network | [optional] 
**NetworkDomain** | Pointer to [**NetworkCreateNetworkDomain**](NetworkCreateNetworkDomain.md) |  | [optional] 
**SearchDomains** | Pointer to **string** | Search Domains | [optional] 
**NetworkProxy** | Pointer to [**NetworkCreateNetworkProxy**](NetworkCreateNetworkProxy.md) |  | [optional] 
**ApplianceUrlProxyBypass** | Pointer to **bool** | Bypass Proxy for Appliance URL | [optional] 
**NoProxy** | Pointer to **NullableString** | Comma-separated list of ip addresses or name servers to exclude proxy traversal for. Typically locally routable servers are excluded. | [optional] 
**Visibility** | Pointer to **string** | Visibility, private or public. | [optional] [default to "private"]
**Config** | Pointer to [**NetworkCreateConfig**](NetworkCreateConfig.md) |  | [optional] 
**Tenants** | Pointer to [**[]NetworkCreateTenantsInner**](NetworkCreateTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 
**ResourcePermission** | Pointer to [**NetworkCreateResourcePermission**](NetworkCreateResourcePermission.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkCreate{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### GatewayIPv6 (Nullable)

Use the Nullable wrapper methods:
- `obj.GatewayIPv6.IsSet()` — check if set
- `obj.GatewayIPv6.Get()` — get the inner value (returns pointer)
- `obj.GatewayIPv6.Set(&val)` — set the value
- `obj.GatewayIPv6.Unset()` — clear the value
### NetmaskIPv6 (Nullable)

Use the Nullable wrapper methods:
- `obj.NetmaskIPv6.IsSet()` — check if set
- `obj.NetmaskIPv6.Get()` — get the inner value (returns pointer)
- `obj.NetmaskIPv6.Set(&val)` — set the value
- `obj.NetmaskIPv6.Unset()` — clear the value
### DnsPrimaryIPv6 (Nullable)

Use the Nullable wrapper methods:
- `obj.DnsPrimaryIPv6.IsSet()` — check if set
- `obj.DnsPrimaryIPv6.Get()` — get the inner value (returns pointer)
- `obj.DnsPrimaryIPv6.Set(&val)` — set the value
- `obj.DnsPrimaryIPv6.Unset()` — clear the value
### DnsSecondaryIPv6 (Nullable)

Use the Nullable wrapper methods:
- `obj.DnsSecondaryIPv6.IsSet()` — check if set
- `obj.DnsSecondaryIPv6.Get()` — get the inner value (returns pointer)
- `obj.DnsSecondaryIPv6.Set(&val)` — set the value
- `obj.DnsSecondaryIPv6.Unset()` — clear the value
### Pool (Nullable)

Use the Nullable wrapper methods:
- `obj.Pool.IsSet()` — check if set
- `obj.Pool.Get()` — get the inner value (returns pointer)
- `obj.Pool.Set(&val)` — set the value
- `obj.Pool.Unset()` — clear the value
### PoolIPv6 (Nullable)

Use the Nullable wrapper methods:
- `obj.PoolIPv6.IsSet()` — check if set
- `obj.PoolIPv6.Get()` — get the inner value (returns pointer)
- `obj.PoolIPv6.Set(&val)` — set the value
- `obj.PoolIPv6.Unset()` — clear the value
### NoProxy (Nullable)

Use the Nullable wrapper methods:
- `obj.NoProxy.IsSet()` — check if set
- `obj.NoProxy.Get()` — get the inner value (returns pointer)
- `obj.NoProxy.Set(&val)` — set the value
- `obj.NoProxy.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


