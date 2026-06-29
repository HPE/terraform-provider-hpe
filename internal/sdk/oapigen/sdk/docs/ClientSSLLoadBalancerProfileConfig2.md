# ClientSSLLoadBalancerProfileConfig2

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProfileType** | Pointer to **string** | Profile category. Automatically set from serviceType (ssl-profile). | [optional] 
**SslSuite** | Pointer to **string** | Predefined cipher suite. Use CUSTOM to specify supportedSslCiphers and supportedSslProtocols explicitly. | [optional] 
**SessionCache** | Pointer to **bool** | Whether SSL session caching is enabled. Defaults to true. | [optional] [default to true]
**SupportedSslCiphers** | Pointer to **[]string** | Supported SSL ciphers. Required when sslSuite is CUSTOM. | [optional] 
**SupportedSslProtocols** | Pointer to **[]string** | Supported SSL/TLS protocols. Required when sslSuite is CUSTOM. | [optional] 
**SessionCacheTimeout** | Pointer to **NullableInt64** | SSL session cache timeout in seconds. Defaults to 300. Range 1-86400. | [optional] [default to 300]
**PreferServerCipher** | Pointer to **bool** | Whether the server cipher order is preferred during SSL negotiation. Defaults to true. | [optional] [default to true]
**Tags** | Pointer to [**[]LoadBalancerProfileTag22**](LoadBalancerProfileTag22.md) | NSX-T tags applied to the profile. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClientSSLLoadBalancerProfileConfig2{
    // Set fields directly
}
```

### SessionCacheTimeout (Nullable)

Use the Nullable wrapper methods:
- `obj.SessionCacheTimeout.IsSet()` — check if set
- `obj.SessionCacheTimeout.Get()` — get the inner value (returns pointer)
- `obj.SessionCacheTimeout.Set(&val)` — set the value
- `obj.SessionCacheTimeout.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


