# GetNetworkServer200ResponseNetworkServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Network Server ID | [optional] 
**Name** | Pointer to **string** | Name | [optional] 
**Description** | Pointer to **NullableString** | Description | [optional] 
**Type** | Pointer to [**GetNetworkServer200ResponseNetworkServerType**](GetNetworkServer200ResponseNetworkServerType.md) |  | [optional] 
**Integration** | Pointer to [**GetNetworkServer200ResponseNetworkServerIntegration**](GetNetworkServer200ResponseNetworkServerIntegration.md) |  | [optional] 
**Account** | Pointer to [**GetNetworkServer200ResponseNetworkServerAccount**](GetNetworkServer200ResponseNetworkServerAccount.md) |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **NullableString** | Internal ID | [optional] 
**ExternalId** | Pointer to **NullableString** | External ID | [optional] 
**ServiceUrl** | Pointer to **NullableString** | Service URL | [optional] 
**ServiceHost** | Pointer to **NullableString** | Service Host | [optional] 
**ServicePort** | Pointer to **NullableInt32** | Service Port | [optional] 
**ServiceMode** | Pointer to **NullableString** | Service Mode | [optional] 
**ServicePath** | Pointer to **NullableString** | Service Path | [optional] 
**ServiceUsername** | Pointer to **NullableString** | Service Username | [optional] 
**ServicePassword** | Pointer to **NullableString** | Service Password | [optional] 
**ServicePasswordHash** | Pointer to **NullableString** |  | [optional] 
**ServiceToken** | Pointer to **NullableString** | Service Token | [optional] 
**ServiceTokenHash** | Pointer to **NullableString** |  | [optional] 
**ApiPort** | Pointer to **NullableInt32** |  | [optional] 
**AdminPort** | Pointer to **NullableInt32** |  | [optional] 
**Status** | Pointer to **string** | Status | [optional] 
**StatusMessage** | Pointer to **NullableString** | Status Message | [optional] 
**StatusDate** | Pointer to **NullableTime** |  | [optional] 
**LastSync** | Pointer to **NullableTime** | Last Sync Date | [optional] 
**NextRunDate** | Pointer to **NullableTime** | Next Run Date | [optional] 
**LastSyncDuration** | Pointer to **NullableInt64** | Last Sync Duration in milliseconds | [optional] 
**Config** | Pointer to **map[string]interface{}** | Config object varies with network server type. | [optional] 
**NetworkFilter** | Pointer to **NullableString** | Network Filter | [optional] 
**TenantMatch** | Pointer to **NullableString** | Tenant Match | [optional] 
**ZoneId** | Pointer to **NullableInt64** | Cloud ID | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Visible** | Pointer to **bool** |  | [optional] 
**Credential** | Pointer to [**GetNetworkServer200ResponseNetworkServerCredential**](GetNetworkServer200ResponseNetworkServerCredential.md) |  | [optional] 
**Tenants** | Pointer to [**[]GetNetworkServer200ResponseNetworkServerTenantsInner**](GetNetworkServer200ResponseNetworkServerTenantsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkServer200ResponseNetworkServer{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### ServiceUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceUrl.IsSet()` — check if set
- `obj.ServiceUrl.Get()` — get the inner value (returns pointer)
- `obj.ServiceUrl.Set(&val)` — set the value
- `obj.ServiceUrl.Unset()` — clear the value
### ServiceHost (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceHost.IsSet()` — check if set
- `obj.ServiceHost.Get()` — get the inner value (returns pointer)
- `obj.ServiceHost.Set(&val)` — set the value
- `obj.ServiceHost.Unset()` — clear the value
### ServicePort (Nullable)

Use the Nullable wrapper methods:
- `obj.ServicePort.IsSet()` — check if set
- `obj.ServicePort.Get()` — get the inner value (returns pointer)
- `obj.ServicePort.Set(&val)` — set the value
- `obj.ServicePort.Unset()` — clear the value
### ServiceMode (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceMode.IsSet()` — check if set
- `obj.ServiceMode.Get()` — get the inner value (returns pointer)
- `obj.ServiceMode.Set(&val)` — set the value
- `obj.ServiceMode.Unset()` — clear the value
### ServicePath (Nullable)

Use the Nullable wrapper methods:
- `obj.ServicePath.IsSet()` — check if set
- `obj.ServicePath.Get()` — get the inner value (returns pointer)
- `obj.ServicePath.Set(&val)` — set the value
- `obj.ServicePath.Unset()` — clear the value
### ServiceUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceUsername.IsSet()` — check if set
- `obj.ServiceUsername.Get()` — get the inner value (returns pointer)
- `obj.ServiceUsername.Set(&val)` — set the value
- `obj.ServiceUsername.Unset()` — clear the value
### ServicePassword (Nullable)

Use the Nullable wrapper methods:
- `obj.ServicePassword.IsSet()` — check if set
- `obj.ServicePassword.Get()` — get the inner value (returns pointer)
- `obj.ServicePassword.Set(&val)` — set the value
- `obj.ServicePassword.Unset()` — clear the value
### ServicePasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.ServicePasswordHash.IsSet()` — check if set
- `obj.ServicePasswordHash.Get()` — get the inner value (returns pointer)
- `obj.ServicePasswordHash.Set(&val)` — set the value
- `obj.ServicePasswordHash.Unset()` — clear the value
### ServiceToken (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceToken.IsSet()` — check if set
- `obj.ServiceToken.Get()` — get the inner value (returns pointer)
- `obj.ServiceToken.Set(&val)` — set the value
- `obj.ServiceToken.Unset()` — clear the value
### ServiceTokenHash (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceTokenHash.IsSet()` — check if set
- `obj.ServiceTokenHash.Get()` — get the inner value (returns pointer)
- `obj.ServiceTokenHash.Set(&val)` — set the value
- `obj.ServiceTokenHash.Unset()` — clear the value
### ApiPort (Nullable)

Use the Nullable wrapper methods:
- `obj.ApiPort.IsSet()` — check if set
- `obj.ApiPort.Get()` — get the inner value (returns pointer)
- `obj.ApiPort.Set(&val)` — set the value
- `obj.ApiPort.Unset()` — clear the value
### AdminPort (Nullable)

Use the Nullable wrapper methods:
- `obj.AdminPort.IsSet()` — check if set
- `obj.AdminPort.Get()` — get the inner value (returns pointer)
- `obj.AdminPort.Set(&val)` — set the value
- `obj.AdminPort.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### StatusDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusDate.IsSet()` — check if set
- `obj.StatusDate.Get()` — get the inner value (returns pointer)
- `obj.StatusDate.Set(&val)` — set the value
- `obj.StatusDate.Unset()` — clear the value
### LastSync (Nullable)

Use the Nullable wrapper methods:
- `obj.LastSync.IsSet()` — check if set
- `obj.LastSync.Get()` — get the inner value (returns pointer)
- `obj.LastSync.Set(&val)` — set the value
- `obj.LastSync.Unset()` — clear the value
### NextRunDate (Nullable)

Use the Nullable wrapper methods:
- `obj.NextRunDate.IsSet()` — check if set
- `obj.NextRunDate.Get()` — get the inner value (returns pointer)
- `obj.NextRunDate.Set(&val)` — set the value
- `obj.NextRunDate.Unset()` — clear the value
### LastSyncDuration (Nullable)

Use the Nullable wrapper methods:
- `obj.LastSyncDuration.IsSet()` — check if set
- `obj.LastSyncDuration.Get()` — get the inner value (returns pointer)
- `obj.LastSyncDuration.Set(&val)` — set the value
- `obj.LastSyncDuration.Unset()` — clear the value
### NetworkFilter (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkFilter.IsSet()` — check if set
- `obj.NetworkFilter.Get()` — get the inner value (returns pointer)
- `obj.NetworkFilter.Set(&val)` — set the value
- `obj.NetworkFilter.Unset()` — clear the value
### TenantMatch (Nullable)

Use the Nullable wrapper methods:
- `obj.TenantMatch.IsSet()` — check if set
- `obj.TenantMatch.Get()` — get the inner value (returns pointer)
- `obj.TenantMatch.Set(&val)` — set the value
- `obj.TenantMatch.Unset()` — clear the value
### ZoneId (Nullable)

Use the Nullable wrapper methods:
- `obj.ZoneId.IsSet()` — check if set
- `obj.ZoneId.Get()` — get the inner value (returns pointer)
- `obj.ZoneId.Set(&val)` — set the value
- `obj.ZoneId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


