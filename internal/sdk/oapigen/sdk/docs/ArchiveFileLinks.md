# ArchiveFileLinks

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**SecretAccessKey** | Pointer to **string** |  | [optional] 
**ArchiveFile** | Pointer to [**ArchiveFileLinksArchiveFile**](ArchiveFileLinksArchiveFile.md) |  | [optional] 
**CreatedBy** | Pointer to [**ArchiveFileLinksCreatedBy**](ArchiveFileLinksCreatedBy.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**LastAccessDate** | Pointer to **NullableTime** |  | [optional] 
**ExpirationDate** | Pointer to **NullableTime** |  | [optional] 
**DownloadCount** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ArchiveFileLinks{
    // Set fields directly
}
```

### LastAccessDate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastAccessDate.IsSet()` — check if set
- `obj.LastAccessDate.Get()` — get the inner value (returns pointer)
- `obj.LastAccessDate.Set(&val)` — set the value
- `obj.LastAccessDate.Unset()` — clear the value
### ExpirationDate (Nullable)

Use the Nullable wrapper methods:
- `obj.ExpirationDate.IsSet()` — check if set
- `obj.ExpirationDate.Get()` — get the inner value (returns pointer)
- `obj.ExpirationDate.Set(&val)` — set the value
- `obj.ExpirationDate.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


