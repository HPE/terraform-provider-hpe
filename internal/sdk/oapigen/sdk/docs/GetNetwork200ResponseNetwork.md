# GetNetwork200ResponseNetwork

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Network ID | [optional] 
**Name** | Pointer to **string** | Name | [optional] 
**DisplayName** | Pointer to **NullableString** | Network Display Name | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Group** | Pointer to [**GetNetwork200ResponseNetworkGroup**](GetNetwork200ResponseNetworkGroup.md) |  | [optional] 
**Zone** | Pointer to [**GetNetwork200ResponseNetworkZone**](GetNetwork200ResponseNetworkZone.md) |  | [optional] 
**Type** | Pointer to [**GetNetwork200ResponseNetworkType**](GetNetwork200ResponseNetworkType.md) |  | [optional] 
**Owner** | Pointer to [**GetNetwork200ResponseNetworkOwner**](GetNetwork200ResponseNetworkOwner.md) |  | [optional] 
**Code** | Pointer to **NullableString** | Network Code | [optional] 
**Ipv4Enabled** | Pointer to **bool** |  | [optional] 
**Ipv6Enabled** | Pointer to **bool** |  | [optional] 
**Category** | Pointer to **NullableString** | Network Category | [optional] 
**InterfaceName** | Pointer to **NullableString** |  | [optional] 
**BridgeName** | Pointer to **NullableString** |  | [optional] 
**BridgeInterface** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **NullableString** | Description | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**UniqueId** | Pointer to **NullableString** |  | [optional] 
**ExternalType** | Pointer to **string** |  | [optional] 
**RefUrl** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**VlanId** | Pointer to **NullableInt64** |  | [optional] 
**VswitchName** | Pointer to **NullableString** |  | [optional] 
**DhcpServer** | Pointer to **bool** |  | [optional] 
**DhcpIp** | Pointer to **NullableString** |  | [optional] 
**DhcpServerIPv6** | Pointer to **bool** |  | [optional] 
**Gateway** | Pointer to **NullableString** | Network Gateway | [optional] 
**Netmask** | Pointer to **NullableString** |  | [optional] 
**Broadcast** | Pointer to **NullableString** |  | [optional] 
**SubnetAddress** | Pointer to **NullableString** |  | [optional] 
**DnsPrimary** | Pointer to **NullableString** | Primary DNS Server | [optional] 
**DnsSecondary** | Pointer to **NullableString** | Secondary DNS Server | [optional] 
**Cidr** | Pointer to **NullableString** | Network CIDR | [optional] 
**GatewayIPv6** | Pointer to **NullableString** | IPv6 Network Gateway | [optional] 
**NetmaskIPv6** | Pointer to **NullableString** |  | [optional] 
**DnsPrimaryIPv6** | Pointer to **NullableString** | Primary IPv6 DNS Server | [optional] 
**DnsSecondaryIPv6** | Pointer to **NullableString** | Secondary IPv6 DNS Server | [optional] 
**CidrIPv6** | Pointer to **NullableString** | IPv6 Network CIDR | [optional] 
**TftpServer** | Pointer to **NullableString** |  | [optional] 
**BootFile** | Pointer to **NullableString** |  | [optional] 
**SwitchId** | Pointer to **NullableString** | Network switch identifier | [optional] 
**FabricId** | Pointer to **NullableString** |  | [optional] 
**NetworkRole** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **NullableString** |  | [optional] 
**AvailabilityZone** | Pointer to **NullableString** |  | [optional] 
**Pool** | Pointer to [**GetNetwork200ResponseNetworkPool**](GetNetwork200ResponseNetworkPool.md) |  | [optional] 
**PoolIPv6** | Pointer to [**GetNetwork200ResponseNetworkPoolIPv6**](GetNetwork200ResponseNetworkPoolIPv6.md) |  | [optional] 
**NetworkProxy** | Pointer to [**GetNetwork200ResponseNetworkNetworkProxy**](GetNetwork200ResponseNetworkNetworkProxy.md) |  | [optional] 
**NetworkDomain** | Pointer to [**GetNetwork200ResponseNetworkNetworkDomain**](GetNetwork200ResponseNetworkNetworkDomain.md) |  | [optional] 
**SearchDomains** | Pointer to **NullableString** |  | [optional] 
**PrefixLength** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**EnableAdmin** | Pointer to **bool** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**DefaultNetwork** | Pointer to **bool** |  | [optional] 
**AssignPublicIp** | Pointer to **bool** |  | [optional] 
**NoProxy** | Pointer to **NullableString** |  | [optional] 
**ApplianceUrlProxyBypass** | Pointer to **bool** |  | [optional] 
**ZonePool** | Pointer to [**GetNetwork200ResponseNetworkZonePool**](GetNetwork200ResponseNetworkZonePool.md) |  | [optional] 
**AllowStaticOverride** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to [**GetNetwork200ResponseNetworkConfig**](GetNetwork200ResponseNetworkConfig.md) |  | [optional] 
**Tenants** | Pointer to [**[]GetNetwork200ResponseNetworkTenantsInner**](GetNetwork200ResponseNetworkTenantsInner.md) |  | [optional] 
**ResourcePermission** | Pointer to [**GetNetwork200ResponseNetworkResourcePermission**](GetNetwork200ResponseNetworkResourcePermission.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetwork200ResponseNetwork{
    // Set fields directly
}
```

### DisplayName (Nullable)

Use the Nullable wrapper methods:
- `obj.DisplayName.IsSet()` — check if set
- `obj.DisplayName.Get()` — get the inner value (returns pointer)
- `obj.DisplayName.Set(&val)` — set the value
- `obj.DisplayName.Unset()` — clear the value
### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### InterfaceName (Nullable)

Use the Nullable wrapper methods:
- `obj.InterfaceName.IsSet()` — check if set
- `obj.InterfaceName.Get()` — get the inner value (returns pointer)
- `obj.InterfaceName.Set(&val)` — set the value
- `obj.InterfaceName.Unset()` — clear the value
### BridgeName (Nullable)

Use the Nullable wrapper methods:
- `obj.BridgeName.IsSet()` — check if set
- `obj.BridgeName.Get()` — get the inner value (returns pointer)
- `obj.BridgeName.Set(&val)` — set the value
- `obj.BridgeName.Unset()` — clear the value
### BridgeInterface (Nullable)

Use the Nullable wrapper methods:
- `obj.BridgeInterface.IsSet()` — check if set
- `obj.BridgeInterface.Get()` — get the inner value (returns pointer)
- `obj.BridgeInterface.Set(&val)` — set the value
- `obj.BridgeInterface.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### UniqueId (Nullable)

Use the Nullable wrapper methods:
- `obj.UniqueId.IsSet()` — check if set
- `obj.UniqueId.Get()` — get the inner value (returns pointer)
- `obj.UniqueId.Set(&val)` — set the value
- `obj.UniqueId.Unset()` — clear the value
### RefUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.RefUrl.IsSet()` — check if set
- `obj.RefUrl.Get()` — get the inner value (returns pointer)
- `obj.RefUrl.Set(&val)` — set the value
- `obj.RefUrl.Unset()` — clear the value
### VlanId (Nullable)

Use the Nullable wrapper methods:
- `obj.VlanId.IsSet()` — check if set
- `obj.VlanId.Get()` — get the inner value (returns pointer)
- `obj.VlanId.Set(&val)` — set the value
- `obj.VlanId.Unset()` — clear the value
### VswitchName (Nullable)

Use the Nullable wrapper methods:
- `obj.VswitchName.IsSet()` — check if set
- `obj.VswitchName.Get()` — get the inner value (returns pointer)
- `obj.VswitchName.Set(&val)` — set the value
- `obj.VswitchName.Unset()` — clear the value
### DhcpIp (Nullable)

Use the Nullable wrapper methods:
- `obj.DhcpIp.IsSet()` — check if set
- `obj.DhcpIp.Get()` — get the inner value (returns pointer)
- `obj.DhcpIp.Set(&val)` — set the value
- `obj.DhcpIp.Unset()` — clear the value
### Gateway (Nullable)

Use the Nullable wrapper methods:
- `obj.Gateway.IsSet()` — check if set
- `obj.Gateway.Get()` — get the inner value (returns pointer)
- `obj.Gateway.Set(&val)` — set the value
- `obj.Gateway.Unset()` — clear the value
### Netmask (Nullable)

Use the Nullable wrapper methods:
- `obj.Netmask.IsSet()` — check if set
- `obj.Netmask.Get()` — get the inner value (returns pointer)
- `obj.Netmask.Set(&val)` — set the value
- `obj.Netmask.Unset()` — clear the value
### Broadcast (Nullable)

Use the Nullable wrapper methods:
- `obj.Broadcast.IsSet()` — check if set
- `obj.Broadcast.Get()` — get the inner value (returns pointer)
- `obj.Broadcast.Set(&val)` — set the value
- `obj.Broadcast.Unset()` — clear the value
### SubnetAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.SubnetAddress.IsSet()` — check if set
- `obj.SubnetAddress.Get()` — get the inner value (returns pointer)
- `obj.SubnetAddress.Set(&val)` — set the value
- `obj.SubnetAddress.Unset()` — clear the value
### DnsPrimary (Nullable)

Use the Nullable wrapper methods:
- `obj.DnsPrimary.IsSet()` — check if set
- `obj.DnsPrimary.Get()` — get the inner value (returns pointer)
- `obj.DnsPrimary.Set(&val)` — set the value
- `obj.DnsPrimary.Unset()` — clear the value
### DnsSecondary (Nullable)

Use the Nullable wrapper methods:
- `obj.DnsSecondary.IsSet()` — check if set
- `obj.DnsSecondary.Get()` — get the inner value (returns pointer)
- `obj.DnsSecondary.Set(&val)` — set the value
- `obj.DnsSecondary.Unset()` — clear the value
### Cidr (Nullable)

Use the Nullable wrapper methods:
- `obj.Cidr.IsSet()` — check if set
- `obj.Cidr.Get()` — get the inner value (returns pointer)
- `obj.Cidr.Set(&val)` — set the value
- `obj.Cidr.Unset()` — clear the value
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
### CidrIPv6 (Nullable)

Use the Nullable wrapper methods:
- `obj.CidrIPv6.IsSet()` — check if set
- `obj.CidrIPv6.Get()` — get the inner value (returns pointer)
- `obj.CidrIPv6.Set(&val)` — set the value
- `obj.CidrIPv6.Unset()` — clear the value
### TftpServer (Nullable)

Use the Nullable wrapper methods:
- `obj.TftpServer.IsSet()` — check if set
- `obj.TftpServer.Get()` — get the inner value (returns pointer)
- `obj.TftpServer.Set(&val)` — set the value
- `obj.TftpServer.Unset()` — clear the value
### BootFile (Nullable)

Use the Nullable wrapper methods:
- `obj.BootFile.IsSet()` — check if set
- `obj.BootFile.Get()` — get the inner value (returns pointer)
- `obj.BootFile.Set(&val)` — set the value
- `obj.BootFile.Unset()` — clear the value
### SwitchId (Nullable)

Use the Nullable wrapper methods:
- `obj.SwitchId.IsSet()` — check if set
- `obj.SwitchId.Get()` — get the inner value (returns pointer)
- `obj.SwitchId.Set(&val)` — set the value
- `obj.SwitchId.Unset()` — clear the value
### FabricId (Nullable)

Use the Nullable wrapper methods:
- `obj.FabricId.IsSet()` — check if set
- `obj.FabricId.Get()` — get the inner value (returns pointer)
- `obj.FabricId.Set(&val)` — set the value
- `obj.FabricId.Unset()` — clear the value
### NetworkRole (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkRole.IsSet()` — check if set
- `obj.NetworkRole.Get()` — get the inner value (returns pointer)
- `obj.NetworkRole.Set(&val)` — set the value
- `obj.NetworkRole.Unset()` — clear the value
### Status (Nullable)

Use the Nullable wrapper methods:
- `obj.Status.IsSet()` — check if set
- `obj.Status.Get()` — get the inner value (returns pointer)
- `obj.Status.Set(&val)` — set the value
- `obj.Status.Unset()` — clear the value
### AvailabilityZone (Nullable)

Use the Nullable wrapper methods:
- `obj.AvailabilityZone.IsSet()` — check if set
- `obj.AvailabilityZone.Get()` — get the inner value (returns pointer)
- `obj.AvailabilityZone.Set(&val)` — set the value
- `obj.AvailabilityZone.Unset()` — clear the value
### SearchDomains (Nullable)

Use the Nullable wrapper methods:
- `obj.SearchDomains.IsSet()` — check if set
- `obj.SearchDomains.Get()` — get the inner value (returns pointer)
- `obj.SearchDomains.Set(&val)` — set the value
- `obj.SearchDomains.Unset()` — clear the value
### PrefixLength (Nullable)

Use the Nullable wrapper methods:
- `obj.PrefixLength.IsSet()` — check if set
- `obj.PrefixLength.Get()` — get the inner value (returns pointer)
- `obj.PrefixLength.Set(&val)` — set the value
- `obj.PrefixLength.Unset()` — clear the value
### NoProxy (Nullable)

Use the Nullable wrapper methods:
- `obj.NoProxy.IsSet()` — check if set
- `obj.NoProxy.Get()` — get the inner value (returns pointer)
- `obj.NoProxy.Set(&val)` — set the value
- `obj.NoProxy.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


