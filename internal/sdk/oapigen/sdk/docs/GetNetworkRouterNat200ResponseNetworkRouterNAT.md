# GetNetworkRouterNat200ResponseNetworkRouterNAT

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int32** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Action** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**SourceNetwork** | Pointer to **string** |  | [optional] 
**DestinationNetwork** | Pointer to **NullableString** |  | [optional] 
**TranslatedNetwork** | Pointer to **string** |  | [optional] 
**SourcePorts** | Pointer to **NullableString** |  | [optional] 
**DestinationPorts** | Pointer to **NullableString** |  | [optional] 
**TranslatedPorts** | Pointer to **NullableString** |  | [optional] 
**Priority** | Pointer to **int32** |  | [optional] 
**Protocol** | Pointer to **NullableString** |  | [optional] 
**MatchIpv6DestinationPrefix** | Pointer to **NullableString** |  | [optional] 
**TranslatedIpv4SourcePrefix** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableString** |  | [optional] 
**SyncSource** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**ProviderId** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkRouterNat200ResponseNetworkRouterNAT{
    // Set fields directly
}
```

### DestinationNetwork (Nullable)

Use the Nullable wrapper methods:
- `obj.DestinationNetwork.IsSet()` — check if set
- `obj.DestinationNetwork.Get()` — get the inner value (returns pointer)
- `obj.DestinationNetwork.Set(&val)` — set the value
- `obj.DestinationNetwork.Unset()` — clear the value
### SourcePorts (Nullable)

Use the Nullable wrapper methods:
- `obj.SourcePorts.IsSet()` — check if set
- `obj.SourcePorts.Get()` — get the inner value (returns pointer)
- `obj.SourcePorts.Set(&val)` — set the value
- `obj.SourcePorts.Unset()` — clear the value
### DestinationPorts (Nullable)

Use the Nullable wrapper methods:
- `obj.DestinationPorts.IsSet()` — check if set
- `obj.DestinationPorts.Get()` — get the inner value (returns pointer)
- `obj.DestinationPorts.Set(&val)` — set the value
- `obj.DestinationPorts.Unset()` — clear the value
### TranslatedPorts (Nullable)

Use the Nullable wrapper methods:
- `obj.TranslatedPorts.IsSet()` — check if set
- `obj.TranslatedPorts.Get()` — get the inner value (returns pointer)
- `obj.TranslatedPorts.Set(&val)` — set the value
- `obj.TranslatedPorts.Unset()` — clear the value
### Protocol (Nullable)

Use the Nullable wrapper methods:
- `obj.Protocol.IsSet()` — check if set
- `obj.Protocol.Get()` — get the inner value (returns pointer)
- `obj.Protocol.Set(&val)` — set the value
- `obj.Protocol.Unset()` — clear the value
### MatchIpv6DestinationPrefix (Nullable)

Use the Nullable wrapper methods:
- `obj.MatchIpv6DestinationPrefix.IsSet()` — check if set
- `obj.MatchIpv6DestinationPrefix.Get()` — get the inner value (returns pointer)
- `obj.MatchIpv6DestinationPrefix.Set(&val)` — set the value
- `obj.MatchIpv6DestinationPrefix.Unset()` — clear the value
### TranslatedIpv4SourcePrefix (Nullable)

Use the Nullable wrapper methods:
- `obj.TranslatedIpv4SourcePrefix.IsSet()` — check if set
- `obj.TranslatedIpv4SourcePrefix.Get()` — get the inner value (returns pointer)
- `obj.TranslatedIpv4SourcePrefix.Set(&val)` — set the value
- `obj.TranslatedIpv4SourcePrefix.Unset()` — clear the value
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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


