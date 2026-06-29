# ProvisioningSettingsUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllowZoneSelection** | Pointer to **bool** | Use this to enable / disable allowing cloud selection | [optional] 
**AllowServerSelection** | Pointer to **bool** | Use this to enable / disable allowing host selection | [optional] 
**RequireEnvironments** | Pointer to **bool** | Use this to enable / disable requiring environment selection | [optional] 
**ShowPricing** | Pointer to **bool** | Use this to enable / disable showing pricing | [optional] 
**HideDatastoreStats** | Pointer to **bool** | Use this to enable / disable hiding datastore stats | [optional] 
**CrossTenantNamingPolicies** | Pointer to **bool** | Use this to enable / disable cross-tenant naming policies | [optional] 
**ReuseSequence** | Pointer to **bool** | Use this to enable / disable reusing naming sequence numbers | [optional] 
**CloudInitUsername** | Pointer to **string** | Cloud-init username | [optional] 
**CloudInitPassword** | Pointer to **string** | Cloud-init password | [optional] 
**CloudInitKeyPair** | Pointer to [**ProvisioningSettingsUpdateCloudInitKeyPair**](ProvisioningSettingsUpdateCloudInitKeyPair.md) |  | [optional] 
**DeployStorageProvider** | Pointer to [**ProvisioningSettingsUpdateDeployStorageProvider**](ProvisioningSettingsUpdateDeployStorageProvider.md) |  | [optional] 
**WindowsPassword** | Pointer to **string** | Windows administrator password | [optional] 
**PxeRootPassword** | Pointer to **string** | PXE Boot default root password | [optional] 
**DefaultTemplateType** | Pointer to [**ProvisioningSettingsUpdateDefaultTemplateType**](ProvisioningSettingsUpdateDefaultTemplateType.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ProvisioningSettingsUpdate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


