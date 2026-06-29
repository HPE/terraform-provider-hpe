# UpdateStorageBucketsRequestStorageBucketConfigOneOf

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessKey** | Pointer to **string** | Access Key | [optional] 
**SecretKey** | Pointer to **string** | Secret Key | [optional] 
**Region** | Pointer to **string** | Optional Amazon region if creating a new bucket | [optional] 
**Endpoint** | Pointer to **string** | Optional endpoint URL if pointing to an object store other than amazon that mimics the Amazon S3 APIs. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateStorageBucketsRequestStorageBucketConfigOneOf{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


