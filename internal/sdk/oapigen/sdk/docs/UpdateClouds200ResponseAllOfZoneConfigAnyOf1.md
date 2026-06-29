# UpdateClouds200ResponseAllOfZoneConfigAnyOf1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApplianceUrl** | Pointer to **NullableString** |  | [optional] 
**DatacenterName** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**InventoryLevel** | Pointer to **NullableString** |  | [optional] 
**ConsoleKeymap** | Pointer to **NullableString** |  | [optional] 
**ApiUrl** | Pointer to **string** |  | [optional] 
**ApiVersion** | Pointer to **NullableString** |  | [optional] 
**BackupMode** | Pointer to **NullableString** |  | [optional] 
**CertificateProvider** | Pointer to **NullableString** |  | [optional] 
**Cluster** | Pointer to **string** |  | [optional] 
**ClusterRef** | Pointer to **string** |  | [optional] 
**ConfigCmId** | Pointer to **NullableString** |  | [optional] 
**ConfigCmdbDiscovery** | Pointer to **bool** |  | [optional] 
**ConfigCmdbId** | Pointer to **NullableString** |  | [optional] 
**ConfigManagementId** | Pointer to **NullableString** |  | [optional] 
**Datacenter** | Pointer to **string** |  | [optional] 
**DatacenterId** | Pointer to **NullableString** |  | [optional] 
**DiskStorageType** | Pointer to **NullableString** |  | [optional] 
**DistributedWorkerId** | Pointer to **NullableString** |  | [optional] 
**DnsIntegrationId** | Pointer to **NullableString** |  | [optional] 
**EnableDiskTypeSelection** | Pointer to **NullableString** |  | [optional] 
**EnableNetworkTypeSelection** | Pointer to **NullableString** |  | [optional] 
**EnableStorageTypeSelection** | Pointer to **NullableString** |  | [optional] 
**EnableVnc** | Pointer to **NullableString** |  | [optional] 
**HideHostSelection** | Pointer to **NullableString** |  | [optional] 
**ImportExisting** | Pointer to **NullableString** |  | [optional] 
**KubeUrl** | Pointer to **NullableString** |  | [optional] 
**NetworkServer** | Pointer to [**AddClouds200ResponseAllOfZoneConfigAnyOf1NetworkServer**](AddClouds200ResponseAllOfZoneConfigAnyOf1NetworkServer.md) |  | [optional] 
**NetworkServerId** | Pointer to **string** |  | [optional] 
**Password** | Pointer to **NullableString** |  | [optional] 
**PasswordHash** | Pointer to **NullableString** |  | [optional] 
**ReplicationMode** | Pointer to **NullableString** |  | [optional] 
**ResourcePool** | Pointer to **string** |  | [optional] 
**ResourcePoolId** | Pointer to **string** |  | [optional] 
**RpcMode** | Pointer to **string** |  | [optional] 
**SecurityMode** | Pointer to **string** |  | [optional] 
**SecurityServer** | Pointer to **NullableString** |  | [optional] 
**ServiceRegistryId** | Pointer to **NullableString** |  | [optional] 
**Username** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateClouds200ResponseAllOfZoneConfigAnyOf1{
    // Set fields directly
}
```

### ApplianceUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.ApplianceUrl.IsSet()` — check if set
- `obj.ApplianceUrl.Get()` — get the inner value (returns pointer)
- `obj.ApplianceUrl.Set(&val)` — set the value
- `obj.ApplianceUrl.Unset()` — clear the value
### DatacenterName (Nullable)

