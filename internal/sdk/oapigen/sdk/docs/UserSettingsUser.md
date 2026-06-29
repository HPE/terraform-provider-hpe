# UserSettingsUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Username** | Pointer to **string** |  | [optional] 
**FirstName** | Pointer to **string** |  | [optional] 
**LastName** | Pointer to **string** |  | [optional] 
**Email** | Pointer to **string** |  | [optional] 
**LinuxUsername** | Pointer to **string** |  | [optional] 
**LinuxPassword** | Pointer to **string** |  | [optional] 
**LinuxKeyPairId** | Pointer to **NullableInt64** |  | [optional] 
**WindowsUsername** | Pointer to **string** |  | [optional] 
**WindowsPassword** | Pointer to **string** |  | [optional] 
**Avatar** | Pointer to **string** |  | [optional] 
**DesktopBackground** | Pointer to **string** |  | [optional] 
**ReceiveNotifications** | Pointer to **bool** |  | [optional] 
**DefaultGroup** | Pointer to [**UserSettingsUserDefaultGroup**](UserSettingsUserDefaultGroup.md) |  | [optional] 
**DefaultCloud** | Pointer to [**UserSettingsUserDefaultCloud**](UserSettingsUserDefaultCloud.md) |  | [optional] 
**DefaultPersona** | Pointer to [**UserSettingsUserDefaultPersona**](UserSettingsUserDefaultPersona.md) |  | [optional] 
**IsUsing2FA** | Pointer to **bool** |  | [optional] 
**Tenant** | Pointer to [**UserSettingsUserTenant**](UserSettingsUserTenant.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UserSettingsUser{
    // Set fields directly
}
```

### LinuxKeyPairId (Nullable)

Use the Nullable wrapper methods:
- `obj.LinuxKeyPairId.IsSet()` — check if set
- `obj.LinuxKeyPairId.Get()` — get the inner value (returns pointer)
- `obj.LinuxKeyPairId.Set(&val)` — set the value
- `obj.LinuxKeyPairId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


