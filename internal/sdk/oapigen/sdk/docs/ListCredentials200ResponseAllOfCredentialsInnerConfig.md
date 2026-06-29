# ListCredentials200ResponseAllOfCredentialsInnerConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientSecret** | Pointer to **NullableString** |  | [optional] 
**ClientId** | Pointer to **NullableString** |  | [optional] 
**ClientAuth** | Pointer to **string** |  | [optional] 
**Scope** | Pointer to **string** |  | [optional] 
**GrantType** | Pointer to **string** |  | [optional] 
**AccessTokenUrl** | Pointer to **string** |  | [optional] 
**ClientSecretHash** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListCredentials200ResponseAllOfCredentialsInnerConfig{
    // Set fields directly
}
```

### ClientSecret (Nullable)

Use the Nullable wrapper methods:
- `obj.ClientSecret.IsSet()` — check if set
- `obj.ClientSecret.Get()` — get the inner value (returns pointer)
- `obj.ClientSecret.Set(&val)` — set the value
- `obj.ClientSecret.Unset()` — clear the value
### ClientId (Nullable)

Use the Nullable wrapper methods:
- `obj.ClientId.IsSet()` — check if set
- `obj.ClientId.Get()` — get the inner value (returns pointer)
- `obj.ClientId.Set(&val)` — set the value
- `obj.ClientId.Unset()` — clear the value
### ClientSecretHash (Nullable)

Use the Nullable wrapper methods:
- `obj.ClientSecretHash.IsSet()` — check if set
- `obj.ClientSecretHash.Get()` — get the inner value (returns pointer)
- `obj.ClientSecretHash.Set(&val)` — set the value
- `obj.ClientSecretHash.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