Use the Nullable wrapper methods:
- `obj.DatacenterName.IsSet()` — check if set
- `obj.DatacenterName.Get()` — get the inner value (returns pointer)
- `obj.DatacenterName.Set(&val)` — set the value
- `obj.DatacenterName.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### InventoryLevel (Nullable)

Use the Nullable wrapper methods:
- `obj.InventoryLevel.IsSet()` — check if set
- `obj.InventoryLevel.Get()` — get the inner value (returns pointer)
- `obj.InventoryLevel.Set(&val)` — set the value
- `obj.InventoryLevel.Unset()` — clear the value
### ConsoleKeymap (Nullable)

Use the Nullable wrapper methods:
- `obj.ConsoleKeymap.IsSet()` — check if set
- `obj.ConsoleKeymap.Get()` — get the inner value (returns pointer)
- `obj.ConsoleKeymap.Set(&val)` — set the value
- `obj.ConsoleKeymap.Unset()` — clear the value
### ApiVersion (Nullable)

Use the Nullable wrapper methods:
- `obj.ApiVersion.IsSet()` — check if set
- `obj.ApiVersion.Get()` — get the inner value (returns pointer)
- `obj.ApiVersion.Set(&val)` — set the value
- `obj.ApiVersion.Unset()` — clear the value
### BackupMode (Nullable)

Use the Nullable wrapper methods:
- `obj.BackupMode.IsSet()` — check if set
- `obj.BackupMode.Get()` — get the inner value (returns pointer)
- `obj.BackupMode.Set(&val)` — set the value
- `obj.BackupMode.Unset()` — clear the value
### CertificateProvider (Nullable)

Use the Nullable wrapper methods:
- `obj.CertificateProvider.IsSet()` — check if set
- `obj.CertificateProvider.Get()` — get the inner value (returns pointer)
- `obj.CertificateProvider.Set(&val)` — set the value
- `obj.CertificateProvider.Unset()` — clear the value
### ConfigCmId (Nullable)

Use the Nullable wrapper methods:
- `obj.ConfigCmId.IsSet()` — check if set
- `obj.ConfigCmId.Get()` — get the inner value (returns pointer)
- `obj.ConfigCmId.Set(&val)` — set the value
- `obj.ConfigCmId.Unset()` — clear the value
### ConfigCmdbId (Nullable)

Use the Nullable wrapper methods:
- `obj.ConfigCmdbId.IsSet()` — check if set
- `obj.ConfigCmdbId.Get()` — get the inner value (returns pointer)
- `obj.ConfigCmdbId.Set(&val)` — set the value
- `obj.ConfigCmdbId.Unset()` — clear the value
### ConfigManagementId (Nullable)

Use the Nullable wrapper methods:
- `obj.ConfigManagementId.IsSet()` — check if set
- `obj.ConfigManagementId.Get()` — get the inner value (returns pointer)
- `obj.ConfigManagementId.Set(&val)` — set the value
- `obj.ConfigManagementId.Unset()` — clear the value
### DatacenterId (Nullable)

Use the Nullable wrapper methods:
- `obj.DatacenterId.IsSet()` — check if set
- `obj.DatacenterId.Get()` — get the inner value (returns pointer)
- `obj.DatacenterId.Set(&val)` — set the value
- `obj.DatacenterId.Unset()` — clear the value
### DiskStorageType (Nullable)

Use the Nullable wrapper methods:
- `obj.DiskStorageType.IsSet()` — check if set
- `obj.DiskStorageType.Get()` — get the inner value (returns pointer)
- `obj.DiskStorageType.Set(&val)` — set the value
- `obj.DiskStorageType.Unset()` — clear the value
### DistributedWorkerId (Nullable)

Use the Nullable wrapper methods:
- `obj.DistributedWorkerId.IsSet()` — check if set
- `obj.DistributedWorkerId.Get()` — get the inner value (returns pointer)
- `obj.DistributedWorkerId.Set(&val)` — set the value
- `obj.DistributedWorkerId.Unset()` — clear the value
### DnsIntegrationId (Nullable)

Use the Nullable wrapper methods:
- `obj.DnsIntegrationId.IsSet()` — check if set
- `obj.DnsIntegrationId.Get()` — get the inner value (returns pointer)
- `obj.DnsIntegrationId.Set(&val)` — set the value
- `obj.DnsIntegrationId.Unset()` — clear the value
### EnableDiskTypeSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableDiskTypeSelection.IsSet()` — check if set
- `obj.EnableDiskTypeSelection.Get()` — get the inner value (returns pointer)
- `obj.EnableDiskTypeSelection.Set(&val)` — set the value
- `obj.EnableDiskTypeSelection.Unset()` — clear the value
### EnableNetworkTypeSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableNetworkTypeSelection.IsSet()` — check if set
- `obj.EnableNetworkTypeSelection.Get()` — get the inner value (returns pointer)
- `obj.EnableNetworkTypeSelection.Set(&val)` — set the value
- `obj.EnableNetworkTypeSelection.Unset()` — clear the value
### EnableStorageTypeSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableStorageTypeSelection.IsSet()` — check if set
- `obj.EnableStorageTypeSelection.Get()` — get the inner value (returns pointer)
- `obj.EnableStorageTypeSelection.Set(&val)` — set the value
- `obj.EnableStorageTypeSelection.Unset()` — clear the value
### EnableVnc (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableVnc.IsSet()` — check if set
- `obj.EnableVnc.Get()` — get the inner value (returns pointer)
- `obj.EnableVnc.Set(&val)` — set the value
- `obj.EnableVnc.Unset()` — clear the value
### HideHostSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.HideHostSelection.IsSet()` — check if set
- `obj.HideHostSelection.Get()` — get the inner value (returns pointer)
- `obj.HideHostSelection.Set(&val)` — set the value
- `obj.HideHostSelection.Unset()` — clear the value
### ImportExisting (Nullable)

Use the Nullable wrapper methods:
- `obj.ImportExisting.IsSet()` — check if set
- `obj.ImportExisting.Get()` — get the inner value (returns pointer)
- `obj.ImportExisting.Set(&val)` — set the value
- `obj.ImportExisting.Unset()` — clear the value
### KubeUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.KubeUrl.IsSet()` — check if set
- `obj.KubeUrl.Get()` — get the inner value (returns pointer)
- `obj.KubeUrl.Set(&val)` — set the value
- `obj.KubeUrl.Unset()` — clear the value
### Password (Nullable)

Use the Nullable wrapper methods:
- `obj.Password.IsSet()` — check if set
- `obj.Password.Get()` — get the inner value (returns pointer)
- `obj.Password.Set(&val)` — set the value
- `obj.Password.Unset()` — clear the value
### PasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.PasswordHash.IsSet()` — check if set
- `obj.PasswordHash.Get()` — get the inner value (returns pointer)
- `obj.PasswordHash.Set(&val)` — set the value
- `obj.PasswordHash.Unset()` — clear the value
### ReplicationMode (Nullable)

Use the Nullable wrapper methods:
- `obj.ReplicationMode.IsSet()` — check if set
- `obj.ReplicationMode.Get()` — get the inner value (returns pointer)
- `obj.ReplicationMode.Set(&val)` — set the value
- `obj.ReplicationMode.Unset()` — clear the value
### SecurityServer (Nullable)

Use the Nullable wrapper methods:
- `obj.SecurityServer.IsSet()` — check if set
- `obj.SecurityServer.Get()` — get the inner value (returns pointer)
- `obj.SecurityServer.Set(&val)` — set the value
- `obj.SecurityServer.Unset()` — clear the value
### ServiceRegistryId (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceRegistryId.IsSet()` — check if set
- `obj.ServiceRegistryId.Get()` — get the inner value (returns pointer)
- `obj.ServiceRegistryId.Set(&val)` — set the value
- `obj.ServiceRegistryId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


