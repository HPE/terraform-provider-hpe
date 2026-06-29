# ProvisioningSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllowZoneSelection** | Pointer to **bool** |  | [optional] 
**AllowServerSelection** | Pointer to **bool** |  | [optional] 
**RequireEnvironments** | Pointer to **bool** |  | [optional] 
**ShowPricing** | Pointer to **bool** |  | [optional] 
**HideDatastoreStats** | Pointer to **bool** |  | [optional] 
**CrossTenantNamingPolicies** | Pointer to **bool** |  | [optional] 
**ReuseSequence** | Pointer to **bool** |  | [optional] 
**CloudInitUsername** | Pointer to **string** |  | [optional] 
**CloudInitPassword** | Pointer to **string** |  | [optional] 
**CloudInitKeyPair** | Pointer to [**ProvisioningSettingsCloudInitKeyPair**](ProvisioningSettingsCloudInitKeyPair.md) |  | [optional] 
**WindowsPassword** | Pointer to **NullableString** |  | [optional] 
**PxeRootPassword** | Pointer to **NullableString** |  | [optional] 
**DefaultTemplateType** | Pointer to [**ProvisioningSettingsDefaultTemplateType**](ProvisioningSettingsDefaultTemplateType.md) |  | [optional] 
**DeployStorageProvider** | Pointer to [**ProvisioningSettingsDeployStorageProvider**](ProvisioningSettingsDeployStorageProvider.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ProvisioningSettings{
    // Set fields directly
}
```

### WindowsPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.WindowsPassword.IsSet()` — check if set
- `obj.WindowsPassword.Get()` — get the inner value (returns pointer)
- `obj.WindowsPassword.Set(&val)` — set the value
- `obj.WindowsPassword.Unset()` — clear the value
### PxeRootPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.PxeRootPassword.IsSet()` — check if set
- `obj.PxeRootPassword.Get()` — get the inner value (returns pointer)
- `obj.PxeRootPassword.Set(&val)` — set the value
- `obj.PxeRootPassword.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


