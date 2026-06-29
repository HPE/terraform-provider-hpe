# ListSubnets200ResponseAllOfSubnetsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**UniqueId** | Pointer to **NullableString** |  | [optional] 
**AddressPrefix** | Pointer to **NullableString** |  | [optional] 
**Cidr** | Pointer to **string** |  | [optional] 
**Gateway** | Pointer to **NullableString** |  | [optional] 
**Netmask** | Pointer to **string** |  | [optional] 
**SubnetAddress** | Pointer to **string** |  | [optional] 
**TftpServer** | Pointer to **NullableString** |  | [optional] 
**BootFile** | Pointer to **NullableString** |  | [optional] 
**Pool** | Pointer to [**ListSubnets200ResponseAllOfSubnetsInnerPool**](ListSubnets200ResponseAllOfSubnetsInnerPool.md) |  | [optional] 
**DhcpServer** | Pointer to **bool** |  | [optional] 
**HasFloatingIps** | Pointer to **bool** |  | [optional] 
**DhcpIp** | Pointer to **NullableString** |  | [optional] 
**DnsPrimary** | Pointer to **NullableString** |  | [optional] 
**DnsSecondary** | Pointer to **NullableString** |  | [optional] 
**DhcpStart** | Pointer to **string** |  | [optional] 
**DhcpEnd** | Pointer to **string** |  | [optional] 
**DhcpRange** | Pointer to **NullableString** |  | [optional] 
**NetworkProxy** | Pointer to [**ListSubnets200ResponseAllOfSubnetsInnerNetworkProxy**](ListSubnets200ResponseAllOfSubnetsInnerNetworkProxy.md) |  | [optional] 
**NetworkDomain** | Pointer to [**ListSubnets200ResponseAllOfSubnetsInnerNetworkDomain**](ListSubnets200ResponseAllOfSubnetsInnerNetworkDomain.md) |  | [optional] 
**SearchDomains** | Pointer to **NullableString** |  | [optional] 
**DefaultNetwork** | Pointer to **bool** |  | [optional] 
**AssignPublicIp** | Pointer to **bool** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Status** | Pointer to [**ListSubnets200ResponseAllOfSubnetsInnerStatus**](ListSubnets200ResponseAllOfSubnetsInnerStatus.md) |  | [optional] 
**Network** | Pointer to [**ListSubnets200ResponseAllOfSubnetsInnerNetwork**](ListSubnets200ResponseAllOfSubnetsInnerNetwork.md) |  | [optional] 
**Zone** | Pointer to [**ListSubnets200ResponseAllOfSubnetsInnerZone**](ListSubnets200ResponseAllOfSubnetsInnerZone.md) |  | [optional] 
**Type** | Pointer to [**ListSubnets200ResponseAllOfSubnetsInnerType**](ListSubnets200ResponseAllOfSubnetsInnerType.md) |  | [optional] 
**Account** | Pointer to [**ListSubnets200ResponseAllOfSubnetsInnerAccount**](ListSubnets200ResponseAllOfSubnetsInnerAccount.md) |  | [optional] 
**SecurityGroups** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Tenants** | Pointer to [**[]ListSubnets200ResponseAllOfSubnetsInnerTenantsInner**](ListSubnets200ResponseAllOfSubnetsInnerTenantsInner.md) |  | [optional] 
**ResourcePermission** | Pointer to [**ListSubnets200ResponseAllOfSubnetsInnerResourcePermission**](ListSubnets200ResponseAllOfSubnetsInnerResourcePermission.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListSubnets200ResponseAllOfSubnetsInner{
    // Set fields directly
}
```

### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
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
### UniqueId (Nullable)

Use the Nullable wrapper methods:
- `obj.UniqueId.IsSet()` — check if set
- `obj.UniqueId.Get()` — get the inner value (returns pointer)
- `obj.UniqueId.Set(&val)` — set the value
- `obj.UniqueId.Unset()` — clear the value
### AddressPrefix (Nullable)

Use the Nullable wrapper methods:
- `obj.AddressPrefix.IsSet()` — check if set
- `obj.AddressPrefix.Get()` — get the inner value (returns pointer)
- `obj.AddressPrefix.Set(&val)` — set the value
- `obj.AddressPrefix.Unset()` — clear the value
### Gateway (Nullable)

Use the Nullable wrapper methods:
- `obj.Gateway.IsSet()` — check if set
- `obj.Gateway.Get()` — get the inner value (returns pointer)
- `obj.Gateway.Set(&val)` — set the value
- `obj.Gateway.Unset()` — clear the value
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
### DhcpIp (Nullable)

Use the Nullable wrapper methods:
- `obj.DhcpIp.IsSet()` — check if set
- `obj.DhcpIp.Get()` — get the inner value (returns pointer)
- `obj.DhcpIp.Set(&val)` — set the value
- `obj.DhcpIp.Unset()` — clear the value
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
### DhcpRange (Nullable)

Use the Nullable wrapper methods:
- `obj.DhcpRange.IsSet()` — check if set
- `obj.DhcpRange.Get()` — get the inner value (returns pointer)
- `obj.DhcpRange.Set(&val)` — set the value
- `obj.DhcpRange.Unset()` — clear the value
### SearchDomains (Nullable)

Use the Nullable wrapper methods:
- `obj.SearchDomains.IsSet()` — check if set
- `obj.SearchDomains.Get()` — get the inner value (returns pointer)
- `obj.SearchDomains.Set(&val)` — set the value
- `obj.SearchDomains.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


