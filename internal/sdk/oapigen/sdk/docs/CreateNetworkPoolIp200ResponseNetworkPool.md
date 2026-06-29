# CreateNetworkPoolIp200ResponseNetworkPool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**NetworkPoolId** | Pointer to **int64** |  | [optional] 
**IpType** | Pointer to **string** |  | [optional] 
**IpAddress** | Pointer to **string** |  | [optional] 
**GatewayAddress** | Pointer to **NullableString** |  | [optional] 
**SubnetMask** | Pointer to **NullableString** |  | [optional] 
**DnsServer** | Pointer to **NullableString** |  | [optional] 
**InterfaceName** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**StaticIp** | Pointer to **bool** |  | [optional] 
**Fqdn** | Pointer to **string** |  | [optional] 
**DomainName** | Pointer to **NullableString** |  | [optional] 
**Hostname** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**PtrId** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableInt64** |  | [optional] 
**SubRefId** | Pointer to **NullableInt64** |  | [optional] 
**NetworkDomain** | Pointer to [**CreateNetworkPoolIp200ResponseNetworkPoolNetworkDomain**](CreateNetworkPoolIp200ResponseNetworkPoolNetworkDomain.md) |  | [optional] 
**CreatedBy** | Pointer to [**CreateNetworkPoolIp200ResponseNetworkPoolCreatedBy**](CreateNetworkPoolIp200ResponseNetworkPoolCreatedBy.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateNetworkPoolIp200ResponseNetworkPool{
    // Set fields directly
}
```

### GatewayAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.GatewayAddress.IsSet()` — check if set
- `obj.GatewayAddress.Get()` — get the inner value (returns pointer)
- `obj.GatewayAddress.Set(&val)` — set the value
- `obj.GatewayAddress.Unset()` — clear the value
### SubnetMask (Nullable)

Use the Nullable wrapper methods:
- `obj.SubnetMask.IsSet()` — check if set
- `obj.SubnetMask.Get()` — get the inner value (returns pointer)
- `obj.SubnetMask.Set(&val)` — set the value
- `obj.SubnetMask.Unset()` — clear the value
### DnsServer (Nullable)

Use the Nullable wrapper methods:
- `obj.DnsServer.IsSet()` — check if set
- `obj.DnsServer.Get()` — get the inner value (returns pointer)
- `obj.DnsServer.Set(&val)` — set the value
- `obj.DnsServer.Unset()` — clear the value
### InterfaceName (Nullable)

Use the Nullable wrapper methods:
- `obj.InterfaceName.IsSet()` — check if set
- `obj.InterfaceName.Get()` — get the inner value (returns pointer)
- `obj.InterfaceName.Set(&val)` — set the value
- `obj.InterfaceName.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### DomainName (Nullable)

Use the Nullable wrapper methods:
- `obj.DomainName.IsSet()` — check if set
- `obj.DomainName.Get()` — get the inner value (returns pointer)
- `obj.DomainName.Set(&val)` — set the value
- `obj.DomainName.Unset()` — clear the value
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
### PtrId (Nullable)

Use the Nullable wrapper methods:
- `obj.PtrId.IsSet()` — check if set
- `obj.PtrId.Get()` — get the inner value (returns pointer)
- `obj.PtrId.Set(&val)` — set the value
- `obj.PtrId.Unset()` — clear the value
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
### SubRefId (Nullable)

Use the Nullable wrapper methods:
- `obj.SubRefId.IsSet()` — check if set
- `obj.SubRefId.Get()` — get the inner value (returns pointer)
- `obj.SubRefId.Set(&val)` — set the value
- `obj.SubRefId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


