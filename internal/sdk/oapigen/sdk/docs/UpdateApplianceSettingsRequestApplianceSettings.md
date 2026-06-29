# UpdateApplianceSettingsRequestApplianceSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApplianceUrl** | Pointer to **string** | Appliance URL | [optional] 
**InternalApplianceUrl** | Pointer to **NullableString** | Internal Appliance URL (PXE) | [optional] 
**CorsAllowed** | Pointer to **NullableString** | API Allowed Origins | [optional] 
**RegistrationEnabled** | Pointer to **bool** | Registration enabled (true, false) | [optional] 
**DefaultRoleId** | Pointer to **int64** | Default tenant role ID | [optional] 
**DefaultUserRoleId** | Pointer to **int64** | Default user role ID | [optional] 
**DockerPrivilegedMode** | Pointer to **bool** | Docker privileged mode (true, false) | [optional] 
**PasswordMinLength** | Pointer to **string** | Min Password Length | [optional] 
**PasswordMinUpperCase** | Pointer to **string** | Min Password Uppercase | [optional] 
**PasswordMinNumbers** | Pointer to **string** | Min Password Numbers | [optional] 
**PasswordMinSymbols** | Pointer to **string** | Min Password Symbols | [optional] 
**UserBrowserSessionTimeout** | Pointer to **string** | User Browser Session Timeout (Minutes) | [optional] 
**UserBrowserSessionWarning** | Pointer to **string** | User Browser Session Warning (Minutes) | [optional] 
**ExpirePwdDays** | Pointer to **int64** | Expire password after days. Setting to 0 disabled this feature | [optional] 
**DisableAfterAttempts** | Pointer to **int64** | Disable user after number of attempts. Set to 0 to disable this feature | [optional] 
**DisableAfterDaysInactive** | Pointer to **int64** | Disable user if inactive for specified days. Set to 0 to disable this feature | [optional] 
**WarnUserDaysBefore** | Pointer to **int64** | Send warning email number of days in advance before deactivating. Set to 0 to disable this feature | [optional] 
**SmtpMailFrom** | Pointer to **string** | From email address | [optional] 
**SmtpServer** | Pointer to **string** | SMTP server / host | [optional] 
**SmtpPort** | Pointer to **int64** | SMTP port | [optional] 
**SmtpSSL** | Pointer to **bool** | Use SSL for SMTP connection | [optional] 
**SmtpTLS** | Pointer to **bool** | Use TLS for SMTP connections | [optional] 
**SmtpUser** | Pointer to **string** | SMTP username | [optional] 
**SmtpPassword** | Pointer to **string** | SMTP password | [optional] 
**ProxyHost** | Pointer to **NullableString** | Proxy host | [optional] 
**ProxyPort** | Pointer to **NullableString** | Proxy port | [optional] 
**ProxyUser** | Pointer to **string** | Proxy username | [optional] 
**ProxyPassword** | Pointer to **string** | Proxy password | [optional] 
**ProxyDomain** | Pointer to **NullableString** | Proxy domain | [optional] 
**ProxyWorkstation** | Pointer to **NullableString** | Proxy workstation | [optional] 
**CurrencyProvider** | Pointer to **string** | Currency provider | [optional] 
**CurrencyKey** | Pointer to **NullableString** | Currency provider API key | [optional] 
**EnableAllZoneTypes** | Pointer to **bool** | Set all cloud types enabled status on, overrides enableZoneTypes and disableZoneTypes parameters | [optional] 
**EnableZoneTypes** | Pointer to **[]int64** | List of cloud type IDs to set enabled status on | [optional] 
**DisableZoneTypes** | Pointer to **[]int64** | List of cloud type IDs to set enabled status off | [optional] 
**DisableAllZoneTypes** | Pointer to **bool** | Set all cloud types enabled status off, can be used in conjunction with enableZoneTypes | [optional] 
**TwilioAccountSid** | Pointer to **string** | Twilio SMS Account SID | [optional] 
**TwilioSmsFrom** | Pointer to **string** | Twilio SMS From | [optional] 
**TwilioAuthToken** | Pointer to **string** | Twilio SMS Auth Token | [optional] 
**CloudSyncIntervalSeconds** | Pointer to **int64** | Cloud Sync Interval (Seconds) | [optional] 
**ClusterSyncIntervalSeconds** | Pointer to **int64** | Cluster Sync Interval (Seconds) | [optional] 
**UsageRetainmentPeriod** | Pointer to **int64** | Usage Retainment (Days) | [optional] 
**InvoiceRetainmentPeriod** | Pointer to **int64** | Invoice Retainment (Days) | [optional] 
**IncidentRetainmentPeriod** | Pointer to **int64** | Incident Retainment (Days) | [optional] 
**StatsRetainmentPeriod** | Pointer to **int64** | The number of days stats will be retained. (30, 60 or 90) | [optional] 
**ReportsRetainmentPeriod** | Pointer to **int64** | The number of days reports will be retained. | [optional] 
**HttpBlacklistHosts** | Pointer to **string** | Provide a comma delimited list of ips/hostnames to be blocked when using HTTP Task Types or REST Datasource Option Lists | [optional] 
**HttpApprovelistHosts** | Pointer to **string** | Provide a comma delimited list of ips/hostnames to be allowed when using HTTP Task Types or REST Datasource Option Lists. If not specified, only deny list is filtered out. | [optional] 
**NoAgent** | Pointer to **bool** | If true, disables Agent installation globally. | [optional] 
**AgentSSLVerify** | Pointer to **bool** | Enable/Disable SSL Verification of Agent | [optional] 
**DisableSSHPasswordAuth** | Pointer to **bool** | Enable/Disable SSH Password Authentication for the Appliance | [optional] 
**DefaultLocale** | Pointer to **string** | Default appliance Locale. Setting a default locale for the application will override user browser preferences. | [optional] 
**DefaultVdiGateway** | Pointer to **int64** | ID of the VDI gateway. | [optional] 
**MaxOptionListSize** | Pointer to **int64** | Max option list size. Units are x10^3 (thousand). Increasing this value can adversely affect Morpheus performance. Increase with caution. | [optional] 
**ExchangeUrl** | Pointer to **string** | The url used for checking if there is an update for plugins. Default https\\://share.morpheusdata.com | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateApplianceSettingsRequestApplianceSettings{
    // Set fields directly
}
```

### InternalApplianceUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalApplianceUrl.IsSet()` — check if set
- `obj.InternalApplianceUrl.Get()` — get the inner value (returns pointer)
- `obj.InternalApplianceUrl.Set(&val)` — set the value
- `obj.InternalApplianceUrl.Unset()` — clear the value
### CorsAllowed (Nullable)

Use the Nullable wrapper methods:
- `obj.CorsAllowed.IsSet()` — check if set
- `obj.CorsAllowed.Get()` — get the inner value (returns pointer)
- `obj.CorsAllowed.Set(&val)` — set the value
- `obj.CorsAllowed.Unset()` — clear the value
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
### CurrencyKey (Nullable)

Use the Nullable wrapper methods:
- `obj.CurrencyKey.IsSet()` — check if set
- `obj.CurrencyKey.Get()` — get the inner value (returns pointer)
- `obj.CurrencyKey.Set(&val)` — set the value
- `obj.CurrencyKey.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


