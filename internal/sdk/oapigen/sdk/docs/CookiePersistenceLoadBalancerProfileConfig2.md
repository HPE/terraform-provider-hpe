# CookiePersistenceLoadBalancerProfileConfig2

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProfileType** | Pointer to **string** | Profile category. Automatically set from serviceType (persistence-profile). | [optional] 
**SharePersistence** | Pointer to **bool** | Whether persistence is shared across virtual servers that use this profile. | [optional] 
**CookieName** | Pointer to **string** | Cookie name. | [optional] 
**CookieFallback** | Pointer to **bool** | Whether to fall back to another pool member when the cookie-pointed member is down. Defaults to true. | [optional] [default to true]
**CookieGarbling** | Pointer to **bool** | Whether the cookie value is encrypted (garbled). Defaults to true. | [optional] [default to true]
**CookieMode** | Pointer to **string** | Cookie persistence mode. | [optional] 
**CookieDomain** | Pointer to **string** | Cookie domain. Applies only when cookieMode is INSERT. | [optional] 
**CookiePath** | Pointer to **string** | Cookie path. Applies only when cookieMode is INSERT. | [optional] 
**CookieType** | Pointer to **string** | Cookie time type. Required when cookieMode is INSERT. | [optional] 
**MaxIdleTime** | Pointer to **NullableInt64** | Maximum idle time in milliseconds the cookie persists. Applies only when cookieMode is INSERT. | [optional] 
**MaxCookieAge** | Pointer to **NullableInt64** | Maximum age in milliseconds the cookie persists. Applies only when cookieMode is INSERT. | [optional] 
**Tags** | Pointer to [**[]LoadBalancerProfileTag19**](LoadBalancerProfileTag19.md) | NSX-T tags applied to the profile. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CookiePersistenceLoadBalancerProfileConfig2{
    // Set fields directly
}
```

### MaxIdleTime (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxIdleTime.IsSet()` — check if set
- `obj.MaxIdleTime.Get()` — get the inner value (returns pointer)
- `obj.MaxIdleTime.Set(&val)` — set the value
- `obj.MaxIdleTime.Unset()` — clear the value
### MaxCookieAge (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxCookieAge.IsSet()` — check if set
- `obj.MaxCookieAge.Get()` — get the inner value (returns pointer)
- `obj.MaxCookieAge.Set(&val)` — set the value
- `obj.MaxCookieAge.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


