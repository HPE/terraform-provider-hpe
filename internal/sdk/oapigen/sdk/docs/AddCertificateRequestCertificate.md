# AddCertificateRequestCertificate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A unique name scoped to your account for the key | [optional] 
**Description** | Pointer to **string** | A description of the certificate | [optional] 
**CertFile** | Pointer to **string** | The contents of the certificate file | [optional] 
**KeyFile** | Pointer to **string** | The contents of the key file | [optional] 
**ChainFile** | Pointer to **string** | The contents of the root certificate file | [optional] 
**DomainName** | Pointer to **string** | The domain name this certificate is tied to | [optional] 
**Wildcard** | Pointer to **bool** | Whether or not this certificate is a wildcard cert | [optional] [default to false]
**Type** | Pointer to **string** | Certificate Type Code to create a certificate of a type other than the default &#39;internal&#39;. | [optional] [default to "internal"]
**IntegrationId** | Pointer to **int64** | ID of the Service (Trust Integration) to create the certificate with, if using a type other than &#39;internal&#39;. eg. Internal, NSXT or Venafi | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCertificateRequestCertificate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


