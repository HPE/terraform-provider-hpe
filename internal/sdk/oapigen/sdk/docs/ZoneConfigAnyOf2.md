# ZoneConfigAnyOf2

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Endpoint** | Pointer to **string** |  | [optional] 
**AccessKey** | Pointer to **string** |  | [optional] 
**SecretKey** | Pointer to **string** |  | [optional] 
**UseHostCredentials** | Pointer to **string** |  | [optional] 
**StsAssumeRole** | Pointer to **string** |  | [optional] 
**IsVpc** | Pointer to **string** |  | [optional] 
**Vpc** | Pointer to **string** |  | [optional] 
**ImageStoreId** | Pointer to **string** |  | [optional] 
**EbsEncryption** | Pointer to **string** |  | [optional] 
**CostingReport** | Pointer to **string** |  | [optional] 
**CostingFolder** | Pointer to **string** |  | [optional] 
**CostingBucket** | Pointer to **string** |  | [optional] 
**CostingBucketName** | Pointer to **string** |  | [optional] 
**CostingRegion** | Pointer to **string** |  | [optional] 
**CostingAccessKey** | Pointer to **string** |  | [optional] 
**CostingSecretKey** | Pointer to **NullableString** |  | [optional] 
**CostingReportName** | Pointer to **string** |  | [optional] 
**ApplianceUrl** | Pointer to **string** |  | [optional] 
**DatacenterName** | Pointer to **string** |  | [optional] 
**NetworkServerId** | Pointer to **string** |  | [optional] 
**NetworkServer** | Pointer to [**AddClouds200ResponseAllOfZoneConfigAnyOf2NetworkServer**](AddClouds200ResponseAllOfZoneConfigAnyOf2NetworkServer.md) |  | [optional] 
**SecurityServer** | Pointer to **NullableString** |  | [optional] 
**CertificateProvider** | Pointer to **string** |  | [optional] 
**BackupMode** | Pointer to **string** |  | [optional] 
**ReplicationMode** | Pointer to **string** |  | [optional] 
**DnsIntegrationId** | Pointer to **string** |  | [optional] 
**ServiceRegistryId** | Pointer to **string** |  | [optional] 
**ConfigManagementId** | Pointer to **string** |  | [optional] 
**ConfigCmdbDiscovery** | Pointer to **bool** |  | [optional] 
**SecretKeyHash** | Pointer to **string** |  | [optional] 
**CostingSecretKeyHash** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ZoneConfigAnyOf2{
    // Set fields directly
}
```

### CostingSecretKey (Nullable)

Use the Nullable wrapper methods:
- `obj.CostingSecretKey.IsSet()` — check if set
- `obj.CostingSecretKey.Get()` — get the inner value (returns pointer)
- `obj.CostingSecretKey.Set(&val)` — set the value
- `obj.CostingSecretKey.Unset()` — clear the value
### SecurityServer (Nullable)

Use the Nullable wrapper methods:
- `obj.SecurityServer.IsSet()` — check if set
- `obj.SecurityServer.Get()` — get the inner value (returns pointer)
- `obj.SecurityServer.Set(&val)` — set the value
- `obj.SecurityServer.Unset()` — clear the value
### CostingSecretKeyHash (Nullable)

Use the Nullable wrapper methods:
- `obj.CostingSecretKeyHash.IsSet()` — check if set
- `obj.CostingSecretKeyHash.Get()` — get the inner value (returns pointer)
- `obj.CostingSecretKeyHash.Set(&val)` — set the value
- `obj.CostingSecretKeyHash.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


