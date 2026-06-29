# AddServicePlansRequestServicePlanPermissionsResourcePermissions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllSites** | Pointer to **bool** | Gives all groups (sites) permission to Provision and Reconfigure with servicePlan. | [optional] [default to false]
**Sites** | Pointer to [**[]AddServicePlansRequestServicePlanPermissionsResourcePermissionsSitesInner**](AddServicePlansRequestServicePlanPermissionsResourcePermissionsSitesInner.md) | Group (site) permissions that determines what Groups the Service Plan will be available in for Provisioning and Reconfigure. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddServicePlansRequestServicePlanPermissionsResourcePermissions{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


