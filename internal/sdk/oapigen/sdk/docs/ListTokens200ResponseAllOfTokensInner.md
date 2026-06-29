# ListTokens200ResponseAllOfTokensInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **NullableString** |  | [optional] 
**ClientId** | Pointer to **string** |  | [optional] 
**Username** | Pointer to **string** | Username associated with the token. Sub-tenant users will have their username formatted as &#x60;subdomain\\username&#x60; with a prefix that is the tenant subdomain or id by default. | [optional] 
**Expiration** | Pointer to **NullableTime** | Expiration date of the token | [optional] 
**TokenType** | Pointer to **string** | Type of the token | [optional] 
**Scope** | Pointer to **string** | Authorized scope(s), separated by spaces. Either &#x60;write&#x60; or &#x60;write openid&#x60;. | [optional] 
**MaskedAccessToken** | Pointer to **string** | Masked Access Token, with all but the first 8 characters replaced by asterisks for security. | [optional] 
**MaskedRefreshToken** | Pointer to **string** | Masked Refresh Token, with all but the first 8 characters replaced by asterisks for security. | [optional] 
**DateCreated** | Pointer to **NullableTime** | Date the token was created | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListTokens200ResponseAllOfTokensInner{
    // Set fields directly
}
```

### Name (Nullable)

Use the Nullable wrapper methods:
- `obj.Name.IsSet()` — check if set
- `obj.Name.Get()` — get the inner value (returns pointer)
- `obj.Name.Set(&val)` — set the value
- `obj.Name.Unset()` — clear the value
### Expiration (Nullable)

Use the Nullable wrapper methods:
- `obj.Expiration.IsSet()` — check if set
- `obj.Expiration.Get()` — get the inner value (returns pointer)
- `obj.Expiration.Set(&val)` — set the value
- `obj.Expiration.Unset()` — clear the value
### DateCreated (Nullable)

Use the Nullable wrapper methods:
- `obj.DateCreated.IsSet()` — check if set
- `obj.DateCreated.Get()` — get the inner value (returns pointer)
- `obj.DateCreated.Set(&val)` — set the value
- `obj.DateCreated.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


