# GenerateSupportBundleRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StorageProviderId** | Pointer to **int64** | ID of the storage provider (storage bucket) where the bundle should be saved. Defaults to the storage bucket configured as the default support bundle target. If no default is configured, the bundle is saved to the appliance&#39;s local storage. If provided, this ID must reference an existing storage provider or the request returns &#x60;404&#x60;. | [optional] 
**StartDate** | **time.Time** | ISO 8601 start of the log collection window (e.g. &#x60;2026-01-15T00:00:00Z&#x60;). Required. Appliance logs are collected starting from this time. | 
**EndDate** | Pointer to **time.Time** | ISO 8601 end of the log collection window (e.g. &#x60;2026-01-15T23:59:59Z&#x60;). Defaults to the time the request is processed when omitted. | [optional] 
**Contents** | Pointer to [**[]GenerateSupportBundleRequestContentsInner**](GenerateSupportBundleRequestContentsInner.md) | Flat list of support bundle content entries to include. Resource-backed entries should include &#x60;resourceId&#x60;. Standalone entries omit it. If omitted or empty, all eligible content entries are included. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GenerateSupportBundleRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


