# GetGuidances200ResponseDiscoveryAnyOf1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**ActionCategory** | Pointer to **string** |  | [optional] 
**ActionMessage** | Pointer to **string** |  | [optional] 
**ActionTitle** | Pointer to **string** |  | [optional] 
**ActionType** | Pointer to **string** |  | [optional] 
**ActionValue** | Pointer to **string** |  | [optional] 
**ActionValueType** | Pointer to **string** |  | [optional] 
**ActionPlanId** | Pointer to **NullableString** |  | [optional] 
**StatusMessage** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**UserId** | Pointer to **NullableString** |  | [optional] 
**SiteId** | Pointer to **int64** |  | [optional] 
**Zone** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOf1Zone**](GetGuidances200ResponseDiscoveryAnyOf1Zone.md) |  | [optional] 
**State** | Pointer to **string** |  | [optional] 
**StateMessage** | Pointer to **NullableString** |  | [optional] 
**Severity** | Pointer to **string** |  | [optional] 
**Resolved** | Pointer to **bool** |  | [optional] 
**ResolvedMessage** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**RefName** | Pointer to **string** |  | [optional] 
**Type** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOf1Type**](GetGuidances200ResponseDiscoveryAnyOf1Type.md) |  | [optional] 
**Savings** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOf1Savings**](GetGuidances200ResponseDiscoveryAnyOf1Savings.md) |  | [optional] 
**Config** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOf1Config**](GetGuidances200ResponseDiscoveryAnyOf1Config.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetGuidances200ResponseDiscoveryAnyOf1{
    // Set fields directly
}
```

### ActionPlanId (Nullable)

Use the Nullable wrapper methods:
- `obj.ActionPlanId.IsSet()` — check if set
- `obj.ActionPlanId.Get()` — get the inner value (returns pointer)
- `obj.ActionPlanId.Set(&val)` — set the value
- `obj.ActionPlanId.Unset()` — clear the value
### UserId (Nullable)

Use the Nullable wrapper methods:
- `obj.UserId.IsSet()` — check if set
- `obj.UserId.Get()` — get the inner value (returns pointer)
- `obj.UserId.Set(&val)` — set the value
- `obj.UserId.Unset()` — clear the value
### StateMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StateMessage.IsSet()` — check if set
- `obj.StateMessage.Get()` — get the inner value (returns pointer)
- `obj.StateMessage.Set(&val)` — set the value
- `obj.StateMessage.Unset()` — clear the value
### ResolvedMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.ResolvedMessage.IsSet()` — check if set
- `obj.ResolvedMessage.Get()` — get the inner value (returns pointer)
- `obj.ResolvedMessage.Set(&val)` — set the value
- `obj.ResolvedMessage.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


