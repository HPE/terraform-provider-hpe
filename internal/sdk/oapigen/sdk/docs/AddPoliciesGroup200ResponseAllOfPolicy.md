# AddPoliciesGroup200ResponseAllOfPolicy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**PolicyType** | Pointer to [**AddPoliciesGroup200ResponseAllOfPolicyPolicyType**](AddPoliciesGroup200ResponseAllOfPolicyPolicyType.md) |  | [optional] 
**Zone** | Pointer to [**AddPoliciesGroup200ResponseAllOfPolicyZone**](AddPoliciesGroup200ResponseAllOfPolicyZone.md) |  | [optional] 
**Site** | Pointer to [**AddPoliciesGroup200ResponseAllOfPolicySite**](AddPoliciesGroup200ResponseAllOfPolicySite.md) |  | [optional] 
**User** | Pointer to [**AddPoliciesGroup200ResponseAllOfPolicyUser**](AddPoliciesGroup200ResponseAllOfPolicyUser.md) |  | [optional] 
**Role** | Pointer to [**AddPoliciesGroup200ResponseAllOfPolicyRole**](AddPoliciesGroup200ResponseAllOfPolicyRole.md) |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableInt64** |  | [optional] 
**EachUser** | Pointer to **NullableBool** |  | [optional] 
**Config** | Pointer to [**AddPoliciesGroup200ResponseAllOfPolicyConfig**](AddPoliciesGroup200ResponseAllOfPolicyConfig.md) |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Owner** | Pointer to [**NullableAddPoliciesGroup200ResponseAllOfPolicyOwner**](AddPoliciesGroup200ResponseAllOfPolicyOwner.md) |  | [optional] 
**Accounts** | Pointer to [**[]ListPolicies200ResponseAllOfPoliciesInnerAccountsInner**](ListPolicies200ResponseAllOfPoliciesInnerAccountsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddPoliciesGroup200ResponseAllOfPolicy{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
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
### EachUser (Nullable)

Use the Nullable wrapper methods:
- `obj.EachUser.IsSet()` — check if set
- `obj.EachUser.Get()` — get the inner value (returns pointer)
- `obj.EachUser.Set(&val)` — set the value
- `obj.EachUser.Unset()` — clear the value
### Owner (Nullable)

Use the Nullable wrapper methods:
- `obj.Owner.IsSet()` — check if set
- `obj.Owner.Get()` — get the inner value (returns pointer)
- `obj.Owner.Set(&val)` — set the value
- `obj.Owner.Unset()` — clear the value
### Accounts (Nullable)

Use the Nullable wrapper methods:
- `obj.Accounts.IsSet()` — check if set
- `obj.Accounts.Get()` — get the inner value (returns pointer)
- `obj.Accounts.Set(&val)` — set the value
- `obj.Accounts.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


