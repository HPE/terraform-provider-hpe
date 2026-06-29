# GenerateSupportBundleRequestContentsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Code** | Pointer to **string** | Support bundle content type code. Preferred identifier for this API. Use &#x60;GET /api/options/supportBundles.contentTypes&#x60; to see available values. If provided, the code must resolve to a valid support bundle content type or the request returns &#x60;400&#x60;. | [optional] 
**TypeId** | Pointer to **int64** | Alternative support bundle content type ID returned by &#x60;GET /api/options/supportBundles.contentTypes&#x60;. Use this instead of &#x60;code&#x60; when the numeric type ID is known. If provided, the type ID must resolve to a valid support bundle content type or the request returns &#x60;400&#x60;. | [optional] 
**ResourceId** | Pointer to **int64** | Resource ID for resource-backed content types. Use &#x60;GET /api/options/supportBundles.contentTypeResources?contentTypeCode&#x3D;...&#x60; to see eligible values for a selected content type. Omit or set to null for standalone content types. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GenerateSupportBundleRequestContentsInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


