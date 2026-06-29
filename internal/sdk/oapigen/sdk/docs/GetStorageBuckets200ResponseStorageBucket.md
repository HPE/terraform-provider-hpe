# GetStorageBuckets200ResponseStorageBucket

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**ProviderType** | Pointer to **string** |  | [optional] 
**StorageServer** | Pointer to [**GetStorageBuckets200ResponseStorageBucketStorageServer**](GetStorageBuckets200ResponseStorageBucketStorageServer.md) |  | [optional] 
**Config** | Pointer to [**GetStorageBuckets200ResponseStorageBucketConfig**](GetStorageBuckets200ResponseStorageBucketConfig.md) |  | [optional] 
**BucketName** | Pointer to **string** |  | [optional] 
**ReadOnly** | Pointer to **bool** |  | [optional] 
**DefaultBackupTarget** | Pointer to **bool** |  | [optional] 
**DefaultDeploymentTarget** | Pointer to **bool** |  | [optional] 
**DefaultVirtualImageTarget** | Pointer to **bool** |  | [optional] 
**CopyToStore** | Pointer to **bool** |  | [optional] 
**RetentionPolicyType** | Pointer to **NullableString** |  | [optional] 
**RetentionPolicyDays** | Pointer to **NullableString** |  | [optional] 
**RetentionProvider** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetStorageBuckets200ResponseStorageBucket{
    // Set fields directly
}
```

### RetentionPolicyType (Nullable)

Use the Nullable wrapper methods:
- `obj.RetentionPolicyType.IsSet()` — check if set
- `obj.RetentionPolicyType.Get()` — get the inner value (returns pointer)
- `obj.RetentionPolicyType.Set(&val)` — set the value
- `obj.RetentionPolicyType.Unset()` — clear the value
### RetentionPolicyDays (Nullable)

Use the Nullable wrapper methods:
- `obj.RetentionPolicyDays.IsSet()` — check if set
- `obj.RetentionPolicyDays.Get()` — get the inner value (returns pointer)
- `obj.RetentionPolicyDays.Set(&val)` — set the value
- `obj.RetentionPolicyDays.Unset()` — clear the value
### RetentionProvider (Nullable)

Use the Nullable wrapper methods:
- `obj.RetentionProvider.IsSet()` — check if set
- `obj.RetentionProvider.Get()` — get the inner value (returns pointer)
- `obj.RetentionProvider.Set(&val)` — set the value
- `obj.RetentionProvider.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


