# GetNetworkProxy200ResponseNetworkProxy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**ProxyHost** | Pointer to **string** |  | [optional] 
**ProxyPort** | Pointer to **int64** |  | [optional] 
**ProxyUser** | Pointer to **NullableString** |  | [optional] 
**ProxyPassword** | Pointer to **NullableString** |  | [optional] 
**ProxyDomain** | Pointer to **string** |  | [optional] 
**ProxyWorkstation** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Account** | Pointer to [**NullableGetNetworkProxy200ResponseNetworkProxyAccount**](GetNetworkProxy200ResponseNetworkProxyAccount.md) |  | [optional] 
**Owner** | Pointer to [**NullableGetNetworkProxy200ResponseNetworkProxyOwner**](GetNetworkProxy200ResponseNetworkProxyOwner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkProxy200ResponseNetworkProxy{
    // Set fields directly
}
```

### ProxyUser (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyUser.IsSet()` — check if set
- `obj.ProxyUser.Get()` — get the inner value (returns pointer)
- `obj.ProxyUser.Set(&val)` — set the value
- `obj.ProxyUser.Unset()` — clear the value
### ProxyPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyPassword.IsSet()` — check if set
- `obj.ProxyPassword.Get()` — get the inner value (returns pointer)
- `obj.ProxyPassword.Set(&val)` — set the value
- `obj.ProxyPassword.Unset()` — clear the value
### ProxyWorkstation (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyWorkstation.IsSet()` — check if set
- `obj.ProxyWorkstation.Get()` — get the inner value (returns pointer)
- `obj.ProxyWorkstation.Set(&val)` — set the value
- `obj.ProxyWorkstation.Unset()` — clear the value
### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
### Owner (Nullable)

Use the Nullable wrapper methods:
- `obj.Owner.IsSet()` — check if set
- `obj.Owner.Get()` — get the inner value (returns pointer)
- `obj.Owner.Set(&val)` — set the value
- `obj.Owner.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


