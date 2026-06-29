# GetBackupType200ResponseBackupType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Backup Type ID | [optional] 
**Code** | Pointer to **string** | Unique code identifier for the backup type | [optional] 
**Name** | Pointer to **string** | Display name of the backup type | [optional] 
**BackupFormat** | Pointer to **NullableString** | The format of the backup output | [optional] 
**ContainerFormat** | Pointer to **NullableString** | The container format for the backup | [optional] 
**ProviderCode** | Pointer to **NullableString** | The code of the backup provider that handles this type (e.g., &#x60;commvault&#x60;, &#x60;veeam&#x60;) | [optional] 
**ContainerType** | Pointer to **NullableString** | Which containers this backup type applies to (&#x60;all&#x60;, &#x60;single&#x60;, or null) | [optional] 
**ContainerCategory** | Pointer to **NullableString** | Comma-separated container categories this type is limited to | [optional] 
**RestoreType** | Pointer to **NullableString** | The restore type identifier | [optional] 
**Active** | Pointer to **bool** | Whether this backup type is active and available for use | [optional] 
**HasCopyToStore** | Pointer to **bool** | Whether this backup type supports copying results to a storage provider | [optional] 
**HasStreamToStore** | Pointer to **bool** | Whether this backup type supports streaming results directly to a storage provider | [optional] 
**CopyToStore** | Pointer to **bool** | Whether backup results are copied to the storage provider by default | [optional] 
**DownloadEnabled** | Pointer to **bool** | Whether backup results can be downloaded | [optional] 
**DownloadFromStoreOnly** | Pointer to **bool** | Whether downloads are only available from the storage provider | [optional] 
**RestoreExistingEnabled** | Pointer to **bool** | Whether restoring to the existing instance/server is supported | [optional] 
**RestoreNewEnabled** | Pointer to **bool** | Whether restoring to a new instance/server is supported | [optional] 
**ViewSet** | Pointer to **NullableString** | The view set used for rendering backup type specific UI | [optional] 
**RestoreNewMode** | Pointer to **NullableString** | The mode used when restoring to a new target | [optional] 
**PruneResultsOnRestoreExisting** | Pointer to **bool** | Whether existing backup results are pruned when restoring to an existing target | [optional] 
**RestrictTargets** | Pointer to **bool** | Whether restore targets are restricted | [optional] 
**IsPlugin** | Pointer to **NullableBool** | Whether this backup type is provided by an external plugin | [optional] 
**IsEmbedded** | Pointer to **bool** | Whether this backup type is embedded (bundled with the platform) | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetBackupType200ResponseBackupType{
    // Set fields directly
}
```

### BackupFormat (Nullable)

Use the Nullable wrapper methods:
- `obj.BackupFormat.IsSet()` — check if set
- `obj.BackupFormat.Get()` — get the inner value (returns pointer)
- `obj.BackupFormat.Set(&val)` — set the value
- `obj.BackupFormat.Unset()` — clear the value
### ContainerFormat (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerFormat.IsSet()` — check if set
- `obj.ContainerFormat.Get()` — get the inner value (returns pointer)
- `obj.ContainerFormat.Set(&val)` — set the value
- `obj.ContainerFormat.Unset()` — clear the value
### ProviderCode (Nullable)

Use the Nullable wrapper methods:
- `obj.ProviderCode.IsSet()` — check if set
- `obj.ProviderCode.Get()` — get the inner value (returns pointer)
- `obj.ProviderCode.Set(&val)` — set the value
- `obj.ProviderCode.Unset()` — clear the value
### ContainerType (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerType.IsSet()` — check if set
- `obj.ContainerType.Get()` — get the inner value (returns pointer)
- `obj.ContainerType.Set(&val)` — set the value
- `obj.ContainerType.Unset()` — clear the value
### ContainerCategory (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerCategory.IsSet()` — check if set
- `obj.ContainerCategory.Get()` — get the inner value (returns pointer)
- `obj.ContainerCategory.Set(&val)` — set the value
- `obj.ContainerCategory.Unset()` — clear the value
### RestoreType (Nullable)

Use the Nullable wrapper methods:
- `obj.RestoreType.IsSet()` — check if set
- `obj.RestoreType.Get()` — get the inner value (returns pointer)
- `obj.RestoreType.Set(&val)` — set the value
- `obj.RestoreType.Unset()` — clear the value
### ViewSet (Nullable)

Use the Nullable wrapper methods:
- `obj.ViewSet.IsSet()` — check if set
- `obj.ViewSet.Get()` — get the inner value (returns pointer)
- `obj.ViewSet.Set(&val)` — set the value
- `obj.ViewSet.Unset()` — clear the value
### RestoreNewMode (Nullable)

Use the Nullable wrapper methods:
- `obj.RestoreNewMode.IsSet()` — check if set
- `obj.RestoreNewMode.Get()` — get the inner value (returns pointer)
- `obj.RestoreNewMode.Set(&val)` — set the value
- `obj.RestoreNewMode.Unset()` — clear the value
### IsPlugin (Nullable)

Use the Nullable wrapper methods:
- `obj.IsPlugin.IsSet()` — check if set
- `obj.IsPlugin.Get()` — get the inner value (returns pointer)
- `obj.IsPlugin.Set(&val)` — set the value
- `obj.IsPlugin.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


