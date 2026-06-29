# GetMigration200ResponseMigration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Migration ID | [optional] 
**Name** | Pointer to **string** | Name | [optional] 
**Status** | Pointer to **string** | Migration Status. The possible status values are: &#39;pending&#39;, &#39;running&#39;, &#39;failed&#39;, &#39;completed&#39; | [optional] 
**StatusMessage** | Pointer to **NullableString** | Status Message | [optional] 
**SkippedPrechecks** | Pointer to **bool** | Indicates if the precheck was skipped | [optional] 
**InstallGuestTools** | Pointer to **bool** | Indicates if guest tools should be installed | [optional] 
**ReInitializeServerOnMigration** | Pointer to **bool** | Indicates if the server should be re-initialized on migration | [optional] 
**LinuxUsername** | Pointer to **NullableString** | Linux Username for migrated servers | [optional] 
**LinuxPassword** | Pointer to **NullableString** | Linux Password for migrated servers | [optional] 
**LinuxKeyPair** | Pointer to [**GetMigration200ResponseMigrationLinuxKeyPair**](GetMigration200ResponseMigrationLinuxKeyPair.md) |  | [optional] 
**WindowsUsername** | Pointer to **NullableString** | Windows Username for migrated servers | [optional] 
**WindowsPassword** | Pointer to **NullableString** | Windows Password for migrated servers | [optional] 
**SourceCloud** | Pointer to [**GetMigration200ResponseMigrationSourceCloud**](GetMigration200ResponseMigrationSourceCloud.md) |  | [optional] 
**TargetCloud** | Pointer to [**GetMigration200ResponseMigrationTargetCloud**](GetMigration200ResponseMigrationTargetCloud.md) |  | [optional] 
**TargetGroup** | Pointer to [**GetMigration200ResponseMigrationTargetGroup**](GetMigration200ResponseMigrationTargetGroup.md) |  | [optional] 
**TargetPool** | Pointer to [**GetMigration200ResponseMigrationTargetPool**](GetMigration200ResponseMigrationTargetPool.md) |  | [optional] 
**Servers** | Pointer to [**[]GetMigration200ResponseMigrationServersInner**](GetMigration200ResponseMigrationServersInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetMigration200ResponseMigration{
    // Set fields directly
}
```

### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
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


