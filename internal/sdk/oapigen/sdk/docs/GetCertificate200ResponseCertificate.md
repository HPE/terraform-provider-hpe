# GetCertificate200ResponseCertificate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**DomainName** | Pointer to **NullableString** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**IntegrationId** | Pointer to **NullableInt64** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Generated** | Pointer to **bool** |  | [optional] 
**Wildcard** | Pointer to **bool** |  | [optional] 
**SelfSigned** | Pointer to **bool** |  | [optional] 
**Type** | Pointer to [**GetCertificate200ResponseCertificateType**](GetCertificate200ResponseCertificateType.md) |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**CommonName** | Pointer to **NullableString** |  | [optional] 
**CertType** | Pointer to **NullableString** |  | [optional] 
**KeyFileMD5** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetCertificate200ResponseCertificate{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### DomainName (Nullable)

Use the Nullable wrapper methods:
- `obj.DomainName.IsSet()` — check if set
- `obj.DomainName.Get()` — get the inner value (returns pointer)
- `obj.DomainName.Set(&val)` — set the value
- `obj.DomainName.Unset()` — clear the value
### IntegrationId (Nullable)

Use the Nullable wrapper methods:
- `obj.IntegrationId.IsSet()` — check if set
- `obj.IntegrationId.Get()` — get the inner value (returns pointer)
- `obj.IntegrationId.Set(&val)` — set the value
- `obj.IntegrationId.Unset()` — clear the value
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### CommonName (Nullable)

Use the Nullable wrapper methods:
- `obj.CommonName.IsSet()` — check if set
- `obj.CommonName.Get()` — get the inner value (returns pointer)
- `obj.CommonName.Set(&val)` — set the value
- `obj.CommonName.Unset()` — clear the value
### CertType (Nullable)

Use the Nullable wrapper methods:
- `obj.CertType.IsSet()` — check if set
- `obj.CertType.Get()` — get the inner value (returns pointer)
- `obj.CertType.Set(&val)` — set the value
- `obj.CertType.Unset()` — clear the value
### KeyFileMD5 (Nullable)

Use the Nullable wrapper methods:
- `obj.KeyFileMD5.IsSet()` — check if set
- `obj.KeyFileMD5.Get()` — get the inner value (returns pointer)
- `obj.KeyFileMD5.Set(&val)` — set the value
- `obj.KeyFileMD5.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


