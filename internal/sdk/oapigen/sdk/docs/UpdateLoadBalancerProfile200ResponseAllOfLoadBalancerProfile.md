# UpdateLoadBalancerProfile200ResponseAllOfLoadBalancerProfile

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**LoadBalancer** | Pointer to [**UpdateLoadBalancerProfile200ResponseAllOfLoadBalancerProfileLoadBalancer**](UpdateLoadBalancerProfile200ResponseAllOfLoadBalancerProfileLoadBalancer.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**ServiceType** | Pointer to **string** |  | [optional] 
**ServiceTypeDisplay** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**ProxyType** | Pointer to **NullableString** |  | [optional] 
**RedirectRewrite** | Pointer to **NullableString** |  | [optional] 
**PersistenceType** | Pointer to **NullableString** |  | [optional] 
**SslEnabled** | Pointer to **NullableString** |  | [optional] 
**SslCert** | Pointer to **NullableString** |  | [optional] 
**AccountCertificate** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**RedirectUrl** | Pointer to **NullableString** |  | [optional] 
**InsertXforwardedFor** | Pointer to **bool** |  | [optional] 
**PersistenceCookieName** | Pointer to **NullableString** |  | [optional] 
**PersistenceExpiresIn** | Pointer to **NullableString** |  | [optional] 
**Editable** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to [**UpdateLoadBalancerProfile200ResponseAllOfLoadBalancerProfileConfig**](UpdateLoadBalancerProfile200ResponseAllOfLoadBalancerProfileConfig.md) |  | [optional] 
**CreatedBy** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateLoadBalancerProfile200ResponseAllOfLoadBalancerProfile{
    // Set fields directly
}
```

### ProxyType (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyType.IsSet()` — check if set
- `obj.ProxyType.Get()` — get the inner value (returns pointer)
- `obj.ProxyType.Set(&val)` — set the value
- `obj.ProxyType.Unset()` — clear the value
### RedirectRewrite (Nullable)

Use the Nullable wrapper methods:
- `obj.RedirectRewrite.IsSet()` — check if set
- `obj.RedirectRewrite.Get()` — get the inner value (returns pointer)
- `obj.RedirectRewrite.Set(&val)` — set the value
- `obj.RedirectRewrite.Unset()` — clear the value
### PersistenceType (Nullable)

Use the Nullable wrapper methods:
- `obj.PersistenceType.IsSet()` — check if set
- `obj.PersistenceType.Get()` — get the inner value (returns pointer)
- `obj.PersistenceType.Set(&val)` — set the value
- `obj.PersistenceType.Unset()` — clear the value
### SslEnabled (Nullable)

Use the Nullable wrapper methods:
- `obj.SslEnabled.IsSet()` — check if set
- `obj.SslEnabled.Get()` — get the inner value (returns pointer)
- `obj.SslEnabled.Set(&val)` — set the value
- `obj.SslEnabled.Unset()` — clear the value
### SslCert (Nullable)

Use the Nullable wrapper methods:
- `obj.SslCert.IsSet()` — check if set
- `obj.SslCert.Get()` — get the inner value (returns pointer)
- `obj.SslCert.Set(&val)` — set the value
- `obj.SslCert.Unset()` — clear the value
### AccountCertificate (Nullable)

Use the Nullable wrapper methods:
- `obj.AccountCertificate.IsSet()` — check if set
- `obj.AccountCertificate.Get()` — get the inner value (returns pointer)
- `obj.AccountCertificate.Set(&val)` — set the value
- `obj.AccountCertificate.Unset()` — clear the value
### RedirectUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.RedirectUrl.IsSet()` — check if set
- `obj.RedirectUrl.Get()` — get the inner value (returns pointer)
- `obj.RedirectUrl.Set(&val)` — set the value
- `obj.RedirectUrl.Unset()` — clear the value
### PersistenceCookieName (Nullable)

Use the Nullable wrapper methods:
- `obj.PersistenceCookieName.IsSet()` — check if set
- `obj.PersistenceCookieName.Get()` — get the inner value (returns pointer)
- `obj.PersistenceCookieName.Set(&val)` — set the value
- `obj.PersistenceCookieName.Unset()` — clear the value
### PersistenceExpiresIn (Nullable)

Use the Nullable wrapper methods:
- `obj.PersistenceExpiresIn.IsSet()` — check if set
- `obj.PersistenceExpiresIn.Get()` — get the inner value (returns pointer)
- `obj.PersistenceExpiresIn.Set(&val)` — set the value
- `obj.PersistenceExpiresIn.Unset()` — clear the value
### CreatedBy (Nullable)

Use the Nullable wrapper methods:
- `obj.CreatedBy.IsSet()` — check if set
- `obj.CreatedBy.Get()` — get the inner value (returns pointer)
- `obj.CreatedBy.Set(&val)` — set the value
- `obj.CreatedBy.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


