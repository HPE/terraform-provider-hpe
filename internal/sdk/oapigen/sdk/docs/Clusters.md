# Clusters

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
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
**ServiceToken** | Pointer to **string** |  | [optional] 
**ServiceTokenHash** | Pointer to **string** |  | [optional] 
**ServiceAccess** | Pointer to **string** |  | [optional] 
**ServiceAccessHash** | Pointer to **string** |  | [optional] 
**ServiceCert** | Pointer to **NullableString** |  | [optional] 
**ServiceCertHash** | Pointer to **NullableString** |  | [optional] 
**ServiceVersion** | Pointer to **string** |  | [optional] 
**SearchDomains** | Pointer to **NullableString** |  | [optional] 
**EnableInternalDns** | Pointer to **bool** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**DatacenterId** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusDate** | Pointer to **time.Time** |  | [optional] 
**StatusMessage** | Pointer to **string** |  | [optional] 
**InventoryLevel** | Pointer to **string** |  | [optional] 
**LastSync** | Pointer to **time.Time** |  | [optional] 
**NextRunDate** | Pointer to **time.Time** |  | [optional] 
**LastSyncDuration** | Pointer to **int64** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Managed** | Pointer to **bool** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**AutoRecoverPowerState** | Pointer to **bool** |  | [optional] 
**CpuPlacementMode** | Pointer to **NullableString** | Cluster CPU placement mode | [optional] 
**ServiceEntry** | Pointer to **NullableString** |  | [optional] 
**CreatedBy** | Pointer to [**ClustersCreatedBy**](ClustersCreatedBy.md) |  | [optional] 
**UserGroup** | Pointer to **NullableString** |  | [optional] 
**Layout** | Pointer to [**ClustersLayout**](ClustersLayout.md) |  | [optional] 
**Owner** | Pointer to [**ClustersOwner**](ClustersOwner.md) |  | [optional] 
**Servers** | Pointer to [**[]ClustersServersInner**](ClustersServersInner.md) |  | [optional] 
**Accounts** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Integrations** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Site** | Pointer to [**ClustersSite**](ClustersSite.md) |  | [optional] 
**Type** | Pointer to [**ClustersType**](ClustersType.md) |  | [optional] 
**Zone** | Pointer to [**ClustersZone**](ClustersZone.md) |  | [optional] 
**WorkerStats** | Pointer to [**ClustersWorkerStats**](ClustersWorkerStats.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &Clusters{
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
### CpuPlacementMode (Nullable)

Use the Nullable wrapper methods:
- `obj.CpuPlacementMode.IsSet()` — check if set
- `obj.CpuPlacementMode.Get()` — get the inner value (returns pointer)
- `obj.CpuPlacementMode.Set(&val)` — set the value
- `obj.CpuPlacementMode.Unset()` — clear the value
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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


