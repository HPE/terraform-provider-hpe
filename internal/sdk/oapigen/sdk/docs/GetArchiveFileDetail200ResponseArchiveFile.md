# GetArchiveFileDetail200ResponseArchiveFile

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**FilePath** | Pointer to **string** |  | [optional] 
**ArchiveBucket** | Pointer to [**AddArchiveFile200ResponseAllOfArchiveFileArchiveBucket**](AddArchiveFile200ResponseAllOfArchiveFileArchiveBucket.md) |  | [optional] 
**CreatedBy** | Pointer to [**AddArchiveFile200ResponseAllOfArchiveFileCreatedBy**](AddArchiveFile200ResponseAllOfArchiveFileCreatedBy.md) |  | [optional] 
**IsDirectory** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**RawSize** | Pointer to **int64** |  | [optional] 
**ContentType** | Pointer to **NullableString** |  | [optional] 
**DownloadCount** | Pointer to **int64** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetArchiveFileDetail200ResponseArchiveFile{
    // Set fields directly
}
```

### ContentType (Nullable)

Use the Nullable wrapper methods:
- `obj.ContentType.IsSet()` — check if set
- `obj.ContentType.Get()` — get the inner value (returns pointer)
- `obj.ContentType.Set(&val)` — set the value
- `obj.ContentType.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


