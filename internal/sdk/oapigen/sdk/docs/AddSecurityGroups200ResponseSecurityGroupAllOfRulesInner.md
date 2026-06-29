# AddSecurityGroups200ResponseSecurityGroupAllOfRulesInner

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
**SourceGroup** | Pointer to **NullableString** |  | [optional] 
**SourceTier** | Pointer to **NullableString** |  | [optional] 
**PortRange** | Pointer to **NullableString** |  | [optional] 
**Protocol** | Pointer to **string** |  | [optional] 
**DestinationType** | Pointer to **string** |  | [optional] 
**Destination** | Pointer to **NullableString** |  | [optional] 
**DestinationGroup** | Pointer to **NullableString** |  | [optional] 
**DestinationTier** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddSecurityGroups200ResponseSecurityGroupAllOfRulesInner{
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
### PortRange (Nullable)

Use the Nullable wrapper methods:
- `obj.PortRange.IsSet()` — check if set
- `obj.PortRange.Get()` — get the inner value (returns pointer)
- `obj.PortRange.Set(&val)` — set the value
- `obj.PortRange.Unset()` — clear the value
### Destination (Nullable)

Use the Nullable wrapper methods:
- `obj.Destination.IsSet()` — check if set
- `obj.Destination.Get()` — get the inner value (returns pointer)
- `obj.Destination.Set(&val)` — set the value
- `obj.Destination.Unset()` — clear the value
### DestinationGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.DestinationGroup.IsSet()` — check if set
- `obj.DestinationGroup.Get()` — get the inner value (returns pointer)
- `obj.DestinationGroup.Set(&val)` — set the value
- `obj.DestinationGroup.Unset()` — clear the value
### DestinationTier (Nullable)

Use the Nullable wrapper methods:
- `obj.DestinationTier.IsSet()` — check if set
- `obj.DestinationTier.Get()` — get the inner value (returns pointer)
- `obj.DestinationTier.Set(&val)` — set the value
- `obj.DestinationTier.Unset()` — clear the value
### Enabled (Nullable)

Use the Nullable wrapper methods:
- `obj.Enabled.IsSet()` — check if set
- `obj.Enabled.Get()` — get the inner value (returns pointer)
- `obj.Enabled.Set(&val)` — set the value
- `obj.Enabled.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


