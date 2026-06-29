# Policy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**PolicyType** | Pointer to [**PolicyPolicyType**](PolicyPolicyType.md) |  | [optional] 
**Zone** | Pointer to [**PolicyZone**](PolicyZone.md) |  | [optional] 
**Site** | Pointer to [**PolicySite**](PolicySite.md) |  | [optional] 
**User** | Pointer to [**PolicyUser**](PolicyUser.md) |  | [optional] 
**Role** | Pointer to [**PolicyRole**](PolicyRole.md) |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableInt64** |  | [optional] 
**EachUser** | Pointer to **NullableBool** |  | [optional] 
**Config** | Pointer to [**PolicyConfig**](PolicyConfig.md) |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Owner** | Pointer to [**NullablePolicyOwner**](PolicyOwner.md) |  | [optional] 
**Accounts** | Pointer to [**[]ListPolicies200ResponseAllOfPoliciesInnerAccountsInner**](ListPolicies200ResponseAllOfPoliciesInnerAccountsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &Policy{
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


