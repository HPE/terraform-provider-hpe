# ZoneCreateConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApplianceUrl** | Pointer to **string** | The URL used by workloads provisioned in the cloud for interacting with the Morpheus appliance. | [optional] 
**DatacenterName** | Pointer to **string** | A custom name used to reference the datacenter for the cloud. | [optional] 
**ExternalId** | Pointer to **string** | The external id of the cloud | [optional] 
**InventoryLevel** | Pointer to **string** | Whether to import existing virtual machines. | [optional] 
**ConsoleKeymap** | Pointer to **string** | The keyboard layout to use for the console | [optional] 
**Endpoint** | **string** | AWS endpoint | 
**AccessKey** | Pointer to **string** | AWS access key | [optional] 
**SecretKey** | Pointer to **string** | AWS secret key | [optional] 
**UseHostCredentials** | Pointer to **string** | Whether to use the IAM profile associated with the Morpheus server or not | [optional] [default to "on"]
**EbsEncryption** | Pointer to **string** | Determines whether to configure default EBS volume encryption or not | [optional] [default to "on"]
**StsAssumeRole** | Pointer to **string** | The AWS IAM role ARN to assume for authentication | [optional] 
**ConfigManagementId** | Pointer to **string** | The id of the configuration management integration associated with the vSphere cloud. | [optional] 
**Vpc** | Pointer to **string** | The VPC ID for a specific VPC | [optional] 
**VdiGateway** | Pointer to **string** | The VDI gateway for this cloud to use for provisioning virtual desktops. | [optional] 
**CmdbConfig** | Pointer to **string** | The CMDB configuration for this cloud to use for syncing with the CMDB. | [optional] 
**ChangeManagementConfig** | Pointer to **string** | The change management configuration for this cloud to use for syncing with the change management system. | [optional] 
**NetworkMode** | Pointer to **string** | Whether to use public or private IP addresses for provisioning and managing instances in this cloud. | [optional] 
**ApiProxy** | Pointer to **string** | The API proxy to use for API calls to the cloud. | [optional] 
**Region** | Pointer to **string** | The AWS region to use for this cloud. | [optional] 
**Credentials** | Pointer to **string** | The credentials to use for this cloud. | [optional] 
**CostingBucket** | Pointer to **string** | The S3 bucket to use for storing costing reports. | [optional] 
**CostingFolder** | Pointer to **string** | The folder within the S3 bucket to use for storing costing reports. | [optional] 
**CostingReportName** | Pointer to **string** | The name of the costing report to generate. | [optional] 
**CostingKey** | Pointer to **string** | The AWS access key to use for generating costing reports. | [optional] 
**CostingSecret** | Pointer to **string** | The AWS secret key to use for generating costing reports. | [optional] 
**Domain** | Pointer to **string** | The domain to use for this cloud. | [optional] 
**Timezone** | Pointer to **string** | The timezone to use for this cloud. | [optional] 
**SecurityServer** | Pointer to **string** | The security server to use for this cloud, or &#39;off&#39; to not use a security server. | [optional] 
**Guidance** | Pointer to **string** | Optional guidance field if you want to put more info there | [optional] 
**Costing** | Pointer to **string** | Whether to enable costing for this cloud or not. | [optional] 
**ConfigCmdbDiscovery** | Pointer to **string** | The CMDB discovery configuration to use for this cloud. | [optional] 
**Logo** | Pointer to **string** | The logo to use for this cloud. | [optional] 
**DarkModeLogo** | Pointer to **string** | The logo to use for this cloud in dark mode. | [optional] 
**Proxy** | Pointer to **string** | The proxy to use for this cloud. | [optional] 
**BypassProxyForCloud** | Pointer to **string** | Whether to bypass the proxy for API calls to the cloud or not. | [optional] 
**NoProxy** | Pointer to **string** | A comma separated list of hosts to bypass the proxy for when making API calls to the cloud. | [optional] 
**UserDataLinux** | Pointer to **string** | The user data script to use for provisioning instances in this cloud. | [optional] 
**CloudType** | Pointer to **string** | The Azure cloud type (global, usgov, german, china). | [optional] [default to "global"]
**ImportExisting** | Pointer to **string** | Whether to import existing resources from the cloud (on, off). | [optional] 
**SubscriberId** | Pointer to **string** | Azure subscriber id | [optional] 
**TenantId** | Pointer to **string** | Azure tenant id | [optional] 
**ClientId** | Pointer to **string** | Azure client id | [optional] 
**ClientSecret** | Pointer to **string** | Azure client secret | [optional] 
**ResourceGroup** | Pointer to **string** | Azure resource group | [optional] 
**StorageAccount** | Pointer to **string** | The Azure storage account to use. | [optional] 
**RpcMode** | Pointer to [**NullableZoneCreateConfigAnyOf3RpcMode**](ZoneCreateConfigAnyOf3RpcMode.md) |  | [optional] 
**CertificateProvider** | Pointer to **string** | Certificate provider | [optional] [default to "internal"]
**EnableNetworkTypeSelection** | Pointer to [**NullableZoneCreateConfigAnyOf3EnableNetworkTypeSelection**](ZoneCreateConfigAnyOf3EnableNetworkTypeSelection.md) |  | [optional] 
**ApiUrl** | **string** | The SDK URL of the vCenter server. | 
**ApiVersion** | **string** | The SDK version of the vCenter server. | 
**Datacenter** | **string** | The vSphere datacenter to add. | 
**Cluster** | Pointer to **string** | The name of the vSphere cluster | [optional] [default to "all"]
**ResourcePool** | Pointer to **string** | The name of the vSphere resource pool | [optional] 
**StorageType** | Pointer to **string** | The default vSphere VMDK type for virtual machines | [optional] [default to "thin"]
**EnableVnc** | Pointer to [**NullableZoneCreateConfigAnyOf3EnableVnc**](ZoneCreateConfigAnyOf3EnableVnc.md) |  | [optional] 
**HideHostSelection** | Pointer to [**NullableZoneCreateConfigAnyOf3HideHostSelection**](ZoneCreateConfigAnyOf3HideHostSelection.md) |  | [optional] 
**EnableDiskTypeSelection** | Pointer to [**NullableZoneCreateConfigAnyOf3EnableDiskTypeSelection**](ZoneCreateConfigAnyOf3EnableDiskTypeSelection.md) |  | [optional] 
**EnableStorageTypeSelection** | Pointer to [**NullableZoneCreateConfigAnyOf3EnableStorageTypeSelection**](ZoneCreateConfigAnyOf3EnableStorageTypeSelection.md) |  | [optional] 
**Username** | Pointer to **string** | Username. | [optional] 
**Password** | Pointer to **string** | Password to apply to the user | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ZoneCreateConfig{
    // Set fields directly
}
```

### RpcMode (Nullable)

Use the Nullable wrapper methods:
- `obj.RpcMode.IsSet()` — check if set
- `obj.RpcMode.Get()` — get the inner value (returns pointer)
- `obj.RpcMode.Set(&val)` — set the value
- `obj.RpcMode.Unset()` — clear the value
### EnableNetworkTypeSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableNetworkTypeSelection.IsSet()` — check if set
- `obj.EnableNetworkTypeSelection.Get()` — get the inner value (returns pointer)
- `obj.EnableNetworkTypeSelection.Set(&val)` — set the value
- `obj.EnableNetworkTypeSelection.Unset()` — clear the value
### EnableVnc (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableVnc.IsSet()` — check if set
- `obj.EnableVnc.Get()` — get the inner value (returns pointer)
- `obj.EnableVnc.Set(&val)` — set the value
- `obj.EnableVnc.Unset()` — clear the value
### HideHostSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.HideHostSelection.IsSet()` — check if set
- `obj.HideHostSelection.Get()` — get the inner value (returns pointer)
- `obj.HideHostSelection.Set(&val)` — set the value
- `obj.HideHostSelection.Unset()` — clear the value
### EnableDiskTypeSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableDiskTypeSelection.IsSet()` — check if set
- `obj.EnableDiskTypeSelection.Get()` — get the inner value (returns pointer)
- `obj.EnableDiskTypeSelection.Set(&val)` — set the value
- `obj.EnableDiskTypeSelection.Unset()` — clear the value
### EnableStorageTypeSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableStorageTypeSelection.IsSet()` — check if set
- `obj.EnableStorageTypeSelection.Get()` — get the inner value (returns pointer)
- `obj.EnableStorageTypeSelection.Set(&val)` — set the value
- `obj.EnableStorageTypeSelection.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


