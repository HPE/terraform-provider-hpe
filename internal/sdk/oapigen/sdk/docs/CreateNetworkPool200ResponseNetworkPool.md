# CreateNetworkPool200ResponseNetworkPool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Type** | Pointer to [**CreateNetworkPool200ResponseNetworkPoolType**](CreateNetworkPool200ResponseNetworkPoolType.md) |  | [optional] 
**Account** | Pointer to [**CreateNetworkPool200ResponseNetworkPoolAccount**](CreateNetworkPool200ResponseNetworkPoolAccount.md) |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**DisplayName** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**DnsDomain** | Pointer to **NullableString** |  | [optional] 
**DnsSearchPath** | Pointer to **NullableString** |  | [optional] 
**HostPrefix** | Pointer to **NullableString** |  | [optional] 
**HttpProxy** | Pointer to **NullableString** |  | [optional] 
**DnsServers** | Pointer to **[]string** |  | [optional] 
**DnsSuffixList** | Pointer to **[]string** |  | [optional] 
**DhcpServer** | Pointer to **bool** |  | [optional] 
**DhcpIp** | Pointer to **NullableString** |  | [optional] 
**Gateway** | Pointer to **NullableString** |  | [optional] 
**Netmask** | Pointer to **NullableString** |  | [optional] 
**SubnetAddress** | Pointer to **NullableString** |  | [optional] 
**IpCount** | Pointer to **int64** |  | [optional] 
**FreeCount** | Pointer to **int64** |  | [optional] 
**PoolEnabled** | Pointer to **bool** |  | [optional] 
**TftpServer** | Pointer to **NullableString** |  | [optional] 
**BootFile** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableString** |  | [optional] 
**ParentType** | Pointer to **NullableString** |  | [optional] 
**ParentId** | Pointer to **NullableString** |  | [optional] 
**PoolGroup** | Pointer to **NullableString** |  | [optional] 
**IpRanges** | Pointer to [**[]CreateNetworkPool200ResponseNetworkPoolIpRangesInner**](CreateNetworkPool200ResponseNetworkPoolIpRangesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateNetworkPool200ResponseNetworkPool{
    // Set fields directly
}
```

### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
### DisplayName (Nullable)

Use the Nullable wrapper methods:
- `obj.DisplayName.IsSet()` — check if set
- `obj.DisplayName.Get()` — get the inner value (returns pointer)
- `obj.DisplayName.Set(&val)` — set the value
- `obj.DisplayName.Unset()` — clear the value
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
### DnsDomain (Nullable)

Use the Nullable wrapper methods:
- `obj.DnsDomain.IsSet()` — check if set
- `obj.DnsDomain.Get()` — get the inner value (returns pointer)
- `obj.DnsDomain.Set(&val)` — set the value
- `obj.DnsDomain.Unset()` — clear the value
### DnsSearchPath (Nullable)

Use the Nullable wrapper methods:
- `obj.DnsSearchPath.IsSet()` — check if set
- `obj.DnsSearchPath.Get()` — get the inner value (returns pointer)
- `obj.DnsSearchPath.Set(&val)` — set the value
- `obj.DnsSearchPath.Unset()` — clear the value
### HostPrefix (Nullable)

Use the Nullable wrapper methods:
- `obj.HostPrefix.IsSet()` — check if set
- `obj.HostPrefix.Get()` — get the inner value (returns pointer)
- `obj.HostPrefix.Set(&val)` — set the value
- `obj.HostPrefix.Unset()` — clear the value
### HttpProxy (Nullable)

Use the Nullable wrapper methods:
- `obj.HttpProxy.IsSet()` — check if set
- `obj.HttpProxy.Get()` — get the inner value (returns pointer)
- `obj.HttpProxy.Set(&val)` — set the value
- `obj.HttpProxy.Unset()` — clear the value
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
### SubnetAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.SubnetAddress.IsSet()` — check if set
- `obj.SubnetAddress.Get()` — get the inner value (returns pointer)
- `obj.SubnetAddress.Set(&val)` — set the value
- `obj.SubnetAddress.Unset()` — clear the value
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
### ParentType (Nullable)

Use the Nullable wrapper methods:
- `obj.ParentType.IsSet()` — check if set
- `obj.ParentType.Get()` — get the inner value (returns pointer)
- `obj.ParentType.Set(&val)` — set the value
- `obj.ParentType.Unset()` — clear the value
### ParentId (Nullable)

Use the Nullable wrapper methods:
- `obj.ParentId.IsSet()` — check if set
- `obj.ParentId.Get()` — get the inner value (returns pointer)
- `obj.ParentId.Set(&val)` — set the value
- `obj.ParentId.Unset()` — clear the value
### PoolGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.PoolGroup.IsSet()` — check if set
- `obj.PoolGroup.Get()` — get the inner value (returns pointer)
- `obj.PoolGroup.Set(&val)` — set the value
- `obj.PoolGroup.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


