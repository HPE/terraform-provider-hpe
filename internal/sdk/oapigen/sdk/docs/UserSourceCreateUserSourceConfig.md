# UserSourceCreateUserSourceConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Url** | Pointer to **string** | Login Redirect URL | [optional] 
**BindingUsername** | Pointer to **string** | Binding Username | [optional] 
**BindingPassword** | Pointer to **string** | Binding Password | [optional] 
**RequiredGroup** | Pointer to **string** | Required Group | [optional] 
**OrganizationId** | Pointer to **bool** | Organization ID | [optional] [default to false]
**RequiredRole** | Pointer to **string** | Required Role | [optional] 
**Domain** | Pointer to **string** | Domain | [optional] 
**UseSSL** | Pointer to **string** | Use SSL (\&quot;on\&quot; or \&quot;off\&quot;) | [optional] 
**SearchMemberGroups** | Pointer to **bool** | Include Member Groups | [optional] [default to false]
**AdministratorAPIToken** | Pointer to **string** | Administrator API Token | [optional] 
**Subdomain** | Pointer to **string** | OneLogin Subdomain | [optional] 
**Region** | Pointer to **string** | OneLogin Region | [optional] 
**ClientSecret** | Pointer to **string** | API Client Secret | [optional] 
**ClientId** | Pointer to **string** | API Client ID | [optional] 
**DoNotIncludeSAMLRequest** | Pointer to **bool** | Do not include SAMLRequest | [optional] [default to false]
**LogoutUrl** | Pointer to **string** | Lougout Post URL | [optional] 
**SAMLSignatureMode** | Pointer to **string** | SAML Request signing mode | [optional] [default to "NoSignature"]
**Request509Certificate** | Pointer to **string** | Only applies when &#x60;SAMLSignatureMode&#x3D;CustomSignature&#x60; | [optional] 
**RequestPrivateKey** | Pointer to **string** | RSA Private Key. Only applies when &#x60;SAMLSignatureMode&#x3D;CustomSignature&#x60; | [optional] 
**DoNotValidateSignature** | Pointer to **bool** | SAML Response Signing Flag | [optional] [default to true]
**PublicKey** | Pointer to **string** | Signing Public Key. Only applies when &#x60;doNotValidateSignature&#x3D;true&#x60; | [optional] 
**PrivateKey** | Pointer to **string** | Encryption RSA Private Key | [optional] 
**GivenNameAttribute** | Pointer to **string** | Given Name Attribute Name | [optional] 
**SurnameAttribute** | Pointer to **string** | Surname Attribute Name | [optional] 
**RoleAttributeName** | Pointer to **string** | Role Attribute Name | [optional] 
**RequiredAttributeValue** | Pointer to **string** | Role Attibute Required Value | [optional] 
**LoginUrl** | Pointer to **string** | External Login URL | [optional] 
**Logout** | Pointer to **string** | External Logout URL | [optional] 
**EncryptionAlgo** | Pointer to **string** | Encryption Algorithm | [optional] 
**EncryptionKey** | Pointer to **string** | Encryption Key | [optional] 
**Endpoint** | Pointer to **string** | API Endpoint | [optional] 
**ApiStyle** | Pointer to **string** | API Style | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UserSourceCreateUserSourceConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


