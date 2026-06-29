# SpecTemplateCreateFile

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SourceType** | **string** | File Source i.e. local, repository, url. | [default to "local"]
**Content** | Pointer to **string** | File content, the template text. Only required when sourceType is &#x60;local&#x60;. | [optional] 
**ContentPath** | Pointer to **string** | Content Path, the repo file location or url. Required when sourceType is repository or url. | [optional] 
**ContentRef** | Pointer to **string** | Content Ref, the branch/tag. Only used when sourceType is repo. | [optional] 
**Repository** | Pointer to [**SpecTemplateCreateFileRepository**](SpecTemplateCreateFileRepository.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SpecTemplateCreateFile{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


