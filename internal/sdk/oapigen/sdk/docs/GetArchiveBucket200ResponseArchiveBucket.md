# GetArchiveBucket200ResponseArchiveBucket

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**StorageProvider** | Pointer to [**AddArchiveBucket200ResponseAllOfArchiveBucketStorageProvider**](AddArchiveBucket200ResponseAllOfArchiveBucketStorageProvider.md) |  | [optional] 
**Owner** | Pointer to [**AddArchiveBucket200ResponseAllOfArchiveBucketOwner**](AddArchiveBucket200ResponseAllOfArchiveBucketOwner.md) |  | [optional] 
**CreatedBy** | Pointer to [**AddArchiveBucket200ResponseAllOfArchiveBucketCreatedBy**](AddArchiveBucket200ResponseAllOfArchiveBucketCreatedBy.md) |  | [optional] 
**IsPublic** | Pointer to **bool** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**FilePath** | Pointer to **string** |  | [optional] 
**RawSize** | Pointer to **NullableInt64** |  | [optional] 
**FileCount** | Pointer to **int64** |  | [optional] 
**Accounts** | Pointer to [**[]AddArchiveBucket200ResponseAllOfArchiveBucketAccountsInner**](AddArchiveBucket200ResponseAllOfArchiveBucketAccountsInner.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetArchiveBucket200ResponseArchiveBucket{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### RawSize (Nullable)

Use the Nullable wrapper methods:
- `obj.RawSize.IsSet()` — check if set
- `obj.RawSize.Get()` — get the inner value (returns pointer)
- `obj.RawSize.Set(&val)` — set the value
- `obj.RawSize.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


