# GetNetworkPool200ResponseNetworkPoolIpRangesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**StartAddress** | Pointer to **NullableString** |  | [optional] 
**EndAddress** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**AddressCount** | Pointer to **int64** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Cidr** | Pointer to **NullableString** |  | [optional] 
**CidrIPv6** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkPool200ResponseNetworkPoolIpRangesInner{
    // Set fields directly
}
```

### StartAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.StartAddress.IsSet()` — check if set
- `obj.StartAddress.Get()` — get the inner value (returns pointer)
- `obj.StartAddress.Set(&val)` — set the value
- `obj.StartAddress.Unset()` — clear the value
### EndAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.EndAddress.IsSet()` — check if set
- `obj.EndAddress.Get()` — get the inner value (returns pointer)
- `obj.EndAddress.Set(&val)` — set the value
- `obj.EndAddress.Unset()` — clear the value
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
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Cidr (Nullable)

Use the Nullable wrapper methods:
- `obj.Cidr.IsSet()` — check if set
- `obj.Cidr.Get()` — get the inner value (returns pointer)
- `obj.Cidr.Set(&val)` — set the value
- `obj.Cidr.Unset()` — clear the value
### CidrIPv6 (Nullable)

Use the Nullable wrapper methods:
- `obj.CidrIPv6.IsSet()` — check if set
- `obj.CidrIPv6.Get()` — get the inner value (returns pointer)
- `obj.CidrIPv6.Set(&val)` — set the value
- `obj.CidrIPv6.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


