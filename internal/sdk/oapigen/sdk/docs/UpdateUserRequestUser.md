# UpdateUserRequestUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FirstName** | Pointer to **NullableString** | First Name | [optional] 
**LastName** | Pointer to **NullableString** | Last Name | [optional] 
**Username** | Pointer to **string** | Username (unique per tenant). | [optional] 
**LinuxUsername** | Pointer to **NullableString** |  | [optional] 
**LinuxPassword** | Pointer to **NullableString** |  | [optional] 
**LinuxKeyPairId** | Pointer to **NullableInt64** |  | [optional] 
**WindowsUsername** | Pointer to **NullableString** |  | [optional] 
**WindowsPassword** | Pointer to **NullableString** |  | [optional] 
**Email** | Pointer to **string** | Email address | [optional] 
**Password** | Pointer to **string** | Password | [optional] 
**Roles** | Pointer to [**[]UpdateUserRequestUserRolesInner**](UpdateUserRequestUserRolesInner.md) | List of Roles | [optional] 
**ReceiveNotifications** | Pointer to **bool** | Receive Notifications? | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateUserRequestUser{
    // Set fields directly
}
```

### FirstName (Nullable)

Use the Nullable wrapper methods:
- `obj.FirstName.IsSet()` — check if set
- `obj.FirstName.Get()` — get the inner value (returns pointer)
- `obj.FirstName.Set(&val)` — set the value
- `obj.FirstName.Unset()` — clear the value
### LastName (Nullable)

Use the Nullable wrapper methods:
- `obj.LastName.IsSet()` — check if set
- `obj.LastName.Get()` — get the inner value (returns pointer)
- `obj.LastName.Set(&val)` — set the value
- `obj.LastName.Unset()` — clear the value
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


