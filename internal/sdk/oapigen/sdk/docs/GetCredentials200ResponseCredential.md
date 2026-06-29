# GetCredentials200ResponseCredential

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Type** | Pointer to [**GetCredentials200ResponseCredentialType**](GetCredentials200ResponseCredentialType.md) |  | [optional] 
**Integration** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Username** | Pointer to **NullableString** |  | [optional] 
**Password** | Pointer to **NullableString** |  | [optional] 
**PasswordHash** | Pointer to **NullableString** |  | [optional] 
**AuthKey** | Pointer to [**GetCredentials200ResponseCredentialAuthKey**](GetCredentials200ResponseCredentialAuthKey.md) |  | [optional] 
**AuthPath** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**Scope** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**StatusDate** | Pointer to **NullableTime** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Account** | Pointer to [**NullableGetCredentials200ResponseCredentialAccount**](GetCredentials200ResponseCredentialAccount.md) |  | [optional] 
**User** | Pointer to [**GetCredentials200ResponseCredentialUser**](GetCredentials200ResponseCredentialUser.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Config** | Pointer to [**GetCredentials200ResponseCredentialConfig**](GetCredentials200ResponseCredentialConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetCredentials200ResponseCredential{
    // Set fields directly
}
```

### Integration (Nullable)

Use the Nullable wrapper methods:
- `obj.Integration.IsSet()` — check if set
- `obj.Integration.Get()` — get the inner value (returns pointer)
- `obj.Integration.Set(&val)` — set the value
- `obj.Integration.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Username (Nullable)

Use the Nullable wrapper methods:
- `obj.Username.IsSet()` — check if set
- `obj.Username.Get()` — get the inner value (returns pointer)
- `obj.Username.Set(&val)` — set the value
- `obj.Username.Unset()` — clear the value
### Password (Nullable)

Use the Nullable wrapper methods:
- `obj.Password.IsSet()` — check if set
- `obj.Password.Get()` — get the inner value (returns pointer)
- `obj.Password.Set(&val)` — set the value
- `obj.Password.Unset()` — clear the value
### PasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.PasswordHash.IsSet()` — check if set
- `obj.PasswordHash.Get()` — get the inner value (returns pointer)
- `obj.PasswordHash.Set(&val)` — set the value
- `obj.PasswordHash.Unset()` — clear the value
### AuthPath (Nullable)

Use the Nullable wrapper methods:
- `obj.AuthPath.IsSet()` — check if set
- `obj.AuthPath.Get()` — get the inner value (returns pointer)
- `obj.AuthPath.Set(&val)` — set the value
- `obj.AuthPath.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### RefType (Nullable)

Use the Nullable wrapper methods:
- `obj.RefType.IsSet()` — check if set
- `obj.RefType.Get()` — get the inner value (returns pointer)
- `obj.RefType.Set(&val)` — set the value
- `obj.RefType.Unset()` — clear the value
### RefId (Nullable)

Use the Nullable wrapper methods:
- `obj.RefId.IsSet()` — check if set
- `obj.RefId.Get()` — get the inner value (returns pointer)
- `obj.RefId.Set(&val)` — set the value
- `obj.RefId.Unset()` — clear the value
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### StatusDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusDate.IsSet()` — check if set
- `obj.StatusDate.Get()` — get the inner value (returns pointer)
- `obj.StatusDate.Set(&val)` — set the value
- `obj.StatusDate.Unset()` — clear the value
### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


