# NetworkDomain

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Fqdn** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**DomainController** | Pointer to **bool** |  | [optional] 
**PublicZone** | Pointer to **bool** |  | [optional] 
**DomainUsername** | Pointer to **NullableString** |  | [optional] 
**DomainPassword** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableInt64** |  | [optional] 
**RefSource** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**OuPath** | Pointer to **NullableString** |  | [optional] 
**DcServer** | Pointer to **NullableString** |  | [optional] 
**ZoneType** | Pointer to **NullableString** |  | [optional] 
**Dnssec** | Pointer to **NullableString** |  | [optional] 
**DomainSerial** | Pointer to **NullableString** |  | [optional] 
**Account** | Pointer to [**CreateNetworkDomain200ResponseNetworkDomainAccount**](CreateNetworkDomain200ResponseNetworkDomainAccount.md) |  | [optional] 
**Owner** | Pointer to [**CreateNetworkDomain200ResponseNetworkDomainOwner**](CreateNetworkDomain200ResponseNetworkDomainOwner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkDomain{
    // Set fields directly
}
```

### Fqdn (Nullable)

Use the Nullable wrapper methods:
- `obj.Fqdn.IsSet()` — check if set
- `obj.Fqdn.Get()` — get the inner value (returns pointer)
- `obj.Fqdn.Set(&val)` — set the value
- `obj.Fqdn.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### DomainUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.DomainUsername.IsSet()` — check if set
- `obj.DomainUsername.Get()` — get the inner value (returns pointer)
- `obj.DomainUsername.Set(&val)` — set the value
- `obj.DomainUsername.Unset()` — clear the value
### DomainPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.DomainPassword.IsSet()` — check if set
- `obj.DomainPassword.Get()` — get the inner value (returns pointer)
- `obj.DomainPassword.Set(&val)` — set the value
- `obj.DomainPassword.Unset()` — clear the value
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
### RefSource (Nullable)

Use the Nullable wrapper methods:
- `obj.RefSource.IsSet()` — check if set
- `obj.RefSource.Get()` — get the inner value (returns pointer)
- `obj.RefSource.Set(&val)` — set the value
- `obj.RefSource.Unset()` — clear the value
### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### OuPath (Nullable)

Use the Nullable wrapper methods:
- `obj.OuPath.IsSet()` — check if set
- `obj.OuPath.Get()` — get the inner value (returns pointer)
- `obj.OuPath.Set(&val)` — set the value
- `obj.OuPath.Unset()` — clear the value
### DcServer (Nullable)

Use the Nullable wrapper methods:
- `obj.DcServer.IsSet()` — check if set
- `obj.DcServer.Get()` — get the inner value (returns pointer)
- `obj.DcServer.Set(&val)` — set the value
- `obj.DcServer.Unset()` — clear the value
### ZoneType (Nullable)

Use the Nullable wrapper methods:
- `obj.ZoneType.IsSet()` — check if set
- `obj.ZoneType.Get()` — get the inner value (returns pointer)
- `obj.ZoneType.Set(&val)` — set the value
- `obj.ZoneType.Unset()` — clear the value
### Dnssec (Nullable)

Use the Nullable wrapper methods:
- `obj.Dnssec.IsSet()` — check if set
- `obj.Dnssec.Get()` — get the inner value (returns pointer)
- `obj.Dnssec.Set(&val)` — set the value
- `obj.Dnssec.Unset()` — clear the value
### DomainSerial (Nullable)

Use the Nullable wrapper methods:
- `obj.DomainSerial.IsSet()` — check if set
- `obj.DomainSerial.Get()` — get the inner value (returns pointer)
- `obj.DomainSerial.Set(&val)` — set the value
- `obj.DomainSerial.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


