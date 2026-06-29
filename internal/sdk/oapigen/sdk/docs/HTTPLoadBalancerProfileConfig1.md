# HTTPLoadBalancerProfileConfig1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProfileType** | Pointer to **string** | Profile category. Automatically set from serviceType (application-profile). | [optional] 
**HttpIdleTimeout** | Pointer to **NullableInt64** | Idle timeout in seconds before an inactive connection is closed. Defaults to 15. Range 1-5400. | [optional] [default to 15]
**RequestHeaderSize** | Pointer to **NullableInt64** | Maximum request header size in bytes. Defaults to 1024. Range 1-65536. | [optional] [default to 1024]
**ResponseHeaderSize** | Pointer to **NullableInt64** | Maximum response header size in bytes. Defaults to 4096. Range 1-65536. | [optional] [default to 4096]
**HttpsRedirect** | Pointer to **string** | HTTP to HTTPS redirection. Stored as \&quot;on\&quot; or \&quot;off\&quot;. | [optional] 
**RedirectAddress** | Pointer to **string** | Redirect address. Required when httpsRedirect is \&quot;off\&quot;. | [optional] 
**XForwardedFor** | Pointer to **string** | X-Forwarded-For header handling. | [optional] 
**RequestBodySize** | Pointer to **NullableInt64** | Maximum request body size in bytes. Range 1-2147483647. | [optional] 
**ResponseTimeout** | Pointer to **NullableInt64** | Response timeout in seconds. Defaults to 60. | [optional] [default to 60]
**NtlmAuthentication** | Pointer to **bool** | Whether NTLM authentication is enabled. | [optional] 
**Tags** | Pointer to [**[]LoadBalancerProfileTag8**](LoadBalancerProfileTag8.md) | NSX-T tags applied to the profile. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &HTTPLoadBalancerProfileConfig1{
    // Set fields directly
}
```

### HttpIdleTimeout (Nullable)

Use the Nullable wrapper methods:
- `obj.HttpIdleTimeout.IsSet()` — check if set
- `obj.HttpIdleTimeout.Get()` — get the inner value (returns pointer)
- `obj.HttpIdleTimeout.Set(&val)` — set the value
- `obj.HttpIdleTimeout.Unset()` — clear the value
### RequestHeaderSize (Nullable)

Use the Nullable wrapper methods:
- `obj.RequestHeaderSize.IsSet()` — check if set
- `obj.RequestHeaderSize.Get()` — get the inner value (returns pointer)
- `obj.RequestHeaderSize.Set(&val)` — set the value
- `obj.RequestHeaderSize.Unset()` — clear the value
### ResponseHeaderSize (Nullable)

Use the Nullable wrapper methods:
- `obj.ResponseHeaderSize.IsSet()` — check if set
- `obj.ResponseHeaderSize.Get()` — get the inner value (returns pointer)
- `obj.ResponseHeaderSize.Set(&val)` — set the value
- `obj.ResponseHeaderSize.Unset()` — clear the value
### RequestBodySize (Nullable)

Use the Nullable wrapper methods:
- `obj.RequestBodySize.IsSet()` — check if set
- `obj.RequestBodySize.Get()` — get the inner value (returns pointer)
- `obj.RequestBodySize.Set(&val)` — set the value
- `obj.RequestBodySize.Unset()` — clear the value
### ResponseTimeout (Nullable)

Use the Nullable wrapper methods:
- `obj.ResponseTimeout.IsSet()` — check if set
- `obj.ResponseTimeout.Get()` — get the inner value (returns pointer)
- `obj.ResponseTimeout.Set(&val)` — set the value
- `obj.ResponseTimeout.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


