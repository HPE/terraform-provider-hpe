# ListLogs200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Sort** | Pointer to [**ListLogs200ResponseAllOfSort**](ListLogs200ResponseAllOfSort.md) |  | [optional] 
**Offset** | Pointer to **int64** |  | [optional] 
**Start** | Pointer to **time.Time** |  | [optional] 
**End** | Pointer to **time.Time** |  | [optional] 
**Data** | Pointer to [**[]ListLogs200ResponseAllOfDataInner**](ListLogs200ResponseAllOfDataInner.md) |  | [optional] 
**Max** | Pointer to **int64** |  | [optional] 
**GrandTotal** | Pointer to **int64** |  | [optional] 
**Total** | Pointer to **int64** |  | [optional] 
**Success** | Pointer to **bool** |  | [optional] 
**Count** | Pointer to **int64** |  | [optional] 
**Meta** | Pointer to [**ListApprovals200ResponseAllOfMeta**](ListApprovals200ResponseAllOfMeta.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListLogs200Response{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


