# GetClusterMasters200ResponseMastersInnerInterfacesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**UniqueId** | Pointer to **NullableString** |  | [optional] 
**PublicIpAddress** | Pointer to **string** |  | [optional] 
**PublicIpv6Address** | Pointer to **NullableString** |  | [optional] 
**IpAddress** | Pointer to **string** |  | [optional] 
**Ipv6Address** | Pointer to **NullableString** |  | [optional] 
**IpSubnet** | Pointer to **NullableString** |  | [optional] 
**Ipv6Subnet** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Dhcp** | Pointer to **bool** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**PoolAssigned** | Pointer to **bool** |  | [optional] 
**PrimaryInterface** | Pointer to **bool** |  | [optional] 
**Network** | Pointer to [**NullableGetClusterMasters200ResponseMastersInnerInterfacesInnerNetwork**](GetClusterMasters200ResponseMastersInnerInterfacesInnerNetwork.md) |  | [optional] 
**Subnet** | Pointer to **NullableString** |  | [optional] 
**NetworkGroup** | Pointer to **NullableString** |  | [optional] 
**NetworkPosition** | Pointer to **NullableString** |  | [optional] 
**NetworkPool** | Pointer to **NullableString** |  | [optional] 
**NetworkDomain** | Pointer to [**GetClusterMasters200ResponseMastersInnerInterfacesInnerNetworkDomain**](GetClusterMasters200ResponseMastersInnerInterfacesInnerNetworkDomain.md) |  | [optional] 
**Type** | Pointer to **NullableString** |  | [optional] 
**IpMode** | Pointer to **NullableString** |  | [optional] 
**MacAddress** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetClusterMasters200ResponseMastersInnerInterfacesInner{
    // Set fields directly
}
```

### RefType (Nullable)

Use the Nullable wrapper methods:
- `obj.RefType.IsSet()` — check if set
- `obj.RefType.Get()` — get the inner value (returns pointer)
- `obj.RefType.Set(&val)` — set the value
- `obj.RefType.Unset()` — clear the value
### RefId (Nullable)

Use the Nullable wrapper methods:
- `obj.RefId.IsSet()` — check if set
- `obj.RefId.Get()` — get the inner value (returns pointer)
- `obj.RefId.Set(&val)` — set the value
- `obj.RefId.Unset()` — clear the value
### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### UniqueId (Nullable)

Use the Nullable wrapper methods:
- `obj.UniqueId.IsSet()` — check if set
- `obj.UniqueId.Get()` — get the inner value (returns pointer)
- `obj.UniqueId.Set(&val)` — set the value
- `obj.UniqueId.Unset()` — clear the value
### PublicIpv6Address (Nullable)

Use the Nullable wrapper methods:
- `obj.PublicIpv6Address.IsSet()` — check if set
- `obj.PublicIpv6Address.Get()` — get the inner value (returns pointer)
- `obj.PublicIpv6Address.Set(&val)` — set the value
- `obj.PublicIpv6Address.Unset()` — clear the value
### Ipv6Address (Nullable)

Use the Nullable wrapper methods:
- `obj.Ipv6Address.IsSet()` — check if set
- `obj.Ipv6Address.Get()` — get the inner value (returns pointer)
- `obj.Ipv6Address.Set(&val)` — set the value
- `obj.Ipv6Address.Unset()` — clear the value
### IpSubnet (Nullable)

Use the Nullable wrapper methods:
- `obj.IpSubnet.IsSet()` — check if set
- `obj.IpSubnet.Get()` — get the inner value (returns pointer)
- `obj.IpSubnet.Set(&val)` — set the value
- `obj.IpSubnet.Unset()` — clear the value
### Ipv6Subnet (Nullable)

Use the Nullable wrapper methods:
- `obj.Ipv6Subnet.IsSet()` — check if set
- `obj.Ipv6Subnet.Get()` — get the inner value (returns pointer)
- `obj.Ipv6Subnet.Set(&val)` — set the value
- `obj.Ipv6Subnet.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Network (Nullable)

Use the Nullable wrapper methods:
- `obj.Network.IsSet()` — check if set
- `obj.Network.Get()` — get the inner value (returns pointer)
- `obj.Network.Set(&val)` — set the value
- `obj.Network.Unset()` — clear the value
### Subnet (Nullable)

Use the Nullable wrapper methods:
- `obj.Subnet.IsSet()` — check if set
- `obj.Subnet.Get()` — get the inner value (returns pointer)
- `obj.Subnet.Set(&val)` — set the value
- `obj.Subnet.Unset()` — clear the value
### NetworkGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkGroup.IsSet()` — check if set
- `obj.NetworkGroup.Get()` — get the inner value (returns pointer)
- `obj.NetworkGroup.Set(&val)` — set the value
- `obj.NetworkGroup.Unset()` — clear the value
### NetworkPosition (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkPosition.IsSet()` — check if set
- `obj.NetworkPosition.Get()` — get the inner value (returns pointer)
- `obj.NetworkPosition.Set(&val)` — set the value
- `obj.NetworkPosition.Unset()` — clear the value
### NetworkPool (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkPool.IsSet()` — check if set
- `obj.NetworkPool.Get()` — get the inner value (returns pointer)
- `obj.NetworkPool.Set(&val)` — set the value
- `obj.NetworkPool.Unset()` — clear the value
### Type (Nullable)

Use the Nullable wrapper methods:
- `obj.Type.IsSet()` — check if set
- `obj.Type.Get()` — get the inner value (returns pointer)
- `obj.Type.Set(&val)` — set the value
- `obj.Type.Unset()` — clear the value
### IpMode (Nullable)

Use the Nullable wrapper methods:
- `obj.IpMode.IsSet()` — check if set
- `obj.IpMode.Get()` — get the inner value (returns pointer)
- `obj.IpMode.Set(&val)` — set the value
- `obj.IpMode.Unset()` — clear the value
### MacAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.MacAddress.IsSet()` — check if set
- `obj.MacAddress.Get()` — get the inner value (returns pointer)
- `obj.MacAddress.Set(&val)` — set the value
- `obj.MacAddress.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


