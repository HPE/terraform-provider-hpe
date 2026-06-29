# AddMigrationRequestMigration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**SkippedPrechecks** | Pointer to **bool** | Indicates if the precheck should be skipped | [optional] [default to false]
**InstallGuestTools** | Pointer to **bool** | Indicates if guest tools should be installed | [optional] [default to true]
**ReInitializeServerOnMigration** | Pointer to **bool** | Indicates if the server should be re-initialized on migration | [optional] [default to false]
**LinuxUsername** | Pointer to **NullableString** | Linux Username for migrated servers | [optional] 
**LinuxPassword** | Pointer to **NullableString** | Linux Password for migrated servers | [optional] 
**LinuxKeyPair** | Pointer to [**AddMigrationRequestMigrationLinuxKeyPair**](AddMigrationRequestMigrationLinuxKeyPair.md) |  | [optional] 
**WindowsUsername** | Pointer to **NullableString** | Windows Username for migrated servers | [optional] 
**WindowsPassword** | Pointer to **NullableString** | Windows Password for migrated servers | [optional] 
**SourceCloudId** | Pointer to **int64** | Source Cloud ID. The API &#x60;/api/migrations/source-clouds&#x60; can be used to find available options.  | [optional] 
**TargetCloudId** | Pointer to **int64** | Target Cloud ID. The API &#x60;/api/migrations/target-clouds?sourceCloudId&#x3D;34&#x60; can be used to find available options.  | [optional] 
**TargetGroupId** | Pointer to **int64** | Target Group ID.  The Options API &#x60;/api/options/targetGroups?sourceCloudId&#x3D;34&amp;targetCloudId&#x3D;129&#x60; can be used to find available options.  | [optional] 
**TargetPoolId** | Pointer to [**AddMigrationRequestMigrationTargetPoolId**](AddMigrationRequestMigrationTargetPoolId.md) |  | [optional] 
**SourceServerIds** | Pointer to **[]int64** | Array of Server IDs to be migrated. The API &#x60;/api/migrations/source-servers?sourceCloudId&#x3D;34&#x60; can be used to find available options.  | [optional] 
**Datastores** | Pointer to [**[]AddMigrationRequestMigrationDatastoresInner**](AddMigrationRequestMigrationDatastoresInner.md) | Array of datastore mappings. | [optional] 
**Networks** | Pointer to [**[]AddMigrationRequestMigrationNetworksInner**](AddMigrationRequestMigrationNetworksInner.md) | Array of network mappings. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddMigrationRequestMigration{
    // Set fields directly
}
```

### LinuxUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.LinuxUsername.IsSet()` — check if set
- `obj.LinuxUsername.Get()` — get the inner value (returns pointer)
- `obj.LinuxUsername.Set(&val)` — set the value
- `obj.LinuxUsername.Unset()` — clear the value
### LinuxPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.LinuxPassword.IsSet()` — check if set
- `obj.LinuxPassword.Get()` — get the inner value (returns pointer)
- `obj.LinuxPassword.Set(&val)` — set the value
- `obj.LinuxPassword.Unset()` — clear the value
### WindowsUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.WindowsUsername.IsSet()` — check if set
- `obj.WindowsUsername.Get()` — get the inner value (returns pointer)
- `obj.WindowsUsername.Set(&val)` — set the value
- `obj.WindowsUsername.Unset()` — clear the value
### WindowsPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.WindowsPassword.IsSet()` — check if set
- `obj.WindowsPassword.Get()` — get the inner value (returns pointer)
- `obj.WindowsPassword.Set(&val)` — set the value
- `obj.WindowsPassword.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


