# AddClientRequestClient

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientId** | **string** | Client ID | 
**ClientSecret** | Pointer to **string** | Client Secret | [optional] 
**AccessTokenValiditySeconds** | **NullableInt32** | Length of time accessToken is valid in seconds. | 
**RefreshTokenValiditySeconds** | **NullableInt32** | Length of time refreshToken is valid in seconds. | 
**RedirectUris** | Pointer to **[]string** | List of Redirect URIs for use with the OpenID Authorization Code Flow | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddClientRequestClient{
    // Set fields directly
}
```

### AccessTokenValiditySeconds (Nullable)

Use the Nullable wrapper methods:
- `obj.AccessTokenValiditySeconds.IsSet()` — check if set
- `obj.AccessTokenValiditySeconds.Get()` — get the inner value (returns pointer)
- `obj.AccessTokenValiditySeconds.Set(&val)` — set the value
- `obj.AccessTokenValiditySeconds.Unset()` — clear the value
### RefreshTokenValiditySeconds (Nullable)

Use the Nullable wrapper methods:
- `obj.RefreshTokenValiditySeconds.IsSet()` — check if set
- `obj.RefreshTokenValiditySeconds.Get()` — get the inner value (returns pointer)
- `obj.RefreshTokenValiditySeconds.Set(&val)` — set the value
- `obj.RefreshTokenValiditySeconds.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


