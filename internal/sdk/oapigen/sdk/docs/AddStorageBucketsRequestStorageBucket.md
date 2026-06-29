# AddStorageBucketsRequestStorageBucket

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A unique name scoped to your account for the storage bucket | 
**ProviderType** | **string** | The type of storage bucket | 
**StorageServer** | Pointer to **int64** | Storage Server ID. This is required when creating buckets tied to a cloud&#39;s storage server such as Amazon S3 and Google Cloud Storage.  The API &#x60;/api/storage-servers?type&#x3D;amazons3&#x60; can be used to see which options are available.  | [optional] 
**DefaultBackupTarget** | Pointer to **bool** | Default Backup Target | [optional] [default to false]
**CopyToStore** | Pointer to **bool** | Archive Snapshots | [optional] 
**DefaultDeploymentTarget** | Pointer to **bool** | Default Deployment Target | [optional] [default to false]
**DefaultVirtualImageTarget** | Pointer to **bool** | Default Virtual Image Store | [optional] [default to false]
**RetentionPolicyType** | Pointer to **string** | Cleanup mode. &#x60;backup&#x60; - Move old files to a backup provider. &#x60;delete&#x60; - Delete old files. &#x60;none&#x60; - Keep all files. | [optional] [default to "none"]
**RetentionPolicyDays** | Pointer to **int64** | The number of days old a file must be before it is deleted. | [optional] 
**RetentionProvider** | Pointer to **string** | The backup Storage Bucket where old files are moved to. | [optional] 
**BucketName** | **string** | The name of the bucket. Only applies to &#x60;Amazon&#x60;, &#x60;Azure&#x60;, &#x60;CIFS&#x60;, &#x60;NFSv3&#x60;, &#x60;Openstack Swift&#x60;, and &#x60;Rackspace CDN&#x60;. | 
**CreateBucket** | Pointer to **bool** | Create the bucket if it does not exist. Only applies to &#x60;Amazon&#x60;, &#x60;Azure&#x60;, &#x60;Openstack Swift&#x60;, and &#x60;Rackspace CDN&#x60;. | [optional] [default to false]
**Config** | [**AddStorageBucketsRequestStorageBucketConfig**](AddStorageBucketsRequestStorageBucketConfig.md) |  | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddStorageBucketsRequestStorageBucket{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


