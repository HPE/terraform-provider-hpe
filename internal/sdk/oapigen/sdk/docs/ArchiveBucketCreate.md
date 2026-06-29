# ArchiveBucketCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the archive bucket. Must be globally unique. | [optional] 
**Description** | Pointer to **string** | A description for the archive bucket | [optional] 
**StorageProvider** | Pointer to [**ArchiveBucketCreateStorageProvider**](ArchiveBucketCreateStorageProvider.md) |  | [optional] 
**Visibility** | Pointer to **string** | Visibility - Set to public to allow all tenants | [optional] [default to "private"]
**IsPublic** | Pointer to **bool** | Public URL - Set to true to allow anonymous access | [optional] [default to false]
**Accounts** | Pointer to [**ArchiveBucketCreateAccounts**](ArchiveBucketCreateAccounts.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ArchiveBucketCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


