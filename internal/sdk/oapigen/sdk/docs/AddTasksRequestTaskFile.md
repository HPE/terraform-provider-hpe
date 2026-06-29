# AddTasksRequestTaskFile

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SourceType** | **string** | File Source i.e. &#x60;local&#x60;, &#x60;repository&#x60;, &#x60;url&#x60;. Default is &#x60;local&#x60;. | [default to "local"]
**Content** | Pointer to **string** | File content, the script text. Only required when sourceType is &#x60;local&#x60;. | [optional] 
**ContentPath** | Pointer to **string** | Content Path, the repo file location or url. Required when sourceType is &#x60;repository&#x60; or &#x60;url&#x60;. | [optional] 
**ContentRef** | Pointer to **string** | Content Ref, the branch/tag. Only used when sourceType is &#x60;repository&#x60;. | [optional] 
**Repository** | Pointer to [**AddTasksRequestTaskFileRepository**](AddTasksRequestTaskFileRepository.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddTasksRequestTaskFile{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


