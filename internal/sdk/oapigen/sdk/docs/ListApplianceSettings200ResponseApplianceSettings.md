# ListApplianceSettings200ResponseApplianceSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Uuid** | Pointer to **string** |  | [optional] 
**ApplianceId** | Pointer to **string** |  | [optional] 
**ApplianceUrl** | Pointer to **string** |  | [optional] 
**InternalApplianceUrl** | Pointer to **string** |  | [optional] 
**CorsAllowed** | Pointer to **string** |  | [optional] 
**RegistrationEnabled** | Pointer to **bool** |  | [optional] 
**DefaultRoleId** | Pointer to **string** |  | [optional] 
**DefaultUserRoleId** | Pointer to **NullableString** |  | [optional] 
**DockerPrivilegedMode** | Pointer to **bool** |  | [optional] 
**ExpirePwdDays** | Pointer to **NullableString** |  | [optional] 
**DisableAfterAttempts** | Pointer to **NullableString** |  | [optional] 
**DisableAfterDaysInactive** | Pointer to **NullableString** |  | [optional] 
**WarnUserDaysBefore** | Pointer to **NullableString** |  | [optional] 
**SmtpMailFrom** | Pointer to **NullableString** |  | [optional] 
**SmtpServer** | Pointer to **NullableString** |  | [optional] 
**SmtpPort** | Pointer to **NullableString** |  | [optional] 
**SmtpSSL** | Pointer to **bool** |  | [optional] 
**SmtpTLS** | Pointer to **bool** |  | [optional] 
**SmtpUser** | Pointer to **NullableString** |  | [optional] 
**SmtpPassword** | Pointer to **NullableString** |  | [optional] 
**SmtpPasswordHash** | Pointer to **NullableString** |  | [optional] 
**ProxyHost** | Pointer to **NullableString** |  | [optional] 
**ProxyPort** | Pointer to **NullableString** |  | [optional] 
**ProxyUser** | Pointer to **NullableString** |  | [optional] 
**ProxyPassword** | Pointer to **NullableString** |  | [optional] 
**ProxyPasswordHash** | Pointer to **NullableString** |  | [optional] 
**ProxyDomain** | Pointer to **NullableString** |  | [optional] 
**ProxyWorkstation** | Pointer to **NullableString** |  | [optional] 
**CurrencyProvider** | Pointer to **NullableString** |  | [optional] 
**CurrencyKey** | Pointer to **NullableString** |  | [optional] 
**EnabledZoneTypes** | Pointer to [**[]ListApplianceSettings200ResponseApplianceSettingsEnabledZoneTypesInner**](ListApplianceSettings200ResponseApplianceSettingsEnabledZoneTypesInner.md) |  | [optional] 
**StatsRetainmentPeriod** | Pointer to **int64** |  | [optional] 
**DashboardsToDisplay** | Pointer to **NullableString** |  | [optional] 
**VmeAppliance** | Pointer to **bool** |  | [optional] 
**LatestWindowsAgentVersion** | Pointer to **string** |  | [optional] 
**LatestLinuxAgentVersion** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListApplianceSettings200ResponseApplianceSettings{
    // Set fields directly
}
```

### DefaultUserRoleId (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultUserRoleId.IsSet()` — check if set
- `obj.DefaultUserRoleId.Get()` — get the inner value (returns pointer)
- `obj.DefaultUserRoleId.Set(&val)` — set the value
- `obj.DefaultUserRoleId.Unset()` — clear the value
### ExpirePwdDays (Nullable)

Use the Nullable wrapper methods:
- `obj.ExpirePwdDays.IsSet()` — check if set
- `obj.ExpirePwdDays.Get()` — get the inner value (returns pointer)
- `obj.ExpirePwdDays.Set(&val)` — set the value
- `obj.ExpirePwdDays.Unset()` — clear the value
### DisableAfterAttempts (Nullable)

Use the Nullable wrapper methods:
- `obj.DisableAfterAttempts.IsSet()` — check if set
- `obj.DisableAfterAttempts.Get()` — get the inner value (returns pointer)
- `obj.DisableAfterAttempts.Set(&val)` — set the value
- `obj.DisableAfterAttempts.Unset()` — clear the value
### DisableAfterDaysInactive (Nullable)

Use the Nullable wrapper methods:
- `obj.DisableAfterDaysInactive.IsSet()` — check if set
- `obj.DisableAfterDaysInactive.Get()` — get the inner value (returns pointer)
- `obj.DisableAfterDaysInactive.Set(&val)` — set the value
- `obj.DisableAfterDaysInactive.Unset()` — clear the value
### WarnUserDaysBefore (Nullable)

Use the Nullable wrapper methods:
- `obj.WarnUserDaysBefore.IsSet()` — check if set
- `obj.WarnUserDaysBefore.Get()` — get the inner value (returns pointer)
- `obj.WarnUserDaysBefore.Set(&val)` — set the value
- `obj.WarnUserDaysBefore.Unset()` — clear the value
### SmtpMailFrom (Nullable)

