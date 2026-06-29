# UpdateIdentitySourceSubdomains200ResponseAllOfUserSourceAnyOf6Config

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Url** | Pointer to **string** |  | [optional] 
**LogoutUrl** | Pointer to **string** |  | [optional] 
**DoNotIncludeSAMLRequest** | Pointer to **bool** |  | [optional] 
**SAMLSignatureMode** | Pointer to **string** |  | [optional] 
**DoNotValidateSignature** | Pointer to **bool** |  | [optional] 
**DoNotValidateStatusCode** | Pointer to **bool** |  | [optional] 
**DoNotValidateDestination** | Pointer to **bool** |  | [optional] 
**DoNotValidateIssueInstants** | Pointer to **bool** |  | [optional] 
**DoNotValidateAssertions** | Pointer to **bool** |  | [optional] 
**GivenNameAttribute** | Pointer to **string** |  | [optional] 
**SurnameAttribute** | Pointer to **string** |  | [optional] 
**EmailAttribute** | Pointer to **string** |  | [optional] 
**RequiredAttributeValue** | Pointer to **string** |  | [optional] 
**RoleAttributeName** | Pointer to **string** |  | [optional] 
**AzureTenantId** | Pointer to **string** |  | [optional] 
**AzureAppId** | Pointer to **string** |  | [optional] 
**AzureAppSecret** | Pointer to **NullableString** |  | [optional] 
**RoleLinkAttributeName** | Pointer to **string** |  | [optional] 
**PublicKey** | Pointer to **string** |  | [optional] 
**AzureAppSecretHash** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateIdentitySourceSubdomains200ResponseAllOfUserSourceAnyOf6Config{
    // Set fields directly
}
```

### AzureAppSecret (Nullable)

Use the Nullable wrapper methods:
- `obj.AzureAppSecret.IsSet()` — check if set
- `obj.AzureAppSecret.Get()` — get the inner value (returns pointer)
- `obj.AzureAppSecret.Set(&val)` — set the value
- `obj.AzureAppSecret.Unset()` — clear the value
### AzureAppSecretHash (Nullable)

Use the Nullable wrapper methods:
- `obj.AzureAppSecretHash.IsSet()` — check if set
- `obj.AzureAppSecretHash.Get()` — get the inner value (returns pointer)
- `obj.AzureAppSecretHash.Set(&val)` — set the value
- `obj.AzureAppSecretHash.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


