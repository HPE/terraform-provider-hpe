# GetNetworkRoutersFirewallRules200ResponseRulesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Priority** | Pointer to **int64** |  | [optional] 
**GroupName** | Pointer to **string** |  | [optional] 
**Direction** | Pointer to **string** |  | [optional] 
**RuleType** | Pointer to **string** |  | [optional] 
**Policy** | Pointer to **string** |  | [optional] 
**Source** | Pointer to **[]string** |  | [optional] 
**SourceType** | Pointer to **string** |  | [optional] 
**Destination** | Pointer to **[]string** |  | [optional] 
**DestinationType** | Pointer to **string** |  | [optional] 
**Profiles** | Pointer to **[]string** |  | [optional] 
**Protocol** | Pointer to **NullableString** |  | [optional] 
**Application** | Pointer to **NullableString** |  | [optional] 
**ApplicationType** | Pointer to **string** |  | [optional] 
**PortRange** | Pointer to **NullableString** |  | [optional] 
**SourcePortRange** | Pointer to **NullableString** |  | [optional] 
**DestinationPortRange** | Pointer to **NullableString** |  | [optional] 
**SourceGroup** | Pointer to **NullableString** |  | [optional] 
**SourceTier** | Pointer to **NullableString** |  | [optional] 
**Applications** | Pointer to [**[]GetNetworkRoutersFirewallRules200ResponseRulesInnerApplicationsInner**](GetNetworkRoutersFirewallRules200ResponseRulesInnerApplicationsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkRoutersFirewallRules200ResponseRulesInner{
    // Set fields directly
}
```

### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
### Protocol (Nullable)

Use the Nullable wrapper methods:
- `obj.Protocol.IsSet()` — check if set
- `obj.Protocol.Get()` — get the inner value (returns pointer)
- `obj.Protocol.Set(&val)` — set the value
- `obj.Protocol.Unset()` — clear the value
### Application (Nullable)

Use the Nullable wrapper methods:
- `obj.Application.IsSet()` — check if set
- `obj.Application.Get()` — get the inner value (returns pointer)
- `obj.Application.Set(&val)` — set the value
- `obj.Application.Unset()` — clear the value
### PortRange (Nullable)

Use the Nullable wrapper methods:
- `obj.PortRange.IsSet()` — check if set
- `obj.PortRange.Get()` — get the inner value (returns pointer)
- `obj.PortRange.Set(&val)` — set the value
- `obj.PortRange.Unset()` — clear the value
### SourcePortRange (Nullable)

Use the Nullable wrapper methods:
- `obj.SourcePortRange.IsSet()` — check if set
- `obj.SourcePortRange.Get()` — get the inner value (returns pointer)
- `obj.SourcePortRange.Set(&val)` — set the value
- `obj.SourcePortRange.Unset()` — clear the value
### DestinationPortRange (Nullable)

Use the Nullable wrapper methods:
- `obj.DestinationPortRange.IsSet()` — check if set
- `obj.DestinationPortRange.Get()` — get the inner value (returns pointer)
- `obj.DestinationPortRange.Set(&val)` — set the value
- `obj.DestinationPortRange.Unset()` — clear the value
### SourceGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.SourceGroup.IsSet()` — check if set
- `obj.SourceGroup.Get()` — get the inner value (returns pointer)
- `obj.SourceGroup.Set(&val)` — set the value
- `obj.SourceGroup.Unset()` — clear the value
### SourceTier (Nullable)

Use the Nullable wrapper methods:
- `obj.SourceTier.IsSet()` — check if set
- `obj.SourceTier.Get()` — get the inner value (returns pointer)
- `obj.SourceTier.Set(&val)` — set the value
- `obj.SourceTier.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


