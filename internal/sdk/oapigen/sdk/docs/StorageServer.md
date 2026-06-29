# StorageServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Type** | Pointer to [**AddStorageServers200ResponseAllOfStorageServerType**](AddStorageServers200ResponseAllOfStorageServerType.md) |  | [optional] 
**Chassis** | Pointer to [**AddStorageServers200ResponseAllOfStorageServerChassis**](AddStorageServers200ResponseAllOfStorageServerChassis.md) |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**ServiceUrl** | Pointer to **NullableString** |  | [optional] 
**ServiceHost** | Pointer to **NullableString** |  | [optional] 
**ServicePath** | Pointer to **NullableString** |  | [optional] 
**ServiceToken** | Pointer to **NullableString** |  | [optional] 
**ServiceTokenHash** | Pointer to **NullableString** |  | [optional] 
**ServiceVersion** | Pointer to **NullableString** |  | [optional] 
**ServiceUsername** | Pointer to **NullableString** |  | [optional] 
**ServicePassword** | Pointer to **NullableString** |  | [optional] 
**ServicePasswordHash** | Pointer to **NullableString** |  | [optional] 
**InternalIp** | Pointer to **NullableString** |  | [optional] 
**ExternalIp** | Pointer to **NullableString** |  | [optional] 
**ApiPort** | Pointer to **NullableInt32** |  | [optional] 
**AdminPort** | Pointer to **NullableInt32** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**ServerVendor** | Pointer to **NullableString** |  | [optional] 
**ServerModel** | Pointer to **NullableString** |  | [optional] 
**SerialNumber** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **NullableString** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**StatusDate** | Pointer to **time.Time** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**MaxStorage** | Pointer to **NullableInt64** |  | [optional] 
**UsedStorage** | Pointer to **NullableInt64** |  | [optional] 
**DiskCount** | Pointer to **NullableInt32** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Groups** | Pointer to [**[]AddStorageServers200ResponseAllOfStorageServerGroupsInner**](AddStorageServers200ResponseAllOfStorageServerGroupsInner.md) |  | [optional] 
**HostGroups** | Pointer to [**[]AddStorageServers200ResponseAllOfStorageServerHostGroupsInner**](AddStorageServers200ResponseAllOfStorageServerHostGroupsInner.md) |  | [optional] 
**Hosts** | Pointer to [**[]AddStorageServers200ResponseAllOfStorageServerHostsInner**](AddStorageServers200ResponseAllOfStorageServerHostsInner.md) |  | [optional] 
**Tenants** | Pointer to [**[]AddStorageServers200ResponseAllOfStorageServerTenantsInner**](AddStorageServers200ResponseAllOfStorageServerTenantsInner.md) |  | [optional] 
**Owner** | Pointer to [**AddStorageServers200ResponseAllOfStorageServerOwner**](AddStorageServers200ResponseAllOfStorageServerOwner.md) |  | [optional] 
**Credential** | Pointer to [**AddStorageServers200ResponseAllOfStorageServerCredential**](AddStorageServers200ResponseAllOfStorageServerCredential.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &StorageServer{
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
### ServicePath (Nullable)

Use the Nullable wrapper methods:
- `obj.ServicePath.IsSet()` — check if set
- `obj.ServicePath.Get()` — get the inner value (returns pointer)
- `obj.ServicePath.Set(&val)` — set the value
- `obj.ServicePath.Unset()` — clear the value
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
### ServiceVersion (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceVersion.IsSet()` — check if set
- `obj.ServiceVersion.Get()` — get the inner value (returns pointer)
- `obj.ServiceVersion.Set(&val)` — set the value
- `obj.ServiceVersion.Unset()` — clear the value
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
### InternalIp (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalIp.IsSet()` — check if set
- `obj.InternalIp.Get()` — get the inner value (returns pointer)
- `obj.InternalIp.Set(&val)` — set the value
- `obj.InternalIp.Unset()` — clear the value
### ExternalIp (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalIp.IsSet()` — check if set
- `obj.ExternalIp.Get()` — get the inner value (returns pointer)
- `obj.ExternalIp.Set(&val)` — set the value
- `obj.ExternalIp.Unset()` — clear the value
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
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### ServerVendor (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerVendor.IsSet()` — check if set
- `obj.ServerVendor.Get()` — get the inner value (returns pointer)
- `obj.ServerVendor.Set(&val)` — set the value
- `obj.ServerVendor.Unset()` — clear the value
### ServerModel (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerModel.IsSet()` — check if set
- `obj.ServerModel.Get()` — get the inner value (returns pointer)
- `obj.ServerModel.Set(&val)` — set the value
- `obj.ServerModel.Unset()` — clear the value
### SerialNumber (Nullable)

Use the Nullable wrapper methods:
- `obj.SerialNumber.IsSet()` — check if set
- `obj.SerialNumber.Get()` — get the inner value (returns pointer)
- `obj.SerialNumber.Set(&val)` — set the value
- `obj.SerialNumber.Unset()` — clear the value
### Status (Nullable)

Use the Nullable wrapper methods:
- `obj.Status.IsSet()` — check if set
- `obj.Status.Get()` — get the inner value (returns pointer)
- `obj.Status.Set(&val)` — set the value
- `obj.Status.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### ErrorMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.ErrorMessage.IsSet()` — check if set
- `obj.ErrorMessage.Get()` — get the inner value (returns pointer)
- `obj.ErrorMessage.Set(&val)` — set the value
- `obj.ErrorMessage.Unset()` — clear the value
### MaxStorage (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxStorage.IsSet()` — check if set
- `obj.MaxStorage.Get()` — get the inner value (returns pointer)
- `obj.MaxStorage.Set(&val)` — set the value
- `obj.MaxStorage.Unset()` — clear the value
### UsedStorage (Nullable)

Use the Nullable wrapper methods:
- `obj.UsedStorage.IsSet()` — check if set
- `obj.UsedStorage.Get()` — get the inner value (returns pointer)
- `obj.UsedStorage.Set(&val)` — set the value
- `obj.UsedStorage.Unset()` — clear the value
### DiskCount (Nullable)

Use the Nullable wrapper methods:
- `obj.DiskCount.IsSet()` — check if set
- `obj.DiskCount.Get()` — get the inner value (returns pointer)
- `obj.DiskCount.Set(&val)` — set the value
- `obj.DiskCount.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


