# ListProvisioningSettings200ResponseProvisioningSettings

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
**CloudInitKeyPair** | Pointer to [**NullableListProvisioningSettings200ResponseProvisioningSettingsCloudInitKeyPair**](ListProvisioningSettings200ResponseProvisioningSettingsCloudInitKeyPair.md) |  | [optional] 
**WindowsPassword** | Pointer to **NullableString** |  | [optional] 
**PxeRootPassword** | Pointer to **NullableString** |  | [optional] 
**DefaultTemplateType** | Pointer to [**ListProvisioningSettings200ResponseProvisioningSettingsDefaultTemplateType**](ListProvisioningSettings200ResponseProvisioningSettingsDefaultTemplateType.md) |  | [optional] 
**DeployStorageProvider** | Pointer to [**NullableListProvisioningSettings200ResponseProvisioningSettingsDeployStorageProvider**](ListProvisioningSettings200ResponseProvisioningSettingsDeployStorageProvider.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListProvisioningSettings200ResponseProvisioningSettings{
    // Set fields directly
}
```

### CloudInitKeyPair (Nullable)

Use the Nullable wrapper methods:
- `obj.CloudInitKeyPair.IsSet()` — check if set
- `obj.CloudInitKeyPair.Get()` — get the inner value (returns pointer)
- `obj.CloudInitKeyPair.Set(&val)` — set the value
- `obj.CloudInitKeyPair.Unset()` — clear the value
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
### DeployStorageProvider (Nullable)

Use the Nullable wrapper methods:
- `obj.DeployStorageProvider.IsSet()` — check if set
- `obj.DeployStorageProvider.Get()` — get the inner value (returns pointer)
- `obj.DeployStorageProvider.Set(&val)` — set the value
- `obj.DeployStorageProvider.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


