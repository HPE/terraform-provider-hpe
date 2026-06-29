# CatalogItemTypeInstanceScribe

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Group** | [**CatalogItemTypeInstanceScribeGroup**](CatalogItemTypeInstanceScribeGroup.md) |  | 
**Cloud** | [**CatalogItemTypeInstanceScribeCloud**](CatalogItemTypeInstanceScribeCloud.md) |  | 
**Type** | **string** | The type of instance by code we want to fetch. | 
**Name** | **string** | Name of the instance to be created. | 
**Config** | [**CatalogItemTypeInstanceScribeConfig**](CatalogItemTypeInstanceScribeConfig.md) |  | 
**Volumes** | [**[]CatalogItemTypeInstanceScribeVolumesInner**](CatalogItemTypeInstanceScribeVolumesInner.md) | The (optional) volumes parameter is for LV configuration, can create additional LVs at provision It should be passed as an array of | 
**HostName** | Pointer to **string** | Hostname of the instance to be created.  Can be the same as the instance name. | [optional] 
**Environment** | Pointer to **string** | Environment code | [optional] 
**Layout** | [**CatalogItemTypeInstanceScribeLayout**](CatalogItemTypeInstanceScribeLayout.md) |  | 
**Plan** | [**CatalogItemTypeInstanceScribePlan**](CatalogItemTypeInstanceScribePlan.md) |  | 
**Version** | Pointer to **string** | Version of the layout to create. | [optional] 
**Evars** | Pointer to [**[]CatalogItemTypeInstanceScribeEvarsInner**](CatalogItemTypeInstanceScribeEvarsInner.md) | Environment Variables, an array of objects that have name and value. | [optional] 
**ServicePlanOptions** | Pointer to [**CatalogItemTypeInstanceScribeServicePlanOptions**](CatalogItemTypeInstanceScribeServicePlanOptions.md) |  | [optional] 
**SecurityGroups** | Pointer to [**[]CatalogItemTypeInstanceScribeSecurityGroupsInner**](CatalogItemTypeInstanceScribeSecurityGroupsInner.md) | Key for security group configuration. It should be passed as an array of objects containing the id of the security group to assign the instance to. | [optional] 
**NetworkInterfaces** | Pointer to [**[]InstancesNetworkInterfaces6**](InstancesNetworkInterfaces6.md) | The networkInterfaces parameter is for network configuration.  The Options API &#x60;/api/options/zoneNetworkOptions?zoneId&#x3D;5&amp;provisionTypeId&#x3D;10&#x60; can be used to see which options are available.  | [optional] 
**Labels** | Pointer to **[]string** | Array of strings (keywords). | [optional] 
**Tags** | Pointer to [**[]CatalogItemTypeInstanceScribeTagsInner**](CatalogItemTypeInstanceScribeTagsInner.md) | Metadata tags, Array of objects having a name and value. | [optional] 
**Metadata** | Pointer to [**[]CatalogItemTypeInstanceScribeMetadataInner**](CatalogItemTypeInstanceScribeMetadataInner.md) | Alias for &#x60;tags&#x60;. | [optional] 
**Ports** | Pointer to [**[]CatalogItemTypeInstanceScribePortsInner**](CatalogItemTypeInstanceScribePortsInner.md) | The ports parameter is for port configuration.  The layout may have default ports, which are defined in node types, that are always configured. This parameter will be for additional custom ports to be opened.  | [optional] 
**TaskSetId** | Pointer to **int64** | The Workflow ID to execute. | [optional] 
**TaskSetName** | Pointer to **string** | The Workflow Name to execute. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CatalogItemTypeInstanceScribe{
    // Set fields directly
}
```

### SecurityGroups (Nullable)

Use the Nullable wrapper methods:
- `obj.SecurityGroups.IsSet()` — check if set
- `obj.SecurityGroups.Get()` — get the inner value (returns pointer)
- `obj.SecurityGroups.Set(&val)` — set the value
- `obj.SecurityGroups.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


