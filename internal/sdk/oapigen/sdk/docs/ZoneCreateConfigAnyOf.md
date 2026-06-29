# ZoneCreateConfigAnyOf

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApplianceUrl** | Pointer to **string** | The URL used by workloads provisioned in the cloud for interacting with the Morpheus appliance. | [optional] 
**DatacenterName** | Pointer to **string** | A custom name used to reference the datacenter for the cloud. | [optional] 
**ExternalId** | Pointer to **NullableString** | The external id of the cloud | [optional] 
**InventoryLevel** | Pointer to **string** | Whether to import existing virtual machines. | [optional] 
**ConsoleKeymap** | Pointer to **string** | The keyboard layout to use for the console | [optional] 
**Endpoint** | **string** | AWS endpoint | 
**AccessKey** | Pointer to **string** | AWS access key | [optional] 
**SecretKey** | Pointer to **string** | AWS secret key | [optional] 
**UseHostCredentials** | Pointer to **string** | Whether to use the IAM profile associated with the Morpheus server or not | [optional] [default to "on"]
**EbsEncryption** | Pointer to **string** | Determines whether to configure default EBS volume encryption or not | [optional] [default to "on"]
**StsAssumeRole** | Pointer to **string** | The AWS IAM role ARN to assume for authentication | [optional] 
**ConfigManagementId** | Pointer to **string** | The id of the configuration management integration associated with the AWS cloud | [optional] 
**Vpc** | Pointer to **string** | The VPC ID for a specific VPC | [optional] 
**VdiGateway** | Pointer to **string** | The VDI gateway for this cloud to use for provisioning virtual desktops. | [optional] 
**CmdbConfig** | Pointer to **string** | The CMDB configuration for this cloud to use for syncing with the CMDB. | [optional] 
**ChangeManagementConfig** | Pointer to **string** | The change management configuration for this cloud to use for syncing with the change management system. | [optional] 
**NetworkMode** | Pointer to **NullableString** | Whether to use public or private IP addresses for provisioning and managing instances in this cloud. | [optional] 
**ApiProxy** | Pointer to **NullableString** | The API proxy to use for API calls to the cloud. | [optional] 
**Region** | Pointer to **string** | The AWS region to use for this cloud. | [optional] 
**Credentials** | Pointer to **string** | The credentials to use for this cloud. | [optional] 
**CostingBucket** | Pointer to **string** | The S3 bucket to use for storing costing reports. | [optional] 
**CostingFolder** | Pointer to **NullableString** | The folder within the S3 bucket to use for storing costing reports. | [optional] 
**CostingReportName** | Pointer to **NullableString** | The name of the costing report to generate. | [optional] 
**CostingKey** | Pointer to **NullableString** | The AWS access key to use for generating costing reports. | [optional] 
**CostingSecret** | Pointer to **NullableString** | The AWS secret key to use for generating costing reports. | [optional] 
**Domain** | Pointer to **string** | The domain to use for this cloud. | [optional] 
**Timezone** | Pointer to **string** | The timezone to use for this cloud. | [optional] 
**SecurityServer** | Pointer to **string** | The security server to use for this cloud, or &#39;off&#39; to not use a security server. | [optional] 
**Guidance** | Pointer to **string** | Optional guidance field if you want to put more info there | [optional] 
**Costing** | Pointer to **string** | Whether to enable costing for this cloud or not. | [optional] 
**ConfigCmdbDiscovery** | Pointer to **string** | The CMDB discovery configuration to use for this cloud. | [optional] 
**Logo** | Pointer to **NullableString** | The logo to use for this cloud. | [optional] 
**DarkModeLogo** | Pointer to **NullableString** | The logo to use for this cloud in dark mode. | [optional] 
**Proxy** | Pointer to **string** | The proxy to use for this cloud. | [optional] 
**BypassProxyForCloud** | Pointer to **string** | Whether to bypass the proxy for API calls to the cloud or not. | [optional] 
**NoProxy** | Pointer to **string** | A comma separated list of hosts to bypass the proxy for when making API calls to the cloud. | [optional] 
**UserDataLinux** | Pointer to **string** | The user data script to use for provisioning instances in this cloud. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ZoneCreateConfigAnyOf{
    // Set fields directly
}
```

### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### NetworkMode (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkMode.IsSet()` — check if set
- `obj.NetworkMode.Get()` — get the inner value (returns pointer)
- `obj.NetworkMode.Set(&val)` — set the value
- `obj.NetworkMode.Unset()` — clear the value
### ApiProxy (Nullable)

Use the Nullable wrapper methods:
- `obj.ApiProxy.IsSet()` — check if set
- `obj.ApiProxy.Get()` — get the inner value (returns pointer)
- `obj.ApiProxy.Set(&val)` — set the value
- `obj.ApiProxy.Unset()` — clear the value
### CostingFolder (Nullable)

Use the Nullable wrapper methods:
- `obj.CostingFolder.IsSet()` — check if set
- `obj.CostingFolder.Get()` — get the inner value (returns pointer)
- `obj.CostingFolder.Set(&val)` — set the value
- `obj.CostingFolder.Unset()` — clear the value
### CostingReportName (Nullable)

Use the Nullable wrapper methods:
- `obj.CostingReportName.IsSet()` — check if set
- `obj.CostingReportName.Get()` — get the inner value (returns pointer)
- `obj.CostingReportName.Set(&val)` — set the value
- `obj.CostingReportName.Unset()` — clear the value
### CostingKey (Nullable)

Use the Nullable wrapper methods:
- `obj.CostingKey.IsSet()` — check if set
- `obj.CostingKey.Get()` — get the inner value (returns pointer)
- `obj.CostingKey.Set(&val)` — set the value
- `obj.CostingKey.Unset()` — clear the value
### CostingSecret (Nullable)

Use the Nullable wrapper methods:
- `obj.CostingSecret.IsSet()` — check if set
- `obj.CostingSecret.Get()` — get the inner value (returns pointer)
- `obj.CostingSecret.Set(&val)` — set the value
- `obj.CostingSecret.Unset()` — clear the value
### Logo (Nullable)

Use the Nullable wrapper methods:
- `obj.Logo.IsSet()` — check if set
- `obj.Logo.Get()` — get the inner value (returns pointer)
- `obj.Logo.Set(&val)` — set the value
- `obj.Logo.Unset()` — clear the value
### DarkModeLogo (Nullable)

Use the Nullable wrapper methods:
- `obj.DarkModeLogo.IsSet()` — check if set
- `obj.DarkModeLogo.Get()` — get the inner value (returns pointer)
- `obj.DarkModeLogo.Set(&val)` — set the value
- `obj.DarkModeLogo.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


