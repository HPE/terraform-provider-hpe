# AddCloudsRequestZone

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A unique name scoped to your account for the cloud | 
**GroupId** | **int64** | Specifies which Server group this cloud should be assigned to | 
**ZoneType** | [**AddCloudsRequestZoneZoneType**](AddCloudsRequestZoneZoneType.md) |  | 
**Config** | [**AddCloudsRequestZoneConfig**](AddCloudsRequestZoneConfig.md) |  | 
**AgentMode** | Pointer to **string** | The method used to install the Morpheus agent on virtual machines provisioned in the cloud (ssh, cloudInit). | [optional] [default to "cloudInit"]
**Description** | Pointer to **string** | Optional description field if you want to put more info there | [optional] 
**Code** | Pointer to **string** | Optional code for use with policies | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Location** | Pointer to **NullableString** | Optional location for your cloud | [optional] 
**Visibility** | Pointer to **string** | The visibility of the cloud (private or public) | [optional] [default to "private"]
**AccountId** | Pointer to **int64** | Specifies which Tenant this cloud should be assigned to | [optional] 
**Enabled** | Pointer to **bool** | Can be used to disable the cloud | [optional] [default to true]
**AutoRecoverPowerState** | Pointer to **bool** | Automatically Power on VMs | [optional] [default to false]
**ScalePriority** | Pointer to **int64** | Scale Priority | [optional] [default to 1]
**DefaultDatastoreSyncActive** | Pointer to **bool** | Sets the default active state during discovery of new datastores. | [optional] 
**DefaultNetworkSyncActive** | Pointer to **bool** | Sets the default active state during discovery of new networks. | [optional] 
**DefaultFolderSyncActive** | Pointer to **bool** | Sets the default active state during discovery of new folders. | [optional] 
**DefaultSecurityGroupSyncActive** | Pointer to **bool** | Sets the default active state during discovery of new security groups. | [optional] 
**DefaultPoolSyncActive** | Pointer to **bool** | Sets the default active state during discovery of new resource pools. | [optional] 
**DefaultPlanSyncActive** | Pointer to **bool** | Sets the default active state during discovery of new plans. | [optional] 
**LinkedAccountId** | Pointer to **int64** | Linked Account ID (enter commercial ID to get costing for AWS Govcloud) | [optional] 
**SecurityMode** | Pointer to **string** | host firewall. &#x60;off&#x60; or &#x60;internal&#x60;. a.k.a. \&quot;local firewall\&quot; | [optional] [default to "off"]
**Credential** | Pointer to [**AddCloudsRequestZoneCredential**](AddCloudsRequestZoneCredential.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCloudsRequestZone{
    // Set fields directly
}
```

### Location (Nullable)

Use the Nullable wrapper methods:
- `obj.Location.IsSet()` — check if set
- `obj.Location.Get()` — get the inner value (returns pointer)
- `obj.Location.Set(&val)` — set the value
- `obj.Location.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


