# ListUserSettings200ResponseAllOfUser

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
**DefaultGroup** | Pointer to [**NullableListUserSettings200ResponseAllOfUserDefaultGroup**](ListUserSettings200ResponseAllOfUserDefaultGroup.md) |  | [optional] 
**DefaultCloud** | Pointer to [**NullableListUserSettings200ResponseAllOfUserDefaultCloud**](ListUserSettings200ResponseAllOfUserDefaultCloud.md) |  | [optional] 
**DefaultPersona** | Pointer to [**ListUserSettings200ResponseAllOfUserDefaultPersona**](ListUserSettings200ResponseAllOfUserDefaultPersona.md) |  | [optional] 
**IsUsing2FA** | Pointer to **bool** |  | [optional] 
**Tenant** | Pointer to [**NullableListUserSettings200ResponseAllOfUserTenant**](ListUserSettings200ResponseAllOfUserTenant.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListUserSettings200ResponseAllOfUser{
    // Set fields directly
}
```

### LinuxKeyPairId (Nullable)

Use the Nullable wrapper methods:
- `obj.LinuxKeyPairId.IsSet()` — check if set
- `obj.LinuxKeyPairId.Get()` — get the inner value (returns pointer)
- `obj.LinuxKeyPairId.Set(&val)` — set the value
- `obj.LinuxKeyPairId.Unset()` — clear the value
### DefaultGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultGroup.IsSet()` — check if set
- `obj.DefaultGroup.Get()` — get the inner value (returns pointer)
- `obj.DefaultGroup.Set(&val)` — set the value
- `obj.DefaultGroup.Unset()` — clear the value
### DefaultCloud (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultCloud.IsSet()` — check if set
- `obj.DefaultCloud.Get()` — get the inner value (returns pointer)
- `obj.DefaultCloud.Set(&val)` — set the value
- `obj.DefaultCloud.Unset()` — clear the value
### Tenant (Nullable)

Use the Nullable wrapper methods:
- `obj.Tenant.IsSet()` — check if set
- `obj.Tenant.Get()` — get the inner value (returns pointer)
- `obj.Tenant.Set(&val)` — set the value
- `obj.Tenant.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


