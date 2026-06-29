# UpdateServerNetworkInterface200ResponseAllOfOneOfNetworkInterface

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Addresses** | Pointer to **[]map[string]interface{}** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**InterfaceId** | Pointer to **NullableString** |  | [optional] 
**DisplayOrder** | Pointer to **int64** |  | [optional] 
**NetworkPool** | Pointer to **map[string]interface{}** |  | [optional] 
**Dhcp** | Pointer to **bool** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**UniqueId** | Pointer to **string** |  | [optional] 
**Subnet** | Pointer to **NullableString** |  | [optional] 
**ReplaceHostRecord** | Pointer to **bool** |  | [optional] 
**IpMode** | Pointer to **string** |  | [optional] 
**Version** | Pointer to **NullableString** |  | [optional] 
**IpSubnet** | Pointer to **NullableString** |  | [optional] 
**Config** | Pointer to **NullableString** |  | [optional] 
**PublicIpAddress** | Pointer to **string** |  | [optional] 
**FabricId** | Pointer to **NullableString** |  | [optional] 
**Ipv6Subnet** | Pointer to **NullableString** |  | [optional] 
**MacAddress** | Pointer to **string** |  | [optional] 
**PublicIpv6Address** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**NetworkGroup** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableString** |  | [optional] 
**NetworkDomain** | Pointer to [**UpdateServerNetworkInterface200ResponseAllOfOneOfNetworkInterfaceNetworkDomain**](UpdateServerNetworkInterface200ResponseAllOfOneOfNetworkInterfaceNetworkDomain.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**PrimaryInterface** | Pointer to **bool** |  | [optional] 
**NetworkPoolIPv6** | Pointer to **map[string]interface{}** |  | [optional] 
**Network** | Pointer to [**UpdateServerNetworkInterface200ResponseAllOfOneOfNetworkInterfaceNetwork**](UpdateServerNetworkInterface200ResponseAllOfOneOfNetworkInterfaceNetwork.md) |  | [optional] 
**VlanId** | Pointer to **NullableString** |  | [optional] 
**Type** | Pointer to [**UpdateServerNetworkInterface200ResponseAllOfOneOfNetworkInterfaceType**](UpdateServerNetworkInterface200ResponseAllOfOneOfNetworkInterfaceType.md) |  | [optional] 
**NetworkPosition** | Pointer to **NullableString** |  | [optional] 
**PoolAssigned** | Pointer to **bool** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**ExternalType** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateServerNetworkInterface200ResponseAllOfOneOfNetworkInterface{
    // Set fields directly
}
```

### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### InterfaceId (Nullable)

Use the Nullable wrapper methods:
- `obj.InterfaceId.IsSet()` — check if set
- `obj.InterfaceId.Get()` — get the inner value (returns pointer)
- `obj.InterfaceId.Set(&val)` — set the value
- `obj.InterfaceId.Unset()` — clear the value
### NetworkPool (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkPool.IsSet()` — check if set
- `obj.NetworkPool.Get()` — get the inner value (returns pointer)
- `obj.NetworkPool.Set(&val)` — set the value
- `obj.NetworkPool.Unset()` — clear the value
### Subnet (Nullable)

Use the Nullable wrapper methods:
- `obj.Subnet.IsSet()` — check if set
- `obj.Subnet.Get()` — get the inner value (returns pointer)
- `obj.Subnet.Set(&val)` — set the value
- `obj.Subnet.Unset()` — clear the value
### Version (Nullable)

Use the Nullable wrapper methods:
- `obj.Version.IsSet()` — check if set
- `obj.Version.Get()` — get the inner value (returns pointer)
- `obj.Version.Set(&val)` — set the value
- `obj.Version.Unset()` — clear the value
### IpSubnet (Nullable)

Use the Nullable wrapper methods:
- `obj.IpSubnet.IsSet()` — check if set
- `obj.IpSubnet.Get()` — get the inner value (returns pointer)
- `obj.IpSubnet.Set(&val)` — set the value
- `obj.IpSubnet.Unset()` — clear the value
### Config (Nullable)

Use the Nullable wrapper methods:
- `obj.Config.IsSet()` — check if set
- `obj.Config.Get()` — get the inner value (returns pointer)
- `obj.Config.Set(&val)` — set the value
- `obj.Config.Unset()` — clear the value
### FabricId (Nullable)

Use the Nullable wrapper methods:
- `obj.FabricId.IsSet()` — check if set
- `obj.FabricId.Get()` — get the inner value (returns pointer)
- `obj.FabricId.Set(&val)` — set the value
- `obj.FabricId.Unset()` — clear the value
### Ipv6Subnet (Nullable)

Use the Nullable wrapper methods:
- `obj.Ipv6Subnet.IsSet()` — check if set
- `obj.Ipv6Subnet.Get()` — get the inner value (returns pointer)
- `obj.Ipv6Subnet.Set(&val)` — set the value
- `obj.Ipv6Subnet.Unset()` — clear the value
### PublicIpv6Address (Nullable)

Use the Nullable wrapper methods:
- `obj.PublicIpv6Address.IsSet()` — check if set
- `obj.PublicIpv6Address.Get()` — get the inner value (returns pointer)
- `obj.PublicIpv6Address.Set(&val)` — set the value
- `obj.PublicIpv6Address.Unset()` — clear the value
### RefType (Nullable)

Use the Nullable wrapper methods:
- `obj.RefType.IsSet()` — check if set
- `obj.RefType.Get()` — get the inner value (returns pointer)
- `obj.RefType.Set(&val)` — set the value
- `obj.RefType.Unset()` — clear the value
### NetworkGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkGroup.IsSet()` — check if set
- `obj.NetworkGroup.Get()` — get the inner value (returns pointer)
- `obj.NetworkGroup.Set(&val)` — set the value
- `obj.NetworkGroup.Unset()` — clear the value
### RefId (Nullable)

Use the Nullable wrapper methods:
- `obj.RefId.IsSet()` — check if set
- `obj.RefId.Get()` — get the inner value (returns pointer)
- `obj.RefId.Set(&val)` — set the value
- `obj.RefId.Unset()` — clear the value
### NetworkPoolIPv6 (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkPoolIPv6.IsSet()` — check if set
- `obj.NetworkPoolIPv6.Get()` — get the inner value (returns pointer)
- `obj.NetworkPoolIPv6.Set(&val)` — set the value
- `obj.NetworkPoolIPv6.Unset()` — clear the value
### VlanId (Nullable)

Use the Nullable wrapper methods:
- `obj.VlanId.IsSet()` — check if set
- `obj.VlanId.Get()` — get the inner value (returns pointer)
- `obj.VlanId.Set(&val)` — set the value
- `obj.VlanId.Unset()` — clear the value
### NetworkPosition (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkPosition.IsSet()` — check if set
- `obj.NetworkPosition.Get()` — get the inner value (returns pointer)
- `obj.NetworkPosition.Set(&val)` — set the value
- `obj.NetworkPosition.Unset()` — clear the value
### ExternalType (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalType.IsSet()` — check if set
- `obj.ExternalType.Get()` — get the inner value (returns pointer)
- `obj.ExternalType.Set(&val)` — set the value
- `obj.ExternalType.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


