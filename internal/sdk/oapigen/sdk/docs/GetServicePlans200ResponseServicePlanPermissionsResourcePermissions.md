# GetServicePlans200ResponseServicePlanPermissionsResourcePermissions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DefaultStore** | Pointer to **bool** |  | [optional] 
**AllPlans** | Pointer to **bool** |  | [optional] 
**DefaultTarget** | Pointer to **bool** |  | [optional] 
**CanManage** | Pointer to **bool** |  | [optional] 
**All** | Pointer to **bool** |  | [optional] 
**Account** | Pointer to [**GetServicePlans200ResponseServicePlanPermissionsResourcePermissionsAccount**](GetServicePlans200ResponseServicePlanPermissionsResourcePermissionsAccount.md) |  | [optional] 
**Sites** | Pointer to [**[]GetServicePlans200ResponseServicePlanPermissionsResourcePermissionsSitesInner**](GetServicePlans200ResponseServicePlanPermissionsResourcePermissionsSitesInner.md) |  | [optional] 
**Plans** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetServicePlans200ResponseServicePlanPermissionsResourcePermissions{
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


