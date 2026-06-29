# UpdateClouds200ResponseAllOfZoneConfigAnyOf3

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SubscriberId** | Pointer to **string** |  | [optional] 
**TenantId** | Pointer to **string** |  | [optional] 
**ClientId** | Pointer to **string** |  | [optional] 
**ClientSecret** | Pointer to **string** |  | [optional] 
**ResourceGroup** | Pointer to **string** |  | [optional] 
**ImportExisting** | Pointer to **string** |  | [optional] 
**InventoryLevel** | Pointer to **string** |  | [optional] 
**ApplianceUrl** | Pointer to **string** |  | [optional] 
**DatacenterName** | Pointer to **string** |  | [optional] 
**NetworkServerId** | Pointer to **string** |  | [optional] 
**NetworkServer** | Pointer to [**AddClouds200ResponseAllOfZoneConfigAnyOf3NetworkServer**](AddClouds200ResponseAllOfZoneConfigAnyOf3NetworkServer.md) |  | [optional] 
**SecurityMode** | Pointer to **string** |  | [optional] 
**CertificateProvider** | Pointer to **string** |  | [optional] 
**BackupMode** | Pointer to **string** |  | [optional] 
**ReplicationMode** | Pointer to **string** |  | [optional] 
**DnsIntegrationId** | Pointer to **string** |  | [optional] 
**ConfigManagementId** | Pointer to **string** |  | [optional] 
**ConfigCmdbId** | Pointer to **string** |  | [optional] 
**SecurityServer** | Pointer to **NullableString** |  | [optional] 
**AccountType** | Pointer to **string** |  | [optional] 
**ServiceRegistryId** | Pointer to **string** |  | [optional] 
**CloudType** | Pointer to **string** |  | [optional] 
**StorageAccount** | Pointer to **string** |  | [optional] 
**RpcMode** | Pointer to **string** |  | [optional] 
**DiskEncryption** | Pointer to **string** |  | [optional] 
**EncryptionSet** | Pointer to **string** |  | [optional] 
**CspTenantId** | Pointer to **string** |  | [optional] 
**CspClientId** | Pointer to **string** |  | [optional] 
**CspClientSecret** | Pointer to **NullableString** |  | [optional] 
**CspCustomer** | Pointer to **NullableString** |  | [optional] 
**ConfigCmdbDiscovery** | Pointer to **bool** |  | [optional] 
**AzureCostingMode** | Pointer to **string** |  | [optional] 
**ClientSecretHash** | Pointer to **string** |  | [optional] 
**CspClientSecretHash** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateClouds200ResponseAllOfZoneConfigAnyOf3{
    // Set fields directly
}
```

### SecurityServer (Nullable)

Use the Nullable wrapper methods:
- `obj.SecurityServer.IsSet()` — check if set
- `obj.SecurityServer.Get()` — get the inner value (returns pointer)
- `obj.SecurityServer.Set(&val)` — set the value
- `obj.SecurityServer.Unset()` — clear the value
### CspClientSecret (Nullable)

Use the Nullable wrapper methods:
- `obj.CspClientSecret.IsSet()` — check if set
- `obj.CspClientSecret.Get()` — get the inner value (returns pointer)
- `obj.CspClientSecret.Set(&val)` — set the value
- `obj.CspClientSecret.Unset()` — clear the value
### CspCustomer (Nullable)

Use the Nullable wrapper methods:
- `obj.CspCustomer.IsSet()` — check if set
- `obj.CspCustomer.Get()` — get the inner value (returns pointer)
- `obj.CspCustomer.Set(&val)` — set the value
- `obj.CspCustomer.Unset()` — clear the value
### CspClientSecretHash (Nullable)

Use the Nullable wrapper methods:
- `obj.CspClientSecretHash.IsSet()` — check if set
- `obj.CspClientSecretHash.Get()` — get the inner value (returns pointer)
- `obj.CspClientSecretHash.Set(&val)` — set the value
- `obj.CspClientSecretHash.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


