# ValidateBackupCreate200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | Pointer to **bool** | Whether the validation passed | [optional] 
**Errors** | Pointer to **map[string]string** | Map of field names to error messages | [optional] 
**Backup** | Pointer to [**ValidateBackupCreate200ResponseBackup**](ValidateBackupCreate200ResponseBackup.md) |  | [optional] 
**ZoneId** | Pointer to **NullableInt64** | The cloud (zone) ID associated with the instance. Present when locationType is &#x60;instance&#x60;. | [optional] 
**Instance** | Pointer to [**ValidateBackupCreate200ResponseInstance**](ValidateBackupCreate200ResponseInstance.md) |  | [optional] 
**Containers** | Pointer to [**[]ValidateBackupCreate200ResponseContainersInner**](ValidateBackupCreate200ResponseContainersInner.md) | List of containers (workloads) belonging to the instance. Present when locationType is &#x60;instance&#x60;. | [optional] 
**Server** | Pointer to [**ValidateBackupCreate200ResponseServer**](ValidateBackupCreate200ResponseServer.md) |  | [optional] 
**BackupTypes** | Pointer to [**[]ValidateBackupCreate200ResponseBackupTypesInner**](ValidateBackupCreate200ResponseBackupTypesInner.md) | List of available backup type codes for the resolved instance or server. | [optional] 
**AccountId** | Pointer to **NullableInt64** | The account ID. Present when locationType is &#x60;instance&#x60;. | [optional] 
**ProviderBackupTypes** | Pointer to [**[]ValidateBackupCreate200ResponseProviderBackupTypesInner**](ValidateBackupCreate200ResponseProviderBackupTypesInner.md) | List of external backup provider types available for this instance (e.g. Veeam, Commvault). Present when locationType is &#x60;instance&#x60;. | [optional] 
**VeeamSupported** | Pointer to **bool** | Whether Veeam backup is available for this instance. | [optional] 
**CommvaultSupported** | Pointer to **bool** | Whether Commvault backup is available for this instance. | [optional] 
**BackupSettings** | Pointer to [**ValidateBackupCreate200ResponseBackupSettings**](ValidateBackupCreate200ResponseBackupSettings.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ValidateBackupCreate200Response{
    // Set fields directly
}
```

### ZoneId (Nullable)

Use the Nullable wrapper methods:
- `obj.ZoneId.IsSet()` — check if set
- `obj.ZoneId.Get()` — get the inner value (returns pointer)
- `obj.ZoneId.Set(&val)` — set the value
- `obj.ZoneId.Unset()` — clear the value
### AccountId (Nullable)

Use the Nullable wrapper methods:
- `obj.AccountId.IsSet()` — check if set
- `obj.AccountId.Get()` — get the inner value (returns pointer)
- `obj.AccountId.Set(&val)` — set the value
- `obj.AccountId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


