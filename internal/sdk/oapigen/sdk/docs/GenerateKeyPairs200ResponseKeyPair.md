# GenerateKeyPairs200ResponseKeyPair

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**PublicKey** | Pointer to **NullableString** |  | [optional] 
**HasPrivateKey** | Pointer to **bool** |  | [optional] 
**PrivateKeyHash** | Pointer to **NullableString** |  | [optional] 
**PrivateKey** | Pointer to **NullableString** | Only present in response to generate | [optional] 
**Fingerprint** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GenerateKeyPairs200ResponseKeyPair{
    // Set fields directly
}
```

### PublicKey (Nullable)

Use the Nullable wrapper methods:
- `obj.PublicKey.IsSet()` — check if set
- `obj.PublicKey.Get()` — get the inner value (returns pointer)
- `obj.PublicKey.Set(&val)` — set the value
- `obj.PublicKey.Unset()` — clear the value
### PrivateKeyHash (Nullable)

Use the Nullable wrapper methods:
- `obj.PrivateKeyHash.IsSet()` — check if set
- `obj.PrivateKeyHash.Get()` — get the inner value (returns pointer)
- `obj.PrivateKeyHash.Set(&val)` — set the value
- `obj.PrivateKeyHash.Unset()` — clear the value
### PrivateKey (Nullable)

Use the Nullable wrapper methods:
- `obj.PrivateKey.IsSet()` — check if set
- `obj.PrivateKey.Get()` — get the inner value (returns pointer)
- `obj.PrivateKey.Set(&val)` — set the value
- `obj.PrivateKey.Unset()` — clear the value
### Fingerprint (Nullable)

Use the Nullable wrapper methods:
- `obj.Fingerprint.IsSet()` — check if set
- `obj.Fingerprint.Get()` — get the inner value (returns pointer)
- `obj.Fingerprint.Set(&val)` — set the value
- `obj.Fingerprint.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


