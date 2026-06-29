# ListGuidances200ResponseAllOfDiscoveriesInnerAnyOf

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
**ActionPlanId** | Pointer to **int64** |  | [optional] 
**StatusMessage** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**UserId** | Pointer to **NullableString** |  | [optional] 
**SiteId** | Pointer to **NullableInt64** |  | [optional] 
**Zone** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfZone**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfZone.md) |  | [optional] 
**State** | Pointer to **string** |  | [optional] 
**StateMessage** | Pointer to **NullableString** |  | [optional] 
**Severity** | Pointer to **string** |  | [optional] 
**Resolved** | Pointer to **bool** |  | [optional] 
**ResolvedMessage** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**RefName** | Pointer to **string** |  | [optional] 
**Type** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfType**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfType.md) |  | [optional] 
**Savings** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfSavings**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfSavings.md) |  | [optional] 
**Resource** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResource**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResource.md) |  | [optional] 
**PlanBeforeAction** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfPlanBeforeAction**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfPlanBeforeAction.md) |  | [optional] 
**PlanAfterAction** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfPlanAfterAction**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfPlanAfterAction.md) |  | [optional] 
**Config** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfConfig**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListGuidances200ResponseAllOfDiscoveriesInnerAnyOf{
    // Set fields directly
}
```

### UserId (Nullable)

Use the Nullable wrapper methods:
- `obj.UserId.IsSet()` — check if set
- `obj.UserId.Get()` — get the inner value (returns pointer)
- `obj.UserId.Set(&val)` — set the value
- `obj.UserId.Unset()` — clear the value
### SiteId (Nullable)

Use the Nullable wrapper methods:
- `obj.SiteId.IsSet()` — check if set
- `obj.SiteId.Get()` — get the inner value (returns pointer)
- `obj.SiteId.Set(&val)` — set the value
- `obj.SiteId.Unset()` — clear the value
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


