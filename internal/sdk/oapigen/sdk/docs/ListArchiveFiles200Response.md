# ListArchiveFiles200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ArchiveBucket** | Pointer to [**ListArchiveFiles200ResponseAllOfArchiveBucket**](ListArchiveFiles200ResponseAllOfArchiveBucket.md) |  | [optional] 
**IsOwner** | Pointer to **bool** |  | [optional] 
**ParentDirectory** | Pointer to **NullableString** |  | [optional] 
**ArchiveFiles** | Pointer to [**[]ListArchiveFiles200ResponseAllOfArchiveFilesInner**](ListArchiveFiles200ResponseAllOfArchiveFilesInner.md) |  | [optional] 
**Meta** | Pointer to [**ListApprovals200ResponseAllOfMeta**](ListApprovals200ResponseAllOfMeta.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListArchiveFiles200Response{
    // Set fields directly
}
```

### ParentDirectory (Nullable)

Use the Nullable wrapper methods:
- `obj.ParentDirectory.IsSet()` — check if set
- `obj.ParentDirectory.Get()` — get the inner value (returns pointer)
- `obj.ParentDirectory.Set(&val)` — set the value
- `obj.ParentDirectory.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


