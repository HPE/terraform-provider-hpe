# AddSecurityGroupRules200ResponseAllOfRule

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **NullableString** |  | [optional] 
**RuleType** | Pointer to **string** |  | [optional] 
**CustomRule** | Pointer to **bool** |  | [optional] 
**InstanceTypeId** | Pointer to **NullableString** |  | [optional] 
**Direction** | Pointer to **string** |  | [optional] 
**Policy** | Pointer to **string** |  | [optional] 
**SourceType** | Pointer to **string** |  | [optional] 
**Source** | Pointer to **NullableString** |  | [optional] 
**SourceGroup** | Pointer to [**AddSecurityGroupRules200ResponseAllOfRuleSourceGroup**](AddSecurityGroupRules200ResponseAllOfRuleSourceGroup.md) |  | [optional] 
**SourceTier** | Pointer to [**AddSecurityGroupRules200ResponseAllOfRuleSourceTier**](AddSecurityGroupRules200ResponseAllOfRuleSourceTier.md) |  | [optional] 
**PortRange** | Pointer to **NullableString** |  | [optional] 
**SourcePortRange** | Pointer to **NullableString** |  | [optional] 
**DestinationPortRange** | Pointer to **NullableString** |  | [optional] 
**Protocol** | Pointer to **string** |  | [optional] 
**DestinationType** | Pointer to **string** |  | [optional] 
**Destination** | Pointer to **NullableString** |  | [optional] 
**DestinationGroup** | Pointer to [**AddSecurityGroupRules200ResponseAllOfRuleDestinationGroup**](AddSecurityGroupRules200ResponseAllOfRuleDestinationGroup.md) |  | [optional] 
**DestinationTier** | Pointer to [**AddSecurityGroupRules200ResponseAllOfRuleDestinationTier**](AddSecurityGroupRules200ResponseAllOfRuleDestinationTier.md) |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddSecurityGroupRules200ResponseAllOfRule{
    // Set fields directly
}
```

### Name (Nullable)

Use the Nullable wrapper methods:
- `obj.Name.IsSet()` — check if set
- `obj.Name.Get()` — get the inner value (returns pointer)
- `obj.Name.Set(&val)` — set the value
- `obj.Name.Unset()` — clear the value
### InstanceTypeId (Nullable)

Use the Nullable wrapper methods:
- `obj.InstanceTypeId.IsSet()` — check if set
- `obj.InstanceTypeId.Get()` — get the inner value (returns pointer)
- `obj.InstanceTypeId.Set(&val)` — set the value
- `obj.InstanceTypeId.Unset()` — clear the value
### Source (Nullable)

Use the Nullable wrapper methods:
- `obj.Source.IsSet()` — check if set
- `obj.Source.Get()` — get the inner value (returns pointer)
- `obj.Source.Set(&val)` — set the value
- `obj.Source.Unset()` — clear the value
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
### Destination (Nullable)

Use the Nullable wrapper methods:
- `obj.Destination.IsSet()` — check if set
- `obj.Destination.Get()` — get the inner value (returns pointer)
- `obj.Destination.Set(&val)` — set the value
- `obj.Destination.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### Enabled (Nullable)

Use the Nullable wrapper methods:
- `obj.Enabled.IsSet()` — check if set
- `obj.Enabled.Get()` — get the inner value (returns pointer)
- `obj.Enabled.Set(&val)` — set the value
- `obj.Enabled.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


