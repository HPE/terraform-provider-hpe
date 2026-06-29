# GetSupportBundle200ResponseSupportBundle

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** | The current status of the support bundle generation | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**StartedAt** | Pointer to **NullableTime** |  | [optional] 
**CompletedAt** | Pointer to **NullableTime** |  | [optional] 
**FilePath** | Pointer to **NullableString** | The file path / filename of the generated bundle archive | [optional] 
**ContentLength** | Pointer to **NullableInt64** | File size in bytes | [optional] 
**ContentType** | Pointer to **NullableString** | MIME type of the bundle file | [optional] 
**Checksum** | Pointer to **NullableString** |  | [optional] 
**StorageProvider** | Pointer to [**GetSupportBundle200ResponseSupportBundleStorageProvider**](GetSupportBundle200ResponseSupportBundleStorageProvider.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetSupportBundle200ResponseSupportBundle{
    // Set fields directly
}
```

### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### StartedAt (Nullable)

Use the Nullable wrapper methods:
- `obj.StartedAt.IsSet()` — check if set
- `obj.StartedAt.Get()` — get the inner value (returns pointer)
- `obj.StartedAt.Set(&val)` — set the value
- `obj.StartedAt.Unset()` — clear the value
### CompletedAt (Nullable)

Use the Nullable wrapper methods:
- `obj.CompletedAt.IsSet()` — check if set
- `obj.CompletedAt.Get()` — get the inner value (returns pointer)
- `obj.CompletedAt.Set(&val)` — set the value
- `obj.CompletedAt.Unset()` — clear the value
### FilePath (Nullable)

Use the Nullable wrapper methods:
- `obj.FilePath.IsSet()` — check if set
- `obj.FilePath.Get()` — get the inner value (returns pointer)
- `obj.FilePath.Set(&val)` — set the value
- `obj.FilePath.Unset()` — clear the value
### ContentLength (Nullable)

Use the Nullable wrapper methods:
- `obj.ContentLength.IsSet()` — check if set
- `obj.ContentLength.Get()` — get the inner value (returns pointer)
- `obj.ContentLength.Set(&val)` — set the value
- `obj.ContentLength.Unset()` — clear the value
### ContentType (Nullable)

Use the Nullable wrapper methods:
- `obj.ContentType.IsSet()` — check if set
- `obj.ContentType.Get()` — get the inner value (returns pointer)
- `obj.ContentType.Set(&val)` — set the value
- `obj.ContentType.Unset()` — clear the value
### Checksum (Nullable)

Use the Nullable wrapper methods:
- `obj.Checksum.IsSet()` — check if set
- `obj.Checksum.Get()` — get the inner value (returns pointer)
- `obj.Checksum.Set(&val)` — set the value
- `obj.Checksum.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


