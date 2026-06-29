# GetLoadBalancerProfile200ResponseLoadBalancerProfileConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProfileType** | Pointer to **string** | Profile category. Automatically set from serviceType (ssl-profile). | [optional] 
**HttpIdleTimeout** | Pointer to **int64** | Idle timeout in seconds before an inactive connection is closed. Defaults to 15. Range 1-5400. | [optional] [default to 15]
**RequestHeaderSize** | Pointer to **int64** | Maximum request header size in bytes. Defaults to 1024. Range 1-65536. | [optional] [default to 1024]
**ResponseHeaderSize** | Pointer to **int64** | Maximum response header size in bytes. Defaults to 4096. Range 1-65536. | [optional] [default to 4096]
**HttpsRedirect** | Pointer to **string** | HTTP to HTTPS redirection. Stored as \&quot;on\&quot; or \&quot;off\&quot;. | [optional] 
**RedirectAddress** | Pointer to **string** | Redirect address. Required when httpsRedirect is \&quot;off\&quot;. | [optional] 
**XForwardedFor** | Pointer to **string** | X-Forwarded-For header handling. | [optional] 
**RequestBodySize** | Pointer to **int64** | Maximum request body size in bytes. Range 1-2147483647. | [optional] 
**ResponseTimeout** | Pointer to **int64** | Response timeout in seconds. Defaults to 60. | [optional] [default to 60]
**NtlmAuthentication** | Pointer to **bool** | Whether NTLM authentication is enabled. | [optional] 
**Tags** | Pointer to [**[]LoadBalancerProfileTag31**](LoadBalancerProfileTag31.md) | NSX-T tags applied to the profile. | [optional] 
**FastTcpIdleTimeout** | Pointer to **int64** | Idle timeout in seconds before an inactive connection is closed. Defaults to 1800. | [optional] [default to 1800]
**HaFlowMirroring** | Pointer to **bool** | Whether all flows to the active member are mirrored to the standby member. | [optional] 
**ConnectionCloseTimeout** | Pointer to **int64** | Timeout in seconds before a closed connection is removed. Defaults to 8. Range 1-60. | [optional] [default to 8]
**FastUdpIdleTimeout** | Pointer to **int64** | Idle timeout in seconds before an inactive flow is closed. Defaults to 300. | [optional] [default to 300]
**SharePersistence** | Pointer to **bool** | Whether persistence is shared across virtual servers that use this profile. | [optional] 
**CookieName** | Pointer to **string** | Cookie name. | [optional] 
**CookieFallback** | Pointer to **bool** | Whether to fall back to another pool member when the cookie-pointed member is down. Defaults to true. | [optional] [default to true]
**CookieGarbling** | Pointer to **bool** | Whether the cookie value is encrypted (garbled). Defaults to true. | [optional] [default to true]
**CookieMode** | Pointer to **string** | Cookie persistence mode. | [optional] 
**CookieDomain** | Pointer to **string** | Cookie domain. Applies only when cookieMode is INSERT. | [optional] 
**CookiePath** | Pointer to **string** | Cookie path. Applies only when cookieMode is INSERT. | [optional] 
**CookieType** | Pointer to **string** | Cookie time type. Required when cookieMode is INSERT. | [optional] 
**MaxIdleTime** | Pointer to **int64** | Maximum idle time in milliseconds the cookie persists. Applies only when cookieMode is INSERT. | [optional] 
**MaxCookieAge** | Pointer to **int64** | Maximum age in milliseconds the cookie persists. Applies only when cookieMode is INSERT. | [optional] 
**PurgeEntries** | Pointer to **bool** | Whether the oldest persistence entries are purged when the table is full. Defaults to true. | [optional] [default to true]
**HaPersistenceMirroring** | Pointer to **bool** | Whether persistence entries are synchronized to the standby node for high availability. | [optional] 
**PersistenceEntryTimeout** | Pointer to **int64** | Persistence entry timeout in seconds. Defaults to 300. | [optional] [default to 300]
**SslSuite** | Pointer to **string** | Predefined cipher suite. Use CUSTOM to specify supportedSslCiphers and supportedSslProtocols explicitly. | [optional] 
**SessionCache** | Pointer to **bool** | Whether SSL session caching is enabled. Defaults to true. | [optional] [default to true]
**SupportedSslCiphers** | Pointer to **[]string** | Supported SSL ciphers. Required when sslSuite is CUSTOM. | [optional] 
**SupportedSslProtocols** | Pointer to **[]string** | Supported SSL/TLS protocols. Required when sslSuite is CUSTOM. | [optional] 
**SessionCacheTimeout** | Pointer to **int64** | SSL session cache timeout in seconds. Defaults to 300. Range 1-86400. | [optional] [default to 300]
**PreferServerCipher** | Pointer to **bool** | Whether the server cipher order is preferred during SSL negotiation. Defaults to true. | [optional] [default to true]

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetLoadBalancerProfile200ResponseLoadBalancerProfileConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


