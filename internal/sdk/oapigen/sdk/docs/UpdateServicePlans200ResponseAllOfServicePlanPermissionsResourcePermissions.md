# UpdateServicePlans200ResponseAllOfServicePlanPermissionsResourcePermissions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DefaultStore** | Pointer to **bool** |  | [optional] 
**AllPlans** | Pointer to **bool** |  | [optional] 
**DefaultTarget** | Pointer to **bool** |  | [optional] 
**CanManage** | Pointer to **bool** |  | [optional] 
**All** | Pointer to **bool** |  | [optional] 
**Account** | Pointer to [**UpdateServicePlans200ResponseAllOfServicePlanPermissionsResourcePermissionsAccount**](UpdateServicePlans200ResponseAllOfServicePlanPermissionsResourcePermissionsAccount.md) |  | [optional] 
**Sites** | Pointer to [**[]UpdateServicePlans200ResponseAllOfServicePlanPermissionsResourcePermissionsSitesInner**](UpdateServicePlans200ResponseAllOfServicePlanPermissionsResourcePermissionsSitesInner.md) |  | [optional] 
**Plans** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateServicePlans200ResponseAllOfServicePlanPermissionsResourcePermissions{
    // Set fields directly
}
```

### Plans (Nullable)

Use the Nullable wrapper methods:
- `obj.Plans.IsSet()` — check if set
- `obj.Plans.Get()` — get the inner value (returns pointer)
- `obj.Plans.Set(&val)` — set the value
- `obj.Plans.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


