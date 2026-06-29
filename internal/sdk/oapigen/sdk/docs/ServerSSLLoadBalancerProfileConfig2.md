# ServerSSLLoadBalancerProfileConfig2

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProfileType** | Pointer to **string** | Profile category. Automatically set from serviceType (ssl-profile). | [optional] 
**SslSuite** | Pointer to **string** | Predefined cipher suite. Use CUSTOM to specify supportedSslCiphers and supportedSslProtocols explicitly. | [optional] 
**SessionCache** | Pointer to **bool** | Whether SSL session caching is enabled. Defaults to true. | [optional] [default to true]
**SupportedSslCiphers** | Pointer to **[]string** | Supported SSL ciphers. Required when sslSuite is CUSTOM. | [optional] 
**SupportedSslProtocols** | Pointer to **[]string** | Supported SSL/TLS protocols. Required when sslSuite is CUSTOM. | [optional] 
**Tags** | Pointer to [**[]LoadBalancerProfileTag23**](LoadBalancerProfileTag23.md) | NSX-T tags applied to the profile. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ServerSSLLoadBalancerProfileConfig2{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


