# InstanceCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Instance** | [**InstanceCreateInstance**](InstanceCreateInstance.md) |  | 
**ZoneId** | Pointer to **int64** | The Cloud ID to provision the instance onto. | [optional] 
**Evars** | Pointer to [**[]InstanceCreateEvarsInner**](InstanceCreateEvarsInner.md) | Environment Variables, an array of objects that have name and value. | [optional] 
**Copies** | Pointer to **int64** | Number of copies to provision. | [optional] [default to 1]
**LayoutSize** | Pointer to **int64** | Apply a multiply factor of containers/vms within the instance. | [optional] [default to 1]
**ServicePlanOptions** | Pointer to [**InstanceCreateServicePlanOptions**](InstanceCreateServicePlanOptions.md) |  | [optional] 
**SecurityGroups** | Pointer to [**[]InstanceCreateSecurityGroupsInner**](InstanceCreateSecurityGroupsInner.md) | Key for security group configuration. It should be passed as an array of objects containing the id of the security group to assign the instance to. | [optional] 
**Volumes** | Pointer to [**[]InstanceCreateVolumesInner**](InstanceCreateVolumesInner.md) | The (optional) volumes parameter is for LV configuration, can create additional LVs at provision It should be passed as an array of | [optional] 
**NetworkInterfaces** | Pointer to [**[]InstancesNetworkInterfaces9**](InstancesNetworkInterfaces9.md) | The networkInterfaces parameter is for network configuration.  The Options API &#x60;/api/options/zoneNetworkOptions?zoneId&#x3D;5&amp;provisionTypeId&#x3D;10&#x60; can be used to see which options are available.  | [optional] 
**Config** | [**InstanceCreateConfig**](InstanceCreateConfig.md) |  | 
**Labels** | Pointer to **[]string** | Array of strings (keywords). | [optional] 
**Tags** | Pointer to [**[]InstanceCreateTagsInner**](InstanceCreateTagsInner.md) | Metadata tags, Array of objects having a name and value. | [optional] 
**Metadata** | Pointer to [**[]InstanceCreateMetadataInner**](InstanceCreateMetadataInner.md) | Alias for &#x60;tags&#x60;. | [optional] 
**Ports** | Pointer to [**[]InstanceCreatePortsInner**](InstanceCreatePortsInner.md) | The ports parameter is for port configuration.  The layout may have default ports, which are defined in node types, that are always configured. This parameter will be for additional custom ports to be opened.  | [optional] 
**TaskSetId** | Pointer to **int64** | The Workflow ID to execute. | [optional] 
**TaskSetName** | Pointer to **string** | The Workflow Name to execute. | [optional] 
**ServerUUIDs** | Pointer to **[]string** | Server UUIDs is an optional array of uuid string values to assign to the new servers being provisioned for this instance. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceCreate{
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


