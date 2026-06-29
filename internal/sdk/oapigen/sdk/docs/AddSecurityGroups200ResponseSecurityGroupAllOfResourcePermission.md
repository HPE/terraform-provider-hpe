# AddSecurityGroups200ResponseSecurityGroupAllOfResourcePermission

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DefaultStore** | Pointer to **bool** |  | [optional] 
**DefaultTarget** | Pointer to **bool** |  | [optional] 
**CanManage** | Pointer to **bool** |  | [optional] 
**All** | Pointer to **bool** |  | [optional] 
**Account** | Pointer to [**AddSecurityGroups200ResponseSecurityGroupAllOfResourcePermissionAccount**](AddSecurityGroups200ResponseSecurityGroupAllOfResourcePermissionAccount.md) |  | [optional] 
**Sites** | Pointer to [**[]AddSecurityGroups200ResponseSecurityGroupAllOfResourcePermissionSitesInner**](AddSecurityGroups200ResponseSecurityGroupAllOfResourcePermissionSitesInner.md) |  | [optional] 
**AllPlans** | Pointer to **bool** |  | [optional] 
**Plans** | Pointer to [**[]AddSecurityGroups200ResponseSecurityGroupAllOfResourcePermissionPlansInner**](AddSecurityGroups200ResponseSecurityGroupAllOfResourcePermissionPlansInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddSecurityGroups200ResponseSecurityGroupAllOfResourcePermission{
    // Set fields directly
}
```

### Sites (Nullable)

Use the Nullable wrapper methods:
- `obj.Sites.IsSet()` — check if set
- `obj.Sites.Get()` — get the inner value (returns pointer)
- `obj.Sites.Set(&val)` — set the value
- `obj.Sites.Unset()` — clear the value
### Plans (Nullable)

Use the Nullable wrapper methods:
- `obj.Plans.IsSet()` — check if set
- `obj.Plans.Get()` — get the inner value (returns pointer)
- `obj.Plans.Set(&val)` — set the value
- `obj.Plans.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


