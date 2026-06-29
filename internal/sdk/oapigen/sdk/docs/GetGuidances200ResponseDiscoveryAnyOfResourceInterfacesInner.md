# GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**UniqueId** | Pointer to **NullableString** |  | [optional] 
**PublicIpAddress** | Pointer to **string** |  | [optional] 
**PublicIpv6Address** | Pointer to **NullableString** |  | [optional] 
**IpAddress** | Pointer to **string** |  | [optional] 
**Ipv6Address** | Pointer to **string** |  | [optional] 
**IpSubnet** | Pointer to **NullableString** |  | [optional] 
**Ipv6Subnet** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Dhcp** | Pointer to **bool** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**PoolAssigned** | Pointer to **bool** |  | [optional] 
**PrimaryInterface** | Pointer to **bool** |  | [optional] 
**Network** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInnerNetwork**](GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInnerNetwork.md) |  | [optional] 
**Subnet** | Pointer to **NullableString** |  | [optional] 
**NetworkGroup** | Pointer to **NullableString** |  | [optional] 
**NetworkPosition** | Pointer to **NullableString** |  | [optional] 
**NetworkPool** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInnerNetworkPool**](GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInnerNetworkPool.md) |  | [optional] 
**NetworkDomain** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInnerNetworkDomain**](GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInnerNetworkDomain.md) |  | [optional] 
**Type** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInnerType**](GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInnerType.md) |  | [optional] 
**IpMode** | Pointer to **string** |  | [optional] 
**MacAddress** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInner{
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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


