# AddCluster200ResponseAllOfCluster

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Location** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**ServiceUrl** | Pointer to **NullableString** |  | [optional] 
**ServiceHost** | Pointer to **NullableString** |  | [optional] 
**ServicePath** | Pointer to **NullableString** |  | [optional] 
**ServiceHostname** | Pointer to **NullableString** |  | [optional] 
**ServicePort** | Pointer to **int64** |  | [optional] 
**ServiceUsername** | Pointer to **NullableString** |  | [optional] 
**ServicePassword** | Pointer to **NullableString** |  | [optional] 
**ServicePasswordHash** | Pointer to **NullableString** |  | [optional] 
**ServiceToken** | Pointer to **NullableString** |  | [optional] 
**ServiceTokenHash** | Pointer to **NullableString** |  | [optional] 
**ServiceAccess** | Pointer to **NullableString** |  | [optional] 
**ServiceAccessHash** | Pointer to **NullableString** |  | [optional] 
**ServiceCert** | Pointer to **NullableString** |  | [optional] 
**ServiceCertHash** | Pointer to **NullableString** |  | [optional] 
**ServiceVersion** | Pointer to **NullableString** |  | [optional] 
**SearchDomains** | Pointer to **NullableString** |  | [optional] 
**EnableInternalDns** | Pointer to **bool** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**DatacenterId** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusDate** | Pointer to **NullableTime** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**InventoryLevel** | Pointer to **string** |  | [optional] 
**LastSync** | Pointer to **NullableTime** |  | [optional] 
**NextRunDate** | Pointer to **NullableTime** |  | [optional] 
**LastSyncDuration** | Pointer to **NullableInt64** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Managed** | Pointer to **bool** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**AutoRecoverPowerState** | Pointer to **bool** | Automatically Power on VMs | [optional] [default to false]
**CpuPlacementMode** | Pointer to **NullableString** | Cluster CPU placement mode | [optional] 
**UseAgent** | Pointer to **NullableString** | Use the Agent to relay communications for the Kubernetes API instead of direct. | [optional] 
**ProvisionComplete** | Pointer to **bool** | Changes from false to true once provisioning is finished. | [optional] 
**ServiceEntry** | Pointer to **NullableString** |  | [optional] 
**CreatedBy** | Pointer to [**AddCluster200ResponseAllOfClusterCreatedBy**](AddCluster200ResponseAllOfClusterCreatedBy.md) |  | [optional] 
**UserGroup** | Pointer to **NullableString** |  | [optional] 
**Layout** | Pointer to [**AddCluster200ResponseAllOfClusterLayout**](AddCluster200ResponseAllOfClusterLayout.md) |  | [optional] 
**Owner** | Pointer to [**NullableAddCluster200ResponseAllOfClusterOwner**](AddCluster200ResponseAllOfClusterOwner.md) |  | [optional] 
**Servers** | Pointer to [**[]AddCluster200ResponseAllOfClusterServersInner**](AddCluster200ResponseAllOfClusterServersInner.md) |  | [optional] 
**Accounts** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Integrations** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Site** | Pointer to [**NullableAddCluster200ResponseAllOfClusterSite**](AddCluster200ResponseAllOfClusterSite.md) |  | [optional] 
**Type** | Pointer to [**NullableAddCluster200ResponseAllOfClusterType**](AddCluster200ResponseAllOfClusterType.md) |  | [optional] 
**Zone** | Pointer to [**AddCluster200ResponseAllOfClusterZone**](AddCluster200ResponseAllOfClusterZone.md) |  | [optional] 
**WorkerStats** | Pointer to [**AddCluster200ResponseAllOfClusterWorkerStats**](AddCluster200ResponseAllOfClusterWorkerStats.md) |  | [optional] 
**ContainersCount** | Pointer to **int64** |  | [optional] 
**DeploymentsCount** | Pointer to **int64** |  | [optional] 
**PodsCount** | Pointer to **int64** |  | [optional] 
**JobsCount** | Pointer to **int64** |  | [optional] 
**VolumesCount** | Pointer to **int64** |  | [optional] 
**NamespacesCount** | Pointer to **int64** |  | [optional] 
**WorkersCount** | Pointer to **int64** |  | [optional] 
**ServicesCount** | Pointer to **int64** |  | [optional] 
**Permissions** | Pointer to [**AddCluster200ResponseAllOfClusterPermissions**](AddCluster200ResponseAllOfClusterPermissions.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCluster200ResponseAllOfCluster{
    // Set fields directly
}
```

### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Location (Nullable)

Use the Nullable wrapper methods:
- `obj.Location.IsSet()` — check if set
- `obj.Location.Get()` — get the inner value (returns pointer)
- `obj.Location.Set(&val)` — set the value
- `obj.Location.Unset()` — clear the value
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
### ServiceHostname (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceHostname.IsSet()` — check if set
- `obj.ServiceHostname.Get()` — get the inner value (returns pointer)
- `obj.ServiceHostname.Set(&val)` — set the value
- `obj.ServiceHostname.Unset()` — clear the value
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
### ServiceAccess (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceAccess.IsSet()` — check if set
- `obj.ServiceAccess.Get()` — get the inner value (returns pointer)
- `obj.ServiceAccess.Set(&val)` — set the value
- `obj.ServiceAccess.Unset()` — clear the value
### ServiceAccessHash (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceAccessHash.IsSet()` — check if set
- `obj.ServiceAccessHash.Get()` — get the inner value (returns pointer)
- `obj.ServiceAccessHash.Set(&val)` — set the value
- `obj.ServiceAccessHash.Unset()` — clear the value
### ServiceCert (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceCert.IsSet()` — check if set
- `obj.ServiceCert.Get()` — get the inner value (returns pointer)
- `obj.ServiceCert.Set(&val)` — set the value
- `obj.ServiceCert.Unset()` — clear the value
### ServiceCertHash (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceCertHash.IsSet()` — check if set
- `obj.ServiceCertHash.Get()` — get the inner value (returns pointer)
- `obj.ServiceCertHash.Set(&val)` — set the value
- `obj.ServiceCertHash.Unset()` — clear the value
### ServiceVersion (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceVersion.IsSet()` — check if set
- `obj.ServiceVersion.Get()` — get the inner value (returns pointer)
- `obj.ServiceVersion.Set(&val)` — set the value
- `obj.ServiceVersion.Unset()` — clear the value
### SearchDomains (Nullable)

Use the Nullable wrapper methods:
- `obj.SearchDomains.IsSet()` — check if set
- `obj.SearchDomains.Get()` — get the inner value (returns pointer)
- `obj.SearchDomains.Set(&val)` — set the value
- `obj.SearchDomains.Unset()` — clear the value
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
### DatacenterId (Nullable)

Use the Nullable wrapper methods:
- `obj.DatacenterId.IsSet()` — check if set
- `obj.DatacenterId.Get()` — get the inner value (returns pointer)
- `obj.DatacenterId.Set(&val)` — set the value
- `obj.DatacenterId.Unset()` — clear the value
### StatusDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusDate.IsSet()` — check if set
- `obj.StatusDate.Get()` — get the inner value (returns pointer)
- `obj.StatusDate.Set(&val)` — set the value
- `obj.StatusDate.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
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
### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### CpuPlacementMode (Nullable)

Use the Nullable wrapper methods:
- `obj.CpuPlacementMode.IsSet()` — check if set
- `obj.CpuPlacementMode.Get()` — get the inner value (returns pointer)
- `obj.CpuPlacementMode.Set(&val)` — set the value
- `obj.CpuPlacementMode.Unset()` — clear the value
### UseAgent (Nullable)

Use the Nullable wrapper methods:
- `obj.UseAgent.IsSet()` — check if set
- `obj.UseAgent.Get()` — get the inner value (returns pointer)
- `obj.UseAgent.Set(&val)` — set the value
- `obj.UseAgent.Unset()` — clear the value
### ServiceEntry (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceEntry.IsSet()` — check if set
- `obj.ServiceEntry.Get()` — get the inner value (returns pointer)
- `obj.ServiceEntry.Set(&val)` — set the value
- `obj.ServiceEntry.Unset()` — clear the value
### UserGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.UserGroup.IsSet()` — check if set
- `obj.UserGroup.Get()` — get the inner value (returns pointer)
- `obj.UserGroup.Set(&val)` — set the value
- `obj.UserGroup.Unset()` — clear the value
### Owner (Nullable)

Use the Nullable wrapper methods:
- `obj.Owner.IsSet()` — check if set
- `obj.Owner.Get()` — get the inner value (returns pointer)
- `obj.Owner.Set(&val)` — set the value
- `obj.Owner.Unset()` — clear the value
### Site (Nullable)

Use the Nullable wrapper methods:
- `obj.Site.IsSet()` — check if set
- `obj.Site.Get()` — get the inner value (returns pointer)
- `obj.Site.Set(&val)` — set the value
- `obj.Site.Unset()` — clear the value
### Type (Nullable)

Use the Nullable wrapper methods:
- `obj.Type.IsSet()` — check if set
- `obj.Type.Get()` — get the inner value (returns pointer)
- `obj.Type.Set(&val)` — set the value
- `obj.Type.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


