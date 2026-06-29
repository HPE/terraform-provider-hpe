# AddUserTenant200ResponseAllOfUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Username** | Pointer to **string** |  | [optional] 
**DisplayName** | Pointer to **string** |  | [optional] 
**Email** | Pointer to **string** |  | [optional] 
**FirstName** | Pointer to **string** |  | [optional] 
**LastName** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**ReceiveNotifications** | Pointer to **bool** |  | [optional] 
**IsUsing2FA** | Pointer to **bool** |  | [optional] 
**AccountExpired** | Pointer to **bool** |  | [optional] 
**AccountLocked** | Pointer to **bool** |  | [optional] 
**PasswordExpired** | Pointer to **bool** |  | [optional] 
**LoginCount** | Pointer to **int64** |  | [optional] 
**LoginAttempts** | Pointer to **int64** |  | [optional] 
**LastLoginDate** | Pointer to **time.Time** |  | [optional] 
**Roles** | Pointer to [**[]AddUserTenant200ResponseAllOfUserRolesInner**](AddUserTenant200ResponseAllOfUserRolesInner.md) |  | [optional] 
**Account** | Pointer to [**NullableAddUserTenant200ResponseAllOfUserAccount**](AddUserTenant200ResponseAllOfUserAccount.md) |  | [optional] 
**LinuxUsername** | Pointer to **NullableString** |  | [optional] 
**LinuxPassword** | Pointer to **NullableString** |  | [optional] 
**LinuxKeyPairId** | Pointer to **NullableInt64** |  | [optional] 
**WindowsUsername** | Pointer to **NullableString** |  | [optional] 
**WindowsPassword** | Pointer to **NullableString** |  | [optional] 
**DefaultPersona** | Pointer to [**AddUserTenant200ResponseAllOfUserDefaultPersona**](AddUserTenant200ResponseAllOfUserDefaultPersona.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Access** | Pointer to [**AddUserTenant200ResponseAllOfUserAccess**](AddUserTenant200ResponseAllOfUserAccess.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddUserTenant200ResponseAllOfUser{
    // Set fields directly
}
```

### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
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
### LinuxKeyPairId (Nullable)

Use the Nullable wrapper methods:
- `obj.LinuxKeyPairId.IsSet()` — check if set
- `obj.LinuxKeyPairId.Get()` — get the inner value (returns pointer)
- `obj.LinuxKeyPairId.Set(&val)` — set the value
- `obj.LinuxKeyPairId.Unset()` — clear the value
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