Use the Nullable wrapper methods:
- `obj.SmtpMailFrom.IsSet()` — check if set
- `obj.SmtpMailFrom.Get()` — get the inner value (returns pointer)
- `obj.SmtpMailFrom.Set(&val)` — set the value
- `obj.SmtpMailFrom.Unset()` — clear the value
### SmtpServer (Nullable)

Use the Nullable wrapper methods:
- `obj.SmtpServer.IsSet()` — check if set
- `obj.SmtpServer.Get()` — get the inner value (returns pointer)
- `obj.SmtpServer.Set(&val)` — set the value
- `obj.SmtpServer.Unset()` — clear the value
### SmtpPort (Nullable)

Use the Nullable wrapper methods:
- `obj.SmtpPort.IsSet()` — check if set
- `obj.SmtpPort.Get()` — get the inner value (returns pointer)
- `obj.SmtpPort.Set(&val)` — set the value
- `obj.SmtpPort.Unset()` — clear the value
### SmtpUser (Nullable)

Use the Nullable wrapper methods:
- `obj.SmtpUser.IsSet()` — check if set
- `obj.SmtpUser.Get()` — get the inner value (returns pointer)
- `obj.SmtpUser.Set(&val)` — set the value
- `obj.SmtpUser.Unset()` — clear the value
### SmtpPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.SmtpPassword.IsSet()` — check if set
- `obj.SmtpPassword.Get()` — get the inner value (returns pointer)
- `obj.SmtpPassword.Set(&val)` — set the value
- `obj.SmtpPassword.Unset()` — clear the value
### SmtpPasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.SmtpPasswordHash.IsSet()` — check if set
- `obj.SmtpPasswordHash.Get()` — get the inner value (returns pointer)
- `obj.SmtpPasswordHash.Set(&val)` — set the value
- `obj.SmtpPasswordHash.Unset()` — clear the value
### ProxyHost (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyHost.IsSet()` — check if set
- `obj.ProxyHost.Get()` — get the inner value (returns pointer)
- `obj.ProxyHost.Set(&val)` — set the value
- `obj.ProxyHost.Unset()` — clear the value
### ProxyPort (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyPort.IsSet()` — check if set
- `obj.ProxyPort.Get()` — get the inner value (returns pointer)
- `obj.ProxyPort.Set(&val)` — set the value
- `obj.ProxyPort.Unset()` — clear the value
### ProxyUser (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyUser.IsSet()` — check if set
- `obj.ProxyUser.Get()` — get the inner value (returns pointer)
- `obj.ProxyUser.Set(&val)` — set the value
- `obj.ProxyUser.Unset()` — clear the value
### ProxyPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyPassword.IsSet()` — check if set
- `obj.ProxyPassword.Get()` — get the inner value (returns pointer)
- `obj.ProxyPassword.Set(&val)` — set the value
- `obj.ProxyPassword.Unset()` — clear the value
### ProxyPasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyPasswordHash.IsSet()` — check if set
- `obj.ProxyPasswordHash.Get()` — get the inner value (returns pointer)
- `obj.ProxyPasswordHash.Set(&val)` — set the value
- `obj.ProxyPasswordHash.Unset()` — clear the value
### ProxyDomain (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyDomain.IsSet()` — check if set
- `obj.ProxyDomain.Get()` — get the inner value (returns pointer)
- `obj.ProxyDomain.Set(&val)` — set the value
- `obj.ProxyDomain.Unset()` — clear the value
### ProxyWorkstation (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyWorkstation.IsSet()` — check if set
- `obj.ProxyWorkstation.Get()` — get the inner value (returns pointer)
- `obj.ProxyWorkstation.Set(&val)` — set the value
- `obj.ProxyWorkstation.Unset()` — clear the value
### CurrencyProvider (Nullable)

Use the Nullable wrapper methods:
- `obj.CurrencyProvider.IsSet()` — check if set
- `obj.CurrencyProvider.Get()` — get the inner value (returns pointer)
- `obj.CurrencyProvider.Set(&val)` — set the value
- `obj.CurrencyProvider.Unset()` — clear the value
### CurrencyKey (Nullable)

Use the Nullable wrapper methods:
- `obj.CurrencyKey.IsSet()` — check if set
- `obj.CurrencyKey.Get()` — get the inner value (returns pointer)
- `obj.CurrencyKey.Set(&val)` — set the value
- `obj.CurrencyKey.Unset()` — clear the value
### EnabledZoneTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.EnabledZoneTypes.IsSet()` — check if set
- `obj.EnabledZoneTypes.Get()` — get the inner value (returns pointer)
- `obj.EnabledZoneTypes.Set(&val)` — set the value
- `obj.EnabledZoneTypes.Unset()` — clear the value
### DashboardsToDisplay (Nullable)

Use the Nullable wrapper methods:
- `obj.DashboardsToDisplay.IsSet()` — check if set
- `obj.DashboardsToDisplay.Get()` — get the inner value (returns pointer)
- `obj.DashboardsToDisplay.Set(&val)` — set the value
- `obj.DashboardsToDisplay.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


